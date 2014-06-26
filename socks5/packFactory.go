package socks5

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

type PackFactory struct {
}

func (this *PackFactory) handshakePack(conn net.Conn) *HandshakePack {
	var ver, nMethods byte
	// handshake
	err := binary.Read(conn, binary.BigEndian, &ver)
	if err != nil {
		// if err == io.EOF {
		// 	return nil
		// }
		log.Println(err.Error())
	}
	err = binary.Read(conn, binary.BigEndian, &nMethods)
	if err != nil {
		log.Println(err.Error())
	}
	methods := make([]byte, nMethods)
	err = binary.Read(conn, binary.BigEndian, methods)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(ver, nMethods, methods)
	pack := HandshakePack{
		Version:     ver,
		MethodCount: nMethods,
		Methods:     methods,
	}
	return &pack
}

func (this *PackFactory) poxyReqPack(conn net.Conn) *RequestPack {
	var ver, cmd, reserved, addrType byte
	err := binary.Read(conn, binary.BigEndian, &ver)
	if err != nil {
		log.Println(err.Error())
	}
	err = binary.Read(conn, binary.BigEndian, &cmd)
	if err != nil {
		log.Println(err.Error())
	}
	err = binary.Read(conn, binary.BigEndian, &reserved)
	if err != nil {
		log.Println(err.Error())
	}
	err = binary.Read(conn, binary.BigEndian, &addrType)
	if err != nil {
		log.Println(err.Error())
	}
	if ver != Version {
		log.Println(err.Error())
		return nil
	}
	if reserved != RESERVED {
		log.Println(err.Error())
		return nil
	}
	//地址类型不支持
	if addrType != ADDR_TYPE_IP && addrType != ADDR_TYPE_DOMAIN && addrType != ADDR_TYPE_IPV6 {
		return nil
	}

	// var DSTAddr string
	var address []byte
	if addrType == ADDR_TYPE_IP {
		address = make([]byte, 4)
	}
	if addrType == ADDR_TYPE_DOMAIN {
		var domainLength byte
		err := binary.Read(conn, binary.BigEndian, &domainLength)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		address = make([]byte, domainLength)
	}
	if addrType == ADDR_TYPE_IPV6 {
		address = make([]byte, 16)
	}
	err = binary.Read(conn, binary.BigEndian, address)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	var port uint16
	err = binary.Read(conn, binary.BigEndian, &port)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var hostPort string
	if addrType == ADDR_TYPE_IP || addrType == ADDR_TYPE_IPV6 {
		ip := net.IP(address)
		hostPort = net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	} else if addrType == ADDR_TYPE_DOMAIN {
		hostPort = net.JoinHostPort(string(address), strconv.Itoa(int(port)))
	}
	log.Println(ver, cmd, reserved, addrType, address, port, hostPort)
	pack := RequestPack{
		Version: ver,
		Cmd:     cmd,
		Rsv:     reserved,
		Atype:   addrType,
		DSTAddr: hostPort,
		DSTPort: int(port),
	}
	return &pack
}
