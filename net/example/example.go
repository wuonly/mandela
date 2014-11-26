package main

import (
	"fmt"
	"github.com/prestonTao/mandela/net"
)

func main() {
	example1()
}

func example1() {
	server := net.NewNet("tao")
	// server
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello World")
}

func nimei(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "nimei")
}
