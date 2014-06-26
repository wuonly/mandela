package main

import (
	// "fmt"
	// "io"
	"log"
	"net"
	// "net/http"
)

func main() {
	sample1()
}

func sample1() {
	ln, _ := net.Listen("tcp", ":1990")
	for {
		conn, _ := ln.Accept()
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		log.Println(buf[:n])

		str := "message from server"
		// conn.Write([]byte(strconv.Itoa(len([]byte(str)))))
		conn.Write([]byte(str))
	}
}
