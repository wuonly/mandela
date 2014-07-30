package service

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/prestonTao/mandela/message"
	"github.com/prestonTao/mandela/nodeStore"
	engine "github.com/prestonTao/messageEngine"
	"math/big"
)

// func init() {
// 	nodeManager := new(NodeManager)
// 	engine.AddRouter(message.FindNodeReqNum, nodeManager.FindNodeReq)
// 	engine.AddRouter(message.FindNodeRspNum, nodeManager.FindNodeRsp)
// }

type NodeManager struct {
}

//连接到本机后，目标机器会给自己发送它的名片
func (this *NodeManager) IntroduceSelfRsp(c engine.Controller, msg engine.GetPacket) {
	recvNode := new(message.FindNodeRsp)
	proto.Unmarshal(msg.Date, recvNode)
	// fmt.Println("接收到：", *recvNode.NodeId)
	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)

	nodeIdInt, _ := new(big.Int).SetString(*recvNode.NodeId, 10)
	newNode := &nodeStore.Node{
		NodeId:  nodeIdInt,
		Addr:    *recvNode.Addr,
		IsSuper: !*recvNode.IsProxy,
		TcpPort: int32(*recvNode.TcpPort),
		UdpPort: int32(*recvNode.UdpPort),
	}

	isNeed, replace := store.CheckNeedNode(newNode.NodeId.String())
	// fmt.Println("这个节点是否需要：", isNeed)
	if isNeed {
		store.AddNode(newNode)
		c.GetNet().AddClientConn(nodeIdInt.String(), newNode.Addr, store.GetRootId(), newNode.TcpPort, false)
		if replace != "" {
			//删除原来的连接
			fmt.Println("替换原有的连接：", replace)
			if session, ok := c.GetSession(replace); ok {
				session.Close()

				delNode := new(nodeStore.Node)
				delNode.NodeId, _ = new(big.Int).SetString(replace, 10)
				store.DelNode(delNode)
				fmt.Println("替换成功")
			}
		}
	}
	//--------------------------------------------
	//    互相介绍自己
	//--------------------------------------------
	rspMsg := message.FindNodeRsp{
		NodeId:  recvNode.NodeId,
		FindId:  proto.String(store.GetRootId()),
		Addr:    proto.String(store.Root.Addr),
		IsProxy: proto.Bool(!store.Root.IsSuper),
		TcpPort: proto.Int32(int32(store.Root.TcpPort)),
		UdpPort: proto.Int32(int32(store.Root.UdpPort)),
	}
	resultBytes, _ := proto.Marshal(&rspMsg)
	session, ok := c.GetSession(msg.Name)
	if !ok {
		fmt.Println("这个session已经不存在了")
		return
	}
	err := session.Send(message.FindNodeRspNum, &resultBytes)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}
}

//
//处理查找节点请求
func (this *NodeManager) FindNodeReq(c engine.Controller, msg engine.GetPacket) {
	findNode := new(message.FindNodeReq)
	proto.Unmarshal(msg.Date, findNode)
	nodeStore := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)

	//--------------------------------------------
	//    查找节点
	//--------------------------------------------
	targetNode := nodeStore.Get(findNode.GetFindId(), true, findNode.GetNodeId())

	if targetNode.NodeId.String() == nodeStore.GetRootId() {
		//找到了
		// fmt.Println("查找到：", targetNode.NodeId.String())
		rspMsg := message.FindNodeRsp{
			NodeId:  findNode.NodeId,
			FindId:  proto.String(targetNode.NodeId.String()),
			Addr:    proto.String(targetNode.Addr),
			IsProxy: proto.Bool(!targetNode.IsSuper),
			TcpPort: proto.Int32(int32(targetNode.TcpPort)),
			UdpPort: proto.Int32(int32(targetNode.UdpPort)),
		}

		resultBytes, _ := proto.Marshal(&rspMsg)
		nodeResult := nodeStore.Get(findNode.GetNodeId(), false, "")
		if nodeResult == nil {
			return
		}
		session, ok := c.GetSession(nodeResult.NodeId.String())
		if !ok {
			fmt.Println("这个session已经不存在了")
			return
		}
		err := session.Send(message.FindNodeRspNum, &resultBytes)
		if err != nil {
			fmt.Println("node发送数据出错：", err.Error())
		}
		// c.GetNet().Send(msg., message.FindNodeRspNum, resultBytes)
		return
	}

	if targetNode.NodeId.String() == msg.Name {
		// fmt.Println("忽略这个查找")
		return
	}

	//转发出去
	session, ok := c.GetSession(targetNode.NodeId.String())
	if !ok {
		fmt.Println("这个session已经不存在了")
		return
	}
	err := session.Send(message.FindNodeReqNum, &msg.Date)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}
}

//处理查找节点返回请求
func (this *NodeManager) FindNodeRsp(c engine.Controller, msg engine.GetPacket) {
	recvNode := new(message.FindNodeRsp)
	proto.Unmarshal(msg.Date, recvNode)
	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)

	if recvNode.GetNodeId() == store.GetRootId() {
		//自己发出的查找请求
		if recvNode.GetFindId() == store.GetRootId() {
			//擦，把自己给找到老
			return
		}
		// nodeIdInt, _ := new(big.Int).SetString(*recvNode.NodeId, 10)
		shouldNodeInt, _ := new(big.Int).SetString(*recvNode.FindId, 10)
		newNode := &nodeStore.Node{
			NodeId:  shouldNodeInt,
			Addr:    *recvNode.Addr,
			IsSuper: !*recvNode.IsProxy,
			TcpPort: int32(*recvNode.TcpPort),
			UdpPort: int32(*recvNode.UdpPort),
		}
		// fmt.Println(*recvNode.NodeId)
		// _, ok := c.GetSession(*recvNode.NodeId)

		// if !ok {
		// 	fmt.Println(*recvNode.NodeId)
		// 	c.GetNet().AddClientConn(*recvNode.NodeId, *recvNode.Addr, store.GetRootId(), *recvNode.TcpPort, false)
		// }

		// // newNode
		// store.AddNode(newNode)

		isNeed, replace := store.CheckNeedNode(newNode.NodeId.String())
		// fmt.Println("这个节点是否需要：", isNeed)
		if isNeed {
			fmt.Println("接收请求:", *recvNode.FindId)
			store.AddNode(newNode)
			c.GetNet().AddClientConn(shouldNodeInt.String(), newNode.Addr, store.GetRootId(), newNode.TcpPort, false)
			if replace != "" {
				fmt.Println("替换原有的连接：", replace)
				//删除原来的连接
				if session, ok := c.GetSession(replace); ok {
					session.Close()
					delNode := new(nodeStore.Node)
					delNode.NodeId, _ = new(big.Int).SetString(replace, 10)
					store.DelNode(delNode)
					fmt.Println("替换成功")
				}
			}
		}
		return
	}

	//不是自己的回复，转发请求
	nodeForword := store.Get(recvNode.GetNodeId(), false, "")
	if nodeForword == nil {
		return
	}
	session, ok := c.GetSession(nodeForword.NodeId.String())
	if !ok {
		fmt.Println("这个session已经不存在了")
		return
	}
	err := session.Send(message.FindNodeRspNum, &msg.Date)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}
}

//查询相邻节点请求
func (this *NodeManager) FindRecentNodeReq(c engine.Controller, msg engine.GetPacket) {
	//--------------------------------------------
	//    查找邻居节点
	//--------------------------------------------
	recvNode := new(message.FindNodeRsp)
	proto.Unmarshal(msg.Date, recvNode)

	nodeIdInt, _ := new(big.Int).SetString(*recvNode.FindId, 10)
	var node *nodeStore.Node
	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	if recvNode.GetFindId() == "left" {
		node = store.GetLeftNode(*nodeIdInt)
		fmt.Println("-+-", node)
		if node == nil {
			return
		}
	}
	fmt.Println("---", node)
	rspMsg := message.FindNodeRsp{
		NodeId:  recvNode.NodeId,
		FindId:  proto.String(node.NodeId.String()),
		Addr:    proto.String(node.Addr),
		IsProxy: proto.Bool(!node.IsSuper),
		TcpPort: proto.Int32(int32(node.TcpPort)),
		UdpPort: proto.Int32(int32(node.UdpPort)),
	}
	resultBytes, _ := proto.Marshal(&rspMsg)
	session, ok := c.GetSession(msg.Name)
	if !ok {
		fmt.Println("这个session已经不存在了")
		return
	}
	err := session.Send(message.FindNodeRspNum, &resultBytes)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}
	if recvNode.GetFindId() == "right" {
		node = store.GetRightNode(*nodeIdInt)
		if node == nil {
			return
		}
	}
	rspMsg = message.FindNodeRsp{
		NodeId:  recvNode.NodeId,
		FindId:  proto.String(node.NodeId.String()),
		Addr:    proto.String(node.Addr),
		IsProxy: proto.Bool(!node.IsSuper),
		TcpPort: proto.Int32(int32(node.TcpPort)),
		UdpPort: proto.Int32(int32(node.UdpPort)),
	}
	resultBytes, _ = proto.Marshal(&rspMsg)
	err = session.Send(message.FindNodeRspNum, &resultBytes)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}

	// recvNode := new(message.FindRecentNodeReq)
	// proto.Unmarshal(msg.Date, recvNode)
	// store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	// allIds := store.GetRecentNodes()
	// if len(allIds) == 1 {
	// 	//这种情况一般是根节点
	// 	rspMsg := message.FindNodeRsp{
	// 		NodeId:  recvNode.NodeId,
	// 		FindId:  proto.String(store.GetRootId()),
	// 		Addr:    proto.String(store.Root.Addr),
	// 		IsProxy: proto.Bool(!store.Root.IsSuper),
	// 		TcpPort: proto.Int32(int32(store.Root.TcpPort)),
	// 		UdpPort: proto.Int32(int32(store.Root.UdpPort)),
	// 	}

	// 	resultBytes, _ := proto.Marshal(&rspMsg)
	// 	nodeResult := store.Get(recvNode.GetNodeId(), false, "")
	// 	if nodeResult == nil {
	// 		return
	// 	}
	// 	session, ok := c.GetSession(nodeResult.NodeId.String())
	// 	if !ok {
	// 		fmt.Println("这个session已经不存在了")
	// 		return
	// 	}
	// 	err := session.Send(message.FindNodeRspNum, &resultBytes)
	// 	if err != nil {
	// 		fmt.Println("node发送数据出错：", err.Error())
	// 	}
	// 	// c.GetNet().Send(msg., message.FindNodeRspNum, resultBytes)
	// 	return
	// }

	// switch recvNode.Cmp(allIds[0]) {
	// case 0:
	// case -1:
	// case 1:
	// }

	//转发出去
	// targetNode := store.Get(recvNode.GetNodeId(), false, "")
	// if targetNode == nil {
	// 	return
	// }
	// if recvNode.GetNodeId() == msg.Name {
	// 	// fmt.Println("忽略这个查找")
	// 	return
	// }
	// session, ok := c.GetSession(targetNode.NodeId.String())
	// if !ok {
	// 	fmt.Println("这个session已经不存在了")
	// 	return
	// }
	// err := session.Send(message.FindRecentNodeReqNum, &msg.Date)
	// if err != nil {
	// 	fmt.Println("node发送数据出错：", err.Error())
	// }

}

//注册节点请求
func (this *NodeManager) RegisterNodeReq(c engine.Controller, msg engine.GetPacket) {

}

//注册节点返回
func (this *NodeManager) RegisterNodeRsp(c engine.Controller, msg engine.GetPacket) {

}
