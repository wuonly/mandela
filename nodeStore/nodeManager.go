package nodeStore

import (
	// "bytes"
	// "crypto/rand"
	// "crypto/rsa"
	// "encoding/gob"
	"fmt"
	// "io/ioutil"
	"math/big"
	// "os"
	"sort"
	"time"
)

type NodeManager struct {
	Root           *Node            //
	isNew          bool             //是否是新节点
	nodes          map[string]*Node //十进制字符串为键
	consistentHash *ConsistentHash  //一致性hash表
	InNodes        chan *Node       //需要更新的节点
	OutFindNode    chan *Node       //需要查询是否在线的节点
	OutRecentNode  chan *Node       //需要查询相邻节点
	Groups         *NodeGroup       //组
	NodeIdLevel    int              //节点id长度
	// recentNode     *RecentNode      //
	// OverTime       time.Duration    `1 * 60 * 60` //超时时间，单位为秒
	// SelectTime     time.Duration    `5 * 60`      //查询时间，单位为秒
}

//定期检查所有节点状态
//一个小时查询所有应该有的节点
//5分钟清理一次已经不在线的节点
func (this *NodeManager) Run() {
	go this.recv()
	//向网络中查找自己，通知相关节点自己上线了
	// this.OutFindNode <- this.Root
	for {
		for _, idOne := range this.getNodeNetworkNum() {
			this.OutFindNode <- &Node{NodeId: idOne}
		}
		this.OutRecentNode <- &Node{NodeId: this.Root.NodeId}
		//清理离线的节点
		// for _, nodeOne := range this.nodes {
		// 	if time.Now().Sub(nodeOne.LastContactTimestamp) > time.Hour {
		// 		this.DelNode(nodeOne)
		// 	}
		// }
		time.Sleep(time.Minute * 1)
	}
}

//需要更新的节点
func (this *NodeManager) recv() {
	for node := range this.InNodes {
		this.AddNode(node)
	}
}

//定期检查所有节点状态
// func (this *NodeStore) checkSelf() {

// }

//添加一个节点
//不保存本节点
func (this *NodeManager) AddNode(node *Node) {
	//是本身节点不添加
	if node.NodeId.Cmp(this.Root.NodeId) == 0 {
		// this.nodes[this.root.NodeId.String()] = this.root
		return
	}

	node.LastContactTimestamp = time.Now()
	this.nodes[node.NodeId.String()] = node
	this.consistentHash.Add(node.NodeId)
	// this.recentNode.Add(node.NodeId)
}

//删除一个节点
func (this *NodeManager) DelNode(node *Node) {
	this.consistentHash.Del(node.NodeId)
	// this.recentNode.Del(node.NodeId)
	delete(this.nodes, node.NodeId.String())
}

//根据节点id得到一个节点的信息，id为十进制字符串
//@nodeId         要查找的节点
//@includeSelf    是否包括自己
//@outId          排除一个节点
//@return         查找到底节点id，可能为空
func (this *NodeManager) Get(nodeId string, includeSelf bool, outId string) *Node {
	nodeIdInt, b := new(big.Int).SetString(nodeId, 10)
	if !b {
		fmt.Println("节点id格式不正确，应该为十进制字符串")
		return nil
	}

	consistentHash := NewHash()
	if includeSelf {
		// fmt.Println("添加根节点：", this.Root)
		consistentHash.Add(this.Root.NodeId)
	}
	for key, value := range this.GetAllNodes() {
		if outId != "" && key == outId {
			continue
		}
		consistentHash.Add(value.NodeId)
	}
	// for _, id := range this.recentNode.GetAll() {
	// 	if outId != "" && id.String() == outId {
	// 		continue
	// 	}
	// 	consistentHash.Add(id)
	// }
	targetId := consistentHash.Get(nodeIdInt)

	if targetId == nil {
		return nil
	}
	if targetId.String() == this.GetRootId() {
		return this.Root
	}
	return this.nodes[targetId.String()]
}

//得到所有的节点，不包括本节点
func (this *NodeManager) GetAllNodes() map[string]*Node {
	return this.nodes
}

//获取所有相邻节点，包括本节点
func (this *NodeManager) GetRecentNodes() []*big.Int {
	var ids IdDESC = this.recentNode.GetAll()
	ids = append(ids, this.Root.NodeId)
	sort.Sort(ids)
	return this.recentNode.GetAll()
}

//检查节点是否是本节点的逻辑节点
func (this *NodeManager) CheckNeedNode(nodeId string) (isNeed bool, replace string) {
	/*
		1.找到已有节点中与本节点最近的节点
		2.计算两个节点是否在同一个网络
		3.若在同一个网络，计算谁的值最小
	*/
	nodeIdInt, b := new(big.Int).SetString(nodeId, 10)
	if !b {
		fmt.Println("节点id格式不正确，应该为十进制字符串")
		return
	}
	if len(this.GetAllNodes()) == 0 {
		return true, ""
	}
	consHash := NewHash()
	for _, value := range this.GetAllNodes() {
		consHash.Add(value.NodeId)
	}
	targetId := consHash.Get(nodeIdInt)

	consHash = NewHash()
	for _, value := range this.getNodeNetworkNum() {
		consHash.Add(value)
	}
	//在同一个网络
	if consHash.Get(targetId).Cmp(consHash.Get(nodeIdInt)) == 0 {
		switch targetId.Cmp(nodeIdInt) {
		case 0:
			// return false, ""
		case -1:
			// return false, ""
		case 1:
			return true, targetId.String()
		}
		return this.recentNode.CheckIn(nodeIdInt)
	} else {
		//不在同一个网络
		return true, ""
	}
}

//得到本节点id十进制字符串
func (this *NodeManager) GetRootId() string {
	return this.Root.NodeId.String()
}

//得到每个节点网络的网络号，不包括本节点
func (this *NodeManager) getNodeNetworkNum() map[string]*big.Int {
	// rootInt, _ := new(big.Int).SetString(, 10)
	networkNums := make(map[string]*big.Int, 3000)
	for i := 0; i < this.NodeIdLevel; i++ {
		//---------------------------------
		//将后面的i位置零
		//---------------------------------
		startInt := new(big.Int).Lsh(new(big.Int).Rsh(this.Root.NodeId, uint(i)), uint(i))
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

func NewNodeManager(node *Node, bits int) *NodeManager {
	// node := Node{
	// 	NodeId:  nodeId,
	// 	IsSuper: IsSuper,
	// 	Addr:    Addr,
	// 	TcpPort: TcpPort,
	// 	UdpPort: UdpPort,
	// 	// Key:     privateKey,
	// }

	// fmt.Println("本次创建的nodeid为：", node.NodeId, "私钥：", node.Key)
	//节点长度为512,深度为513
	nodeManager := &NodeManager{
		Root:           node,
		consistentHash: new(ConsistentHash),
		// recentNode:     NewRecentNode(node.NodeId, 2),
		nodes:         make(map[string]*Node, 1000),
		OutFindNode:   make(chan *Node, 1000),
		OutRecentNode: make(chan *Node, 1000),
		InNodes:       make(chan *Node, 1000),
		Groups:        NewNodeGroup(),
		NodeIdLevel:   bits,
	}

	go nodeManager.Run()
	return nodeManager
}

//=====================================================

// var (
// 	NodeIdLevel       = 256         //节点id二进制长度
// 	IsSuper           = false       //是否是超级节点，并且提供代理功能
// 	Addr              = "127.0.0.1" //外网地址
// 	TcpPort     int32 = 0           //外网端口
// 	UdpPort     int32 = 0           //外网端口
// )

// type NodeStore struct {
// 	root           *Node            //
// 	isNew          bool             //是否是新节点
// 	nodes          map[string]*Node //十进制字符串为键
// 	consistentHash *ConsistentHash  //一致性hash表
// 	pipeNode       chan *Node       //收集各个节点更新，下线消息的管道
// 	InNodes        chan *Node       //需要更新的节点
// 	OutFindNode    chan *Node       //需要查询是否在线的节点
// 	Groups         *NodeGroup       //组
// }

// func (this *NodeStore) Run() {
// 	go this.recv()
// 	go this.checkSelf()
// 	go this.nodeComm()
// 	//向网络中查找自己，通知相关节点自己上线了
// 	this.OutFindNode <- this.root
// 	idsInt := this.getNodeNetworkNum()
// 	for _, idOne := range idsInt {
// 		this.OutFindNode <- &Node{NodeIdShould: idOne}
// 	}
// }

// //检查下线的节点，把他们移除掉
// func (this *NodeStore) checkSelf() {
// 	//160分钟检查一次
// 	time.Sleep(time.Minute * 160)
// 	for _, value := range this.nodes {
// 		if value.Status == 3 {
// 			this.delNode(value)
// 		}
// 	}
// 	go this.checkSelf()
// }

// //需要更新的节点
// func (this *NodeStore) recv() {
// 	for {
// 		node := <-this.InNodes
// 		this.AddNode(node)
// 	}
// }

// //负责处理仓库中节点发来的更新，下线消息
// func (this *NodeStore) nodeComm() {
// 	for {
// 		node := <-this.pipeNode
// 		switch node.Status {
// 		case 1:
// 		case 2:
// 			//需要查询是否在线的节点
// 			this.OutFindNode <- node
// 		case 3:
// 			//需要下线的节点
// 		}
// 	}
// }

// //添加一个节点
// func (this *NodeStore) AddNode(node *Node) {
// 	//添加本节点
// 	if node.NodeId.Cmp(this.root.NodeId) == 0 {
// 		this.nodes[this.root.NodeId.String()] = this.root
// 		return
// 	}
// 	node.Out = this.pipeNode
// 	node.LastContactTimestamp = time.Now()
// 	node.OverTime = 1 * 60 * 60
// 	node.SelectTime = 5 * 60
// 	go node.ticker()
// 	this.nodes[node.NodeId.String()] = node
// 	this.consistentHash.Add(node.NodeId)
// }

// //删除一个节点
// func (this *NodeStore) delNode(node *Node) {
// 	this.consistentHash.Del(node.NodeId)
// 	delete(this.nodes, node.NodeId.String())
// }

// //得到每个节点网络的网络号，不包括本节点
// func (this *NodeStore) getNodeNetworkNum() map[string]*big.Int {
// 	// rootInt, _ := new(big.Int).SetString(, 10)
// 	networkNums := make(map[string]*big.Int, 3000)
// 	for i := 0; i < NodeIdLevel; i++ {
// 		//---------------------------------
// 		//将后面的i位置零
// 		//---------------------------------
// 		startInt := new(big.Int).Lsh(new(big.Int).Rsh(this.root.NodeId, uint(i)), uint(i))
// 		//---------------------------------
// 		//最后一位取反
// 		//---------------------------------
// 		networkNum := new(big.Int).Xor(startInt, new(big.Int).Lsh(big.NewInt(1), uint(i)))
// 		// fmt.Println("haha", i)
// 		// Print(networkNum)
// 		// networkNums = append(networkNums, networkNum)
// 		networkNums[networkNum.String()] = networkNum
// 	}
// 	return networkNums
// }

// //根据节点id得到一个节点的信息，id为十进制字符串
// func (this *NodeStore) Get(nodeId string) *Node {
// 	nodeIdInt, b := new(big.Int).SetString(nodeId, 10)
// 	if !b {
// 		fmt.Println("节点id格式不正确，应该为十进制字符串")
// 	}
// 	targetNodeId := this.consistentHash.Get(nodeIdInt)
// 	if targetNodeId != nil {
// 		return this.nodes[targetNodeId.String()]
// 	}
// 	return this.root
// }

// //得到所有的节点
// func (this *NodeStore) GetAllNodes() map[string]*Node {
// 	return this.nodes
// }

// //得到本节点id
// func (this *NodeStore) GetRootId() string {
// 	return this.root.NodeId.String()
// }

// //根据域名和帐号名称创建一个节点
// //域名为p2p网络中唯一域名
// //域名加帐号名称才能够确定一个节点id，为了把域名和帐号绑定
// func NewNodeStore(nodeId *big.Int) *NodeStore {

// 	node := Node{
// 		NodeId:  nodeId,
// 		IsSuper: IsSuper,
// 		Addr:    Addr,
// 		TcpPort: TcpPort,
// 		UdpPort: UdpPort,
// 		// Key:     privateKey,
// 	}

// 	// fmt.Println("本次创建的nodeid为：", node.NodeId, "私钥：", node.Key)
// 	//节点长度为512,深度为513
// 	nodeStore := &NodeStore{
// 		root:           &node,
// 		consistentHash: new(ConsistentHash),
// 		nodes:          make(map[string]*Node, 1000),
// 		OutFindNode:    make(chan *Node, 1000),
// 		pipeNode:       make(chan *Node, 1000),
// 		InNodes:        make(chan *Node, 1000),
// 		Groups:         NewNodeGroup(),
// 	}

// 	go nodeStore.Run()
// 	return nodeStore
// }

// type PrivateKey struct {
// 	NodeId     string
// 	PrivateKey rsa.PrivateKey
// }
