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

		//先判断是否在自己的代理节点中
		if targetNode, ok := store.GetProxyNode(messageRecv.TargetId); ok {
			if session, ok := c.GetSession(string(targetNode.IdInfo.Build())); ok {
				err := session.Send(message.SendMessage, &msg.Date)
				if err != nil {
					fmt.Println("message发送数据出错：", err.Error())
				}
			} else {
				//这个节点离线了，想办法处理下
			}
			return
		}
		// fmt.Println("把消息转发出去")
		//最后转发出去
		targetNode := store.Get(messageRecv.TargetId, true, "")
		if targetNode == nil {
			return
		}
		// session, ok := c.GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
		session, ok := c.GetSession(string(targetNode.IdInfo.Build()))
		if !ok {
			return
		}
		// fmt.Println(session.GetName())
		err := session.Send(message.SendMessage, &msg.Date)
		if err != nil {
			fmt.Println("message发送数据出错：", err.Error())
		}
	}
}
