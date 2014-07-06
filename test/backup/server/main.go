package main

import (
	"fmt"
	// "io"
	"log"
	"net"
	// "net/http"
)

func main() {
	Run()
}

var ips = []string{}

func Run() {
	addr, err := net.ResolveTCPAddr("tcp", ":9981")
	if err != nil {
		log.Fatal(err)
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 50)
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}
		fmt.Println(string(buf[:n]))
		ips = append(ips, string(buf[:n]))

		conn.Write([]byte(ips[0]))
		conn.Close()
	}

}
