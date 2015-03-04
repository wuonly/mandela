package message_center

import (
	"fmt"
	engine "github.com/prestonTao/mandela/net"
	"sync"
)

type MsgHandler func(c engine.Controller, packet engine.GetPacket, msg *Message)

var handlersMapping = make(map[int32]MsgHandler)
var router_lock = new(sync.RWMutex)

func addRouter(msgId int32, handler MsgHandler) {
	router_lock.Lock()
	defer router_lock.Unlock()
	handlersMapping[msgId] = handler
}

func getHandler(msgId int32) MsgHandler {
	router_lock.Lock()
	defer router_lock.Unlock()
	handler := handlersMapping[msgId]
	return handler
}

/*
	消息分发程序
*/
func handlerProcess(c engine.Controller, packet engine.GetPacket, msg *Message) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				fmt.Println("网络库：", e.Error())
			}
		}
	}()
	handler := getHandler(msg.ProtoId)
	if handler == nil {
		fmt.Println("消息中心：未注册的消息编号-", packet.MsgID)
		return
	}
	handler(c, packet, msg)
}
