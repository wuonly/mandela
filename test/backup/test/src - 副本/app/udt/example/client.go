package main

import (
	"../../udt"
	"log"
	"net"
	"time"
)

func main() {
	if addr, err := net.ResolveUDPAddr("udp", "localhost:47008"); err != nil {
		log.Fatalf("Unable to resolve address: %s", err)
	} else {
		go client(addr)

		time.Sleep(5 * time.Second)
	}
}

func client(addr *net.UDPAddr) {
	if _, err := udt.DialUDT("udp", nil, addr); err != nil {
		log.Fatalf("Unable to dial: %s", err)
	}
}
