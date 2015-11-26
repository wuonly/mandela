package net

import (
// "fmt"
)

const (
	HoldConn  = iota //心跳连接
	CloseConn        //关闭连接
)

var zaro_bytes = []byte{0x00}

func init() {
	AddRouter(CloseConn, CloseConnMsg)
}

/*
	关闭连接消息
*/
func CloseConnMsg(c Controller, msg GetPacket) {
	session, ok := c.GetSession(msg.Name)
	if ok {
		session.Close()
	}
}
