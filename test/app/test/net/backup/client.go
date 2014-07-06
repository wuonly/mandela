package net

import (
	"net"
)

type SocketClient struct {
	Addr    string //IP地址
	Port    int    //端口
	started bool
}

func (this *SocketClient) StartUP() {
}

func (this *SocketClient) handleLisener(conn *net.Conn) {
	//nodeId有154位
}
