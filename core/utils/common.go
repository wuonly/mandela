package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
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

/*
	获取本机能联网的ip地址
	@return    string    获得的ip地址
	@return    bool      是否能联网
*/
func GetLocalIntenetIp() (string, bool) {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		log.Println(err.Error())
		return "", false
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip, true
}

/*
	不联网的情况下，得到本机ip地址
*/
func GetLocalHost() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for i, one := range addrs {
		fmt.Println(i, one)
	}
	return addrs[0].String()
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
	获得一个可用的UDP端口
*/
func GetAvailablePortForUDP() int {
	startPort := 9981
	for i := 0; i < 1000; i++ {
		_, err := net.ListenPacket("udp", "127.0.0.1:"+strconv.Itoa(startPort))
		if err != nil {
			startPort = startPort + 1
		} else {
			return startPort
		}
	}
	return 0
}

/*
	获得一个可用的TCP端口
*/
func GetAvailablePortForTCP() int {
	startPort := 9981
	for i := 0; i < 1000; i++ {
		_, err := net.ListenPacket("tcp", "127.0.0.1:"+strconv.Itoa(startPort))
		if err != nil {
			startPort = startPort + 1
		} else {
			return startPort
		}
	}
	return 0
}

/*
	随机获取一个域名
*/
func GetRandomDomain() string {
	str := "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(int64(time.Now().Nanosecond()))
	result := ""
	r := 0
	for i := 0; i < 8; i++ {
		r = rand.Intn(25)
		result = result + str[r:r+1]
	}
	return result
}
