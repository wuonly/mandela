package nodeStore

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestNodeManager(t *testing.T) {
	// nodeManagerTest()
}

func nodeManagerTest() {
	rootId := RandNodeId(4)
	fmt.Println("本节点id为：", rootId.String())

	node := &Node{
		NodeId:  rootId,
		IsSuper: true, //是超级节点
		Addr:    "1111",
		TcpPort: 8080,
		UdpPort: 0,
	}

	nodeManager := NewNodeManager(node, 4)
	ok, _ := nodeManager.CheckNeedNode("9")
	fmt.Println(ok)

}

//得到指定长度的节点id
//@return 10进制字符串
func RandNodeId(lenght int) *big.Int {
	min := rand.New(rand.NewSource(99))
	timens := int64(time.Now().Nanosecond())
	min.Seed(timens)
	maxId := new(big.Int).Lsh(big.NewInt(1), uint(lenght))
	randInt := new(big.Int).Rand(min, maxId)
	return randInt
}
