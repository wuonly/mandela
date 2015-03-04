package message_center

import (
	// "fmt"
	engine "github.com/prestonTao/mandela/net"
)

func init() {
	// addRouter(, handler)
}

type Domaim struct {
	IdInfo string `json:"id_info"`
}

/*
	保存一个域名
*/
func saveDoMaim(c engine.Controller, packet engine.GetPacket, msg *Message) {

}
