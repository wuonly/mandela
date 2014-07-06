package main

import (
	// "fmt"
	// "io"
	"log"
	"net"
	"time"
	// "net/http"
)

func main() {
	simple1()
}

func simple1() {
	conn, _ := net.Dial("tcp", "127.0.0.1:1991")
	time.Sleep(time.Second * 1)
	conn.Write([]byte("nihao wo shi client"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	log.Println(string(buf[:n]))

}
