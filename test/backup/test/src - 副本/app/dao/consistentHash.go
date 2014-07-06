package dao

import (
	"fmt"
	"math/big"
)

type ConsistentHash struct {
	root     Node   //自己节点
	maxSize  int    //最大数量
	preNodes []Node //前面的节点 保存的节点,下标越小，离自己越近
	sufNodes []Node //后面的节点
}

func (this *ConsistentHash) UpdeteNode(node Node) {
	rootInt := new(big.Int).SetBytes([]byte(this.root.NodeId))
	updateInt := new(big.Int).SetBytes([]byte(node.NodeId))
	qu := new(big.Int).Sub(rootInt, updateInt)
	//负数在本节点前面，正数在节点后面，数越大离本节点越远
	if qu.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("要插入后面的节点")
		index := -1
		//要插入到后面的节点
		for i, node := range this.sufNodes {
			nodeInt := new(big.Int).SetBytes([]byte(node.NodeId))
			quIn := new(big.Int).Sub(nodeInt, updateInt)
			tempInt := quIn.Cmp(big.NewInt(0))
			if tempInt > 0 {
				break
			} else if tempInt == 0 {
				//节点存在，并且位置正确
				return
			} else {
				index = i
				break
			}
		}
		//插入到后节点中
		if index == -1 {
			//插入到最后
			if this.sufNodes == nil {
				this.sufNodes = []Node{node}
			} else {
				this.sufNodes = append(this.sufNodes, node)
			}
		} else {
			pre := this.sufNodes[:index]
			suf := this.sufNodes[index:]
			this.sufNodes = append(pre, node)
			for _, sufNode := range suf {
				this.sufNodes = append(this.sufNodes, sufNode)
			}
		}
	} else {
		fmt.Println("要插入前面的节点")
		//要插入到前面的节点
		index := -1
		for i, node := range this.preNodes {
			nodeInt := new(big.Int).SetBytes([]byte(node.NodeId))
			quIn := new(big.Int).Sub(nodeInt, updateInt)
			tempInt := quIn.Cmp(big.NewInt(0))
			if tempInt < 0 {
				break
			} else if tempInt == 0 {
				//节点存在，并且位置正确
				return
			} else {
				index = i
				break
			}
		}
		//插入到前节点中
		if index == -1 {
			//插入到最后
			if this.preNodes == nil {
				this.preNodes = []Node{node}
			} else {
				this.preNodes = append(this.preNodes, node)
			}
		} else {
			pre := this.preNodes[:index]
			suf := this.preNodes[index:]
			this.preNodes = append(pre, node)
			for _, sufNode := range suf {
				this.preNodes = append(this.preNodes, sufNode)
			}
		}
	}
}

func (this *ConsistentHash) FindNode(nodeId string) Node {
	rootInt := new(big.Int).SetBytes([]byte(this.root.NodeId))
	findInt := new(big.Int).SetBytes([]byte(nodeId))
	qu := new(big.Int).Sub(rootInt, findInt)
	//负数在本节点前面，正数在节点后面，数越大离本节点越远
	if qu.Cmp(big.NewInt(0)) > 0 {
		//在后面节点中查找
		for _, node := range this.sufNodes {
			if node.NodeId == nodeId {
				return node
			}
		}
	} else {
		//在前面节点中查找
		for _, node := range this.preNodes {
			if node.NodeId == nodeId {
				return node
			}
		}
	}
	return Node{}
}

func NewConsistentHash(node Node) *ConsistentHash {
	return &ConsistentHash{root: node}
}
