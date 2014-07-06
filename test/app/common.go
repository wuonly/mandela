package app

import (
	"log"
	"net"
	"strings"
)

const (
	FindNode = 1 //查找节点
)

//获取本机能联网的ip地址
func getLocalIntenetIp() string {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}
