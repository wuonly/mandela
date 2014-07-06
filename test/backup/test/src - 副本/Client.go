package main

import (
	"fmt"
	"net"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	conn, err := net.Dial("tcp", "100.64.211.233:1990")
	chk(err)
	_, err = conn.Write([]byte("hello"))
	chk(err)
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	runes := []rune(string(buf[:n+1]))
	start := 0
	for i, r := range runes {
		if r == 0 {
			command := string(runes[start:i])
			fmt.Println(command, i)

		}
	}

}
