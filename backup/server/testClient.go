package main

import (
	"fmt"
	"net"
)

func main() {
	StartUP()
}

func StartUP() {
	conn, _ := net.Dial("tcp4", "127.0.0.1:9981")
	conn.Write([]byte("127.0.0.1:9990"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
}
