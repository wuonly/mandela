package upnp

import (
// "log"
)

type Header struct {
	proto     string            //请求协议
	host      string            //请求地址
	headerMap map[string]string //请求头参数
}

type Taget struct {
	url    string //请求url
	method string //请求方法
}

type SearchGatewayReq struct{}

func (this SearchGatewayReq) send() {

}
func (this SearchGatewayReq) BuildBody() {

}
