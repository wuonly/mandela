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
	// "sort"
	"time"
)

type NodeManager struct {
	Root           *Node            //
	isNew          bool             //是否是新节点
	nodes          map[string]*Node //十进制字符串为键
	consistentHash *ConsistentHash  //一致性hash表
	InNodes        chan *Node       //需要更新的节点
	OutFindNode    chan *Node       //需要查询是否在线的节点
	Groups         *NodeGroup       //组
	NodeIdLevel    int              //节点id长度
	MaxRecentCount int              //最多多少个邻居节点
	Proxys         map[string]*Node //被代理的节点，十进制字符串为键
	// OutRecentNode  chan *Node       //需要查询相邻节点
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
		this.OutFindNode <- &Node{NodeId: this.Root.NodeId}
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

//添加一个被代理的节点
func (this *NodeManager) AddProxyNode(node *Node) {
	this.Proxys[node.NodeId.String()] = node
}

//得到一个被代理的节点
func (this *NodeManager) GetProxyNode(id string) (*Node, bool) {
	return this.Proxys[id]
}

//删除一个被代理的节点
func (this *NodeManager) DelProxyNode(id string) {
	delete(this.Proxys, id)
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
		consistentHash.Add(this.Root.NodeId)
	}
	for key, value := range this.GetAllNodes() {
		if outId != "" && key == outId {
			continue
		}
		consistentHash.Add(value.NodeId)
	}
	targetId := consistentHash.Get(nodeIdInt)

	if targetId == nil {
		return nil
	}
	if targetId.String() == this.GetRootId() {
		return this.Root
	}
	return this.nodes[targetId.String()]
}

//得到左邻节点
func (this *NodeManager) GetLeftNode(id big.Int, count int) []*Node {
	ids := this.consistentHash.GetLeftLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, id := range ids {
		temp = append(temp, this.nodes[id.String()])
	}
	return temp
}

//得到右邻节点
func (this *NodeManager) GetRightNode(id big.Int, count int) []*Node {
	ids := this.consistentHash.GetRightLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, id := range ids {
		temp = append(temp, this.nodes[id.String()])
	}
	return temp
}

//得到所有的节点，不包括本节点
func (this *NodeManager) GetAllNodes() map[string]*Node {
	return this.nodes
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
			return false, ""
		case -1:
			// return false, ""
		case 1:
			for _, idOne := range this.consistentHash.GetLeftLow(this.Root.NodeId, this.MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			for _, idOne := range this.consistentHash.GetRightLow(this.Root.NodeId, this.MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			return true, targetId.String()
		}
		//判断是否是左边最近的临近节点
		id := this.consistentHash.GetLeftLow(this.Root.NodeId, 1)[0]
		distanceA := new(big.Int).Xor(id, this.Root.NodeId)
		distanceB := new(big.Int).Xor(nodeIdInt, this.Root.NodeId)
		if distanceA.Cmp(distanceB) == 1 {
			return true, id.String()
		}
		//判断是否是右边最近的临近节点
		id = this.consistentHash.GetRightLow(this.Root.NodeId, 1)[0]
		distanceA = new(big.Int).Xor(id, this.Root.NodeId)
		if distanceA.Cmp(distanceB) == 1 {
			return true, id.String()
		}
		return false, ""
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
		networkNums[networkNum.String()] = networkNum
	}
	return networkNums
}

func NewNodeManager(node *Node, bits int) *NodeManager {
	//节点长度为512,深度为513
	nodeManager := &NodeManager{
		Root:           node,
		consistentHash: new(ConsistentHash),
		nodes:          make(map[string]*Node, 0),
		OutFindNode:    make(chan *Node, 1000),
		InNodes:        make(chan *Node, 1000),
		Groups:         NewNodeGroup(),
		NodeIdLevel:    bits,
		MaxRecentCount: 2,
		Proxys:         make(map[string]*Node, 0),
	}

	go nodeManager.Run()
	return nodeManager
}
