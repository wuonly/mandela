package nodeStore

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

var (
	lock           *sync.Mutex      = new(sync.Mutex)           //锁
	Root           *Node                                        //
	isNew          bool                                         //是否是新节点
	nodes          map[string]*Node = make(map[string]*Node, 0) //id符串为键
	consistentHash *ConsistentHash  = new(ConsistentHash)       //一致性hash表
	InNodes        chan *Node       = make(chan *Node, 1000)    //需要更新的节点
	OutFindNode    chan string      = make(chan string, 1000)   //需要查询是否在线的节点
	Groups         *NodeGroup       = NewNodeGroup()            //组
	NodeIdLevel    uint             = 256                       //节点id长度
	MaxRecentCount int              = 2                         //存放相邻节点个数(左半边个数或者右半边个数)
	Proxys         map[string]*Node = make(map[string]*Node, 0) //被代理的节点，id字符串为键
	SuperName      string                                       //超级节点名称
)

//超级节点之间查询的间隔时间
var SpacingInterval time.Duration = time.Second * 30

//id字符串格式为16进制字符串
var IdStrBit int = 16

func InitNodeStore(node *Node) {
	//节点长度为512,深度为513
	Root = node
	Run()
}

/*
	定期检查所有节点状态
	一个小时查询所有应该有的节点
	5分钟清理一次已经不在线的节点
*/
func Run() {
	go recv()
	go func() {
		//查询自己节点和邻居节点
		for {
			//向网络中查找自己
			OutFindNode <- Root.IdInfo.GetId()
			time.Sleep(SpacingInterval)
		}
	}()
	go func() {
		//查询和自己相关的逻辑节点
		for {
			for _, idOne := range getNodeNetworkNum() {
				if idOne.Cmp(big.NewInt(0)) == 0 {
					OutFindNode <- "0"
				} else {
					OutFindNode <- hex.EncodeToString(idOne.Bytes())
				}
				time.Sleep(time.Second * 5)
			}

		}
	}()
}

/*
	需要更新的节点
*/
func recv() {
	for node := range InNodes {
		AddNode(node)
	}
}

//定期检查所有节点状态
// func (*NodeStore) checkSelf() {

// }

//添加一个被代理的节点
func AddProxyNode(node *Node) {
	Proxys[node.IdInfo.GetId()] = node
}

//得到一个被代理的节点
func GetProxyNode(id string) (node *Node, ok bool) {
	node, ok = Proxys[id]
	return
}

//删除一个被代理的节点
func DelProxyNode(id string) {
	delete(Proxys, id)
}

/*
	添加一个节点
	不保存本节点
*/
func AddNode(node *Node) {
	//是本身节点不添加
	if node.IdInfo.GetId() == Root.IdInfo.GetId() {
		return
	}
	node.LastContactTimestamp = time.Now()
	nodes[node.IdInfo.GetId()] = node
	consistentHash.add(node.IdInfo.GetBigIntId())
	//	fmt.Println("add node ", node.IdInfo.GetId())
}

/*
	删除一个节点
*/
func DelNode(idStr string) {
	idBitInt, _ := new(big.Int).SetString(idStr, IdStrBit)
	consistentHash.del(idBitInt)
	// recentNode.Del(node.NodeId)
	delete(nodes, idStr)
	delete(Proxys, idStr)
}

/*
	根据节点id得到一个距离最短节点的信息，不包括代理节点
	@nodeId         要查找的节点
	@includeSelf    是否包括自己
	@outId          排除一个节点
	@return         查找到的节点id，可能为空
*/
func Get(nodeId string, includeSelf bool, outId string) *Node {
	nodeIdInt, b := new(big.Int).SetString(nodeId, IdStrBit)
	if !b {
		fmt.Println("节点id格式不正确，应该为十六进制字符串:")
		fmt.Println(nodeId)
		return nil
	}

	tempHash := NewHash()
	if includeSelf {
		tempHash.add(Root.IdInfo.GetBigIntId())
	}
	for key, value := range GetAllNodes() {
		if outId != "" && key == outId {
			continue
		}
		tempHash.add(value.IdInfo.GetBigIntId())
	}
	targetId := tempHash.get(nodeIdInt)

	if targetId == nil {
		return nil
	}
	if hex.EncodeToString(targetId.Bytes()) == ParseId(GetRootIdInfoString()) {
		return Root
	}
	return nodes[hex.EncodeToString(targetId.Bytes())]
}

/*
	得到一个距离最短的节点信息，包括代理节点
	@nodeId         要查找的节点
	@includeSelf    是否包括自己
	@outId          排除一个节点
	@return         查找到的节点id，可能为空
*/
func GetInAll(nodeId string, includeSelf bool, outId string) *Node {
	nodeIdInt, b := new(big.Int).SetString(nodeId, IdStrBit)
	if !b {
		fmt.Println("节点id格式不正确，应该为十六进制字符串:")
		fmt.Println(nodeId)
		return nil
	}
	tempHash := NewHash()
	if includeSelf {
		tempHash.add(Root.IdInfo.GetBigIntId())
	}
	for key, value := range GetAllNodes() {
		if outId != "" && key == outId {
			continue
		}
		tempHash.add(value.IdInfo.GetBigIntId())
	}
	for key, value := range Proxys {
		if outId != "" && key == outId {
			continue
		}
		tempHash.add(value.IdInfo.GetBigIntId())
	}
	targetId := tempHash.get(nodeIdInt)

	if targetId == nil {
		return nil
	}
	if hex.EncodeToString(targetId.Bytes()) == ParseId(GetRootIdInfoString()) {
		return Root
	}
	if node, ok := nodes[hex.EncodeToString(targetId.Bytes())]; ok {
		return node
	}
	return Proxys[hex.EncodeToString(targetId.Bytes())]
}

//得到左邻节点
//@id         要查询的节点id
//@count      查询的id数量
func GetLeftNode(id big.Int, count int) []*Node {
	ids := consistentHash.getLeftLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, one := range ids {
		temp = append(temp, nodes[hex.EncodeToString(one.Bytes())])
	}
	return temp
}

//得到右邻节点
//@id         要查询的节点id
//@count      查询的id数量
func GetRightNode(id big.Int, count int) []*Node {
	ids := consistentHash.getRightLow(&id, count)
	if ids == nil {
		return nil
	}
	temp := make([]*Node, 0)
	for _, one := range ids {
		temp = append(temp, nodes[hex.EncodeToString(one.Bytes())])
	}
	return temp
}

//得到所有的节点，不包括本节点，也不包括代理节点
func GetAllNodes() map[string]*Node {
	return nodes
}

//检查节点是否是本节点的逻辑节点
func CheckNeedNode(nodeId string) (isNeed bool, replace string) {
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
	if len(GetAllNodes()) == 0 || len(GetAllNodes()) < MaxRecentCount*2 {
		return true, ""
	}
	consHash := NewHash()
	for _, value := range GetAllNodes() {
		consHash.add(value.IdInfo.GetBigIntId())
	}
	targetId := consHash.get(nodeIdInt)
	consHash = NewHash()
	for _, value := range getNodeNetworkNum() {
		consHash.add(value)
	}
	//在同一个网络
	if consHash.get(targetId).Cmp(consHash.get(nodeIdInt)) == 0 {
		switch targetId.Cmp(nodeIdInt) {
		case 0:
			return false, ""
		case -1:
			// return false, ""
		case 1:
			for _, idOne := range consistentHash.getLeftLow(Root.IdInfo.GetBigIntId(), MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			for _, idOne := range consistentHash.getRightLow(Root.IdInfo.GetBigIntId(), MaxRecentCount) {
				if idOne.Cmp(targetId) == 0 {
					return true, ""
				}
			}
			return true, hex.EncodeToString(targetId.Bytes())
		}
		//判断是否是左边最近的临近节点
		ids := consistentHash.getLeftLow(Root.IdInfo.GetBigIntId(), MaxRecentCount)
		distanceB := new(big.Int).Xor(nodeIdInt, Root.IdInfo.GetBigIntId())
		for _, idOne := range ids {
			// fmt.Println("左边最邻近的节点：", hex.EncodeToString(idOne.Bytes()))
			distanceA := new(big.Int).Xor(idOne, Root.IdInfo.GetBigIntId())
			if distanceA.Cmp(distanceB) == 1 {
				return true, hex.EncodeToString(idOne.Bytes())
			}
		}
		//判断是否是右边最近的临近节点
		ids = consistentHash.getRightLow(Root.IdInfo.GetBigIntId(), MaxRecentCount)
		for _, idOne := range ids {
			// fmt.Println("右边最邻近的节点：", hex.EncodeToString(idOne.Bytes()))
			distanceA := new(big.Int).Xor(idOne, Root.IdInfo.GetBigIntId())
			if distanceA.Cmp(distanceB) == 1 {
				return true, hex.EncodeToString(idOne.Bytes())
			}
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
func GetRootIdInfoString() string {
	return string(Root.IdInfo.Build())
}

//得到本节点id十六进制字符串
// func  GetRootId() string {
// 	return Root.IdInfo.GetId()
// 	// return hex.EncodeToString(Root.NodeId.Bytes())
// }

//得到每个节点网络的网络号，不包括本节点
func getNodeNetworkNum() map[string]*big.Int {
	// rootInt, _ := new(big.Int).SetString(, IdStrBit)
	networkNums := make(map[string]*big.Int, 3000)
	for i := 0; i < int(NodeIdLevel); i++ {
		//---------------------------------
		//将后面的i位置零
		//---------------------------------
		startInt := new(big.Int).Lsh(new(big.Int).Rsh(Root.IdInfo.GetBigIntId(), uint(i)), uint(i))
		//---------------------------------
		//最后一位取反
		//---------------------------------
		networkNum := new(big.Int).Xor(startInt, new(big.Int).Lsh(big.NewInt(1), uint(i)))
		networkNums[hex.EncodeToString(networkNum.Bytes())] = networkNum
	}
	return networkNums
}
