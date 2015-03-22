package message_center

const (
	MSGID_Text        = iota + 101 //显示文本消息
	MSGID_findDomain               //查找这个域名是否存在
	MSGID_recv_domain              //返回这个域名是否存在
)
