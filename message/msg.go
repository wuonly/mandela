package message

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

type Message struct {
	TargetId string `json:"target_id"`
	Content  []byte `json:"content"`
}
