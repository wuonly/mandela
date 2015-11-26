package main

import (
	// "fmt"
	"github.com/prestonTao/mandela/net"
	"time"
)

func main() {
	example1()
}

func example1() {
	server := net.NewNet("tao")
	server.RegisterMsg(101, hello)
	time.Sleep(time.Minute * 10)
	// server
}

func hello(c net.Controller, msg FindNode) {

}

type FindNode struct {
	Name string `json:"name"`
}
