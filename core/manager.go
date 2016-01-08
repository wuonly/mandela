package core

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	addrm "github.com/prestonTao/mandela/core/addr_manager"
	"github.com/prestonTao/mandela/core/config"
	msg "github.com/prestonTao/mandela/core/message_center"
	engine "github.com/prestonTao/mandela/core/net"
	"github.com/prestonTao/mandela/core/nodeStore"
	"github.com/prestonTao/mandela/core/utils"
)

var (
	privateKey  *rsa.PrivateKey
	isStartCore = false
)

func init() {
	go startUp()
}

func startUp() {
	one := make(chan string, 0)
	addrm.AddSubscribe(one)
	for {
		//接收到超级节点地址消息
		addr := <-one
		utils.Log.Debug("有新的地址")
		host, portStr, _ := net.SplitHostPort(addr)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			// return "", 0, errors.New("IP地址解析失败")

			continue
		}
		if !isStartCore {
			connectNet(host, port)
		}
	}
}

func StartService() {
	//启动核心组件
	StartUpCore()
	//开启web服务
	// go StartWeb()
}

/*
	启动核心组件
*/
func StartUpCore() {
	if len(Init_IdInfo.Id) == 0 {
		GetId()
		if len(Init_IdInfo.Id) == 0 {
			return
		}
	}
	utils.Log.Debug("启动服务器核心组件")
	utils.Log.Debug("本机id为：\n%s", Init_IdInfo.GetId())

	isSuperPeer := config.CheckIsSuperPeer()
	//是超级节点
	node := &nodeStore.Node{
		IdInfo:  Init_IdInfo,
		IsSuper: isSuperPeer, //是否是超级节点
		UdpPort: 0,
	}
	addr, port := config.GetHost()
	node.Addr = addr
	node.TcpPort = int32(port)

	/*
		启动消息服务器
	*/
	engine.InitEngine(string(Init_IdInfo.Build()))
	/*
		生成密钥文件
	*/
	var err error
	//生成密钥
	privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		fmt.Println("生成密钥错误", err.Error())
		return
	}
	/*
		启动分布式哈希表
	*/
	nodeStore.InitNodeStore(node)
	/*
		设置关闭连接回调函数后监听
	*/
	engine.SetAuth(new(Auth))
	engine.SetCloseCallback(closeConnCallback)
	engine.Listen(config.TCPListener)
	//自己是超级节点就把自己添加到超级节点地址列表中去
	if isSuperPeer {
		addrm.AddSuperPeerAddr(addr + ":" + strconv.Itoa(port))
	}

	isStartCore = true

	/*
		连接到超级节点
	*/
	ip, port, err := addrm.GetSuperAddrOne(false)
	if err == nil {
		connectNet(ip, port)
	}

	go read()

}

/*
	链接到网络中去
*/
func connectNet(ip string, port int) {
	if !isStartCore {
		StartUpCore()
		//启动失败
		if !isStartCore {
			return
		}
	}
	utils.Log.Debug("链接到网络中去")

	nodeStore.SuperName = engine.AddClientConn(ip, int32(port), false)
	utils.Log.Debug("超级节点为: %s", nodeStore.SuperName)
	// config.SuperNodeIp = ip
	// config.SuperNodePort = port
	//给目标机器发送自己的名片
	introduceSelf()

}

/*
	关闭服务器回调函数
*/
func shutdownCallback() {
	//回收映射的端口
	config.Reclaim()
	// addrm.CloseBroadcastServer()
	fmt.Println("Close over")
}

/*
	连接超级节点后，向超级节点介绍自己
	第一次连接超级节点，用代理方式查找离自己最近的节点
*/
func introduceSelf() {
	session, _ := engine.GetController().GetSession(nodeStore.SuperName)
	//用代理方式查找最近的超级节点
	nodeMsg := msg.FindNode{
		NodeId:  session.GetName(),
		WantId:  nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
		IsProxy: true,
		ProxyId: nodeStore.GetRootIdInfoString(),
		IsSuper: config.CheckIsSuperPeer(),
		Addr:    nodeStore.Root.Addr,
		TcpPort: nodeStore.Root.TcpPort,
		UdpPort: nodeStore.Root.UdpPort,
	}

	// resultBytes, _ := proto.Marshal(&nodeMsg)
	resultBytes, _ := json.Marshal(nodeMsg)

	session.Send(msg.FindNodeNum, &resultBytes)
}

/*
	一个连接断开后的回调方法
*/
func closeConnCallback(name string) {
	// utils.Log.Debug("节点下线 %s", nodeStore.ParseId(name))
	node := nodeStore.Get(nodeStore.ParseId(name), false, "")
	// fmt.Println("节点下线", node)

	// utils.Log.Debug("目前超级节点是 %s", nodeStore.ParseId(nodeStore.SuperName))

	if name == nodeStore.SuperName {
		fmt.Println("超级节点断开连接:", name)
		targetNode := nodeStore.Get(nodeStore.Root.IdInfo.GetId(), false, nodeStore.Root.IdInfo.GetId())
		if targetNode == nil {
			return
		}

		if config.Init_role == config.C_role_client {
			nodeStore.SuperName = engine.AddClientConn(targetNode.Addr, targetNode.TcpPort, false)
		} else {
			session, _ := engine.GetController().GetSession(string(targetNode.IdInfo.Build()))
			nodeStore.SuperName = session.GetName()
		}
		return
	}
	// session, ok := engine.GetController().GetSession(name)
	// if err != nil {
	// 	fmt.Println("客户端离线，但找不到这个session")
	// }

	if node != nil && !node.IsSuper {
		fmt.Println("自己代理的节点下线:", nodeStore.ParseId(name))
	}
	nodeStore.DelNode(nodeStore.ParseId(name))
}

//处理查找节点的请求
//本节点定期查询已知节点是否在线，更新节点信息
func read() {
	for {
		nodeIdStr := <-nodeStore.OutFindNode
		session, ok := engine.GetController().GetSession(nodeStore.SuperName)
		//root节点刚启动就没有超级节点
		if !ok {
			continue
		}
		findNodeOne := &msg.FindNode{
			NodeId:  nodeStore.GetRootIdInfoString(),
			IsProxy: false,
			ProxyId: nodeStore.GetRootIdInfoString(),
			WantId:  nodeIdStr,
		}
		/*
			当查找id等于自己的时候：
			超级节点：查找邻居节点
			普通节点：查找离自己最近的超级节点，查找邻居节点做备用超级节点
		*/
		if nodeIdStr == nodeStore.ParseId(nodeStore.GetRootIdInfoString()) {
			//普通节点查找最近的超级节点
			if !nodeStore.Root.IsSuper {
				findNodeOne.NodeId = session.GetName()
				findNodeOne.IsProxy = true
				// findNodeOne.WantId = nodeIdStr
				findNodeOne.IsSuper = nodeStore.Root.IsSuper
				findNodeOne.Addr = nodeStore.Root.Addr
				findNodeOne.TcpPort = nodeStore.Root.TcpPort
				findNodeOne.UdpPort = nodeStore.Root.UdpPort

				// fmt.Println("1---------", findNodeOne)
				resultBytes, _ := json.Marshal(findNodeOne)
				session.Send(msg.FindNodeNum, &resultBytes)

				findNodeOne.WantId = "left"
				findNodeBytes, _ := json.Marshal(findNodeOne)
				err := session.Send(msg.FindNodeNum, &findNodeBytes)
				if err != nil {
					fmt.Println("manager发送数据出错：", err.Error())
				}

				findNodeOne.WantId = "right"
				findNodeBytes, _ = json.Marshal(findNodeOne)
				err = session.Send(msg.FindNodeNum, &findNodeBytes)
				if err != nil {
					fmt.Println("manager发送数据出错：", err.Error())
				}
				continue
			}

			//先发送左邻居节点查找请求
			findNodeOne.WantId = "left"
			id := nodeStore.GetLeftNode(*nodeStore.Root.IdInfo.GetBigIntId(), 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ := json.Marshal(findNodeOne)
			ok := false
			var clientConn engine.Session
			if nodeStore.Root.IsSuper {
				clientConn, ok = engine.GetController().GetSession(string(id[0].IdInfo.Build()))
			} else {
				clientConn, ok = engine.GetController().GetSession(nodeStore.SuperName)
			}
			if !ok {
				continue
			}
			err := clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			//发送右邻居节点查找请求
			findNodeOne.WantId = "right"
			id = nodeStore.GetRightNode(*nodeStore.Root.IdInfo.GetBigIntId(), 1)
			if id == nil || len(id) == 0 {
				continue
			}
			findNodeBytes, _ = json.Marshal(findNodeOne)
			if nodeStore.Root.IsSuper {
				if clientConn, ok = engine.GetController().GetSession(string(id[0].IdInfo.Build())); !ok {
					continue
				}
			}
			err = clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			continue
		}
		//自己不是超级节点，就不需要保存逻辑节点
		// if !nodeStore.Root.IsSuper {
		// 	continue
		// }

		//--------------------------------------------
		//    查找普通节点，只有超级节点才需要查找
		//--------------------------------------------
		//这里临时加上去
		//去掉后有性能问题

		findNodeBytes, _ := json.Marshal(findNodeOne)

		remote := nodeStore.Get(nodeIdStr, false, "")
		if remote == nil {
			continue
		}
		session, _ = engine.GetController().GetSession(string(remote.IdInfo.Build()))
		if session == nil {
			continue
		}
		err := session.Send(msg.FindNodeNum, &findNodeBytes)
		if err != nil {
			fmt.Println("manager发送数据出错：", err.Error())
		}
	}
}
