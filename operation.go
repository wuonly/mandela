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
			clientConn.Send(msgc.SendMessage, &data)
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
	SendMsg(messageSend)
}

func SendMsg(msg *msgc.Message) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		e, ok := err.(error)
	// 		if ok {
	// 			fmt.Println("网络库：", e.Error())
	// 		}
	// 	}
	// }()
	sendBytes, _ := json.Marshal(msg)
	//发送给自己的
	if nodeStore.ParseId(nodeStore.GetRootIdInfoString()) == msg.TargetId {
		// handler := getHandler(msg.ProtoId)
		// if handler == nil {
		// 	fmt.Println("消息中心：未注册的消息编号-")
		// 	return
		// }
		// packet := net.GetPacket{
		// 	MsgID: msgc.SendMessage,
		// 	Date:  sendBytes,
		// }
		// ok, str := handler(engine.GetController(), packet, msg)
		// fire(msg, ok, str)
		fmt.Println(msg.Content)
		return
	}

	var session engine.Session
	var ok bool
	//本机是超级节点
	if Init_IsSuperPeer {
		//是发给自己的弱节点
		if targetNode, ok := nodeStore.GetProxyNode(msg.TargetId); ok {
			if session, ok := engine.GetController().GetSession(string(targetNode.IdInfo.Build())); ok {
				err := session.Send(msgc.SendMessage, &sendBytes)
				if err != nil {
					fmt.Println("message发送数据出错：", err.Error())
				}
			} else {
				//这个节点离线了，想办法处理下
			}
			return
		}
		//转发出去
		targetNode := nodeStore.Get(msg.TargetId, true, "")
		if targetNode == nil {
			fmt.Println("本机未连入mandela网络")
			return
		}
		session, ok = engine.GetController().GetSession(string(targetNode.IdInfo.Build()))
	} else {
		//本机是普通节点
		//获得超级节点
		session, ok = engine.GetController().GetSession(nodeStore.SuperName)
	}
	if !ok {
		return
	}
	err := session.Send(msgc.SendMessage, &sendBytes)
	if err != nil {
		fmt.Println("message发送数据出错：", err.Error())
	}
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
