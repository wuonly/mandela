package mandela

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	msg "github.com/prestonTao/mandela/message_center"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"github.com/prestonTao/upnp"
	"net"
	"strconv"
)

var (
	Init_IsSuperPeer               = false //是超级节点
	Init_GlobalUnicastAddress      = ""    //公网地址
	Init_GlobalUnicastAddress_port = 9981  //

	Sys_mapping = new(upnp.Upnp) //端口映射程序

	Init_LocalIP     = ""   //本地ip地址
	Init_LocalPort   = 9981 //本地监听端口
	Init_ExternalIP  = ""   //
	Init_MappingPort = 9981 //映射到路由器的端口

	// Mode_dev = false //是否是开发者模式
)

/*
	判断自己是否有公网ip地址
	若支持upnp协议，则添加一个端口映射
*/
func portMapping() {
	//得到本地ip地址
	Init_LocalIP = GetLocalIntenetIp()
	/*
		获得一个可用的端口
	*/
	for i := 0; i < 1000; i++ {
		_, err := net.ListenPacket("udp", Init_LocalIP+":"+strconv.Itoa(Init_LocalPort))
		if err != nil {
			Init_LocalPort = Init_LocalPort + 1
		} else {
			break
		}
	}
	fmt.Println("监听一个本地地址：", Init_LocalIP, ":", Init_LocalPort)
	//本地地址是全球唯一公网地址
	if IsOnlyIp(Init_LocalIP) {
		Init_IsSuperPeer = true
		Init_GlobalUnicastAddress = Init_LocalIP
		Init_GlobalUnicastAddress_port = Init_LocalPort
		fmt.Println("本机ip是公网全球唯一地址")
		return
	}
	//获得网关公网地址
	err := Sys_mapping.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		Init_ExternalIP = Sys_mapping.GatewayOutsideIP
		Init_GlobalUnicastAddress = Init_ExternalIP
	}
	for i := 0; i < 1000; i++ {
		if err := Sys_mapping.AddPortMapping(Init_LocalPort, Init_MappingPort, "TCP"); err == nil {
			Init_IsSuperPeer = true
			Init_GlobalUnicastAddress_port = Init_MappingPort
			fmt.Println("映射到公网地址：", Init_ExternalIP, ":", Init_MappingPort)
			return
		}
		Init_MappingPort = Init_MappingPort + 1
	}
	fmt.Println("端口映射失败")
}

var (
	IsRoot bool //是否是第一个节点
	// nodeManager   *nodeStore.NodeManager
	superNodeIp   string
	superNodePort int
	privateKey    *rsa.PrivateKey
	// engine        *msgE.Engine
	// auth *msgE.Auth
)

/*
	根据网络环境启动程序
*/
func StartUpAuto() {
	loadIdInfo()
	if !IsRoot {
		startLoadSuperPeer()
	}
	//尝试端口映射
	portMapping()
	//没有idinfo的新节点
	if !Init_HaveId {
		//连接网络并得到一个idinfo
		idInfo, err := GetId(getSuperAddrOne())
		if err == nil {
			Init_IdInfo = *idInfo
			// saveIdInfo(Path_Id)
		} else {
			fmt.Println("从网络中获得idinfo失败")
			return
		}
	}
	//是超级节点
	var node *nodeStore.Node
	if Init_IsSuperPeer || IsRoot {
		node = &nodeStore.Node{
			IdInfo:  Init_IdInfo,
			IsSuper: Init_IsSuperPeer, //是否是超级节点
			Addr:    Init_GlobalUnicastAddress,
			TcpPort: int32(Init_GlobalUnicastAddress_port),
			UdpPort: 0,
		}
	} else {
		node = &nodeStore.Node{
			IdInfo:  Init_IdInfo,
			IsSuper: Init_IsSuperPeer, //是否是超级节点
			Addr:    Init_LocalIP,
			TcpPort: int32(Init_LocalPort),
			UdpPort: 0,
		}
	}
	startUp(node)
	if IsRoot {
		// StartRootPeer()
		startLoadSuperPeer()
	}
}

func startUp(node *nodeStore.Node) {
	fmt.Println("本机id为：", Init_IdInfo.GetId())
	/*
		启动消息服务器
	*/
	// engine = msgE.NewEngine(string(Init_IdInfo.Build()))
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
	// nodeManager = nodeStore.NewNodeManager(node)
	nodeStore.InitNodeStore(node)
	/*
		设置关闭连接回调函数后监听
	*/
	// auth :=
	engine.SetAuth(new(Auth))
	engine.SetCloseCallback(closeConnCallback)
	engine.Listen(Init_LocalIP, int32(Init_LocalPort))
	// engine.GetController().SetAttribute("nodeStore", nodeManager)

	if !IsRoot {
		/*
			连接到超级节点
		*/
		host, portStr, _ := net.SplitHostPort(getSuperAddrOne())
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
	Sys_mapping.Reclaim()
}

//连接超级节点后，向超级节点介绍自己
//第一次连接超级节点，用代理方式查找离自己最近的节点
func introduceSelf() {
	session, _ := engine.GetController().GetSession(nodeStore.SuperName)
	//用代理方式查找最近的超级节点
	nodeMsg := msg.FindNode{
		NodeId:  session.GetName(),
		WantId:  nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
		IsProxy: true,
		ProxyId: nodeStore.GetRootIdInfoString(),
		IsSuper: Init_IsSuperPeer,
		Addr:    nodeStore.Root.Addr,
		TcpPort: nodeStore.Root.TcpPort,
		UdpPort: nodeStore.Root.UdpPort,
	}

	// resultBytes, _ := proto.Marshal(&nodeMsg)
	resultBytes, _ := json.Marshal(nodeMsg)

	session.Send(msg.FindNodeNum, &resultBytes)
}

//一个连接断开后的回调方法
func closeConnCallback(name string) {
	fmt.Println("客户端离线：", name)
	if name == nodeStore.SuperName {
		return
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

				resultBytes, _ := json.Marshal(findNodeOne)
				session.Send(msg.FindNodeNum, &resultBytes)
				// continue
			}

			//先发送左邻居节点查找请求
			findNodeOne.WantId = "left"
			id := nodeStore.GetLeftNode(*nodeStore.Root.IdInfo.GetBigIntId(), 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ := json.Marshal(findNodeOne)
			clientConn, ok := engine.GetController().GetSession(string(id[0].IdInfo.Build()))
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
			clientConn, ok = engine.GetController().GetSession(string(id[0].IdInfo.Build()))
			if !ok {
				continue
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
