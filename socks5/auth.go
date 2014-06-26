package socks5

// import (
// 	"net"
// 	"strconv"
// )

// //socks代理客户端用户验证
// type Auth struct {
// 	conn net.Conn
// }

// //检查方法是否被我们支持
// func (this *Auth) CheckMethod(methods []uint8) bool {
// 	for _, method := range methods {
// 		//自定义方法
// 		if method == authCustom {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (this *Auth) Valid() bool {
// 	buf := make([]byte, 1024)
// 	this.conn.Read(buf)
// 	return true
// }
