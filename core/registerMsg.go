package core

import (
	"github.com/prestonTao/mandela/core/message_center"
	engine "github.com/prestonTao/mandela/core/net"
)

func init() {
	//注册节点查找服务
	nodeManager := new(message_center.NodeManager)
	engine.RegisterMsg(message_center.FindNodeNum, nodeManager.FindNode)

	//注册发送消息服务
	engine.RegisterMsg(message_center.SendMessageNum, message_center.RecvMsg)

	// dataStore := new(service.DataStore)
	// engine.RegisterMsg(message.SaveKeyValueReqNum, dataStore.SaveDataReq)

}
