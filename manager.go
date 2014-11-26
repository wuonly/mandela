package mandela

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	// "github.com/prestonTao/mandela/cache"
	msg "github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/nodeStore"
	msgE "github.com/prestonTao/messageEngine"
	"github.com/prestonTao/upnp"
	"math/big"
	"net"
	"strconv"
	"strings"
)

type Manager struct {
	IsRoot           bool //是否是第一个节点
	nodeStoreManager *NodeStoreManager
	nodeManager      *nodeStore.NodeManager
	superNodeIp      string
	superNodePort    int
	hostIp           string
	HostPort         int32
	rootId           *big.Int
	privateKey       *rsa.PrivateKey
	upnp             *upnp.Upnp
	engine           *msgE.Engine
	// cache            *cache.Memcache
	auth *msgE.Auth
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
func (this *Manager) Run() error {
	if this.IsRoot {
		//随机产生一个nodeid
		this.rootId = nodeStore.RandNodeId()
	} else {
		//随机产生一个nodeid
		this.rootId = nodeStore.RandNodeId()
		this.nodeStoreManager = new(NodeStoreManager)
		this.nodeStoreManager.loadPeerEntry()
	}
	fmt.Println("本客户端随机id为：", this.rootId.String())
	//---------------------------------------------------------------
	//   启动消息服务器
	//---------------------------------------------------------------
	// this.initMsgEngine(this.rootId.String())
	this.hostIp = GetLocalIntenetIp()
	l, err := net.ListenPacket("udp", this.hostIp+":")
	if err != nil {
		fmt.Println("获取端口失败")
		return err
	}
	hostPort, _ := strconv.Atoi(strings.Split(l.LocalAddr().String(), ":")[1])
	this.HostPort = int32(hostPort)

	this.engine = msgE.NewEngine(this.rootId.String())
	//注册所有的消息
	this.registerMsg()
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	// var err error
	//生成密钥
	this.privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		fmt.Println("生成密钥错误", err.Error())
		return nil
	}

	//---------------------------------------------------------------
	//  启动分布式哈希表
	//---------------------------------------------------------------
	// this.initPeerNode()
	node := &nodeStore.Node{
		NodeId:  this.rootId,
		IsSuper: true, //是超级节点
		Addr:    this.hostIp,
		TcpPort: this.HostPort,
		UdpPort: 0,
	}
	this.nodeManager = nodeStore.NewNodeManager(node)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	//---------------------------------------------------------------
	//  设置关闭连接回调函数后监听
	//---------------------------------------------------------------
	auth := new(Auth)
	auth.nodeManager = this.nodeManager
	this.engine.SetAuth(auth)
	this.engine.SetCloseCallback(this.closeConnCallback)
	this.engine.Listen(this.hostIp, this.HostPort)
	this.engine.GetController().SetAttribute("nodeStore", this.nodeManager)
	//---------------------------------------------------------------
	//  end
	//---------------------------------------------------------------
	if this.IsRoot {
		//自己连接自己
		// this.engine.AddClientConn(this.rootId.String(), this.hostIp, this.HostPort, false)
	} else {
		//连接到超级节点
		host, portStr, _ := net.SplitHostPort(this.nodeStoreManager.superNodeEntry[0])
		// hotsAndPost := strings.Split(this.nodeStoreManager.superNodeEntry[0], ":")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		this.nodeManager.SuperName = this.engine.AddClientConn(host, int32(port), false)
		//给目标机器发送自己的名片
		this.introduceSelf()
	}
	//这里启动存储系统
	// this.cache = cache.NewMencache()
	// this.engine.GetController().SetAttribute("cache", this.cache)
	go this.read()
	return nil
}

//连接超级节点后，向超级节点介绍自己
//第一次连接超级节点，用代理方式查找离自己最近的节点
func (this *Manager) introduceSelf() {
	session, _ := this.engine.GetController().GetSession(this.nodeManager.SuperName)

	//用代理方式查找最近的超级节点
	nodeMsg := msg.FindNode{
		NodeId:  proto.String(session.GetName()),
		WantId:  proto.String(this.nodeManager.GetRootId()),
		IsProxy: proto.Bool(true),
		ProxyId: proto.String(this.nodeManager.GetRootId()),
		IsSuper: proto.Bool(true),
		Addr:    proto.String(this.nodeManager.Root.Addr),
		TcpPort: proto.Int32(this.nodeManager.Root.TcpPort),
		UdpPort: proto.Int32(this.nodeManager.Root.UdpPort),
	}
	resultBytes, _ := proto.Marshal(&nodeMsg)

	session.Send(msg.FindNodeNum, &resultBytes)
}

//一个连接断开后的回调方法
func (this *Manager) closeConnCallback(name string) {
	fmt.Println("客户端离线：", name)
	if name == this.nodeManager.SuperName {
		return
	}
	delNode := new(nodeStore.Node)
	delNode.NodeId, _ = new(big.Int).SetString(name, 10)
	this.nodeManager.DelNode(delNode)
}

//处理查找节点的请求
//本节点定期查询已知节点是否在线，更新节点信息
func (this *Manager) read() {
	for {
		node := <-this.nodeManager.OutFindNode
		session, _ := this.engine.GetController().GetSession(this.nodeManager.SuperName)

		findNodeOne := &msg.FindNode{
			NodeId:  proto.String(this.nodeManager.GetRootId()),
			IsProxy: proto.Bool(false),
			ProxyId: proto.String(this.nodeManager.GetRootId()),
		}
		//普通节点只需要定时查找最近的超级节点
		if !this.nodeManager.Root.IsSuper {
			if node.NodeId.String() == this.nodeManager.GetRootId() {
				findNodeOne.NodeId = proto.String(session.GetName())
				findNodeOne.IsProxy = proto.Bool(true)
				findNodeOne.WantId = proto.String(node.NodeId.String())
				findNodeOne.IsSuper = proto.Bool(true)
				findNodeOne.Addr = proto.String(this.nodeManager.Root.Addr)
				findNodeOne.TcpPort = proto.Int32(this.nodeManager.Root.TcpPort)
				findNodeOne.UdpPort = proto.Int32(this.nodeManager.Root.UdpPort)

				resultBytes, _ := proto.Marshal(findNodeOne)
				session.Send(msg.FindNodeNum, &resultBytes)
			}
			continue
		}
		//--------------------------------------------
		//    查找邻居节点，只有超级节点才需要查找
		//--------------------------------------------
		if node.NodeId.String() == this.nodeManager.GetRootId() {
			//先发送左邻居节点查找请求
			findNodeOne.WantId = proto.String("left")
			id := this.nodeManager.GetLeftNode(*this.nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ := proto.Marshal(findNodeOne)
			clientConn, ok := this.engine.GetController().GetSession(id[0].NodeId.String())
			if !ok {
				continue
			}
			err := clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			//发送右邻居节点查找请求
			findNodeOne.WantId = proto.String("right")
			id = this.nodeManager.GetRightNode(*this.nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ = proto.Marshal(findNodeOne)
			clientConn, ok = this.engine.GetController().GetSession(id[0].NodeId.String())
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
		if this.nodeManager.Root.IsSuper {
			continue
		}
		findNodeOne.WantId = proto.String(node.NodeId.String())
		findNodeBytes, _ := proto.Marshal(findNodeOne)

		remote := this.nodeManager.Get(node.NodeId.String(), false, "")
		if remote == nil {
			continue
		}
		session, _ = this.engine.GetController().GetSession(remote.NodeId.String())
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
func (this *Manager) SaveData(key, value string) {
	clientConn, _ := this.engine.GetController().GetSession(this.nodeManager.SuperName)
	data := []byte(key + "!" + value)
	clientConn.Send(msg.SaveKeyValueReqNum, &data)
}

//给所有客户端发送消息
func (this *Manager) SendMsgForAll(message string) {
	for idOne, _ := range this.nodeManager.GetAllNodes() {
		if clientConn, ok := this.engine.GetController().GetSession(idOne); ok {
			data := []byte(message)
			clientConn.Send(msg.SendMessage, &data)
		}
	}
}

//给某个人发送消息
func (this *Manager) SendMsgForOne(target, message string) {
	if this.nodeManager.GetRootId() == target {
		//发送给自己的
		fmt.Println(message)
		return
	}
	targetNode := this.nodeManager.Get(target, true, "")
	if targetNode == nil {
		fmt.Println("本节点未连入网络")
		return
	}
	session, ok := this.engine.GetController().GetSession(targetNode.NodeId.String())
	if !ok {
		return
	}

	messageSend := msg.Message{
		TargetId: proto.String(target),
		Content:  []byte(message),
	}
	// proto.
	sendBytes, _ := proto.Marshal(&messageSend)
	err := session.Send(msg.SendMessage, &sendBytes)
	if err != nil {
		fmt.Println("message发送数据出错：", err.Error())
	}
}

//注册一个域名帐号
func (this *Manager) CreateAccount(account string) {
	// id := GetHashKey(account)
}

func (this *Manager) See() {
	allNodes := this.nodeManager.GetAllNodes()
	for key, _ := range allNodes {
		fmt.Println(key)
	}
}

func (this *Manager) SeeLeftNode() {
	nodes := this.nodeManager.GetLeftNode(*this.nodeManager.Root.NodeId, this.nodeManager.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(id.NodeId.String())
	}
}

func (this *Manager) SeeRightNode() {
	nodes := this.nodeManager.GetRightNode(*this.nodeManager.Root.NodeId, this.nodeManager.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(id.NodeId.String())
	}
}
