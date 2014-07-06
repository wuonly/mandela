package socks5

import (
	"encoding/binary"
	"net"
)

type Conn struct {
	atype    string //udp/tcp
	tcpConn  net.Conn
	udpConn  *net.UDPConn
	hostPort string
}

func (this *Conn) run() {
	switch this.atype {
	case "tcp":
		go this.readTCP()
	case "udp":
		go this.readUDP()
		go this.holdTCP()
	}
}

func (this *Conn) readUDP() {
	var rsv [2]byte
	var frag, atype byte
	for {
		err := binary.Read(this.udpConn, binary.BigEndian, rsv)
		if err != nil {
			return
		}
		err = binary.Read(this.udpConn, binary.BigEndian, frag)
		if err != nil {
			return
		}
		err = binary.Read(this.udpConn, binary.BigEndian, atype)
		if err != nil {
			return
		}
	}
}
func (this *Conn) readTCP() {

}

func (this *Conn) holdTCP() {
	// buf := make([]byte, 10)
	// n, err := this.tcpConn.Read(buf)
	// if err != nil {

	// }
}
