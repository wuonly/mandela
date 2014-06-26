package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

//通过一个域名和用户名得到节点的id
//@return 16进制字符串
func GetHashKey(uri, account string) string {
	hash := sha256.New()
	hash.Write([]byte(uri + account))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

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

//得到指定长度的节点id
func RandNodeId(lenght int) string {
	min := rand.New(rand.NewSource(99))
	timens := int64(time.Now().Nanosecond())
	min.Seed(timens)
	maxId := new(big.Int).Lsh(big.NewInt(1), uint(lenght+1))
	randInt := new(big.Int).Rand(min, maxId)
	return randInt.String()
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
