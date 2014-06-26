package main

import (
	"encoding/binary"
	"net"
	// "time"
)

func main() {
	simple2()
}

func simple1() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	var VERSION byte = byte(5)
	binary.Write(conn, binary.BigEndian, VERSION)
	var METHOD = byte(6)
	binary.Write(conn, binary.BigEndian, METHOD)
}

func simple2() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	conn.Write([]byte{byte(5), byte(6)})
}
