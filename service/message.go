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

func (this *Message) RecvMsg(c engine.Controller, msg engine.GetPacket) {
	messageRecv := new(message.Message)
	json.Unmarshal(msg.Date, messageRecv)
	// proto.Unmarshal(msg.Date, messageRecv)

	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	if store.GetRootId() == messageRecv.TargetId {
		fmt.Println(string(messageRecv.Content))
	} else {
		targetNode := store.Get(messageRecv.TargetId, true, "")
		if targetNode == nil {
			return
		}
		session, ok := c.GetSession(targetNode.NodeId.String())
		if !ok {
			return
		}
		err := session.Send(message.SendMessage, &msg.Date)
		if err != nil {
			fmt.Println("message发送数据出错：", err.Error())
		}
	}
}
