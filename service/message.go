package service

import (
	"encoding/hex"
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
	fmt.Println(msg.Date)

	err := json.Unmarshal(msg.Date, messageRecv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(messageRecv)
	// proto.Unmarshal(msg.Date, messageRecv)

	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	fmt.Println("== 1", store.GetRootId())
	fmt.Println("== 2", messageRecv.TargetId)
	if store.GetRootId() == messageRecv.TargetId {
		fmt.Println(string(messageRecv.Content))
	} else {
		targetNode := store.Get(messageRecv.TargetId, true, "")
		if targetNode == nil {
			return
		}
		session, ok := c.GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
		if !ok {
			return
		}
		err := session.Send(message.SendMessage, &msg.Date)
		if err != nil {
			fmt.Println("message发送数据出错：", err.Error())
		}
	}
}
