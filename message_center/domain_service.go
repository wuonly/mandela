package message_center

import (
	// "fmt"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"time"
)

var dataStore = make(map[string]string)

func init() {
	addRouter(MSGID_findDomain, findDomain)
}

type Domaim struct {
	IdInfo string `json:"id_info"`
}

/*
	查询一个域名是否存在
*/
func findDomain(c engine.Controller, packet engine.GetPacket, msg *Message) (bool, string) {
	if msg.ReplyHash == "" {
		newMsg := Message{
			TargetId:   msg.Sender,
			ProtoId:    msg.ProtoId,
			CreateTime: time.Now().Unix(),
			Sender:     nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
			ReplyTime:  time.Now().Unix(),
			Content:    []byte("false"),
			ReplyHash:  msg.Hash,
			Accurate:   true,
		}
		newMsg.Hash = GetHash(&newMsg)
		if _, ok := dataStore[string(msg.Content)]; ok {
			newMsg.Content = []byte("true")
		}
		SendMessage(&newMsg)
		return true, ""
	} else {
		return true, string(msg.Content)
	}
}

/*
	保存一个域名
*/
func saveDomain(c engine.Controller, packet engine.GetPacket, msg *Message) {

}
