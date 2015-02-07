package nodeStore

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type NodeManager struct {
	lock           *sync.Mutex      //锁
	Root           *Node            //
	isNew          bool             //是否是新节点
	nodes          map[string]*Node //id符串为键
	consistentHash *ConsistentHash  //一致性hash表
	InNodes        chan *Node       //需要更新的节点
	OutFindNode    chan string      //需要查询是否在线的节点
	Groups         *NodeGroup       //组
	NodeIdLevel    uint             //节点id长度
	MaxRecentCount int              //最多存放多少个邻居节点
	Proxys         map[string]*Node //被代理的节点，id字符串为键
	SuperName      string           //超级节点名称
	// OutRecentNode  chan *Node       //需要查询相邻节点
	// recentNode     *RecentNode      //
	// OverTime       time.Duration    `1 * 60 * 60` //超时时间，单位为秒
	// SelectTime     time.Duration    `5 * 60`      //查询时间，单位为秒
}

/*
	定期检查所有节点状态
	一个小时查询所有应该有的节点
	5分钟清理一次已经不在线的节点
*/
func (this *NodeManager) Run() {
	go this.recv()
	for {
		for _, idOne := range this.getNodeNetworkNum() {
			if idOne.Cmp(big.NewInt(0)) == 0 {
				this.OutFindNode <- "0"
			} else {
				this.OutFindNode <- hex.EncodeToString(idOne.Bytes())
			}
		}
		//向网络中查找自己
		this.OutFindNode <- this.Root.IdInfo.GetId()
		time.Sleep(SpacingInterval)
	}
}

/*
	需要更新的节点
*/
func (this *NodeManager) recv() {
	for node := range this.InNodes {
		this.AddNode(node)
	}
}

//定期检查所有节点状态
// func (this *NodeStore) checkSelf() {

// }

//添加一个被代理的节点
func (this *NodeManager) AddProxyNode(node *Node) {
	this.Proxys[node.IdInfo.GetId()] = node
}

//得到一个被代理的节点
func (this *NodeManager) GetProxyNode(id string) (node *Node, ok bool) {
	node, ok = this.Proxys[id]
	return
}

//删除一个被代理的节点
func (this *NodeManager) DelProxyNode(id string) {
	delete(this.Proxys, id)
}

/*
	添加一个节点
	不保存本节点
*/
func (this *NodeManager) AddNode(node *Node) {
	//是本身节点不添加
	if node.IdInfo.GetId() == this.Root.IdInfo.GetId() {
		return
	}
	node.LastContactTimestamp = time.Now()
	this.nodes[node.IdInfo.GetId()] = node
	this.consistentHash.Add(node.IdInfo.GetBigIntId())
}

/*
	删除一个节点
*/
func (this *NodeManager) DelNode(idStr string) {
	idBitInt, _ := new(big.Int).SetString(idStr, IdStrBit)
	this.consistentHash.Del(idBitInt)
	// this.recentNode.Del(node.NodeId)
	delete(this.nodes, idStr)
	delete(this.Proxys, idStr)
}

/*
	根据节点id得到一个节点的信息，id为十进制字符串
	@nodeId         要查找的节点
	@includeSelf    是否包括自己
	@outId          排除一个节点
	@return         查找到的节点id，可能为空
*/
func (this *NodeManager) Get(nodeId string, includeSelf bool, outId string) *Node {
	nodeIdInt, b := new(big.Int).SetString(nodeId, IdStrBit)
	if !b {
		fmt.Println("节点id格式不正确，应该为十六进制字符串:")
		fmt.Println(nodeId)
		return nil
	}

	consistentHash := NewHash()
	if includeSelf {
		consistentHash.Add(this.Root.IdInfo.GetBigIntId())
	}
	for key, value := range this.GetAllNodes() {
		if outId != "" && key == outId {
			continue
		}
		consistentHash.Add(value.IdInfo.GetBigIntId())
	}
	targetId := consistentHash.Get(nodeIdInt)

	if targetId == nil {
		return nil
	}
	if hex.EncodeToString(targetId.Bytes()) == ParseId(this.GetRootIdInfoString()) {
		return this.Root
	}
	return this.nodes[hex.EncodeToString(targetId.Bytes())]
}

//得到左邻节点
//@id         要查询的节点id
//@count      查询的id数量
func (this *NodeManager) GetLeftNode(id big.Int, count int) []*Node {
	ids := this.consistentHash.GetLeftLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, id := range ids {
		temp = append(temp, this.nodes[hex.EncodeToString(id.Bytes())])
	}
	return temp
}

//得到右邻节点
//@id         要查询的节点id
//@count      查询的id数量
func (this *NodeManager) GetRightNode(id big.Int, count int) []*Node {
	ids := this.consistentHash.GetRightLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, id := range ids {
		temp = append(temp, this.nodes[hex.EncodeToString(id.Bytes())])
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
	nodeIdInt, b := new(big.Int).SetString(nodeId, IdStrBit)
	if !b {
		fmt.Println("节点id格式不正确，应该为十六进制字符串:")
		fmt.Println(nodeId)
		return
	}
	if len(this.GetAllNodes()) == 0 {
		return true, ""
	}
	consHash := NewHash()
	for _, value := range this.GetAllNodes() {
		consHash.Add(value.IdInfo.GetBigIntId())
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
			for _, idOne := range this.consistentHash.GetLeftLow(this.Root.IdInfo.GetBigIntId(), this.MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			for _, idOne := range this.consistentHash.GetRightLow(this.Root.IdInfo.GetBigIntId(), this.MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			return true, hex.EncodeToString(targetId.Bytes())
		}
		//判断是否是左边最近的临近节点
		id := this.consistentHash.GetLeftLow(this.Root.IdInfo.GetBigIntId(), 1)[0]
		distanceA := new(big.Int).Xor(id, this.Root.IdInfo.GetBigIntId())
		distanceB := new(big.Int).Xor(nodeIdInt, this.Root.IdInfo.GetBigIntId())
		if distanceA.Cmp(distanceB) == 1 {
			return true, hex.EncodeToString(id.Bytes())
		}
		//判断是否是右边最近的临近节点
		id = this.consistentHash.GetRightLow(this.Root.IdInfo.GetBigIntId(), 1)[0]
		distanceA = new(big.Int).Xor(id, this.Root.IdInfo.GetBigIntId())
		if distanceA.Cmp(distanceB) == 1 {
			return true, hex.EncodeToString(id.Bytes())
		}
		return false, ""
	} else {
		//不在同一个网络
		return true, ""
	}
}

/*
	得到IdInfo字符串
*/
func (this *NodeManager) GetRootIdInfoString() string {
	return string(this.Root.IdInfo.Build())
}

//得到本节点id十六进制字符串
// func (this *NodeManager) GetRootId() string {
// 	return this.Root.IdInfo.GetId()
// 	// return hex.EncodeToString(this.Root.NodeId.Bytes())
// }

//得到每个节点网络的网络号，不包括本节点
func (this *NodeManager) getNodeNetworkNum() map[string]*big.Int {
	// rootInt, _ := new(big.Int).SetString(, IdStrBit)
	networkNums := make(map[string]*big.Int, 3000)
	for i := 0; i < int(this.NodeIdLevel); i++ {
		//---------------------------------
		//将后面的i位置零
		//---------------------------------
		startInt := new(big.Int).Lsh(new(big.Int).Rsh(this.Root.IdInfo.GetBigIntId(), uint(i)), uint(i))
		//---------------------------------
		//最后一位取反
		//---------------------------------
		networkNum := new(big.Int).Xor(startInt, new(big.Int).Lsh(big.NewInt(1), uint(i)))
		networkNums[hex.EncodeToString(networkNum.Bytes())] = networkNum
	}
	return networkNums
}

func NewNodeManager(node *Node) *NodeManager {
	//节点长度为512,深度为513
	nodeManager := &NodeManager{
		lock:           new(sync.Mutex),
		Root:           node,
		consistentHash: new(ConsistentHash),
		nodes:          make(map[string]*Node, 0),
		OutFindNode:    make(chan string, 1000),
		InNodes:        make(chan *Node, 1000),
		Groups:         NewNodeGroup(),
		NodeIdLevel:    NodeIdLevel,
		MaxRecentCount: MaxRecentCount,
		Proxys:         make(map[string]*Node, 0),
	}

	go nodeManager.Run()
	return nodeManager
}
