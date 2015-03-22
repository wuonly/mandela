package message_center

import (
	// "fmt"
	engine "github.com/prestonTao/mandela/net"
	// "github.com/prestonTao/mandela/nodeStore"
)

func init() {
	// addRouter(, handler)
}

type Domaim struct {
	IdInfo string `json:"id_info"`
}

/*
	查询一个域名是否存在
*/
func findDomain(c engine.Controller, packet engine.GetPacket, msg *Message) (bool, string) {
	//检查这个域名是否归自己管
	// store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	// store.GetRootId()
	return true, ""
}

/*
	查询一个域名是否存在返回
*/
func findDomainRecv(c engine.Controller, packet engine.GetPacket, msg *Message) (bool, string) {
	return true, ""
}

/*
	保存一个域名
*/
func saveDomain(c engine.Controller, packet engine.GetPacket, msg *Message) {

}
