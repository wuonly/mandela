package main

import (
	"../../udt"
	// "io"
	"log"
	"net"
	"time"
)

func main() {
	if addr, err := net.ResolveUDPAddr("udp", "localhost:47008"); err != nil {
		log.Fatalf("Unable to resolve address: %s", err)
	} else {
		go server(addr)
		// time.Sleep(200 * time.Millisecond)
		// go client(addr)

		time.Sleep(50 * time.Second)
	}
}

func server(addr *net.UDPAddr) {
	_, err := udt.ListenUDT("udp", addr)
	if err != nil {
		log.Fatalf("Unable to listen: %s", err)
	}

	// for {
	// 	_, err = listener.Accept()
	// 	if err != nil {
	// 		log.Fatalf("1111111111111111111111")
	// 	} else {
	// 		log.Fatalf("2222222222222222222222")
	// 	}
	// 	time.Sleep(50 * time.Second)
	// }

	// for {
	// 	var buf [512]byte
	// 	_, err := conn.Read(buf[0:])
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 	}
	// 	log.Println(string(buf[:]))
	// }

}
