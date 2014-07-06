package net

import (
	"net"
	"strconv"
)

type connection struct {
	nodeId string   //节点id
	conn   net.Conn //nodeId为键，conn为值
}

func (this *connection) ListenerAndServer() {
	for {
		size := make([]byte, 2)
		this.conn.Read(size)
		dataLen, _ := strconv.Atoi(string(size))
		data := make([]byte, dataLen)
		this.conn.Read(data)
		this.conn.Write([]byte("I'm server"))
	}
}

func (this *connection) Request(str string) {

	this.conn.Write([]byte(strconv.Itoa(len([]byte(str)))))
	this.conn.Write([]byte(str))
}
