package msgServer

import (
	"code.google.com/p/goprotobuf/proto"
	// "fmt"
	"mandela/peerNode"
	"mandela/peerNode/message"
	server "mandela/peerNode/messageEngine"
	"mandela/peerNode/messageEngine/net"
	"math/big"
)

func init() {
	nodeManager := new(NodeManager)
	server.AddRouter(message.FindNodeReqNum, nodeManager.FindNodeReq)
	server.AddRouter(message.FindNodeRspNum, nodeManager.FindNodeRsp)
}

type NodeManager struct {
}

//处理查找节点请求
func (this *NodeManager) FindNodeReq(c server.Controller, msg net.GetPacket) {
	findNode := new(message.FindNodeReq)
	proto.Unmarshal(msg.Date, findNode)
	nodeStore := c.GetAttribute("peerNode").(*peerNode.NodeStore)

	targetNode := nodeStore.Get(findNode.GetFindId())
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
	c.GetNet().Send(msg.ConnId, message.FindNodeRspNum, resultBytes)
}

//处理查找节点返回请求
func (this *NodeManager) FindNodeRsp(c server.Controller, msg net.GetPacket) {
	recvNode := new(message.FindNodeRsp)
	proto.Unmarshal(msg.Date, recvNode)
	// fmt.Println("接收到：", *recvNode.NodeId)
	nodeStore := c.GetAttribute("peerNode").(*peerNode.NodeStore)

	nodeIdInt, _ := new(big.Int).SetString(*recvNode.NodeId, 10)
	shouldNodeInt, _ := new(big.Int).SetString(*recvNode.FindId, 10)
	newNode := &peerNode.Node{
		NodeId:       nodeIdInt,
		NodeIdShould: shouldNodeInt,
		Addr:         *recvNode.Addr,
		IsSuper:      !*recvNode.IsProxy,
		TcpPort:      int(*recvNode.TcpPort),
		UdpPort:      int(*recvNode.UdpPort),
	}

	// newNode
	nodeStore.AddNode(newNode)
}

//注册节点请求
func (this *NodeManager) RegisterNodeReq(c server.Controller, msg net.GetPacket) {

}

//注册节点返回
func (this *NodeManager) RegisterNodeRsp(c server.Controller, msg net.GetPacket) {

}
