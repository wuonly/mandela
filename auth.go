package mandela

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/prestonTao/mandela/nodeStore"
	engine "github.com/prestonTao/messageEngine"
	"io"
	"math/big"
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
func (this *Auth) SendKey(conn net.Conn, session engine.Session, name string) (string, error) {

	lenght := int32(len(name))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, lenght)

	buf.Write([]byte(name))
	conn.Write(buf.Bytes())

	return name, nil
}

//接收
func (this *Auth) RecvKey(conn net.Conn) (name string, err error) {
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	lenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, lenght)

	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	name = string(nameByte[:n])

	node := new(nodeStore.Node)
	nodeIdInt, b := new(big.Int).SetString(name, 10)
	if !b {
		// fmt.Println("节点id格式不正确，应该为十进制字符串")
		err = errors.New("节点id格式不正确，应该为十进制字符串")
		return
	}
	node.NodeId = nodeIdInt
	node.Addr = conn.RemoteAddr().String()
	// node.TcpPort = conn.RemoteAddr()

	this.nodeManager.AddNode(node)
	return
}
