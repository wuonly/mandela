package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/astaxie/beego"
	msg "mandela/message"
	msgE "mandela/messageEngine"
	_ "mandela/msgServer"
	"mandela/peerNode"
	"mandela/upnp"
	"net"
	"strconv"
	"strings"
)

var operation *Operation

func init() {
	operation = new(Operation)

	beego.AddAPPStartHook(operation.Run)
}

type Operation struct {
	superNodeIp   string
	superNodePort int
	hostIp        string
	hostPort      int
	upnp          *upnp.Upnp
	serverManager *msgE.ServerManager
	nodeStore     *peerNode.NodeStore
}

func (this *Operation) Run() error {
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

func (this *Operation) loadPeerEntry() {
	filePath := "conf/nodeEntry.json"
	s := new(SuperNodeEntry)
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
		json.Unmarshal(fileBytes, s)
	}
}

func (this *Operation) loadPrivateKey() {
	filePath := "conf/private.key"
	var nodeId *big.Int
	var privateKey *rsa.PrivateKey
	var keyBytes bytes.Buffer

	//先加载本地的私钥
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		//本地没有私钥，创建一个
		//生成密钥
		privateKey, err = rsa.GenerateKey(rand.Reader, 512)
		if err != nil {
			fmt.Println("生成密钥错误", err.Error())
			return nil
		}
		//  创建id
		if name == "" || account == "" {
			//随机产生一个nodeid
			nodeId = RandNodeId(NodeIdLevel)
		}
	} else {
		keyBytes.Write(fileBytes)
		dec := gob.NewDecoder(&keyBytes)
		key := &PrivateKey{}
		err := dec.Decode(key)
		if err != nil {
			fmt.Println("密钥文件损坏")
			return nil
		}
		nodeId, _ = new(big.Int).SetString(key.NodeId, 10)
		privateKey = &key.PrivateKey
	}

}

//启动消息服务器
func (this *Operation) initMsgEngine() {
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
// func (this *Operation) getSuperPeer() {
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
func (this *Operation) initUpnp() {

	this.upnp = new(upnp.Upnp)
	if ok := this.upnp.AddPortMapping(this.hostPort, this.hostPort, "TCP"); ok {
		fmt.Println("端口映射成功")
	} else {
		fmt.Println("不支持upnp协议")
	}
}

//启动分布式哈希表
func (this *Operation) initPeerNode() {
	peerNode.IsSuper = true //是超级节点
	peerNode.Addr = this.hostIp
	peerNode.TcpPort = this.hostPort
	this.nodeStore = peerNode.NewNodeStore("", "")
	this.serverManager.GetController().SetAttribute("peerNode", this.nodeStore)
	// this.serverManager.GetController().SetAttribute("nodeInQueue", this.nodeStore.InNodes)
	msgE.Name = this.nodeStore.GetRootId()
}

//连接超级节点
func (this *Operation) connSuperPeer() {
	this.serverManager.AddClientConn("firstConnPeer", this.superNodeIp, int32(this.superNodePort))
	clientConn := this.serverManager.GetController().GetClientByName("firstConnPeer")
	fmt.Println("++", clientConn)
}

//处理查找节点的请求
func (this *Operation) read() {
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
