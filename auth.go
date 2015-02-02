package mandela

import (
	"bytes"
	"encoding/binary"
	engineE "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"io"
	"net"
)

const (
	version = 1
)

type Auth struct {
	nodeManager *nodeStore.NodeManager
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
//@name                 本机服务器的名称
//@return  remoteName   对方服务器的名称
func (this *Auth) SendKey(conn net.Conn, session engineE.Session, name string) (remoteName string, err error) {
	//第一次连接，向对方发送自己的名称
	lenght := int32(len(name))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, lenght)
	buf.Write([]byte(name))
	conn.Write(buf.Bytes())

	//对方服务器验证成功后发送给自己的名称
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	nameLenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, nameLenght)
	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	//得到对方名称
	remoteName = string(nameByte[:n])
	return
}

//接收
//@return  remoteName   对方服务器的名称
func (this *Auth) RecvKey(conn net.Conn, name string) (remoteName string, err error) {
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	lenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, lenght)

	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	//得到对方名称
	remoteName = string(nameByte[:n])
	//开始验证对方客户端名称

	//验证成功后，向对方发送自己的名称
	nameLenght := int32(len(name))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, nameLenght)
	buf.Write([]byte(name))
	conn.Write(buf.Bytes())
	return
}

/*

*/
// type NewPeerAuth struct {
// }

// func (this *Auth) SendKey(conn net.Conn, session engineE.Session, name string) (remoteName string, err error) {

// }

// func (this *Auth) RecvKey(conn net.Conn, name string) (remoteName string, err error) {

// }

/*
	连接超级节点，得到一个id
	@ addr   超级节点ip地址
*/
// func GetId(addr string) (idInfo *IdInfo, err error) {

// 	idInfo, err = NewIdInfo("", "", "", zaro)
// 	if err != nil {
// 		fmt.Println(err)
// 		err = errors.New("生成id错误")
// 		return
// 	}

// 	conn, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		err = errors.New("连接超级节点失败")
// 		return
// 	}

// 	//第一次连接，向对方发送自己的名称
// 	lenght := int32(len(name))
// 	buf := bytes.NewBuffer([]byte{})
// 	binary.Write(buf, binary.BigEndian, lenght)
// 	buf.Write([]byte(name))
// 	conn.Write(buf.Bytes())

// 	//对方服务器验证成功后发送给自己的名称
// 	lenghtByte := make([]byte, 4)
// 	io.ReadFull(conn, lenghtByte)
// 	nameLenght := binary.BigEndian.Uint32(lenghtByte)
// 	nameByte := make([]byte, nameLenght)
// 	n, e := conn.Read(nameByte)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	//得到对方名称
// 	remoteName = string(nameByte[:n])

// }
