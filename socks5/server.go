package socks5

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"strings"
)

func NewServer(host string, port int) {
	sp := new(SocksProxy)
	sp.host = host
	sp.tcpPort = port
	sp.packFactory = new(PackFactory)
	sp.Run(host, port)
}

type SocksProxy struct {
	host        string
	tcpPort     int
	packFactory *PackFactory
}

func (this *SocksProxy) Run(host string, port int) {
	addrPort := host + ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", addrPort)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("start")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go this.handshake(conn)
	}
}

func (this *SocksProxy) handshake(conn net.Conn) {
	handshakePack := this.packFactory.handshakePack(conn)
	//不支持的版本
	if handshakePack.Version != Version || handshakePack.MethodCount < byte(1) {
		log.Println("不支持的版本")
		binary.Write(conn, binary.BigEndian, []byte{MethodNoAcceptable})
		return
	}
	//不支持身份验证的方法
	if bytes.IndexByte(handshakePack.Methods, MethodNoRequired) == -1 {
		binary.Write(conn, binary.BigEndian, []byte{MethodNoAcceptable})
		return
	}
	//发送服务器版本
	err := binary.Write(conn, binary.BigEndian, Version)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//不需要身份验证
	err = binary.Write(conn, binary.BigEndian, []byte{MethodNoRequired})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("发送不需要验证请求")

	//第二次请求
	requestPack := this.packFactory.poxyReqPack(conn)
	//
	if requestPack.Cmd == CMD_CONNECT {

	}
	if requestPack.Cmd == CMD_BIND {

	}
	if requestPack.Cmd == CMD_UDP_ASSOCIATE {
		log.Println("请求代理类型为UDP")
		go this.read(requestPack)
		this.createUDPProxy(requestPack)
		portBuf := bytes.NewBuffer([]byte{})
		binary.Write(portBuf, binary.BigEndian, int16(9980))
		port := portBuf.Bytes()
		binary.Write(conn, binary.BigEndian, []byte{5, 0, 0, 1, 127, 0, 0, 1, port[0], port[1]})

		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		log.Println("tcp-------", buf[:n])
	}

}

func (this *SocksProxy) createUDPProxy(pack *RequestPack) {
	// remotAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(pack.DSTPort))
	// locaAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(pack.DSTPort))
	// locaAddr := net.UDPAddr{
	// 	IP:   net.IPv4zero,
	// 	Port: 9980,
	// }
	// conn, err := net.ListenUDP("udp", &locaAddr)

	// addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9980")
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	log.Println("123454", err.Error())
	// }
	// portBuf := bytes.NewBuffer([]byte{})
	// binary.Write(portBuf, binary.LittleEndian, int16(9980))
	// port := portBuf.Bytes()
	// // conn.WriteToUDP([]byte{0, 0, 0, 1, 0, 0, 0, 0, port[0], port[1]}, remotAddr)
	// conn.WriteToUDP([]byte{5, 0, 0, 1, 127, 0, 0, 1, port[0], port[1]}, remotAddr)

	// // conn.Close()
	// buf := make([]byte, 1024)
	// n, _, err := conn.ReadFromUDP(buf)
	// if err != nil {
	// 	log.Println("haha", err.Error())
	// }
	// log.Println("udp-------", buf[:n])
}

func (this *SocksProxy) read(pack *RequestPack) {
	log.Println(pack.DSTAddr, pack.DSTPort)
	//获取一个可用的端口
	packConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	localPortStr := strings.Split(packConn.LocalAddr().String(), ":")[1]
	log.Println(localPortStr)

	remotAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+localPortStr)
	locaAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(pack.DSTPort))
	conn, err := net.ListenUDP("udp", locaAddr)
	if err != nil {
		log.Println("连接客户端失败", err.Error())
	}
	portBuf := bytes.NewBuffer([]byte{})
	localPortInt, _ := strconv.Atoi(localPortStr)
	binary.Write(portBuf, binary.LittleEndian, int16(localPortInt))
	port := portBuf.Bytes()
	_, err = conn.WriteToUDP([]byte{5, 0, 0, 1, 127, 0, 0, 1, port[0], port[1]}, remotAddr)
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Println("haha", err.Error())
	}
	log.Println("udp----client---", buf[:n])

}

type UDPio struct {
	conn *net.UDPConn
}

func (this *UDPio) Run() {
	defer func() {

	}()

}

func (this *UDPio) read() {

}
