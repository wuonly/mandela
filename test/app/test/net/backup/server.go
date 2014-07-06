package net

import (
	"log"
	"net"
	"strconv"
)

type SocketServer struct {
	Port    int //本地socket服务器端口，和upnp映射端口一致
	started bool
}

func (this *SocketServer) StartUP() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(this.Port))
	defer listener.Close()
	if err != nil {
		log.Println(err)
	}
	this.started = true
	for this.started {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			this.started = false
			return
		}
		go this.handleLisener(conn)
	}
}

func (this *SocketServer) handleLisener(conn net.Conn) {
	//nodeId有154位
	buf := make([]byte, 160)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		//判断错误类型是否是net.OpError
		// if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
		// 	log.Println("就是这种")
		// }
	}
	//-----------------------
	//这里验证连接的合法性
	//-----------------------
	connection := connection{nodeId: string(buf[:n]), conn: conn}
	go connection.ListenerAndServer()
	// if n == 160 {
	// } else {
	// 	//不合法
	// 	conn.Close()
	// }
}
