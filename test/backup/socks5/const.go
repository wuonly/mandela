package socks5

const (
	Version = byte(5) //socks协议版本，只支持socks5

	MethodNoRequired   = 0x00 //无需认证
	MethodGSSAPI       = 0x01 //GSSAPI
	MethodUserNamePwd  = 0x02 //用户名/口令验证机制
	MethodCustom       = 0x81 //0x80-0xFE   RESERVED FOR PRIVATE METHODS(私有认证机制)
	MethodNoAcceptable = 0xff //不兼容的版本

	CMD_CONNECT       = byte(1) //CONNECT
	CMD_BIND          = byte(2) //BIND
	CMD_UDP_ASSOCIATE = byte(3) //UDP ASSOCIATE

	RESERVED = byte(0) //保留字段，必须为0x00

	ADDR_TYPE_IP     = byte(1) //IPv4地址
	ADDR_TYPE_IPV6   = byte(4) //IPv6地址
	ADDR_TYPE_DOMAIN = byte(3) //FQDN(全称域名)

	REP_SUCCEED                    = byte(0)
	REP_SERVER_FAILURE             = byte(1)
	REP_CONNECTION_NOT_ALLOW       = byte(2)
	REP_NETWORK_UNREACHABLE        = byte(3)
	REP_HOST_UNREACHABLE           = byte(4)
	REP_CONNECTION_REFUSED         = byte(5)
	REP_TTL_EXPIRED                = byte(6)
	REP_COMMAND_NOT_SUPPORTED      = byte(7)
	REP_ADDRESS_TYPE_NOT_SUPPORTED = byte(8)
)

type HandshakePack struct {
	Version     byte
	MethodCount byte
	Methods     []byte
}

type RequestPack struct {
	Version byte
	Cmd     byte
	Rsv     byte
	Atype   byte
	DSTAddr string
	DSTPort int
}

type ResponsePack struct {
	Version  int
	Rep      int
	Rsv      int
	Atype    int
	BindAddr string
	BindPort int
}

/**

client send
+----+----------+----------+
|VER | NMETHODS | METHODS  |
+----+----------+----------+
| 1  |    1     | 1 to 255 |
+----+----------+----------+

server send
+----+--------+
|VER | METHOD |
+----+--------+
| 1  |   1    |
+----+--------+

client send
+----+-----+-------+------+----------+----------+
|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
+----+-----+-------+------+----------+----------+
| 1  |  1  | X'00' |  1   | Variable |    2     |
+----+-----+-------+------+----------+----------+


server send
+----+-----+-------+------+----------+----------+
|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
+----+-----+-------+------+----------+----------+
| 1  |  1  | X'00' |  1   | Variable |    2     |
+----+-----+-------+------+----------+----------+


**/
