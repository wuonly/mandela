package mandela

import (
	"encoding/json"
	"fmt"
	msgc "github.com/prestonTao/mandela/message_center"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
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

//注册一个域名帐号
func CreateAccount(account string) {
	// id := GetHashKey(account)
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
	创建一个id
*/
func CreateIdInfo() {

}
