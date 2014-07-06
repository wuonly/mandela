package main

import (
	"../../socks5"
)

func main() {
	socks5.NewServer("127.0.0.1", 1080)
}
