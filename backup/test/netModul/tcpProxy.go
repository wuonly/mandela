package main

import (
	// "fmt"
	// "io"
	"log"
	"net"
	// "net/http"
)

func main() {
	c := Connection{}
	c.startup()
}

type Connection struct {
	conn  net.Conn
	pconn net.Conn
}

func (this *Connection) startup() {
	ln, _ := net.Listen("tcp", ":1991")
	this.conn, _ = ln.Accept()
	log.Println("一个客户端连接")
	this.pconn, _ = net.Dial("tcp", "127.0.0.1:1990")
	log.Println("代理连接上了服务器")
	go this.read()
	go this.write()
}

func (this *Connection) read() {
	for {
		log.Println("代理等待读操作")
		buf := make([]byte, 1024)
		n, _ := this.conn.Read(buf)
		log.Println(string(buf[:n]))
		this.pconn.Write(buf[:n])
	}
}

func (this *Connection) write() {
	for {
		buf := make([]byte, 1024)
		n, _ := this.pconn.Read(buf)
		log.Println(string(buf[:n]))
		this.conn.Write(buf[:n])
	}
}
