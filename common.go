package mandela

import (
	"crypto/sha256"
	// "encoding/hex"
	"fmt"
	"log"
	"math/big"
	// "math/rand"
	"net"
	"strings"
	// "time"
)

//通过一个域名和用户名得到节点的id
//@return 10进制字符串
func GetHashKey(account string) *big.Int {
	hash := sha256.New()
	hash.Write([]byte(account))
	md := hash.Sum(nil)
	// str16 := hex.EncodeToString(md)
	// resultInt, _ := new(big.Int).SetString(str16, 16)
	resultInt := new(big.Int).SetBytes(md)
	return resultInt
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

//获取本机能联网的ip地址
func GetLocalIntenetIp() string {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/

	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}
