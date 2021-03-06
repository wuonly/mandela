package nodeStore

import (
	// "bytes"
	// "crypto/rand"
	"crypto/rsa"
	// "encoding/gob"
	"fmt"
	// "io/ioutil"
	"math/big"
	// "os"
	"time"
)

var (
	NodeIdLevel       = 256         //节点id二进制长度
	IsSuper           = false       //是否是超级节点，并且提供代理功能
	Addr              = "127.0.0.1" //外网地址
	TcpPort     int32 = 0           //外网端口
	UdpPort     int32 = 0           //外网端口
)

type NodeStore struct {
	root           *Node            //
	isNew          bool             //是否是新节点
	nodes          map[string]*Node //十进制字符串为键
	consistentHash *ConsistentHash  //一致性hash表
	pipeNode       chan *Node       //收集各个节点更新，下线消息的管道
	InNodes        chan *Node       //需要更新的节点
	OutFindNode    chan *Node       //需要查询是否在线的节点
	Groups         *NodeGroup       //组
}

func (this *NodeStore) Run() {
	go this.recv()
	go this.checkSelf()
	go this.nodeComm()
	//向网络中查找自己，通知相关节点自己上线了
	this.OutFindNode <- this.root
	idsInt := this.getNodeNetworkNum()
	for _, idOne := range idsInt {
		this.OutFindNode <- &Node{NodeIdShould: idOne}
	}
}

//检查下线的节点，把他们移除掉
func (this *NodeStore) checkSelf() {
	//160分钟检查一次
	time.Sleep(time.Minute * 160)
	for _, value := range this.nodes {
		if value.Status == 3 {
			this.delNode(value)
		}
	}
	go this.checkSelf()
}

//需要更新的节点
func (this *NodeStore) recv() {
	for {
		node := <-this.InNodes
		this.AddNode(node)
	}
}

//负责处理仓库中节点发来的更新，下线消息
func (this *NodeStore) nodeComm() {
	for {
		node := <-this.pipeNode
		switch node.Status {
		case 1:
		case 2:
			//需要查询是否在线的节点
			this.OutFindNode <- node
		case 3:
			//需要下线的节点
		}
	}
}

//添加一个节点
func (this *NodeStore) AddNode(node *Node) {
	//添加本节点
	if node.NodeId.Cmp(this.root.NodeId) == 0 {
		this.nodes[this.root.NodeId.String()] = this.root
		return
	}
	node.Out = this.pipeNode
	node.LastContactTimestamp = time.Now()
	node.OverTime = 1 * 60 * 60
	node.SelectTime = 5 * 60
	go node.timeOut()
	this.nodes[node.NodeId.String()] = node
	this.consistentHash.Add(node.NodeId)
}

//删除一个节点
func (this *NodeStore) delNode(node *Node) {
	this.consistentHash.Del(node.NodeId)
	delete(this.nodes, node.NodeId.String())
}

//得到每个节点网络的网络号，不包括本节点
func (this *NodeStore) getNodeNetworkNum() map[string]*big.Int {
	// rootInt, _ := new(big.Int).SetString(, 10)
	networkNums := make(map[string]*big.Int, 3000)
	for i := 0; i < NodeIdLevel; i++ {
		//---------------------------------
		//将后面的i位置零
		//---------------------------------
		startInt := new(big.Int).Lsh(new(big.Int).Rsh(this.root.NodeId, uint(i)), uint(i))
		//---------------------------------
		//最后一位取反
		//---------------------------------
		networkNum := new(big.Int).Xor(startInt, new(big.Int).Lsh(big.NewInt(1), uint(i)))
		// fmt.Println("haha", i)
		// Print(networkNum)
		// networkNums = append(networkNums, networkNum)
		networkNums[networkNum.String()] = networkNum
	}
	return networkNums
}

//根据节点id得到一个节点的信息，id为十进制字符串
func (this *NodeStore) Get(nodeId string) *Node {
	nodeIdInt, b := new(big.Int).SetString(nodeId, 10)
	if !b {
		fmt.Println("节点id格式不正确，应该为十进制字符串")
	}
	targetNodeId := this.consistentHash.Get(nodeIdInt)
	if targetNodeId != nil {
		return this.nodes[targetNodeId.String()]
	}
	return this.root
}

//得到所有的节点
func (this *NodeStore) GetAllNodes() map[string]*Node {
	return this.nodes
}

//得到本节点id
func (this *NodeStore) GetRootId() string {
	return this.root.NodeId.String()
}

//根据域名和帐号名称创建一个节点
//域名为p2p网络中唯一域名
//域名加帐号名称才能够确定一个节点id，为了把域名和帐号绑定
func NewNodeStore(nodeId *big.Int) *NodeStore {

	node := Node{
		NodeId:  nodeId,
		IsSuper: IsSuper,
		Addr:    Addr,
		TcpPort: TcpPort,
		UdpPort: UdpPort,
		// Key:     privateKey,
	}

	// fmt.Println("本次创建的nodeid为：", node.NodeId, "私钥：", node.Key)
	//节点长度为512,深度为513
	nodeStore := &NodeStore{
		root:           &node,
		consistentHash: new(ConsistentHash),
		nodes:          make(map[string]*Node, 1000),
		OutFindNode:    make(chan *Node, 1000),
		pipeNode:       make(chan *Node, 1000),
		InNodes:        make(chan *Node, 1000),
		Groups:         NewNodeGroup(),
	}

	go nodeStore.Run()
	return nodeStore
}

type PrivateKey struct {
	NodeId     string
	PrivateKey rsa.PrivateKey
}
