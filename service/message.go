package service

import (
	"encoding/json"
	"fmt"
	"github.com/prestonTao/mandela/message"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
)

type Message struct {
}

/*
	接收消息并显示或转发
*/
func (this *Message) RecvMsg(c engine.Controller, msg engine.GetPacket) {
	messageRecv := new(message.Message)
	err := json.Unmarshal(msg.Date, messageRecv)
	if err != nil {
		fmt.Println(err)
	}

	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	if nodeStore.ParseId(store.GetRootIdInfoString()) == messageRecv.TargetId {
		fmt.Println(string(messageRecv.Content))
	} else {
		targetNode := store.Get(messageRecv.TargetId, true, "")
		if targetNode == nil {
			return
		}
		// session, ok := c.GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
		session, ok := c.GetSession(string(targetNode.IdInfo.Build()))
		if !ok {
			return
		}
		err := session.Send(message.SendMessage, &msg.Date)
		if err != nil {
			fmt.Println("message发送数据出错：", err.Error())
		}
	}
}
