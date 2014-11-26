package mandela

import (
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/service"
)

func (this *Manager) registerMsg() {
	//注册节点查找服务
	nodeManager := new(service.NodeManager)
	this.engine.RegisterMsg(message.FindNodeNum, nodeManager.FindNode)

	//注册发送消息服务
	messageService := new(service.Message)
	this.engine.RegisterMsg(message.SendMessage, messageService.RecvMsg)

	dataStore := new(service.DataStore)
	this.engine.RegisterMsg(message.SaveKeyValueReqNum, dataStore.SaveDataReq)

}
