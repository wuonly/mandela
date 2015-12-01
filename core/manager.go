package core

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	addrm "github.com/prestonTao/mandela/core/addr_manager"
	"github.com/prestonTao/mandela/core/config"
	msg "github.com/prestonTao/mandela/core/message_center"
	engine "github.com/prestonTao/mandela/core/net"
	"github.com/prestonTao/mandela/core/nodeStore"
	"github.com/prestonTao/mandela/core/utils"
	// "github.com/prestonTao/upnp"
	"net"
	"strconv"
	"time"
)

var (
	privateKey *rsa.PrivateKey
)

func init() {
	utils.GlobalInit("console", "", "debug", 1)
	// utils.GlobalInit("file", `{"filename":"/var/log/gd/gd.log"}`, "", 1000)
	// utils.Log.Debug("session handle receive, %d, %v", msg.Code(), msg.Content())
	utils.Log.Debug("test debug")
	utils.Log.Warn("test warn")
	utils.Log.Error("test error")
}

/*
	根据网络情况自己确定节点角色
*/
func AutoRole() {
	//尝试端口映射
	if !config.Mode_local {
		portMapping()
	}
	//得到本地ip地址
	if config.Mode_local {
		config.Init_LocalIP = utils.GetLocalHost()
	} else {
		config.Init_LocalIP = utils.GetLocalIntenetIp()
	}
	//得到本机可用端口
	config.Init_LocalPort = utils.GetAvailablePortForTCP()
	if config.Mode_local && config.Init_role == config.C_role_super {
		config.Init_IsSuperPeer = true
		config.Init_GlobalUnicastAddress = config.Init_LocalIP
		config.Init_GlobalUnicastAddress_port = config.Init_LocalPort
		config.Init_ExternalIP = config.Init_LocalIP
		config.Init_MappingPort = config.Init_LocalPort
	}
	//自己是根节点
	if config.Init_role == config.C_role_root {
		config.Init_IsSuperPeer = true
	}
	if config.Init_role == config.C_role_auto {

	}
	utils.Log.Debug("本机角色为：%s", config.Init_role)
	utils.Log.Debug("本机监听地址：%s:%d", config.Init_LocalIP, config.Init_LocalPort)
}

/*
	根据网络环境启动程序
*/
func StartUpAuto() {

	if config.Mode_local {
		utils.Log.Debug("局域网模式")
		config.Init_LocalIP = "127.0.0.1"
	}
	AutoRole()
	addrm.LoadIdInfo()
	addrm.InitSuperPeer()

	if config.Init_role != config.C_role_root {
		addrm.StartLoadSuperPeer()
	}
	//开启web服务
	// go StartWeb()

	//没有idinfo的新节点
	if len(addrm.Init_IdInfo.Id) == 0 {
		//连接网络并得到一个idinfo
		newnode := nodeStore.IdInfo{
			Id:          addrm.Str_zaro,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05.999999999"),
			Domain:      utils.GetRandomDomain(),
			Name:        "",
			Email:       "",
			SuperNodeId: addrm.Str_zaro,
		}
		idInfo, err := addrm.GetId(newnode)
		if err == nil {
			addrm.Init_IdInfo = *idInfo
			// saveIdInfo(Path_Id)
		} else {
			fmt.Println("从网络中获得idinfo失败")
			return
		}
	}
	//是超级节点
	var node *nodeStore.Node
	node = &nodeStore.Node{
		IdInfo:  addrm.Init_IdInfo,
		IsSuper: config.Init_IsSuperPeer, //是否是超级节点
		UdpPort: 0,
	}
	if config.Init_role == config.C_role_root {

		node.Addr = config.Init_GlobalUnicastAddress
		node.TcpPort = int32(config.Init_GlobalUnicastAddress_port)
	} else if config.Init_IsSuperPeer {
		node.Addr = config.Init_ExternalIP
		node.TcpPort = int32(config.Init_MappingPort)
	} else {
		node.Addr = config.Init_LocalIP
		node.TcpPort = int32(config.Init_LocalPort)
	}

	startUp(node)
	if config.Init_role == config.C_role_root {
		// StartRootPeer()
		addrm.StartLoadSuperPeer()
	}
}

/*
	开始启动服务器
*/
func StartUp() {
	//尝试端口映射
	portMapping()
	//是超级节点
	var node *nodeStore.Node
	if config.Init_IsSuperPeer || config.Init_role == config.C_role_root {
		node = &nodeStore.Node{
			IdInfo:  addrm.Init_IdInfo,
			IsSuper: config.Init_IsSuperPeer, //是否是超级节点
			Addr:    config.Init_GlobalUnicastAddress,
			TcpPort: int32(config.Init_GlobalUnicastAddress_port),
			UdpPort: 0,
		}
	} else {
		node = &nodeStore.Node{
			IdInfo:  addrm.Init_IdInfo,
			IsSuper: config.Init_IsSuperPeer, //是否是超级节点
			Addr:    config.Init_LocalIP,
			TcpPort: int32(config.Init_LocalPort),
			UdpPort: 0,
		}
	}
	startUp(node)
	if config.Init_role == config.C_role_root {
		// StartRootPeer()
		addrm.StartLoadSuperPeer()
	}
}

/*
	判断自己是否有公网ip地址
	若支持upnp协议，则添加一个端口映射
*/
func portMapping() {
	// utils.Log.Debug("监听一个本地地址：%s:%d", Init_LocalIP, Init_LocalPort)

	// fmt.Println("监听一个本地地址：", Init_LocalIP, ":", Init_LocalPort)
	//本地地址是全球唯一公网地址
	if utils.IsOnlyIp(config.Init_LocalIP) {
		config.Init_IsSuperPeer = true
		config.Init_GlobalUnicastAddress = config.Init_LocalIP
		config.Init_GlobalUnicastAddress_port = config.Init_LocalPort
		// fmt.Println("本机ip是公网全球唯一地址")
		utils.Log.Debug("本机ip是公网全球唯一地址")
		return
	}
	//获得网关公网地址
	err := config.Sys_mapping.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
		utils.Log.Warn("网关不支持端口映射")
		return
	}
	config.Init_ExternalIP = config.Sys_mapping.GatewayOutsideIP
	config.Init_GlobalUnicastAddress = config.Init_ExternalIP
	utils.Log.Debug("正在尝试端口映射")
	for i := 0; i < 1000; i++ {
		if err := config.Sys_mapping.AddPortMapping(config.Init_LocalPort, config.Init_MappingPort, "TCP"); err == nil {
			config.Init_IsSuperPeer = true
			config.Init_GlobalUnicastAddress_port = config.Init_MappingPort
			// fmt.Println("映射到公网地址：", Init_ExternalIP, ":", Init_MappingPort)
			utils.Log.Debug("映射到公网地址：%s:%d", config.Init_ExternalIP, config.Init_MappingPort)
			return
		}
		config.Init_MappingPort = config.Init_MappingPort + 1
	}
	utils.Log.Warn("端口映射失败")
	// fmt.Println("端口映射失败")
}

func startUp(node *nodeStore.Node) {
	utils.Log.Debug("本机id为：%s", addrm.Init_IdInfo.GetId())
	/*
		启动消息服务器
	*/
	engine.InitEngine(string(addrm.Init_IdInfo.Build()))
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
	engine.Listen(config.Init_LocalIP, int32(config.Init_LocalPort))
	if config.Init_role != config.C_role_root {
		/*
			连接到超级节点
		*/
		host, portStr, _ := net.SplitHostPort(addrm.GetSuperAddrOne())
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return
		}
		nodeStore.SuperName = engine.AddClientConn(host, int32(port), false)
		//给目标机器发送自己的名片
		introduceSelf()
	}
	go read()
}

/*
	关闭服务器回调函数
*/
func shutdownCallback() {
	//回收映射的端口
	config.Sys_mapping.Reclaim()
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
		IsSuper: config.Init_IsSuperPeer,
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
	node := nodeStore.Get(nodeStore.ParseId(name), false, "")
	fmt.Println("节点下线", node)
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
			if id == nil {
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
