package dao

import (
	"fmt"
	"math/big"
	"testing"
)

func TestInsertNode(t *testing.T) {
	rootNode := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39f")).String(),
		Addr: "100.45.6.37", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}

	insetNode := Node{NodeId: new(big.Int).SetBytes([]byte("r75y1e620b42e8fbec37edf4bfc765d490f326137740a51837184f61b8aae39e")).String(),
		Addr: "110.45.6.37", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNodeOne := Node{NodeId: new(big.Int).SetBytes([]byte("s75y1e620b42e8fbec37edf4bfc765d490f326137740a51837184f61b8aae39e")).String(),
		Addr: "110.45.6.38", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNodeTwo := Node{NodeId: new(big.Int).SetBytes([]byte("t75y1e620b42e8fbec37edf4bfc765d490f326137740a51837184f61b8aae39e")).String(),
		Addr: "110.45.6.39", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNodeThree := Node{NodeId: new(big.Int).SetBytes([]byte("u75y1e620b42e8fbec37edf4bfc765d490f326137740a51837184f61b8aae39e")).String(),
		Addr: "110.45.6.40", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNode4 := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39e")).String(),
		Addr: "110.45.6.41", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNode5 := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39d")).String(),
		Addr: "110.45.6.42", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNode6 := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39c")).String(),
		Addr: "110.45.6.43", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	insetNode7 := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39b")).String(),
		Addr: "110.45.6.40", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}

	// treeBucket.InsertNode(nodeId, insetNode)
	// bucket := TreeBucket{TreeBucket: make([]Node, MaxSize)}

	treeBucket := NewTreeBucket(rootNode, 512)

	treeBucket.InsertNode(insetNode)
	treeBucket.InsertNode(insetNodeOne)
	treeBucket.InsertNode(insetNodeTwo)
	treeBucket.InsertNode(insetNodeThree)
	treeBucket.InsertNode(insetNode4)
	treeBucket.InsertNode(insetNode5)
	// treeBucket.InsertNode(nodeId, insetNode6)
	// treeBucket.InsertNode(nodeId, insetNode7)
	// fmt.Println("获得根节点右子节点 \r\n")
	// fmt.Println(treeBucket.GetRight().GetRight().GetLeft().GetLeft().GetRight().GetRight().GetLeft(), "\r\n")
	// Print("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39e")

	node := treeBucket.FindRecentNode(new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39c")).String())
	fmt.Println("查找到的节点id为：", node.NodeId, "查找到的节点ip地址为：", node.Addr)

	hashRoot := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39f")).String(),
		Addr: "100.45.6.37", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}
	updateNode := Node{NodeId: new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39e")).String(),
		Addr: "100.45.6.37", PortFromScheme: map[string]int{"UDP": 1990, "TCP": 1990}}

	hashStore := NewHashNodeStore(hashRoot)
	hashStore.UpdeteNode(updateNode)
	hashStore.UpdeteNode(insetNode)
	hashStore.UpdeteNode(insetNodeOne)
	hashStore.UpdeteNode(insetNodeTwo)
	hashStore.UpdeteNode(insetNodeThree)
	hashStore.UpdeteNode(insetNode4)
	hashStore.UpdeteNode(insetNode5)
	hashStore.UpdeteNode(insetNode6)
	hashStore.UpdeteNode(insetNode7)
	findNode := hashStore.FindNode(insetNode7.NodeId)
	fmt.Println("查找到的节点id为：", findNode.NodeId)
	fmt.Println("\r\n++++++++++++++++++++++++++++++++++++++++++++\r\n")
	fmt.Println("根节点id二进制为：\r\n")
	Print(new(big.Int).SetBytes([]byte("27511e620b42e8fbec37edf4bfc765d490f326137a40a51837184f61b8aae39f")))
	treeBucket.GetAllBucketNodeId()

}

func TestGetAllBucketNodeId(t *testing.T) {
	//GetAllBucketNodeId()
}

func Print(findInt *big.Int) {
	fmt.Println("==================================\r\n")
	bi := ""

	// findInt := new(big.Int).SetBytes([]byte(nodeId))
	lenght := findInt.BitLen()
	for i := 0; i < lenght; i++ {
		tempInt := findInt
		findInt = new(big.Int).Div(tempInt, big.NewInt(2))
		mod := new(big.Int).Mod(tempInt, big.NewInt(2))
		bi = mod.String() + bi
	}
	fmt.Println(bi, "\r\n")
	fmt.Println("==================================\r\n")
}
