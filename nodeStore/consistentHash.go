package nodeStore

import (
	// "fmt"
	// "log"
	"math/big"
	"sort"
	"sync"
)

type ConsistentHash struct {
	lock  sync.Mutex
	nodes IdDESC
}

//10进制字符串
func (this *ConsistentHash) Add(nodes ...*big.Int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, idOne := range nodes {
		//判断重复的
		for _, node := range this.nodes {
			if node.Cmp(idOne) == 0 {
				continue
			}
		}
		this.nodes = append(this.nodes, idOne)
	}
	sort.Sort(this.nodes)
}

func (this *ConsistentHash) Get(nodeId *big.Int) *big.Int {
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(this.nodes) == 0 {
		return nil
	}
	isFind := false
	for i, idOne := range this.nodes {
		switch nodeId.Cmp(idOne) {
		case 0:
			return idOne
		case -1:
			// fmt.Println("haha")
			isFind = true
		case 1:
			if i == 0 {
				firstDistanceInt := new(big.Int).Xor(nodeId, idOne)
				lastDistanceInt := new(big.Int).Xor(nodeId, this.nodes[len(this.nodes)-1])
				switch firstDistanceInt.Cmp(lastDistanceInt) {
				case 0:
					return idOne
				case -1:
					return idOne
				case 1:
					return this.nodes[len(this.nodes)-1]
				}
			}
			if isFind {
				startDistanceInt := new(big.Int).Xor(nodeId, this.nodes[i-1])
				lastDistanceInt := new(big.Int).Xor(nodeId, idOne)
				switch startDistanceInt.Cmp(lastDistanceInt) {
				case 0:
					return idOne
				case -1:
					return this.nodes[i-1]
				case 1:
					return idOne
				}
			}
		}
	}

	firstDistanceInt := new(big.Int).Xor(this.nodes[0], nodeId)
	lastDistanceInt := new(big.Int).Xor(nodeId, this.nodes[len(this.nodes)-1])
	switch firstDistanceInt.Cmp(lastDistanceInt) {
	case 0:
		return this.nodes[0]
	case -1:
		return this.nodes[0]
	case 1:
		return this.nodes[len(this.nodes)-1]
	}
	return nil
}

//获得左边或右边最近的节点
func (this *ConsistentHash) GetLeftLow(isLeft bool, nodeId *big.Int) *big.Int {
	// this.lock.Lock()
	// defer this.lock.Unlock()
	// if len(this.nodes) == 0 {
	// 	return nil
	// }
	// for i, idOne := range this.nodes {

	// }
}

//删除一个节点
func (this *ConsistentHash) Del(node *big.Int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	//判断重复的
	for i, nodeOne := range this.nodes {
		if nodeOne.Cmp(node) == 0 {
			this.nodes = append(this.nodes[:i], this.nodes[i+1:]...)
			return
		}
	}
}

//得到hash表中保存的所有节点
func (this *ConsistentHash) GetAll() []*big.Int {
	return this.nodes
}

//创建一个新的一致性hash表
func NewHash() *ConsistentHash {
	chash := &ConsistentHash{nodes: []*big.Int{}}
	return chash
}
