package mandela

import (
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/service"
)

func init() {
	//注册节点查找服务
	nodeManager := new(service.NodeManager)
	engine.RegisterMsg(message.FindNodeNum, nodeManager.FindNode)

	//注册发送消息服务
	messageService := new(service.Message)
	engine.RegisterMsg(message.SendMessage, messageService.RecvMsg)

	dataStore := new(service.DataStore)
	engine.RegisterMsg(message.SaveKeyValueReqNum, dataStore.SaveDataReq)

}
