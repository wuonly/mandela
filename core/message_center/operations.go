package message_center

import (
	"encoding/json"
	// "fmt"
	engine "github.com/prestonTao/mandela/core/net"
	// "github.com/prestonTao/mandela/nodeStore"
)

func SendMessage(msg *Message) {
	data, _ := json.Marshal(msg)
	pkg := engine.GetPacket{
		MsgID: SendMessageNum,
		Size:  uint32(len(data)),
		Date:  data,
	}
	if IsSendToSelf(engine.GetController(), pkg) {
		handlerProcess(engine.GetController(), pkg)
	}
}
