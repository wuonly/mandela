package app

import (
	"fmt"
	"math/big"
)

//nodeOne节点和nodeTwo节点的距离
func NodeDiscern(nodeOne, nodeTwo string) string {
	x := new(big.Int).SetBytes([]byte(nodeOne))
	y := new(big.Int).SetBytes([]byte(nodeTwo))
	z := new(big.Int).Xor(x, y)
	return z.String()
}

//根据一个节点id，算出所有网络逻辑节点id
func GetNodes(nodeId string) []string {
	return nil
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