package peerNode

import (
	engine "mandela/peerNode/messageEngine"
	"net"
)

const (
	version = 1
)

type Auth struct {
}

/*
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
| version   | ctp        | size      | name           |
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
| 版本      | 连接类型   | 数据长度  | 连接名称       |
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
| 2 byte    | 2 byte     | 4 byte    |                |
+++++++++++++++++++++++++++++++++++++++++++++++++++++++

version：版本
	1：第一个版本

ctp：连接类型
	1：带name的连接
	2：不带name的连接

name：连接名称
	区分每一个客户端的名称

*/

//发送
func (this *Auth) SendKey(conn net.Conn, session engine.Session) (err error) {
	session.GetName()
	return
}

//接收
func (this *Auth) RecvKey(conn net.Conn) (name string, err error) {
	this.session++
	// name = strconv.ParseInt(this.session, 10, )
	// name = strconv.Itoa(this.session)
	name = strconv.FormatInt(this.session, 10)
	return
}
