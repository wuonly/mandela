package net

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Net struct {
	name         string        //本服务器名称
	timeout      int           //设置一个ping超时时间，单位为毫秒
	Name         string        //本机名称
	sessionStore *sessionStore //
	router       *RouterStore  //请求路径路由表
}

// func (this *Net) Router(url string, handler MsgHandler) {
// 	this.router.AddRouter(url, handler)
// }

func (this *Net) RegisterMsg(msgId int, handler MsgHandler) {
	if msgId <= 100 {
		fmt.Println("该消息不能注册，消息编号0-100被系统占用。")
		return
	}
	this.router.AddRouter(msgId, handler)
}

func (this *Net) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/hello" {
		return
	}
	if r.URL.Path == "/nimei" {
		return
	}
	http.NotFound(w, r)
	return
}

func (this *Net) Listen(ip string, port int) {
	go http.ListenAndServe(ip+":"+strconv.Itoa(port), this)

	fmt.Println("webServer startup...")

}

//关闭连接
// func (this *Net) CloseClient(name string) bool {

// }

//@serverName   给客户端发送的自己的名字
func (this *Net) AddClientConn(ip, serverName string, port int, powerful bool) (session Session, b error) {
	return
}

func (this *Net) GetSession(name string) (Session, bool) {
	return this.sessionStore.getSession(name)
}

//发送数据
func (this *Net) Send(name string, msgID int, data []byte) bool {
	session, ok := this.sessionStore.getSession(name)
	if ok {
		session.Send(msgID, &data)
		return true
	} else {
		return false
	}
}

func (this *Net) Ping(address string) {
	c, err := net.Dial("ip4:icmp", address)
	if err != nil {
		return
	}
	c.SetDeadline(time.Now().Add(time.Duration(this.timeout) * time.Millisecond))
	defer c.Close()

}

//@name   本服务器名称
func NewNet(name string) *Net {
	net := Net{
		name:         name,
		timeout:      400,
		sessionStore: NewSessionStore(),
		router:       NewRouter(),
	}
	return &net
}
