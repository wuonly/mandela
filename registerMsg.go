package mandela

import (
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/service"
)

func (this *Manager) registerMsg() {
	nodeManager := new(service.NodeManager)
	this.serverManager.RegisterMsg(message.FindNodeReqNum, nodeManager.FindNodeReq)
	this.serverManager.RegisterMsg(message.FindNodeRspNum, nodeManager.FindNodeRsp)
}
