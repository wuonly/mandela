package net

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Net struct {
	Recv          chan *GetPacket //获得数据
	Name          string          //本机名称
	sessionStore  *sessionStore
	closecallback CloseCallback
}

func (this *Net) Listen(ip string, port int32) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(int(port)))
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("监听一个地址：", ip+":"+strconv.Itoa(int(port)))
	// fmt.Println(ip + ":" + strconv.Itoa(int(port)) + "成功启动服务器")
	go this.listener(listener)
}

func (this *Net) listener(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go this.newConnect(conn)
	}
}

//创建一个新的连接
func (this *Net) newConnect(conn net.Conn) {
	remoteName, err := defaultAuth.RecvKey(conn, this.Name)
	if err != nil {
		return
	}

	// this.lock.Lock()
	// defer this.lock.Unlock()

	// this.session++

	serverConn := &ServerConn{
		// attrbutes:      make(map[string]interface{}),
		conn: conn,
		// session:        this.session,
		Ip:             conn.RemoteAddr().String(),
		Connected_time: time.Now().String(),
		// outData:        make(chan *[]byte, 1000),
		inPack: this.Recv,
		net:    this,
	}
	serverConn.name = remoteName
	serverConn.attrbutes = make(map[string]interface{})
	serverConn.run()
	this.sessionStore.addSession(remoteName, serverConn)

	fmt.Println(time.Now().String(), "建立连接", conn.RemoteAddr().String())

}

//关闭连接
func (this *Net) CloseClient(name string) bool {
	session, ok := this.sessionStore.getSession(name)
	if ok {
		if this.closecallback != nil {
			this.closecallback(name)
		}
		this.sessionStore.removeSession(name)
		session.Close()
		return true
	}
	return false
}

/*
	连接一个服务器
	@serverName   给客户端发送的自己的名字
	@powerful     是否是强连接，是强连接断开后自动重连
*/
func (this *Net) AddClientConn(ip, serverName string, port int32, powerful bool) (Session, error) {
	clientConn := &Client{
		serverName: serverName,
		inPack:     this.Recv,
		net:        this,
		isPowerful: powerful,
	}
	clientConn.attrbutes = make(map[string]interface{})
	remoteName, err := clientConn.Connect(ip, port)
	if err == nil {
		clientConn.name = remoteName
		this.sessionStore.addSession(remoteName, clientConn)
		return clientConn, nil
	}
	return nil, err
}

func (this *Net) GetSession(name string) (Session, bool) {
	return this.sessionStore.getSession(name)
}

//发送数据
func (this *Net) Send(name string, msgID uint32, data []byte) bool {
	session, ok := this.sessionStore.getSession(name)
	if ok {
		session.Send(msgID, &data)
		return true
	} else {
		return false
	}
}

//@name   本服务器名称
func NewNet(name string) *Net {
	net := new(Net)
	net.Name = name
	net.sessionStore = NewSessionStore()
	net.Recv = make(chan *GetPacket, 5000)
	return net
}
