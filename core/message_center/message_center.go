package message_center

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	engine "github.com/prestonTao/mandela/core/net"
	"github.com/prestonTao/mandela/core/nodeStore"
	"github.com/prestonTao/mandela/core/utils"
)

const (
	FindNodeNum    = iota + 101 //查找结点服务id
	SendMessageNum              //发送消息服务id

	SaveKeyValueReqNum
	SaveKeyValueRspNum
)

/*
	查找节点序列化对象
*/
type FindNode struct {
	Timeout int32  `json:"timeout"`  //这个节点的超时时间
	NodeId  string `json:"node_id"`  //本机的idinfo字符串
	WantId  string `json:"want_id"`  //想要查找的id 16进制字符串
	FindId  string `json:"find_id"`  //找到后返回的idinfo字符串
	IsProxy bool   `json:"is_proxy"` //这个查找是否是代理查找
	ProxyId string `json:"proxy_id"` //被代理的节点idinfo字符串
	Addr    string `json:"addr"`     //查找到的节点ip地址
	IsSuper bool   `json:"id_super"` //查找到的节点是否是超级节点
	TcpPort int32  `json:"tcp_port"` //查找到的节点tcp端口号
	UdpPort int32  `json:"udp_port"` //查找到的节点udp端口号
	Status  int32  `json:"status"`   //查找到的节点状态
}

/*
	发送消息序列化对象
*/
type Message struct {
	TargetId   string `json:"target_id"`   //接收者id
	ProtoId    int    `json:"proto_id"`    //协议编号
	CreateTime int64  `json:"create_time"` //消息创建时间unix
	Sender     string `json:"sender"`      //发送者id
	ReplyTime  int64  `json:"reply_time"`  //消息回复时间unix
	Content    []byte `json:"content"`     //发送的内容
	Hash       string `json:"hash"`        //消息的hash值
	ReplyHash  string `json:"reply_hash"`  //回复消息的hash
	Accurate   bool   `json:"accurate"`    //是否准确发送给一个节点
}

/*
	检查该消息是否是自己的
	不是自己的则自动转发出去
*/
func IsSendToSelf(c engine.Controller, msg engine.GetPacket) bool {
	messageRecv := new(Message)
	err := json.Unmarshal(msg.Date, messageRecv)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if nodeStore.ParseId(nodeStore.GetRootIdInfoString()) == messageRecv.TargetId {
		//是自己的消息，则处理这个消息
		fmt.Println("11111111111")
		return true
	} else {
		//先判断是否在自己的代理节点中，普通节点没有代理节点
		if targetNode, ok := nodeStore.GetProxyNode(messageRecv.TargetId); ok {
			if session, ok := c.GetSession(string(targetNode.IdInfo.Build())); ok {
				err := session.Send(SendMessageNum, &msg.Date)
				if err != nil {
					//fmt.Println("message发送数据出错：", err.Error())
					utils.Log.Debug("message发送数据出错：%v", err)
				}
			} else {
				//这个节点离线了，想办法处理下
				utils.Log.Debug("这个节点离线了，想办法处理下")
			}
			return false
		}
		//先判断自己是不是超级节点
		if !nodeStore.Root.IsSuper {
			if session, ok := c.GetSession(nodeStore.SuperName); ok {
				err := session.Send(SendMessageNum, &msg.Date)
				if err != nil {
					fmt.Println("message发送数据出错：", err.Error())
				}
			} else {
				fmt.Println("超级节点不在了")
				//超级节点都不在了，搞个锤子
			}
			return false
		}
		//检查这个消息是否因该自己处理，还是转发出去
		targetNode := nodeStore.Get(messageRecv.TargetId, true, "")
		if targetNode == nil {
			utils.Log.Error("邻近节点查询出错")
			return false
		}
		//这个消息因该自己处理
		if string(targetNode.IdInfo.Build()) == nodeStore.GetRootIdInfoString() {
			//指定节点处理
			if messageRecv.Accurate {
				//可能是普通节点未连接，先把消息保存起来
				utils.Log.Warn("可能是普通节点未连接，先把消息保存起来")
				return false
			}
			//未指定节点处理，则自己处理
			fmt.Println("2222222222")
			return true
		}
		//这个消息不该自己处理，则转发出去
		// session, ok := c.GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
		session, ok := c.GetSession(string(targetNode.IdInfo.Build()))
		if !ok {
			//转发失败
			return false
		}
		// fmt.Println(session.GetName())
		err := session.Send(SendMessageNum, &msg.Date)
		if err != nil {
			fmt.Println("message发送数据出错：", err.Error())
		}
		return false

		//		//---------------------------------
		//		//先判断自己是不是超级节点
		//		if !nodeStore.Root.IsSuper {
		//			if session, ok := c.GetSession(nodeStore.SuperName); ok {
		//				err := session.Send(SendMessageNum, &msg.Date)
		//				if err != nil {
		//					fmt.Println("message发送数据出错：", err.Error())
		//				}
		//			} else {
		//				fmt.Println("超级节点不在了")
		//				//超级节点都不在了，搞个锤子
		//			}
		//			return false
		//		}
		//		/*
		//			XXXXXXXXXXXXXXXXXXXXXXXXX
		//			存在一种情况，被代理节点还未连接到该节点
		//			XXXXXXXXXXXXXXXXXXXXXXXXX
		//		*/
		//		//先判断是否在自己的代理节点中
		//		if targetNode, ok := nodeStore.GetProxyNode(messageRecv.TargetId); ok {
		//			if session, ok := c.GetSession(string(targetNode.IdInfo.Build())); ok {
		//				err := session.Send(SendMessageNum, &msg.Date)
		//				if err != nil {
		//					fmt.Println("message发送数据出错：", err.Error())
		//				}
		//			} else {
		//				//这个节点离线了，想办法处理下
		//			}
		//			return false
		//		}
		//		// fmt.Println("把消息转发出去")
		//		//最后转发出去
		//		targetNode := nodeStore.Get(messageRecv.TargetId, true, "")
		//		if targetNode == nil {
		//			return false
		//		}
		//		if string(targetNode.IdInfo.Build()) == nodeStore.GetRootIdInfoString() {
		//			targetNode = nodeStore.GetInAll(messageRecv.TargetId, true, "")
		//			if string(targetNode.IdInfo.Build()) == nodeStore.GetRootIdInfoString() {
		//				if !messageRecv.Accurate {
		//					fmt.Println("看下是不是这里的问题")
		//					return true
		//				} else {
		//					fmt.Println("这个精确发送的消息没人接收")
		//				}
		//				return false
		//			}
		//		}
		//		// session, ok := c.GetSession(hex.EncodeToString(targetNode.NodeId.Bytes()))
		//		session, ok := c.GetSession(string(targetNode.IdInfo.Build()))
		//		if !ok {
		//			return false
		//		}
		//		// fmt.Println(session.GetName())
		//		err := session.Send(SendMessageNum, &msg.Date)
		//		if err != nil {
		//			fmt.Println("message发送数据出错：", err.Error())
		//		}
		//		return false
	}
}

/*
	接收消息并显示或转发
*/
func RecvMsg(c engine.Controller, msg engine.GetPacket) {
	if IsSendToSelf(c, msg) {
		handlerProcess(c, msg)
	}
}

type NodeManager struct {
}

//查找结点消息
func (this *NodeManager) FindNode(c engine.Controller, msg engine.GetPacket) {
	findNode := new(FindNode)
	json.Unmarshal(msg.Date, findNode)

	findNodeIdInfo := new(nodeStore.IdInfo)
	json.Unmarshal([]byte(findNode.NodeId), findNodeIdInfo)

	// proto.Unmarshal(msg.Date, findNode)
	// store := c.GetAttribute("nodeStore").(*nodeStore.NodeManager)
	//--------------------------------------------
	//    接收查找请求
	//--------------------------------------------
	if findNode.FindId != "" {
		//普通节点收到自己发出的代理查找请求
		if findNode.IsProxy && (findNode.ProxyId == nodeStore.GetRootIdInfoString()) {
			// fmt.Println("自己保存这个节点", findNode)
			//自己保存这个节点
			this.saveNode(findNode, c)
			// if nodeStore.SuperName != findNode.FindId {
			// 	fmt.Println("自己获得的超级节点:", findNode.FindId)
			// }
			return
		}
		//是自己发出的非代理查找请求
		if findNode.NodeId == nodeStore.GetRootIdInfoString() {
			//是自己的代理节点发的请求
			if findNode.IsProxy {
				this.sendMsg(findNode.ProxyId, &msg.Date, c)
				return
			}
			//不是代理查找，自己接收这个
			// fmt.Println("找到一个节点：", findNode.GetFindId())
			this.saveNode(findNode, c)
			return
		}
		//不是自己发出的查找请求转发粗去
		//查找除了刚发过来的节点并且包括自己，的临近结点
		targetNode := nodeStore.Get(findNode.WantId, true, nodeStore.ParseId(msg.Name))
		//查找的就是自己，但这个节点已经下线
		// if hex.EncodeToString(targetNode.NodeId.Bytes()) == store.GetRootId() {
		if targetNode.IdInfo.GetId() == nodeStore.ParseId(nodeStore.GetRootIdInfoString()) {
			//这里要想个办法解决下
			fmt.Println("想办法解决下这个问题")
			fmt.Println("from:", nodeStore.ParseId(findNode.ProxyId), "\nwantId: ", findNode.WantId, "\ntargerNodeid: ", targetNode.IdInfo.GetId())
			// fmt.Println(findNode)
			return
		}
		//转发粗去
		this.sendMsg(string(targetNode.IdInfo.Build()), &msg.Date, c)
		return
	}

	//--------------------------------------------
	//    发出查找请求
	//--------------------------------------------
	//--------------------------------------------
	//    查找邻居节点
	//--------------------------------------------
	if findNode.WantId == "left" || findNode.WantId == "right" {
		//需要查找的节点id
		nodeIdInt, _ := new(big.Int).SetString(nodeStore.ParseId(findNode.ProxyId), nodeStore.IdStrBit)
		var nodes []*nodeStore.Node
		//查找左邻居节点
		if findNode.WantId == "left" {
			nodes = nodeStore.GetLeftNode(*nodeIdInt, nodeStore.MaxRecentCount)
			if nodes == nil {
				return
			}
		}
		//查找右邻居节点
		if findNode.WantId == "right" {
			nodes = nodeStore.GetRightNode(*nodeIdInt, nodeStore.MaxRecentCount)
			if nodes == nil {
				return
			}
		}
		//把找到的邻居节点返回给查找者
		for _, nodeOne := range nodes {
			rspMsg := FindNode{
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
			if findNode.IsProxy == true {
				// fmt.Println("查找邻居节点：", findNode.WantId, "  返回结果：", rspMsg)
			}
			resultBytes, _ := json.Marshal(&rspMsg)
			this.sendMsg(msg.Name, &resultBytes, c)
		}
		return
	}

	//自己的代理节点发过来的代理查找请求
	if findNode.IsProxy && (msg.Name == findNode.ProxyId) {
		if findNode.WantId == "left" || findNode.WantId == "right" {
			fmt.Println(findNode.WantId)
		}
		//超级节点刚上线
		if findNode.IsSuper {
			newNode := FindNode{
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
			this.saveNode(&newNode, c)
		}
		//查找超级节点并且包括自己，的临近结点
		targetNode := nodeStore.Get(findNode.WantId, true, nodeStore.ParseId(findNode.ProxyId))
		//要查找的节点就是自己，则发送给自己的代理节点
		if targetNode.IdInfo.GetId() == nodeStore.ParseId(nodeStore.GetRootIdInfoString()) {
			// fmt.Println("自己的代理节点发出的查找请求查找到临近结点：", targetNode.NodeId.String())
			rspMsg := FindNode{
				NodeId:  findNode.NodeId,
				WantId:  findNode.WantId,
				FindId:  nodeStore.GetRootIdInfoString(),
				IsProxy: findNode.IsProxy,
				ProxyId: findNode.ProxyId,
				Addr:    nodeStore.Root.Addr,
				IsSuper: nodeStore.Root.IsSuper,
				TcpPort: nodeStore.Root.TcpPort,
				UdpPort: nodeStore.Root.UdpPort,
			}
			resultBytes, _ := json.Marshal(&rspMsg)
			this.sendMsg(msg.Name, &resultBytes, c)
			//是自己的代理节点
			if !findNode.IsSuper {
				if _, ok := c.GetSession(findNode.ProxyId); ok {
					// fmt.Println("添加一个代理节点")
					findNodeIdInfo := new(nodeStore.IdInfo)
					json.Unmarshal([]byte(findNode.ProxyId), findNodeIdInfo)
					newNode := &nodeStore.Node{
						IdInfo:  *findNodeIdInfo,
						IsSuper: findNode.IsSuper,
						Addr:    findNode.Addr,
						TcpPort: findNode.TcpPort,
						UdpPort: findNode.UdpPort,
					}
					nodeStore.AddProxyNode(newNode)
				}
			}
			return
		}
		//转发代理查找请求
		rspMsg := FindNode{
			// NodeId:  nodeStore.GetRootIdInfoString(),
			NodeId:  findNode.NodeId,
			WantId:  findNode.WantId,
			IsProxy: findNode.IsProxy,
			ProxyId: findNode.ProxyId,
		}
		resultBytes, _ := json.Marshal(&rspMsg)
		this.sendMsg(string(targetNode.IdInfo.Build()), &resultBytes, c)
		return
	}

	//查找除了客户端节点并且包括自己的临近结点
	targetNode := nodeStore.Get(findNode.WantId, true, nodeStore.ParseId(msg.Name))
	//要查找的节点就是自己
	if targetNode.IdInfo.GetId() == nodeStore.ParseId(nodeStore.GetRootIdInfoString()) {
		rspMsg := FindNode{
			NodeId:  findNode.NodeId,
			WantId:  findNode.WantId,
			FindId:  nodeStore.GetRootIdInfoString(),
			IsProxy: findNode.IsProxy,
			ProxyId: findNode.ProxyId,
			Addr:    nodeStore.Root.Addr,
			IsSuper: nodeStore.Root.IsSuper,
			TcpPort: nodeStore.Root.TcpPort,
			UdpPort: nodeStore.Root.UdpPort,
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
		// fmt.Println("这个session已经不存在了:", nodeId)
		return
	}
	err := session.Send(FindNodeNum, data)
	if err != nil {
		// fmt.Println("node发送数据出错：", err.Error())
	}
}

/*
	自己保存这个节点，只能保存超级节点
*/
func (this *NodeManager) saveNode(findNode *FindNode, c engine.Controller) {
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

	//自己不会连接自己
	if nodeStore.GetRootIdInfoString() == string(newNode.IdInfo.Build()) {
		return
	}

	//是否需要这个节点
	if isNeed, replace := nodeStore.CheckNeedNode(findNodeIdInfo.GetId()); isNeed {
		nodeStore.AddNode(newNode)
		if replace != "" {
			nodeStore.DelNode(replace)
		}
		//自己是超级节点，就连接这个节点
		if nodeStore.Root.IsSuper {
			//检查这个session是否存在，不存在就连接它
			if _, ok := c.GetNet().GetSession(string(newNode.IdInfo.Build())); !ok {
				session, err := c.GetNet().AddClientConn(newNode.Addr, nodeStore.GetRootIdInfoString(), newNode.TcpPort, false)
				if err != nil {
					fmt.Println(newNode)
					fmt.Println("连接客户端出错")
				} else {
					introduceSelf(session)
				}
			}
		}
		//找到自己最近的节点
		recentNode := nodeStore.Get(nodeStore.ParseId(nodeStore.GetRootIdInfoString()), false, "")
		if recentNode == nil {
			return
		}
		//是否需要替换超级节点
		if string(recentNode.IdInfo.Build()) != nodeStore.SuperName {
			//需要替换

			//非超级节点的操作
			if !nodeStore.Root.IsSuper {
				oldSuperName := nodeStore.SuperName
				fmt.Println("旧的超级节点id ", oldSuperName)
				// fmt.Println("链接新的超级节点")
				session, _ := c.GetNet().AddClientConn(recentNode.Addr,
					nodeStore.GetRootIdInfoString(), recentNode.TcpPort, false)
				nodeStore.SuperName = session.GetName()
				fmt.Println("新的超级节点id ", nodeStore.SuperName)
				if session, ok := c.GetNet().GetSession(oldSuperName); ok {
					fmt.Println("close -------------------1")
					session.Close()
					// fmt.Println("关闭旧链接：", oldSuperName)
				}
				return
			}
			//超级节点的操作
			//把旧超级节点连接删除
			if session, ok := c.GetNet().GetSession(nodeStore.SuperName); ok {
				session.Close()
			}
			//设置新的连接为超级节点连接
			nodeStore.SuperName = string(recentNode.IdInfo.Build()) // session.GetName()

		}

		//		if !nodeStore.Root.IsSuper {
		//			//			fmt.Println("需要 , ", replace)

		//			nodeStore.DelNode(replace)
		//			fmt.Println("findNode.FindId  ", findNode.FindId)
		//			//代理节点查找的备用超级节点
		//			// if findNode.WantId == "left" || findNode.WantId == "right" {
		//			// 	fmt.Println("添加备用节点：", nodeStore.ParseId(findNode.FindId))
		//			// 	// store.AddNode(node)
		//			// 	// return
		//			// }
		//			recentNode := nodeStore.Get(nodeStore.ParseId(nodeStore.GetRootIdInfoString()), false, "")
		//			if recentNode == nil {
		//				return
		//			}
		//			//是否需要替换超级节点
		//			if string(recentNode.IdInfo.Build()) != nodeStore.SuperName {
		//				oldSuperName := nodeStore.SuperName
		//				fmt.Println("旧的超级节点id ", oldSuperName)
		//				// fmt.Println("链接新的超级节点")
		//				session, _ := c.GetNet().AddClientConn(recentNode.Addr,
		//					nodeStore.GetRootIdInfoString(), recentNode.TcpPort, false)

		//				nodeStore.SuperName = session.GetName()

		//				fmt.Println("新的超级节点id ", nodeStore.SuperName)

		//				//				if _, ok := c.GetNet().GetSession(nodeStore.SuperName); ok {
		//				//					// fmt.Println("链接成功：", nodeStore.SuperName)
		//				//				}
		//				if session, ok := c.GetNet().GetSession(oldSuperName); ok {
		//					fmt.Println("close -------------------1")
		//					session.Close()
		//					// fmt.Println("关闭旧链接：", oldSuperName)
		//				}
		//			}
		//			return
		//		}

		//		//把替换的节点连接删除
		//		if replace != "" {
		//			//是否要替换超级节点
		//			if session, ok := c.GetNet().GetSession(nodeStore.SuperName); ok {
		//				utils.Log.Debug("需要替换的超级节点 %s", nodeStore.ParseId(nodeStore.SuperName))
		//				utils.Log.Debug("要取代的超级节点 %s", replace)
		//				if replace == session.GetName() {
		//					// fmt.Println("close -------------------2")
		//					session.Close()
		//					session, _ := c.GetNet().AddClientConn(newNode.Addr, nodeStore.GetRootIdInfoString(), newNode.TcpPort, false)
		//					nodeStore.SuperName = session.GetName()
		//					introduceSelf(session)
		//				}
		//			}
		//			if session, ok := c.GetSession(replace); ok {
		//				// fmt.Println("close -------------------3")
		//				session.Close()
		//				nodeStore.DelNode(replace)
		//			}
		//		}
		//		if nodeStore.Root.IsSuper {

		//		}
	}
}

/*
	连接节点后，向节点介绍自己
*/
func introduceSelf(session engine.Session) {
	//用代理方式查找最近的超级节点
	nodeMsg := FindNode{
		NodeId:  session.GetName(),
		WantId:  nodeStore.ParseId(nodeStore.GetRootIdInfoString()),
		IsProxy: true,
		ProxyId: nodeStore.GetRootIdInfoString(),
		IsSuper: nodeStore.Root.IsSuper,
		Addr:    nodeStore.Root.Addr,
		TcpPort: nodeStore.Root.TcpPort,
		UdpPort: nodeStore.Root.UdpPort,
	}

	// resultBytes, _ := proto.Marshal(&nodeMsg)
	resultBytes, _ := json.Marshal(nodeMsg)

	session.Send(FindNodeNum, &resultBytes)
}

/*
	得到一条消息的hash值
*/
func GetHash(msg *Message) string {
	hash := sha256.New()
	hash.Write([]byte(msg.TargetId))
	binary.Write(hash, binary.BigEndian, uint64(msg.ProtoId))
	binary.Write(hash, binary.BigEndian, msg.CreateTime)
	// hash.Write([]byte(int64(msg.ProtoId)))
	// hash.Write([]byte(msg.CreateTime))
	hash.Write([]byte(msg.Sender))
	// hash.Write([]byte(msg.RecvTime))
	binary.Write(hash, binary.BigEndian, msg.ReplyTime)
	hash.Write(msg.Content)
	hash.Write([]byte(msg.ReplyHash))
	return hex.EncodeToString(hash.Sum(nil))
}
