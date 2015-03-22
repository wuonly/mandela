package message_center

// import (
// 	"fmt"
// )

// func SendMessage(msg *msgc.Message) {
// 	// defer func() {
// 	// 	if err := recover(); err != nil {
// 	// 		e, ok := err.(error)
// 	// 		if ok {
// 	// 			fmt.Println("网络库：", e.Error())
// 	// 		}
// 	// 	}
// 	// }()
// 	sendBytes, _ := json.Marshal(msg)
// 	//发送给自己的
// 	if nodeStore.ParseId(nodeManager.GetRootIdInfoString()) == msg.TargetId {
// 		// handler := getHandler(msg.ProtoId)
// 		// if handler == nil {
// 		// 	fmt.Println("消息中心：未注册的消息编号-")
// 		// 	return
// 		// }
// 		// packet := net.GetPacket{
// 		// 	MsgID: msgc.SendMessage,
// 		// 	Date:  sendBytes,
// 		// }
// 		// ok, str := handler(engine.GetController(), packet, msg)
// 		// fire(msg, ok, str)
// 		fmt.Println(msg.Content)
// 		return
// 	}

// 	var session engine.Session
// 	var ok bool
// 	//本机是超级节点
// 	if Init_IsSuperPeer {
// 		//是发给自己的弱节点
// 		if targetNode, ok := nodeManager.GetProxyNode(msg.TargetId); ok {
// 			if session, ok := engine.GetController().GetSession(string(targetNode.IdInfo.Build())); ok {
// 				err := session.Send(msgc.SendMessage, &sendBytes)
// 				if err != nil {
// 					fmt.Println("message发送数据出错：", err.Error())
// 				}
// 			} else {
// 				//这个节点离线了，想办法处理下
// 			}
// 			return
// 		}
// 		//转发出去
// 		targetNode := nodeManager.Get(msg.TargetId, true, "")
// 		if targetNode == nil {
// 			fmt.Println("本机未连入mandela网络")
// 			return
// 		}
// 		session, ok = engine.GetController().GetSession(string(targetNode.IdInfo.Build()))
// 	} else {
// 		//本机是普通节点
// 		//获得超级节点
// 		session, ok = engine.GetController().GetSession(nodeManager.SuperName)
// 	}
// 	if !ok {
// 		return
// 	}
// 	err := session.Send(msgc.SendMessage, &sendBytes)
// 	if err != nil {
// 		fmt.Println("message发送数据出错：", err.Error())
// 	}
// }
