package addr_manager

import (
	"fmt"
	"github.com/prestonTao/mandela/core/config"
	"github.com/prestonTao/mandela/core/utils"
	"log"
	"net"
	"strconv"
)

const (
	broadcastStartPort = 8980
)

func init() {
	//startBroadcastServer()
}

/*
	启动一个局域网广播服务器
*/
func startBroadcastServer() {
	utils.Log.Debug("开始启动局域网广播服务器")
	addr, err := net.ResolveUDPAddr("udp", "192.168.1.128:9981")
	if err != nil {
		log.Panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
		// utils.Log.Debug("开始启动局域网广播服务器")
	}
	go func() {
		for {
			if one, err := GetSuperAddrOne(); err != nil {
				for i := 0; i < 10; i++ {
					udpaddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(broadcastStartPort+i))
					if err != nil {
						continue
					}
					_, err = conn.WriteToUDP([]byte(one), udpaddr)
					if err != nil {
						continue
					}
				}
			}
		}
	}()
}

/*
	通过组播方式获取地址列表
*/
func LoadByMulticast() {

	LoadByBroadcast()
}

/*
	通过广播获取地址
*/
func LoadByBroadcast() {
	utils.Log.Debug("通过局域网广播获得超级节点地址")
	count := 10
	for i := 0; i < count; i++ {
		addr, err := net.ResolveUDPAddr("udp", config.Init_LocalIP+":"+strconv.Itoa(broadcastStartPort+i))
		if err != nil {
			log.Panic(err)
		}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			// log.Panic(err)
			count++
			continue
		}
		var b [512]byte
		go func() {
			for {
				n, _, err := conn.ReadFromUDP(b[:])
				if err != nil {
					log.Panic(err)
				}
				if n != 0 {
					fmt.Printf("%s\n", b[0:n])
				}
			}
		}()

	}

}
