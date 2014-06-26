package main

import "net"

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:1990")
	chk(err)
	_, err = conn.Write([]byte("SHOUTDOWN THE APP"))
	chk(err)

}
