package upnp

import (
	// "log"
	"net/http"
	"strconv"
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

type SearchGatewayReq struct {
	host string
}

func (this SearchGatewayReq) send(host string) {
	request := this.BuildRequest()
}
func (this SearchGatewayReq) BuildRequest() http.Request {
	//请求头
	header := http.Header{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"SOAPAction":     `"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"`,
		"Content-Type":   "text/xml",
		"Connection":     "Close",
		"Content-Length": ""}
	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetStatusInfo`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	//请求
	request := http.Request{Method: "POST", Proto: "HTTP/1.1",
		Host: host, Url: "http://" + host + "/ipc", Header: header}
	reqest.Header.Set("Content-Length", strconv.Itoa(len([]byte(body))))
	return request
}
