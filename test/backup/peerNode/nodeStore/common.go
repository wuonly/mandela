package nodeStore

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"math/big"
// 	"math/rand"
// 	"time"
// )

// //得到指定长度的节点id
// //@return 10进制字符串
// func RandNodeId(lenght int) *big.Int {
// 	min := rand.New(rand.NewSource(99))
// 	timens := int64(time.Now().Nanosecond())
// 	min.Seed(timens)
// 	maxId := new(big.Int).Lsh(big.NewInt(1), uint(lenght))
// 	randInt := new(big.Int).Rand(min, maxId)
// 	return randInt
// }

// //通过一个域名和用户名得到节点的id
// //@return 10进制字符串
// func GetHashKey(uri, account string) *big.Int {
// 	hash := sha256.New()
// 	hash.Write([]byte(uri + account))
// 	md := hash.Sum(nil)
// 	str16 := hex.EncodeToString(md)
// 	resultInt, _ := new(big.Int).SetString(str16, 16)
// 	return resultInt
// }

// func Print(findInt *big.Int) {
// 	fmt.Println("==================================\r\n")
// 	bi := ""

// 	// findInt := new(big.Int).SetBytes([]byte(nodeId))
// 	lenght := findInt.BitLen()
// 	for i := 0; i < lenght; i++ {
// 		tempInt := findInt
// 		findInt = new(big.Int).Div(tempInt, big.NewInt(2))
// 		mod := new(big.Int).Mod(tempInt, big.NewInt(2))
// 		bi = mod.String() + bi
// 	}
// 	fmt.Println(bi, "\r\n")
// 	fmt.Println("==================================\r\n")
// }
