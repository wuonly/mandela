package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	StartUp(34567)
}

func StartUp(port int) {
	conn, _ := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
	conn.Write([]byte{0x5, 1, 0})
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(buf[:n])
}
