package nodeStore

import (
	"encoding/hex"
	"math/big"
	"sort"
)

type RecentNode struct {
	root     *big.Int //自己节点
	maxSize  int      //保存的数量
	preNodes IdASC    //前面的节点 保存的节点,下标越小，离自己越近
	sufNodes IdDESC   //后面的节点
}

//检查一个节点是否需要
func (this *RecentNode) CheckIn(nodeId *big.Int) (bool, string) {
	if nodeId.Cmp(this.root) == 0 {
		return false, ""
	}
	if len(append(this.preNodes, this.sufNodes...)) < this.maxSize*2 {
		return true, ""
	}
	temp := NewRecentNode(this.root, this.maxSize)
	for _, idOne := range this.GetAll() {
		temp.Add(idOne)
	}
	switch temp.root.Cmp(nodeId) {
	case 0: //和root节点相等
	case -1: //在节点前面
		temp.preNodes = append(temp.preNodes, nodeId)
		sort.Sort(temp.preNodes)
		passId := temp.preNodes[temp.maxSize]
		if passId.Cmp(nodeId) == 0 {
			return false, ""
		}
		return true, hex.EncodeToString(passId.Bytes())
	case 1: //在节点后面
		temp.sufNodes = append(temp.sufNodes, nodeId)
		sort.Sort(temp.sufNodes)
		passId := temp.sufNodes[temp.maxSize]
		if passId.Cmp(nodeId) == 0 {
			return false, ""
		}
		return true, hex.EncodeToString(passId.Bytes())
	}
	return false, ""
	// if nodeId.Cmp(this.root) == 0 {
	// 	return false
	// }
	// var allIds IdDESC = append(this.preNodes, this.sufNodes...)
	// sort.Sort(allIds)
	// if len(allIds) < this.maxSize*2 {
	// 	return true
	// }
	// switch allIds[0].Cmp(nodeId) {
	// case 0:
	// 	return false
	// case -1:
	// case 1:
	// 	return false
	// }
	// switch allIds[len(allIds)-1].Cmp(nodeId) {
	// case 0:
	// 	return false
	// case -1:
	// 	return false
	// case 1:
	// }
	// return true
}

func (this *RecentNode) Add(nodeId *big.Int) {
	switch this.root.Cmp(nodeId) {
	case 0: //和root节点相等
	case -1: //在节点前面
		this.preNodes = append(this.preNodes, nodeId)
		sort.Sort(this.preNodes)
		this.preNodes = this.preNodes[:this.maxSize]
	case 1: //在节点后面
		this.sufNodes = append(this.sufNodes, nodeId)
		sort.Sort(this.sufNodes)
		this.sufNodes = this.sufNodes[:this.maxSize]
	}

}

func (this *RecentNode) GetAll() []*big.Int {
	return append(this.preNodes, this.sufNodes...)
}

func (this *RecentNode) Del(nodeId *big.Int) {
	switch this.root.Cmp(nodeId) {
	case 0: //和root节点相等
	case -1: //在节点前面
		for i, id := range this.preNodes {
			if id.Cmp(nodeId) == 0 {
				this.preNodes = append(this.preNodes[:i], this.preNodes[i+1:]...)
			}
		}
	case 1: //在节点后面
		for i, id := range this.sufNodes {
			if id.Cmp(nodeId) == 0 {
				this.sufNodes = append(this.sufNodes[:i], this.sufNodes[i+1:]...)
			}
		}
	}
}

func NewRecentNode(nodeId *big.Int, maxSize int) *RecentNode {
	recentNode := new(RecentNode)
	recentNode.root = nodeId
	recentNode.maxSize = maxSize
	return recentNode
}

//------------------------------------------------------
// type RecentNode struct {
// 	root     Node   //自己节点
// 	maxSize  int    //最大数量
// 	preNodes []Node //前面的节点 保存的节点,下标越小，离自己越近
// 	sufNodes []Node //后面的节点
// }

// func (this *RecentNode) UpdeteNode(node Node) {
// 	rootInt := new(big.Int).SetBytes([]byte(this.root.NodeId))
// 	updateInt := new(big.Int).SetBytes([]byte(node.NodeId))
// 	qu := new(big.Int).Sub(rootInt, updateInt)
// 	//负数在本节点前面，正数在节点后面，数越大离本节点越远
// 	if qu.Cmp(big.NewInt(0)) > 0 {
// 		fmt.Println("要插入后面的节点")
// 		index := -1
// 		//要插入到后面的节点
// 		for i, node := range this.sufNodes {
// 			nodeInt := new(big.Int).SetBytes([]byte(node.NodeId))
// 			quIn := new(big.Int).Sub(nodeInt, updateInt)
// 			tempInt := quIn.Cmp(big.NewInt(0))
// 			if tempInt > 0 {
// 				break
// 			} else if tempInt == 0 {
// 				//节点存在，并且位置正确
// 				return
// 			} else {
// 				index = i
// 				break
// 			}
// 		}
// 		//插入到后节点中
// 		if index == -1 {
// 			//插入到最后
// 			if this.sufNodes == nil {
// 				this.sufNodes = []Node{node}
// 			} else {
// 				this.sufNodes = append(this.sufNodes, node)
// 			}
// 		} else {
// 			pre := this.sufNodes[:index]
// 			suf := this.sufNodes[index:]
// 			this.sufNodes = append(pre, node)
// 			for _, sufNode := range suf {
// 				this.sufNodes = append(this.sufNodes, sufNode)
// 			}
// 		}
// 	} else {
// 		fmt.Println("要插入前面的节点")
// 		//要插入到前面的节点
// 		index := -1
// 		for i, node := range this.preNodes {
// 			nodeInt := new(big.Int).SetBytes([]byte(node.NodeId))
// 			quIn := new(big.Int).Sub(nodeInt, updateInt)
// 			tempInt := quIn.Cmp(big.NewInt(0))
// 			if tempInt < 0 {
// 				break
// 			} else if tempInt == 0 {
// 				//节点存在，并且位置正确
// 				return
// 			} else {
// 				index = i
// 				break
// 			}
// 		}
// 		//插入到前节点中
// 		if index == -1 {
// 			//插入到最后
// 			if this.preNodes == nil {
// 				this.preNodes = []Node{node}
// 			} else {
// 				this.preNodes = append(this.preNodes, node)
// 			}
// 		} else {
// 			pre := this.preNodes[:index]
// 			suf := this.preNodes[index:]
// 			this.preNodes = append(pre, node)
// 			for _, sufNode := range suf {
// 				this.preNodes = append(this.preNodes, sufNode)
// 			}
// 		}
// 	}
// }

// func (this *RecentNode) FindNode(nodeId string) Node {
// 	rootInt := new(big.Int).SetBytes([]byte(this.root.NodeId))
// 	findInt := new(big.Int).SetBytes([]byte(nodeId))
// 	qu := new(big.Int).Sub(rootInt, findInt)
// 	//负数在本节点前面，正数在节点后面，数越大离本节点越远
// 	if qu.Cmp(big.NewInt(0)) > 0 {
// 		//在后面节点中查找
// 		for _, node := range this.sufNodes {
// 			if node.NodeId == nodeId {
// 				return node
// 			}
// 		}
// 	} else {
// 		//在前面节点中查找
// 		for _, node := range this.preNodes {
// 			if node.NodeId == nodeId {
// 				return node
// 			}
// 		}
// 	}
// 	return Node{}
// }

// func NewRecentNode(node Node) *RecentNode {
// 	return &RecentNode{root: node}
// }
