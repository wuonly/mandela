package message_center

import (
	"fmt"
	engine "github.com/prestonTao/mandela/net"
)

/*
	显示文本消息
*/
func showTextMsg(c engine.Controller, packet engine.GetPacket, msg *Message) {
	fmt.Println(string(msg.Content))
}
