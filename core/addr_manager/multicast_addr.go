package addr_manager

import (
	"fmt"
	"github.com/prestonTao/mandela/core/config"
	"github.com/prestonTao/mandela/core/utils"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	broadcastStartPort  = 8980
	broadcastServerPort = 9981 //广播服务器起始端口号
)

var (
	broadcastClientIsStart = false
)

func init() {
	startBroadcastServer()
}

/*
	启动一个局域网广播服务器
*/
func startBroadcastServer() {
	utils.Log.Debug("开始启动局域网广播服务器")
	var conn *net.UDPConn
	var err error
	count := 10
	for i := 0; i < count; i++ {
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp", config.Init_LocalIP+":"+strconv.Itoa(broadcastServerPort+i))
		if err != nil {
			// log.Panic(err)
			continue
		}
		fmt.Println(addr)
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			// log.Panic(err)
			// utils.Log.Debug("开始启动局域网广播服务器")
			continue
		} else {
			break
		}
	}
	if err != nil {
		log.Panic("广播服务器启动失败")
		return
	}

	go func() {
		for {
			fmt.Println("111111111111111")
			time.Sleep(time.Second * 3)
			// if len(Sys_superNodeEntry) == 0 {
			// 	continue
			// }
			if ip, port, err := GetSuperAddrOne(true); err == nil {
				fmt.Println("22222222")
				for i := 0; i < 10; i++ {
					// fmt.P "255.255.255.255:" + strconv.Itoa(broadcastStartPort+i)
					udpaddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(broadcastStartPort+i))
					if err != nil {
						fmt.Println("失败")
						continue
					}
					_, err = conn.WriteToUDP([]byte(ip+":"+strconv.Itoa(port)), udpaddr)
					if err != nil {
						fmt.Println("广播失败")
						continue
					} else {
						fmt.Println("广播成功")
					}
				}
			} else {
				fmt.Println("33333333333")
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
	if broadcastClientIsStart {
		utils.Log.Debug("局域网广播客户端正在运行")
		return
	}
	utils.Log.Debug("正在启动局域网广播客户端")
	conns := make([]*net.UDPConn, 0)
	//开始启动监听
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
		conns = append(conns, conn)

		var b [512]byte
		go func() {
			for {
				n, _, err := conn.ReadFromUDP(b[:])
				if err != nil {
					// log.Panic(err)
					return
				}
				if n != 0 {
					fmt.Printf("%s\n", b[0:n])
				}
			}
		}()
	}
	//启动失败
	if len(conns) == 0 {
		return
	}
	broadcastClientIsStart = true
	go func() {
		c := make(chan string, 1)
		AddSubscribe(c)
		<-c
		utils.Log.Debug("开始关闭局域网广播客户端")
		broadcastClientIsStart = false
		for _, one := range conns {
			one.Close()
		}
	}()

}
