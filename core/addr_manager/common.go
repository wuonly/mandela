package addr_manager

import (
	"net"
	"time"
)

const (
	Str_zaro          = "0000000000000000000000000000000000000000000000000000000000000000" //字符串0
	Str_maxNumber     = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //256位的最大数十六进制表示id
	Str_halfNumber    = "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //最大id的二分之一
	Str_quarterNumber = "3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //最大id的四分之一
)

/*
	检查一个地址的计算机是否在线
	@return idOnline    是否在线
*/
func CheckOnline(addr string) (isOnline bool) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

/*
	检查地址是否可用
*/
func CheckAddr() {
	/*
		先获得一个拷贝
	*/
	oldSuperPeerEntry := make(map[string]string)
	for key, value := range Sys_superNodeEntry {
		oldSuperPeerEntry[key] = value
	}
	/*
		一个地址一个地址的判断是否可用
	*/
	for key, _ := range oldSuperPeerEntry {
		if CheckOnline(key) {
			AddSuperPeerAddr(key)
		} else {
			delete(Sys_superNodeEntry, key)
		}
	}
}
