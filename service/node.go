package service

import (
	// "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"github.com/prestonTao/mandela/message"
	engine "github.com/prestonTao/mandela/net"
	"github.com/prestonTao/mandela/nodeStore"
	"math/big"
)

// func init() {
// 	nodeManager := new(NodeManager)
// 	engine.AddRouter(message.FindNodeReqNum, nodeManager.FindNodeReq)
// 	engine.AddRouter(message.FindNodeRspNum, nodeManager.FindNodeRsp)
// }

type NodeManager struct {
}

//查找结点消息
func (this *NodeManager) FindNode(c engine.Controller, msg engine.GetPacket) {
	findNode := new(message.FindNode)
	json.Unmarshal(msg.Date, findNode)

	findNodeIdInfo := new(nodeStore.IdInfo)
	json.Unmarshal([]byte(findNode.NodeId), findNodeIdInfo)
	// proto.Unmarshal(msg.Date, findNode)
	store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	//--------------------------------------------
	//    接收查找请求
	//--------------------------------------------
	if findNode.FindId != "" {
		//普通节点收到自己发出的代理查找请求
		if findNode.IsProxy && (findNode.ProxyId == store.GetRootIdInfoString()) {
			//自己保存这个节点
			fmt.Println("自己获得的超级节点:", findNode.FindId)
			this.saveNode(findNode, store, c)
			return
		}
		//是自己发出的非代理查找请求
		if findNode.NodeId == store.GetRootIdInfoString() {
			//是自己的代理节点发的请求
			if findNode.IsProxy {
				this.sendMsg(findNode.ProxyId, &msg.Date, c)
				return
			}
			//不是代理查找，自己接收这个
			// fmt.Println("找到一个节点：", findNode.GetFindId())
			this.saveNode(findNode, store, c)
			return
		}
		//不是自己发出的查找请求转发粗去
		//查找除了刚发过来的节点并且包括自己，的临近结点
		targetNode := store.Get(findNode.WantId, true, nodeStore.ParseId(msg.Name))
		//查找的就是自己，但这个节点已经下线
		// if hex.EncodeToString(targetNode.NodeId.Bytes()) == store.GetRootId() {
		if targetNode.IdInfo.GetId() == nodeStore.ParseId(store.GetRootIdInfoString()) {
			//这里要想个办法解决下
			fmt.Println("想办法解决下这个问题")
			return
		}
		//转发粗去
		this.sendMsg(string(targetNode.IdInfo.Build()), &msg.Date, c)
		return
	}

	//--------------------------------------------
	//    发出查找请求
	//--------------------------------------------
	//自己的代理节点发过来的代理查找请求
	if findNode.IsProxy && (msg.Name == findNode.ProxyId) {
		//超级节点刚上线
		if findNode.IsSuper {
			newNode := message.FindNode{
				NodeId:  findNode.ProxyId,
				WantId:  findNode.WantId,
				FindId:  findNode.ProxyId,
				IsProxy: findNode.IsProxy,
				ProxyId: findNode.ProxyId,
				Addr:    findNode.Addr,
				IsSuper: findNode.IsSuper,
				TcpPort: findNode.TcpPort,
				UdpPort: findNode.UdpPort,
			}
			this.saveNode(&newNode, store, c)
		}
		//查找除了被代理的节点并且包括自己，的临近结点
		targetNode := store.Get(findNode.WantId, true, nodeStore.ParseId(findNode.ProxyId))
		//要查找的节点就是自己，则发送给自己的代理节点
		if targetNode.IdInfo.GetId() == nodeStore.ParseId(store.GetRootIdInfoString()) {
			// fmt.Println("自己的代理节点发出的查找请求查找到临近结点：", targetNode.NodeId.String())
			rspMsg := message.FindNode{
				NodeId:  findNode.NodeId,
				WantId:  findNode.WantId,
				FindId:  store.GetRootIdInfoString(),
				IsProxy: findNode.IsProxy,
				ProxyId: findNode.ProxyId,
				Addr:    store.Root.Addr,
				IsSuper: store.Root.IsSuper,
				TcpPort: store.Root.TcpPort,
				UdpPort: store.Root.UdpPort,
			}
			resultBytes, _ := json.Marshal(&rspMsg)
			this.sendMsg(msg.Name, &resultBytes, c)
			return
		}
		//转发代理查找请求
		rspMsg := message.FindNode{
			NodeId:  store.GetRootIdInfoString(),
			WantId:  findNode.WantId,
			IsProxy: findNode.IsProxy,
			ProxyId: findNode.ProxyId,
		}
		resultBytes, _ := json.Marshal(&rspMsg)
		this.sendMsg(string(targetNode.IdInfo.Build()), &resultBytes, c)
		return
	}

	//--------------------------------------------
	//    查找邻居节点，只有超级节点才会找邻居节点
	//--------------------------------------------
	if findNode.WantId == "left" || findNode.WantId == "right" {
		//不是代理查找
		nodeIdInt, _ := new(big.Int).SetString(nodeStore.ParseId(findNode.NodeId), nodeStore.IdStrBit)
		var nodes []*nodeStore.Node
		//查找左邻居节点
		if findNode.WantId == "left" {
			nodes = store.GetLeftNode(*nodeIdInt, store.MaxRecentCount)
			if nodes == nil {
				return
			}
		}
		//查找右邻居节点
		if findNode.WantId == "right" {
			nodes = store.GetRightNode(*nodeIdInt, store.MaxRecentCount)
			if nodes == nil {
				return
			}
		}
		//把找到的邻居节点返回给查找者
		for _, nodeOne := range nodes {
			rspMsg := message.FindNode{
				NodeId:  findNode.NodeId,
				WantId:  findNode.WantId,
				FindId:  string(nodeOne.IdInfo.Build()),
				IsProxy: findNode.IsProxy,
				ProxyId: findNode.ProxyId,
				Addr:    nodeOne.Addr,
				IsSuper: nodeOne.IsSuper,
				TcpPort: int32(nodeOne.TcpPort),
				UdpPort: int32(nodeOne.UdpPort),
			}
			resultBytes, _ := json.Marshal(&rspMsg)
			this.sendMsg(msg.Name, &resultBytes, c)
		}
		return
	}

	//查找除了客户端节点并且包括自己的临近结点
	targetNode := store.Get(findNode.WantId, true, nodeStore.ParseId(msg.Name))
	//要查找的节点就是自己
	if targetNode.IdInfo.GetId() == nodeStore.ParseId(store.GetRootIdInfoString()) {
		rspMsg := message.FindNode{
			NodeId:  findNode.NodeId,
			WantId:  findNode.WantId,
			FindId:  store.GetRootIdInfoString(),
			IsProxy: findNode.IsProxy,
			ProxyId: findNode.ProxyId,
			Addr:    store.Root.Addr,
			IsSuper: store.Root.IsSuper,
			TcpPort: store.Root.TcpPort,
			UdpPort: store.Root.UdpPort,
		}
		resultBytes, _ := json.Marshal(&rspMsg)
		this.sendMsg(msg.Name, &resultBytes, c)
		return
	}
	//要找的不是自己，则转发出去
	this.sendMsg(string(targetNode.IdInfo.Build()), &msg.Date, c)
}

func (this *NodeManager) sendMsg(nodeId string, data *[]byte, c engine.Controller) {
	session, ok := c.GetSession(nodeId)
	if !ok {
		fmt.Println("这个session已经不存在了")
		return
	}
	err := session.Send(message.FindNodeNum, data)
	if err != nil {
		fmt.Println("node发送数据出错：", err.Error())
	}
}

//自己保存这个节点，可以保存超级节点，也可以保存代理节点
func (this *NodeManager) saveNode(findNode *message.FindNode, store *nodeStore.NodeManager, c engine.Controller) {
	//自己不是超级节点
	if !store.Root.IsSuper {
		//查找到的节点和自己的超级节点不一样，则连接新的超级节点
		if store.SuperName != findNode.FindId {
			if session, ok := c.GetNet().GetSession(store.SuperName); ok {
				session.Close()
			}
			session, _ := c.GetNet().AddClientConn(findNode.Addr, store.GetRootIdInfoString(), findNode.TcpPort, false)
			store.SuperName = session.GetName()
		}
		return
	}

	findNodeIdInfo := new(nodeStore.IdInfo)
	json.Unmarshal([]byte(findNode.FindId), findNodeIdInfo)
	// nodeStore.Parse(findNode.FindId)
	// shouldNodeInt, _ := new(big.Int).SetString(findNode.FindId, nodeStore.IdStrBit)
	newNode := &nodeStore.Node{
		IdInfo:  *findNodeIdInfo,
		IsSuper: findNode.IsSuper,
		Addr:    findNode.Addr,
		TcpPort: findNode.TcpPort,
		UdpPort: findNode.UdpPort,
	}

	//是否需要这个节点
	if isNeed, replace := store.CheckNeedNode(findNodeIdInfo.GetId()); isNeed {
		store.AddNode(newNode)
		//把替换的节点连接删除
		if replace != "" {
			//是否要替换超级节点
			if session, ok := c.GetNet().GetSession(store.SuperName); ok {
				if replace == session.GetName() {
					session.Close()
					session, _ := c.GetNet().AddClientConn(newNode.Addr, store.GetRootIdInfoString(), newNode.TcpPort, false)
					store.SuperName = session.GetName()
				}
			}
			if session, ok := c.GetSession(replace); ok {
				session.Close()
				// delNode := new(nodeStore.Node)
				// delNode.NodeId, _ = new(big.Int).SetString(replace, nodeStore.IdStrBit)
				// delNode := &nodeStore.Node{IdInfo: nodeStore.IdInfo{Id: replace}}
				store.DelNode(replace)
			}
		}
		if store.Root.IsSuper {
			//检查这个session是否存在
			if _, ok := c.GetNet().GetSession(string(newNode.IdInfo.Build())); !ok {
				_, err := c.GetNet().AddClientConn(newNode.Addr, store.GetRootIdInfoString(), newNode.TcpPort, false)
				if err != nil {
					fmt.Println("连接客户端出错")
				}
			}
		}
	}
}
