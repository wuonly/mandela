package service

import (
	"code.google.com/p/goprotobuf/proto"
	// "fmt"
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
			if session, ok := c.GetSession(replace); ok {
				session.Close()
			}
		}
	}
}

//
//处理查找节点请求
func (this *NodeManager) FindNodeReq(c engine.Controller, msg engine.GetPacket) {
	findNode := new(message.FindNodeReq)
	proto.Unmarshal(msg.Date, findNode)
	nodeStore := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)

	targetNode := nodeStore.Get(findNode.GetFindId(), true, msg.Name)

	if targetNode == nil {
		targetNode = nodeStore.Root
	}

	// fmt.Println("查找到：", targetNode.NodeId.String())
	rspMsg := message.FindNodeRsp{
		NodeId:  proto.String(targetNode.NodeId.String()),
		FindId:  findNode.FindId,
		Addr:    proto.String(targetNode.Addr),
		IsProxy: proto.Bool(!targetNode.IsSuper),
		TcpPort: proto.Int32(int32(targetNode.TcpPort)),
		UdpPort: proto.Int32(int32(targetNode.UdpPort)),
	}

	resultBytes, _ := proto.Marshal(&rspMsg)
	session, _ := c.GetSession(msg.Name)
	session.Send(message.FindNodeRspNum, &resultBytes)
	// c.GetNet().Send(msg., message.FindNodeRspNum, resultBytes)
}

//处理查找节点返回请求
func (this *NodeManager) FindNodeRsp(c engine.Controller, msg engine.GetPacket) {
	recvNode := new(message.FindNodeRsp)
	proto.Unmarshal(msg.Date, recvNode)
	// fmt.Println("接收到：", *recvNode.NodeId)
	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)

	nodeIdInt, _ := new(big.Int).SetString(*recvNode.NodeId, 10)
	shouldNodeInt, _ := new(big.Int).SetString(*recvNode.FindId, 10)
	newNode := &nodeStore.Node{
		NodeId:       nodeIdInt,
		NodeIdShould: shouldNodeInt,
		Addr:         *recvNode.Addr,
		IsSuper:      !*recvNode.IsProxy,
		TcpPort:      int32(*recvNode.TcpPort),
		UdpPort:      int32(*recvNode.UdpPort),
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
		store.AddNode(newNode)
		c.GetNet().AddClientConn(nodeIdInt.String(), newNode.Addr, store.GetRootId(), newNode.TcpPort, false)
		if replace != "" {
			//删除原来的连接
			if session, ok := c.GetSession(replace); ok {
				session.Close()
			}
		}
	}
}

//注册节点请求
func (this *NodeManager) RegisterNodeReq(c engine.Controller, msg engine.GetPacket) {

}

//注册节点返回
func (this *NodeManager) RegisterNodeRsp(c engine.Controller, msg engine.GetPacket) {

}
