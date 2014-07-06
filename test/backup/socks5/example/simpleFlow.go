package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func main() {
	s := Server{}
	s.StartUp(34567)
}

type Server struct {
}

func (this *Server) StartUp(port int) {
	addrPort := ":" + strconv.Itoa(port)
	l, _ := net.Listen("tcp", addrPort)
	for {
		conn, _ := l.Accept()
		go this.handler(conn)
	}
}

func (this *Server) handler(conn net.Conn) {
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("服务器收到：", buf[:n])
	conn.Write([]byte{5, 0})
	n, _ = conn.Read(buf)
	fmt.Println("服务器收到：", buf[:n])
	portBuf := bytes.NewBuffer(buf[8:10])
	var dstPort uint16
	binary.Read(portBuf, binary.BigEndian, &dstPort)
	fmt.Println("版本", buf[:1], "UDP连接类型", buf[1:2], "保留字段", buf[2:3],
		"目标地址类型", buf[3:4], "目标地址", buf[4:8], "端口", int32(dstPort))
	//版本，成功，保留，绑定地址类型，绑定地址，绑定端口
	conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0})

}
