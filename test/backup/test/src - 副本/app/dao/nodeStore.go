package dao

import (
// "fmt"
)

type NodeStore struct {
	kademlia       *Bucket
	consistentHash *ConsistentHash
}

func (this *NodeStore) FindNode() {}

var nodeStore *NodeStore

func NewNodeStore(node Node) *NodeStore {
	//节点长度为512,深度为513
	bucket := Bucket{Bucket: []Node{node}, level: 513}
	consistentHash := ConsistentHash{root: node}
	nodeStore = &NodeStore{kademlia: &bucket, consistentHash: &consistentHash}
	return nodeStore
}
