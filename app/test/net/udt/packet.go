package udt

const (
	// Control packet types
	handshake    = 0x0 //协议连接握手
	keepalive    = 0x1 //保持连接
	ack          = 0x2 //应答，位16-31是应答序号
	nak          = 0x3 //Negative应答（NAK），丢失信息的32位整数数组
	unused       = 0x4 //保留
	shutdown     = 0x5 //关闭
	ack2         = 0x6 //应答一个应答（ACK2），16-31位，应答序号。
	msg_drop_req = 0x7 //保留将来使用
)
