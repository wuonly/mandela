package mandela

import (
	// "code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	msg "github.com/prestonTao/mandela/message"
	msgE "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"github.com/prestonTao/upnp"
	"math/big"
	"net"
	"strconv"
)

var (
	Init_IsSuperPeer          = false //是超级节点
	Init_GlobalUnicastAddress = ""    //公网地址

	Init_LocalIP     = ""   //本地ip地址
	Init_LocalPort   = 9981 //本地监听端口
	Init_ExternalIP  = ""   //
	Init_MappingPort = 9981 //映射到路由器的端口

)

/*
	判断自己是否有公网ip地址
	是否支持upnp协议，添加一个端口映射
*/
func init() {

	Init_LocalIP := GetLocalIntenetIp()
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
		fmt.Println("本机ip是全球唯一公网地址")
		return
	}
	mapping := new(upnp.Upnp)
	err := mapping.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		Init_ExternalIP = mapping.GetewayOutsideIP
	}
	for i := 0; i < 1000; i++ {
		if err := mapping.AddPortMapping(Init_LocalPort, Init_MappingPort, "TCP"); err == nil {
			Init_IsSuperPeer = true
			fmt.Println("映射到公网地址：", Init_ExternalIP, ":", Init_MappingPort)
			return
		}
		Init_MappingPort = Init_MappingPort + 1
	}
	fmt.Println("端口映射失败")
}

var (
	IsRoot        bool //是否是第一个节点
	nodeManager   *nodeStore.NodeManager
	superNodeIp   string
	superNodePort int
	hostIp        string
	HostPort      int32
	rootId        *big.Int
	privateKey    *rsa.PrivateKey
	engine        *msgE.Engine
	auth          *msgE.Auth
)

/*
	根据网络环境启动程序
*/
func StartUp() {
	//最新的节点
	if Init_NewPeer {
		StartNewPeer()
		return
	}
	//是超级节点
	if Init_IsSuperPeer {

	}
}

/*
	启动新的节点
*/
func StartNewPeer() {

}

/*
	启动超级节点
*/
func StartSuperPeer() {

}

/*
	启动弱节点
*/
func StartWeak() {

}

/*
	启动根节点
*/
func StartRootPeer() {
	//随机产生一个nodeid
	rootId = nodeStore.RandNodeId()
	fmt.Println("本机id为：", hex.EncodeToString(rootId.Bytes()))
	//---------------------------------------------------------------
	//   启动消息服务器
	//---------------------------------------------------------------
	HostPort = int32(Init_LocalPort)

	engine = msgE.NewEngine(hex.EncodeToString(rootId.Bytes()))
	//注册所有的消息
	registerMsg()
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	var err error
	//生成密钥
	privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		fmt.Println("生成密钥错误", err.Error())
		return
	}

	//---------------------------------------------------------------
	//  启动分布式哈希表
	//---------------------------------------------------------------
	node := &nodeStore.Node{
		NodeId:  rootId,
		IsSuper: true, //是超级节点
		Addr:    Init_LocalIP,
		TcpPort: HostPort,
		UdpPort: 0,
	}
	nodeManager = nodeStore.NewNodeManager(node)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	//---------------------------------------------------------------
	//  设置关闭连接回调函数后监听
	//---------------------------------------------------------------
	auth := new(Auth)
	auth.nodeManager = nodeManager
	engine.SetAuth(auth)
	engine.SetCloseCallback(closeConnCallback)
	engine.Listen(Init_LocalIP, HostPort)
	engine.GetController().SetAttribute("nodeStore", nodeManager)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	go read()
}

//-------------------------------------------------------
// 1.加载本地超级节点列表，
//   启动消息服务器，
//   连接超级节点发布服务器，得到超级节点的ip地址及端口
//   加载本地密钥和节点id，或随机生成节点id
// 3.连接超级节点
//   使用upnp添加一个端口映射
// 4.注册节点id
//   处理查找节点的请求
//-------------------------------------------------------
func Run() error {

	if IsRoot {
		//随机产生一个nodeid
		rootId = nodeStore.RandNodeId()
	} else {
		//随机产生一个nodeid
		rootId = nodeStore.RandNodeId()
	}
	fmt.Println("本机id为：", hex.EncodeToString(rootId.Bytes()))
	//---------------------------------------------------------------
	//   启动消息服务器
	//---------------------------------------------------------------
	// initMsgEngine(rootId.String())
	hostIp = Init_LocalIP
	HostPort = int32(Init_LocalPort)

	engine = msgE.NewEngine(hex.EncodeToString(rootId.Bytes()))
	//注册所有的消息
	registerMsg()
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	var err error
	//生成密钥
	privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		fmt.Println("生成密钥错误", err.Error())
		return nil
	}

	//---------------------------------------------------------------
	//  启动分布式哈希表
	//---------------------------------------------------------------
	// initPeerNode()
	node := &nodeStore.Node{
		NodeId:  rootId,
		IsSuper: true, //是超级节点
		Addr:    hostIp,
		TcpPort: HostPort,
		UdpPort: 0,
	}
	nodeManager = nodeStore.NewNodeManager(node)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	//---------------------------------------------------------------
	//  设置关闭连接回调函数后监听
	//---------------------------------------------------------------
	auth := new(Auth)
	auth.nodeManager = nodeManager
	engine.SetAuth(auth)
	engine.SetCloseCallback(closeConnCallback)
	engine.Listen(hostIp, HostPort)
	engine.GetController().SetAttribute("nodeStore", nodeManager)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	if IsRoot {
		//自己连接自己
		// engine.AddClientConn(rootId.String(), hostIp, HostPort, false)
	} else {
		//连接到超级节点
		host, portStr, _ := net.SplitHostPort(Sys_superNodeEntry[0])
		// hotsAndPost := strings.Split(nodeStoreManager.superNodeEntry[0], ":")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		nodeManager.SuperName = engine.AddClientConn(host, int32(port), false)
		//给目标机器发送自己的名片
		introduceSelf()
	}
	//这里启动存储系统
	// cache = cache.NewMencache()
	// engine.GetController().SetAttribute("cache", cache)
	go read()
	return nil
}

//连接超级节点后，向超级节点介绍自己
//第一次连接超级节点，用代理方式查找离自己最近的节点
func introduceSelf() {
	session, _ := engine.GetController().GetSession(nodeManager.SuperName)

	//用代理方式查找最近的超级节点
	nodeMsg := msg.FindNode{
		NodeId:  session.GetName(),
		WantId:  nodeManager.GetRootId(),
		IsProxy: true,
		ProxyId: nodeManager.GetRootId(),
		IsSuper: true,
		Addr:    nodeManager.Root.Addr,
		TcpPort: nodeManager.Root.TcpPort,
		UdpPort: nodeManager.Root.UdpPort,
	}
	// resultBytes, _ := proto.Marshal(&nodeMsg)
	resultBytes, _ := json.Marshal(nodeMsg)

	session.Send(msg.FindNodeNum, &resultBytes)
}

//一个连接断开后的回调方法
func closeConnCallback(name string) {
	fmt.Println("客户端离线：", name)
	if name == nodeManager.SuperName {
		return
	}
	delNode := new(nodeStore.Node)
	delNode.NodeId, _ = new(big.Int).SetString(name, nodeStore.IdStrBit)
	nodeManager.DelNode(delNode)
}

//处理查找节点的请求
//本节点定期查询已知节点是否在线，更新节点信息
func read() {
	for {
		node := <-nodeManager.OutFindNode
		session, _ := engine.GetController().GetSession(nodeManager.SuperName)

		findNodeOne := &msg.FindNode{
			NodeId:  nodeManager.GetRootId(),
			IsProxy: false,
			ProxyId: nodeManager.GetRootId(),
		}
		//普通节点只需要定时查找最近的超级节点
		if !nodeManager.Root.IsSuper {
			if hex.EncodeToString(node.NodeId.Bytes()) == nodeManager.GetRootId() {
				findNodeOne.NodeId = session.GetName()
				findNodeOne.IsProxy = true
				findNodeOne.WantId = hex.EncodeToString(node.NodeId.Bytes())
				findNodeOne.IsSuper = true
				findNodeOne.Addr = nodeManager.Root.Addr
				findNodeOne.TcpPort = nodeManager.Root.TcpPort
				findNodeOne.UdpPort = nodeManager.Root.UdpPort

				// resultBytes, _ := proto.Marshal(findNodeOne)
				resultBytes, _ := json.Marshal(findNodeOne)
				session.Send(msg.FindNodeNum, &resultBytes)
			}
			continue
		}
		//--------------------------------------------
		//    查找邻居节点，只有超级节点才需要查找
		//--------------------------------------------
		if hex.EncodeToString(node.NodeId.Bytes()) == nodeManager.GetRootId() {
			//先发送左邻居节点查找请求
			findNodeOne.WantId = "left"
			id := nodeManager.GetLeftNode(*nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ := json.Marshal(findNodeOne)
			clientConn, ok := engine.GetController().GetSession(hex.EncodeToString(id[0].NodeId.Bytes()))
			if !ok {
				continue
			}
			err := clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			//发送右邻居节点查找请求
			findNodeOne.WantId = "right"
			id = nodeManager.GetRightNode(*nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ = json.Marshal(findNodeOne)
			clientConn, ok = engine.GetController().GetSession(hex.EncodeToString(id[0].NodeId.Bytes()))
			if !ok {
				continue
			}
			err = clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			continue
		}
		//--------------------------------------------
		//    查找普通节点，只有超级节点才需要查找
		//--------------------------------------------
		//这里临时加上去
		if nodeManager.Root.IsSuper {
			continue
		}
		findNodeOne.WantId = hex.EncodeToString(node.NodeId.Bytes())
		findNodeBytes, _ := json.Marshal(findNodeOne)

		remote := nodeManager.Get(hex.EncodeToString(node.NodeId.Bytes()), false, "")
		if remote == nil {
			continue
		}
		session, _ = engine.GetController().GetSession(hex.EncodeToString(remote.NodeId.Bytes()))
		if session == nil {
			continue
		}

		err := session.Send(msg.FindNodeNum, &findNodeBytes)
		if err != nil {
			fmt.Println("manager发送数据出错：", err.Error())
		}
	}
}

//保存一个键值对
func SaveData(key, value string) {
	clientConn, _ := engine.GetController().GetSession(nodeManager.SuperName)
	data := []byte(key + "!" + value)
	clientConn.Send(msg.SaveKeyValueReqNum, &data)
}

//给所有客户端发送消息
func SendMsgForAll(message string) {
	messageSend := msg.Message{
		Content: []byte(message),
	}
	for idOne, _ := range nodeManager.GetAllNodes() {
		if clientConn, ok := engine.GetController().GetSession(idOne); ok {
			messageSend.TargetId = idOne
			data, _ := json.Marshal(messageSend)
			clientConn.Send(msg.SendMessage, &data)
		}
	}
}

//给某个人发送消息
func SendMsgForOne(target, message string) {
	if nodeManager.GetRootId() == target {
		//发送给自己的
		fmt.Println(message)
		return
	}
	targetNode := nodeManager.Get(target, true, "")
	if targetNode == nil {
		fmt.Println("本节点未连入网络")
		return
	}
	session, ok := engine.GetController().GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
	if !ok {
		return
	}

	messageSend := msg.Message{
		TargetId: target,
		Content:  []byte(message),
	}
	// proto.
	// sendBytes, _ := proto.Marshal(&messageSend)
	sendBytes, _ := json.Marshal(&messageSend)
	err := session.Send(msg.SendMessage, &sendBytes)
	if err != nil {
		fmt.Println("message发送数据出错：", err.Error())
	}
}

//注册一个域名帐号
func CreateAccount(account string) {
	// id := GetHashKey(account)
}

func See() {
	allNodes := nodeManager.GetAllNodes()
	for key, _ := range allNodes {
		fmt.Println(key)
	}
}

func SeeLeftNode() {
	nodes := nodeManager.GetLeftNode(*nodeManager.Root.NodeId, nodeManager.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(hex.EncodeToString(id.NodeId.Bytes()))
	}
}

func SeeRightNode() {
	nodes := nodeManager.GetRightNode(*nodeManager.Root.NodeId, nodeManager.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(hex.EncodeToString(id.NodeId.Bytes()))
	}
}
