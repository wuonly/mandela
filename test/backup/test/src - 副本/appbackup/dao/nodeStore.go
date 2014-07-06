package dao

import (
// "fmt"
)

type NodeStore struct {
	kademlia       *Bucket
	consistentHash *ConsistentHash
}

func (this *NodeStore) FindNode() {}

func NewNodeStore(node Node) *NodeStore {
	//节点长度为512
	Bucket := Bucket{Bucket: []Node{node}, level: 513}
	ConsistentHash := ConsistentHash{root: node}
	return &NodeStore{kademlia: &Bucket, consistentHash: &ConsistentHash}
}
