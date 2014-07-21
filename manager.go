package mandela

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/prestonTao/mandela/cache"
	msg "github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/nodeStore"
	msgE "github.com/prestonTao/messageEngine"
	"github.com/prestonTao/upnp"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// type NodeInfo struct {
// 	hostIp  string
// 	tcpPort int
// 	udpPort int
// }

type Manager struct {
	IsRoot           bool
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
	cache            *cache.Memcache
	auth             *msgE.Auth
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
		this.rootId = RandNodeId(256)
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
		//  设置回调函数后监听
		//---------------------------------------------------------------
		auth := new(Auth)
		auth.nodeManager = this.nodeManager
		this.engine.SetAuth(auth)
		this.engine.Listen(this.hostIp, this.HostPort)
		this.engine.GetController().SetAttribute("nodeStore", this.nodeManager)
		//---------------------------------------------------------------
		//  end
		//---------------------------------------------------------------
		//自己连接自己
		this.engine.AddClientConn(this.rootId.String(), this.hostIp, this.HostPort, false)
	} else {
		//加载本地超级节点列表，
		// this.nodeStore = NewNodeStoreManager()
		this.nodeStoreManager = new(NodeStoreManager)
		this.nodeStoreManager.loadPeerEntry()

		//随机产生一个nodeid
		this.rootId = RandNodeId(256)
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
		//  设置回调函数后监听
		//---------------------------------------------------------------
		auth := new(Auth)
		auth.nodeManager = this.nodeManager
		this.engine.SetAuth(auth)
		this.engine.Listen(this.hostIp, this.HostPort)
		this.engine.GetController().SetAttribute("nodeStore", this.nodeManager)
		//---------------------------------------------------------------
		//  end
		//---------------------------------------------------------------
		//连接到超级节点
		hotsAndPost := strings.Split(this.nodeStoreManager.superNodeEntry[0], ":")
		port, err := strconv.Atoi(hotsAndPost[1])
		if err != nil {
			return err
		}
		this.engine.AddClientConn("superNode", hotsAndPost[0], int32(port), false)
	}
	this.cache = cache.NewMencache()
	this.engine.GetController().SetAttribute("cache", this.cache)
	go this.read()
	return nil
}

//启动消息服务器
func (this *Manager) initMsgEngine(name string) {
	this.hostIp = GetLocalIntenetIp()
	l, err := net.ListenPacket("udp", this.hostIp+":")
	if err != nil {
		fmt.Println("获取端口失败")
		return
	}
	hostPort, _ := strconv.Atoi(strings.Split(l.LocalAddr().String(), ":")[1])
	this.HostPort = int32(hostPort)

	this.engine = msgE.NewEngine(name)
	//注册所有的消息
	this.registerMsg()

	auth := new(Auth)
	this.engine.SetAuth(auth)
	this.engine.Listen(this.hostIp, this.HostPort)
}

// // //连接服务器，获得超级节点的ip地址
// // func (this *Manager) getSuperPeer() {
// // 	conn, err := net.Dial("tcp4", "127.0.0.1:9981")
// // 	if err != nil {
// // 		fmt.Println("连接超级节点发布服务器失败")
// // 		return
// // 	}
// // 	conn.Write([]byte(this.hostIp + ":" + strconv.Itoa(this.hostPort)))
// // 	buf := make([]byte, 1024)
// // 	n, _ := conn.Read(buf)
// // 	//得到超级节点的ip地址和端口
// // 	addrPort := strings.Split(string(buf[:n]), ":")

// // 	this.superNodeIp = addrPort[0]
// // 	portInt, _ := strconv.Atoi(addrPort[1])
// // 	this.superNodePort = portInt

// // }

// //对消息服务器的端口做映射
// func (this *Manager) initUpnp() {

// 	this.upnp = new(upnp.Upnp)
// 	if ok := this.upnp.AddPortMapping(this.hostPort, this.hostPort, "TCP"); ok {
// 		fmt.Println("端口映射成功")
// 	} else {
// 		fmt.Println("不支持upnp协议")
// 	}
// }

//启动分布式哈希表
func (this *Manager) initPeerNode() {
	// nodeStore.IsSuper = true //是超级节点
	// nodeStore.Addr = this.hostIp
	// nodeStore.TcpPort = this.HostPort

	node := &nodeStore.Node{
		NodeId:  this.rootId,
		IsSuper: true, //是超级节点
		Addr:    this.hostIp,
		TcpPort: this.HostPort,
		UdpPort: 0,
	}

	this.nodeManager = nodeStore.NewNodeManager(node)
	this.engine.GetController().SetAttribute("nodeStore", this.nodeManager)
	// this.engine.GetController().SetAttribute("nodeInQueue", this.nodeStore.InNodes)
	// msgE.Name = this.nodeStore.GetRootId()
}

// //连接超级节点
// func (this *Manager) connSuperPeer() {
// 	this.engine.AddClientConn("firstConnPeer", this.superNodeIp, int32(this.superNodePort))
// 	clientConn := this.engine.GetController().GetClientByName("firstConnPeer")
// 	fmt.Println("++", clientConn)
// }

//处理查找节点的请求
//本节点定期查询已知节点是否在线，更新节点信息
func (this *Manager) read() {
	// clientConn, _ := this.engine.GetController().GetSession(this.rootId.String())
	for {
		node := <-this.nodeManager.OutFindNode
		if node.NodeId != nil {
			remote := this.nodeManager.Get(node.NodeId.String())
			var clientConn msgE.Session
			if remote == nil {
				clientConn, _ = this.engine.GetController().GetSession("superNode")
				if clientConn == nil {
					continue
				}
			} else {
				clientConn, _ = this.engine.GetController().GetSession(remote.NodeId.String())
			}
			findNodeOne := &msg.FindNodeReq{
				NodeId: proto.String(this.nodeManager.GetRootId()),
				FindId: proto.String(node.NodeId.String()),
			}
			findNodeBytes, _ := proto.Marshal(findNodeOne)
			// clientConn := this.engine.GetController().GetClientByName("firstConnPeer")
			// fmt.Println(clientConn)
			clientConn.Send(msg.FindNodeReqNum, &findNodeBytes)
		}
		if node.NodeIdShould != nil {
			remote := this.nodeManager.Get(node.NodeIdShould.String())
			var clientConn msgE.Session
			if remote == nil {
				clientConn, _ = this.engine.GetController().GetSession("superNode")
				if clientConn == nil {
					continue
				}
			} else {
				// fmt.Println(remote.NodeId.String())
				clientConn, _ = this.engine.GetController().GetSession(remote.NodeId.String())
			}
			// fmt.Println(remote.NodeId.String())
			// clientConn, _ := this.engine.GetController().GetSession(remote.NodeId.String())
			findNodeOne := &msg.FindNodeReq{
				NodeId: proto.String(this.nodeManager.GetRootId()),
				FindId: proto.String(node.NodeIdShould.String()),
			}
			findNodeBytes, _ := proto.Marshal(findNodeOne)
			// clientConn := this.engine.GetController().GetClientByName("firstConnPeer")
			// fmt.Println(clientConn, "-0-\n")
			clientConn.Send(msg.FindNodeReqNum, &findNodeBytes)
		}
	}
}

//保存一个键值对
func (this *Manager) SaveData(key, value string) {
	clientConn, _ := this.engine.GetController().GetSession("superNode")
	data := []byte(key + "!" + value)
	clientConn.Send(msg.SaveKeyValueReqNum, &data)
}
