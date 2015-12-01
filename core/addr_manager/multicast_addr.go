package addr_manager

import (
	"log"
	"strconv"
)

const (
	broadcastStartPort = 8980
)

func init() {

	addr, err := net.ResolveUDPAddr("udp", "192.168.1.128:9981")
	if err != nil {
		log.Panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
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
	for i := 0; i < 10; i++ {
		addr, err := net.ResolveUDPAddr("udp", "192.168.1.106:"+strconv.Itoa(broadcastStartPort+i))
		if err != nil {
			log.Panic(err)
		}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Panic(err)
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
