package mandela

import (
	"encoding/json"
	"fmt"
	msgc "github.com/prestonTao/mandela/message_center"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"strings"
	"time"
)

//保存一个键值对
func SaveData(key, value string) {
	clientConn, _ := engine.GetController().GetSession(nodeStore.SuperName)
	data := []byte(key + "!" + value)
	clientConn.Send(msgc.SaveKeyValueReqNum, &data)
}

//给所有客户端发送消息
func SendMsgForAll(message string) {
	messageSend := msgc.Message{
		ProtoId:    msgc.MSGID_Text,
		CreateTime: time.Now().Unix(),
		Content:    []byte(message),
	}
	for idOne, nodeOne := range nodeStore.GetAllNodes() {
		if clientConn, ok := engine.GetController().GetSession(string(nodeOne.IdInfo.Build())); ok {
			messageSend.TargetId = idOne
			data, _ := json.Marshal(messageSend)
			clientConn.Send(msgc.SendMessageNum, &data)
		}
	}
}

//给某个人发送消息
func SendMsgForOne(target, message string) {
	messageSend := &msgc.Message{
		TargetId:   target,
		Sender:     nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
		ProtoId:    msgc.MSGID_Text,
		CreateTime: time.Now().Unix(),
		Content:    []byte(message),
	}
	msgc.SendMessage(messageSend)
}

/*
	发送消息给一个域名
*/
func SendMsgForDomain(tdomain, msg string) {

}

/*
	查看本地保存的所有节点id
*/
func See() {
	allNodes := nodeStore.GetAllNodes()
	for key, _ := range allNodes {
		fmt.Println(key)
	}
}

/*
	查看本地保存的节点中，小于本节点id的所有节点
*/
func SeeLeftNode() {
	nodes := nodeStore.GetLeftNode(*nodeStore.Root.IdInfo.GetBigIntId(), nodeStore.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(id.IdInfo.GetId())
	}
}

/*
	查看本地保存的节点中，大于本节点id的所有节点
*/
func SeeRightNode() {
	nodes := nodeStore.GetRightNode(*nodeStore.Root.IdInfo.GetBigIntId(), nodeStore.MaxRecentCount)
	for _, id := range nodes {
		fmt.Println(id.IdInfo.GetId())
	}
}

/*
	添加一个超级节点ip地址
	@addr   例如：121.45.6.157:8076
*/
func AddAddr(addr string) {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return
	}
	if !IsOnlyIp(addrs[0]) {
		return
	}
	if CheckOnline(addr) {
		addSuperPeerAddr(addr)
	}
}

/*
	注册一个域名帐号
	@name     姓名
	@email    邮箱
	@domain   网络唯一标识
*/
func CreateAccount(name, email, domain string) {
	//连接网络并得到一个idinfo
	idInfo, err := GetId(nodeStore.NewIdInfo(name, email, domain, Str_zaro))
	if err == nil {
		Init_IdInfo = *idInfo
	} else {
		fmt.Println("从网络中获得idinfo失败")
		return
	}
}
