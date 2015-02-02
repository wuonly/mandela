package mandela

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
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

/*
	是全球唯一ip
*/
func IsOnlyIp(ip string) bool {
	ips := strings.Split(ip, ".")
	if ips[0] == "127" && ips[1] == "0" && ips[2] == "0" && ips[3] == "1" {
		return false
	}
	if ips[0] == "192" && ips[1] == "168" {
		return false
	}
	if ips[0] == "10" {
		return false
	}
	return true
}

/*
	检查一个地址的计算机是否在线
	@return idOnline    是否在线
*/
func CheckOnline(addr string) (isOnline bool) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
