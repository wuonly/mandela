package message

type FindNode struct {
	Timeout int32  `json:"timeout"`
	NodeId  string `json:"node_id"`
	WantId  string `json:"want_id"`
	FindId  string `json:"find_id"`
	IsProxy bool   `json:"is_proxy"`
	ProxyId string `json:"proxy_id"`
	Addr    string `json:"addr"`
	IsSuper bool   `json:"id_super"`
	TcpPort int32  `json:"tcp_port"`
	UdpPort int32  `json:"udp_port"`
	Status  int32  `json:"status"`
}

type Message struct {
	TargetId string `json:"target_id"`
	Content  []byte `json:"content"`
}
