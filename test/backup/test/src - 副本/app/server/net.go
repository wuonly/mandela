package server

import (
	"log"
	"net"
)

type Net interface {
}

type TCPServer struct {
	in chan 
}
