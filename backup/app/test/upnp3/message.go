package upnp

import (
	"bytes"
)

type RequestInfo struct {
	Method string `POST` //请求方法
	Host   string //请求ip及端口
	Url    string //请求路径
	Proto  string `HTTP/1.1` //协议版本
}

//消息类
type Msg struct {
	requestInfo RequestInfo
	headerMap   map[string]string
	body        Node
}

func (msg *Msg) BuildString() string {
	buf := bytes.NewBufferString(msg.requestInfo.Method)
	buf.WriteString(" " + msg.requestInfo.Url + " " + msg.requestInfo.Proto + "\r\n")
	for key, value := range msg.headerMap {
		buf.WriteString(key + ": " + value + "\r\n")
	}
	buf.WriteString("\r\n")
	return buf.String()
}

type XMLNode interface {
	AddChild(node XMLNode)
	BuildXML() string
}

type Node struct {
	Name    string
	Content string
	Attr    map[string]string
	Child   []Node
}

func (n *Node) AddChild(node Node) {
	n.Child = append(n.Child, node)
}
func (n *Node) BuildXML() string {
	buf := bytes.NewBufferString("<")
	buf.WriteString(n.Name)
	for key, value := range n.Attr {
		buf.WriteString(" ")
		buf.WriteString(key + "=" + value)
	}
	buf.WriteString(">" + n.Content)

	for _, node := range n.Child {
		buf.WriteString(node.BuildXML())
	}
	buf.WriteString("</" + n.Name + ">")
	return buf.String()
}
