package mandela

import (
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/service"
)

func (this *Manager) registerMsg() {
	nodeManager := new(service.NodeManager)
	// this.engine.RegisterMsg(message.IntroduceSelf, nodeManager.IntroduceSelfRsp)
	// this.engine.RegisterMsg(message.FindNodeReqNum, nodeManager.FindNodeReq)
	// this.engine.RegisterMsg(message.FindNodeRspNum, nodeManager.FindNodeRsp)
	// this.engine.RegisterMsg(message.FindRecentNodeReqNum, nodeManager.FindRecentNodeReq)

	this.engine.RegisterMsg(message.FindNodeNum, nodeManager.FindNode)

	dataStore := new(service.DataStore)
	this.engine.RegisterMsg(message.SaveKeyValueReqNum, dataStore.SaveDataReq)

	messageService := new(service.Message)
	this.engine.RegisterMsg(message.SendMessage, messageService.RecvMsg)

}
