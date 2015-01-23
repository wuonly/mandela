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
	"strings"
)

var (
	Listen_port = 9981
)

type Manager struct {
	IsRoot        bool //是否是第一个节点
	nodeManager   *nodeStore.NodeManager
	superNodeIp   string
	superNodePort int
	hostIp        string
	HostPort      int32
	rootId        *big.Int
	privateKey    *rsa.PrivateKey
	upnp          *upnp.Upnp
	engine        *msgE.Engine
	auth          *msgE.Auth
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
	//是新节点
	if Init_NewPeer {

	}

	if this.IsRoot {
		//随机产生一个nodeid
		this.rootId = nodeStore.RandNodeId()
	} else {
		//随机产生一个nodeid
		this.rootId = nodeStore.RandNodeId()
	}
	fmt.Println("本机id为：", hex.EncodeToString(this.rootId.Bytes()))
	//---------------------------------------------------------------
	//   启动消息服务器
	//---------------------------------------------------------------
	// this.initMsgEngine(this.rootId.String())
	this.hostIp = GetLocalIntenetIp()
tag:
	l, err := net.ListenPacket("udp", this.hostIp+":"+strconv.Itoa(Listen_port))
	if err != nil {
		Listen_port = Listen_port + 1
		goto tag
	}
	hostPort, _ := strconv.Atoi(strings.Split(l.LocalAddr().String(), ":")[1])
	this.HostPort = int32(hostPort)

	this.engine = msgE.NewEngine(hex.EncodeToString(this.rootId.Bytes()))
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
		host, portStr, _ := net.SplitHostPort(Sys_superNodeEntry[0])
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
		NodeId:  session.GetName(),
		WantId:  this.nodeManager.GetRootId(),
		IsProxy: true,
		ProxyId: this.nodeManager.GetRootId(),
		IsSuper: true,
		Addr:    this.nodeManager.Root.Addr,
		TcpPort: this.nodeManager.Root.TcpPort,
		UdpPort: this.nodeManager.Root.UdpPort,
	}
	// resultBytes, _ := proto.Marshal(&nodeMsg)
	resultBytes, _ := json.Marshal(nodeMsg)

	session.Send(msg.FindNodeNum, &resultBytes)
}

//一个连接断开后的回调方法
func (this *Manager) closeConnCallback(name string) {
	fmt.Println("客户端离线：", name)
	if name == this.nodeManager.SuperName {
		return
	}
	delNode := new(nodeStore.Node)
	delNode.NodeId, _ = new(big.Int).SetString(name, nodeStore.IdStrBit)
	this.nodeManager.DelNode(delNode)
}

//处理查找节点的请求
//本节点定期查询已知节点是否在线，更新节点信息
func (this *Manager) read() {
	for {
		node := <-this.nodeManager.OutFindNode
		session, _ := this.engine.GetController().GetSession(this.nodeManager.SuperName)

		findNodeOne := &msg.FindNode{
			NodeId:  this.nodeManager.GetRootId(),
			IsProxy: false,
			ProxyId: this.nodeManager.GetRootId(),
		}
		//普通节点只需要定时查找最近的超级节点
		if !this.nodeManager.Root.IsSuper {
			if hex.EncodeToString(node.NodeId.Bytes()) == this.nodeManager.GetRootId() {
				findNodeOne.NodeId = session.GetName()
				findNodeOne.IsProxy = true
				findNodeOne.WantId = hex.EncodeToString(node.NodeId.Bytes())
				findNodeOne.IsSuper = true
				findNodeOne.Addr = this.nodeManager.Root.Addr
				findNodeOne.TcpPort = this.nodeManager.Root.TcpPort
				findNodeOne.UdpPort = this.nodeManager.Root.UdpPort

				// resultBytes, _ := proto.Marshal(findNodeOne)
				resultBytes, _ := json.Marshal(findNodeOne)
				session.Send(msg.FindNodeNum, &resultBytes)
			}
			continue
		}
		//--------------------------------------------
		//    查找邻居节点，只有超级节点才需要查找
		//--------------------------------------------
		if hex.EncodeToString(node.NodeId.Bytes()) == this.nodeManager.GetRootId() {
			//先发送左邻居节点查找请求
			findNodeOne.WantId = "left"
			id := this.nodeManager.GetLeftNode(*this.nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ := json.Marshal(findNodeOne)
			clientConn, ok := this.engine.GetController().GetSession(hex.EncodeToString(id[0].NodeId.Bytes()))
			if !ok {
				continue
			}
			err := clientConn.Send(msg.FindNodeNum, &findNodeBytes)
			if err != nil {
				fmt.Println("manager发送数据出错：", err.Error())
			}
			//发送右邻居节点查找请求
			findNodeOne.WantId = "right"
			id = this.nodeManager.GetRightNode(*this.nodeManager.Root.NodeId, 1)
			if id == nil {
				continue
			}
			findNodeBytes, _ = json.Marshal(findNodeOne)
			clientConn, ok = this.engine.GetController().GetSession(hex.EncodeToString(id[0].NodeId.Bytes()))
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
		findNodeOne.WantId = hex.EncodeToString(node.NodeId.Bytes())
		findNodeBytes, _ := json.Marshal(findNodeOne)

		remote := this.nodeManager.Get(hex.EncodeToString(node.NodeId.Bytes()), false, "")
		if remote == nil {
			continue
		}
		session, _ = this.engine.GetController().GetSession(hex.EncodeToString(remote.NodeId.Bytes()))
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
	messageSend := msg.Message{
		Content: []byte(message),
	}
	for idOne, _ := range this.nodeManager.GetAllNodes() {
		if clientConn, ok := this.engine.GetController().GetSession(idOne); ok {
			messageSend.TargetId = idOne
			data, _ := json.Marshal(messageSend)
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
	session, ok := this.engine.GetController().GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
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
		fmt.Println(hex.EncodeToString(id.NodeId.Bytes()))
	}
}

func (this *Manager) SeeRightNode() {
	nodes := this.nodeManager.GetRightNode(*this.nodeManager.Root.NodeId, this.nodeManager.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(hex.EncodeToString(id.NodeId.Bytes()))
	}
}
