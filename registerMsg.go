package mandela

import (
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/service"
)

func (this *Manager) registerMsg() {
	nodeManager := new(service.NodeManager)
	this.engine.RegisterMsg(message.FindNodeReqNum, nodeManager.FindNodeReq)
	this.engine.RegisterMsg(message.FindNodeRspNum, nodeManager.FindNodeRsp)
	dataStore := new(service.DataStore)
	this.engine.RegisterMsg(message.SaveKeyValueReqNum, dataStore.SaveDataReq)
}
