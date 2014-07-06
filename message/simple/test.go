package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	msg "mandela/app/message"
)

func main() {
	m := new(msg.FindNode)
	m.MsgType = proto.Int32(1)
	m.Timeout = proto.Int32(1)
	m.NodeId = proto.String("haha")
	m.AccID = proto.Uint64(1)
	buf, _ := proto.Marshal(m)
	fmt.Println(buf)

	base := new(msg.BaseMsg)
	proto.Unmarshal(buf, base)
	fmt.Println(*base.MsgType, *base.Timeout)
}
