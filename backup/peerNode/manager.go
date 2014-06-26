package peerNode

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/astaxie/beego"
	"mandela/peerNode"
	msg "mandela/peerNode/message"
	msgE "mandela/peerNode/messageEngine"
	_ "mandela/peerNode/msgServer"
	"mandela/peerNode/upnp"
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
	superNodeIp   string
	superNodePort int
	hostIp        string
	hostPort      int
	rootId        *big.Int
	privateKey    *rsa.PrivateKey
	upnp          *upnp.Upnp
	serverManager *msgE.ServerManager
	nodeStore     *peerNode.NodeStore
}

func (this *Manager) Run() error {
	//启动消息服务器
	this.initMsgEngine()
	//获得超级节点
	this.getSuperPeer()
	//将端口映射到外网
	this.initUpnp()
	//
	this.initPeerNode()
	//连接超级节点
	this.connSuperPeer()
	//处理查找节点的请求
	go this.read()
	return nil
}

//-------------------------------------------------------
// 1.加载本地超级节点列表，连接超级节点发布服务器，得到超级节点的ip地址及端口
//   加载本地密钥和节点id，或随机生成节点id
// 3.启动消息服务器，连接超级节点
//   使用upnp添加一个端口映射
// 4.注册节点id
//   处理查找节点的请求
//-------------------------------------------------------

type SuperNodeEntry []string

func (this *Manager) loadPeerEntry() {
	filePath := "conf/nodeEntry.json"
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		//本地没有地址列表文件
		//连接服务器获取地址列表
		conn, err := net.Dial("tcp4", "127.0.0.1:9981")
		if err != nil {
			fmt.Println("连接超级节点发布服务器失败")
			return
		}
		conn.Write([]byte(this.hostIp + ":" + strconv.Itoa(this.hostPort)))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		//得到超级节点的ip地址和端口
		addrPort := strings.Split(string(buf[:n]), ":")

		this.superNodeIp = addrPort[0]
		portInt, _ := strconv.Atoi(addrPort[1])
		this.superNodePort = portInt
	} else {
		s := new(SuperNodeEntry)
		json.Unmarshal(fileBytes, s)
		addrPort := strings.Split(s[0], ":")

		this.superNodeIp = addrPort[0]
		portInt, _ := strconv.Atoi(addrPort[1])
		this.superNodePort = portInt
	}
}

func (this *Manager) loadPrivateKey() {
	filePath := "conf/private.key"

	//先加载本地的私钥
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		//本地没有私钥，创建一个
		//随机产生一个nodeid
		this.rootId = RandNodeId(NodeIdLevel)
		//生成密钥
		this.privateKey, err = rsa.GenerateKey(rand.Reader, 512)
		if err != nil {
			fmt.Println("生成密钥错误", err.Error())
			return nil
		}
	} else {
		var keyBytes bytes.Buffer
		keyBytes.Write(fileBytes)
		dec := gob.NewDecoder(&keyBytes)
		key := &PrivateKey{}
		err := dec.Decode(key)
		if err != nil {
			fmt.Println("密钥文件损坏")
			return nil
		}
		this.rootId, _ = new(big.Int).SetString(key.NodeId, 10)
		this.privateKey = &key.PrivateKey
	}

}

//启动消息服务器
func (this *Manager) initMsgEngine() {
	this.hostIp = upnp.GetLocalIntenetIp()
	l, err := net.ListenPacket("udp", this.hostIp+":")
	if err != nil {
		fmt.Println("获取端口失败")
		return
	}
	hostPort := l.LocalAddr().String()
	this.hostPort, _ = strconv.Atoi(strings.Split(hostPort, ":")[1])

	//---------------------------------------
	//  手动设置端口
	//---------------------------------------
	this.hostPort = 9990

	msgE.IP = this.hostIp
	msgE.PORT = this.hostPort
	this.serverManager = new(msgE.ServerManager)
	this.serverManager.Run()

}

// //连接服务器，获得超级节点的ip地址
// func (this *Manager) getSuperPeer() {
// 	conn, err := net.Dial("tcp4", "127.0.0.1:9981")
// 	if err != nil {
// 		fmt.Println("连接超级节点发布服务器失败")
// 		return
// 	}
// 	conn.Write([]byte(this.hostIp + ":" + strconv.Itoa(this.hostPort)))
// 	buf := make([]byte, 1024)
// 	n, _ := conn.Read(buf)
// 	//得到超级节点的ip地址和端口
// 	addrPort := strings.Split(string(buf[:n]), ":")

// 	this.superNodeIp = addrPort[0]
// 	portInt, _ := strconv.Atoi(addrPort[1])
// 	this.superNodePort = portInt

// }

//对消息服务器的端口做映射
func (this *Manager) initUpnp() {

	this.upnp = new(upnp.Upnp)
	if ok := this.upnp.AddPortMapping(this.hostPort, this.hostPort, "TCP"); ok {
		fmt.Println("端口映射成功")
	} else {
		fmt.Println("不支持upnp协议")
	}
}

//启动分布式哈希表
func (this *Manager) initPeerNode() {
	peerNode.IsSuper = true //是超级节点
	peerNode.Addr = this.hostIp
	peerNode.TcpPort = this.hostPort
	this.nodeStore = peerNode.NewNodeStore("", "")
	this.serverManager.GetController().SetAttribute("peerNode", this.nodeStore)
	// this.serverManager.GetController().SetAttribute("nodeInQueue", this.nodeStore.InNodes)
	msgE.Name = this.nodeStore.GetRootId()
}

//连接超级节点
func (this *Manager) connSuperPeer() {
	this.serverManager.AddClientConn("firstConnPeer", this.superNodeIp, int32(this.superNodePort))
	clientConn := this.serverManager.GetController().GetClientByName("firstConnPeer")
	fmt.Println("++", clientConn)
}

//处理查找节点的请求
func (this *Manager) read() {
	for {
		node := <-this.nodeStore.OutFindNode
		if node.NodeId != nil {
			findNodeOne := &msg.FindNodeReq{
				NodeId: proto.String(this.nodeStore.GetRootId()),
				FindId: proto.String(node.NodeId.String()),
			}
			findNodeBytes, _ := proto.Marshal(findNodeOne)
			clientConn := this.serverManager.GetController().GetClientByName("firstConnPeer")
			// fmt.Println(clientConn)
			clientConn.Send(msg.FindNodeReqNum, &findNodeBytes)
		}
		if node.NodeIdShould != nil {
			findNodeOne := &msg.FindNodeReq{
				NodeId: proto.String(this.nodeStore.GetRootId()),
				FindId: proto.String(node.NodeIdShould.String()),
			}
			findNodeBytes, _ := proto.Marshal(findNodeOne)
			clientConn := this.serverManager.GetController().GetClientByName("firstConnPeer")
			clientConn.Send(msg.FindNodeReqNum, &findNodeBytes)
		}
	}
}
