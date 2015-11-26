package message_center

import (
	"fmt"
	engine "github.com/prestonTao/mandela/core/net"
)

func init() {
	addRouter(MSGID_Text, showTextMsg)
}

/*
	显示文本消息
*/
func showTextMsg(c engine.Controller, packet engine.GetPacket, msg *Message) (bool, string) {
	fmt.Println(string(msg.Content))
	return true, ""
}
