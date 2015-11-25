package mandela

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	msg "github.com/prestonTao/mandela/message_center"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"github.com/prestonTao/mandela/utils"
	"github.com/prestonTao/upnp"
	"net"
	"strconv"
	"time"
)

const (
	C_role_auto   = "auto"   //根据网络环境自适应
	C_role_client = "client" //客户端模式
	C_role_super  = "super"  //超级节点模式
	C_role_root   = "root"   //根节点模式
)

var (
	Init_IsSuperPeer               = false //有公网ip或添加了端口映射则是超级节点
	Init_GlobalUnicastAddress      = ""    //公网地址
	Init_GlobalUnicastAddress_port = 9981  //

	Sys_mapping = new(upnp.Upnp) //端口映射程序

	Init_LocalIP     = ""   //本地ip地址
	Init_LocalPort   = 9981 //本地监听端口
	Init_ExternalIP  = ""   //添加端口映射后的网关公网ip地址
	Init_MappingPort = 9981 //映射到路由器的端口

	// Mode_dev   = false       //是否是开发者模式
	Mode_local = false       //是否是局域网开发模式
	Init_role  = C_role_auto //服务器角色

)

var (
	// IsRoot        bool //是否是第一个节点
	superNodeIp   string
	superNodePort int
	privateKey    *rsa.PrivateKey
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
	if !Mode_local {
		portMapping()
	}
	//得到本地ip地址
	if Mode_local {
		Init_LocalIP = GetLocalHost()
	} else {
		Init_LocalIP = GetLocalIntenetIp()
	}
	//得到本机可用端口
	Init_LocalPort = GetAvailablePort()
	if Mode_local && Init_role == C_role_super {
		Init_IsSuperPeer = true
		Init_GlobalUnicastAddress = Init_LocalIP
		Init_GlobalUnicastAddress_port = Init_LocalPort
		Init_ExternalIP = Init_LocalIP
		Init_MappingPort = Init_LocalPort
	}
	//自己是根节点
	if Init_role == C_role_root {
		Init_IsSuperPeer = true
	}
	if Init_role == C_role_auto {

	}
}

/*
	根据网络环境启动程序
*/
func StartUpAuto() {

	if Mode_local {
		Init_LocalIP = "127.0.0.1"
	}
	AutoRole()
	loadIdInfo()
	InitSuperPeer()

	if Init_role != C_role_root {
		startLoadSuperPeer()
	}
	//开启web服务
	// go StartWeb()

	//没有idinfo的新节点
	if len(Init_IdInfo.Id) == 0 {
		//连接网络并得到一个idinfo
		newnode := nodeStore.IdInfo{
			Id:          Str_zaro,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05.999999999"),
			Domain:      GetRandomDomain(),
			Name:        "",
			Email:       "",
			SuperNodeId: Str_zaro,
		}
		idInfo, err := GetId(newnode)
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
	node = &nodeStore.Node{
		IdInfo:  Init_IdInfo,
		IsSuper: Init_IsSuperPeer, //是否是超级节点
		UdpPort: 0,
	}
	if Init_role == C_role_root {

		node.Addr = Init_GlobalUnicastAddress
		node.TcpPort = int32(Init_GlobalUnicastAddress_port)
	} else if Init_IsSuperPeer {
		node.Addr = Init_ExternalIP
		node.TcpPort = int32(Init_MappingPort)
	} else {
		node.Addr = Init_LocalIP
		node.TcpPort = int32(Init_LocalPort)
	}

	startUp(node)
	if Init_role == C_role_root {
		// StartRootPeer()
		startLoadSuperPeer()
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
	if Init_IsSuperPeer || Init_role == C_role_root {
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
	if Init_role == C_role_root {
		// StartRootPeer()
		startLoadSuperPeer()
	}
}

/*
	判断自己是否有公网ip地址
	若支持upnp协议，则添加一个端口映射
*/
func portMapping() {
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

func startUp(node *nodeStore.Node) {
	utils.Log.Debug("本机id为：%s", Init_IdInfo.GetId())
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
	engine.Listen(Init_LocalIP, int32(Init_LocalPort))
	if Init_role != C_role_root {
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
		IsSuper: Init_IsSuperPeer,
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

		if Init_role == C_role_client {
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
