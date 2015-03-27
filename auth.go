package mandela

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	msgcenter "github.com/prestonTao/mandela/message_center"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"github.com/prestonTao/mandela/utils"
	"io"
	"net"
	"time"
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
//@name                 本机服务器的名称
//@return  remoteName   对方服务器的名称
func (this *Auth) SendKey(conn net.Conn, session engine.Session, name string) (remoteName string, err error) {
	//第一次连接，向对方发送自己的名称
	conn.Write(GetBytesForName(name))
	// //得到对方名称
	remoteName, err = GetNameForConn(conn)
	return
}

//接收
//name   自己的名称
//@return  remoteName   对方服务器的名称
func (this *Auth) RecvKey(conn net.Conn, name string) (remoteName string, err error) {
	/*
		获取对方的名称
	*/
	if remoteName, err = GetNameForConn(conn); err != nil {
		return
	}
	/*
		开始验证对方客户端名称
	*/
	clientIdInfo := new(nodeStore.IdInfo)
	json.Unmarshal([]byte(remoteName), clientIdInfo)
	/*
		这是新节点，需要给他生成一个id
	*/
	if clientIdInfo.Id == Str_zaro {
		//生成id之前先检查这个id是否存在
		fmt.Println("在网络中注册一个域名：", clientIdInfo.UserName)
		targetId := utils.GetHashForDomain(clientIdInfo.UserName)
		msgOne := &msgcenter.Message{
			TargetId:   targetId,
			ProtoId:    msgcenter.MSGID_findDomain,
			CreateTime: time.Now().Unix(),
			Sender:     nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
			Content:    []byte(clientIdInfo.UserName),
			Accurate:   false,
		}
		msgOne.Hash = msgcenter.GetHash(msgOne)
		ret := <-msgcenter.SendTimeOutMsg(msgOne, time.Second*10)
		if ret == "timeout" {
			err = errors.New("查询域名超时")
			return
		}
		if ret == "true" {
			err = errors.New("这个域名已经被注册了")
			return
		}

		*clientIdInfo, err = nodeStore.NewIdInfo(clientIdInfo.UserName, clientIdInfo.Email, clientIdInfo.Local, nodeStore.ParseId(name))
		//给服务器发送生成的id
		newName := string(clientIdInfo.Build())
		conn.Write(GetBytesForName(newName))
		err = errors.New("给新节点分配一个idinfo")
		return
	}

	/*
		验证成功后，向对方发送自己的名称
	*/
	//得到对方名称
	conn.Write(GetBytesForName(name))
	return
}

/*
	通过名称字符串获得bytes
	@name   要序列化的name字符串
*/
func GetBytesForName(name string) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(name)))
	buf.Write([]byte(name))
	return buf.Bytes()
}

/*
	通过读连接中的bytes获取name字符串
*/
func GetNameForConn(conn net.Conn) (name string, err error) {
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	nameLenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, nameLenght)
	if n, e := conn.Read(nameByte); e != nil {
		err = e
		return
	} else {
		//得到对方名称
		name = string(nameByte[:n])
		return
	}
}
