package message_center

import (
// "fmt"
)

const (
	MSGID_Text        = iota + 101 //显示文本消息
	MSGID_SendMessage              //
)

func init() {
	addRouter(MSGID_Text, showTextMsg)
}
