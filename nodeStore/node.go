package nodeStore

import (
	// "crypto/rsa"
	// "fmt"
	"math/big"
	"time"
)

//保存节点的id
//ip地址
//不同协议的端口
type Node struct {
	NodeId               *big.Int  //节点id的10进制字符串
	IsSuper              bool      //是不是超级节点，超级节点有外网ip地址，可以为其他节点提供代理服务
	Addr                 string    //外网ip地址
	TcpPort              int32     //TCP端口
	UdpPort              int32     //UDP端口
	LastContactTimestamp time.Time //最后检查的时间戳
	// NodeIdShould         *big.Int  //影子id
	// Status               int       //节点状态，1：在线，2：正在查询中，3：下线
	// Out                  chan *Node //需要查询是否在线的节点
	// OverTime             time.Duration `1 * 60 * 60` //超时时间，单位为秒
	// SelectTime           time.Duration `5 * 60`      //查询时间，单位为秒
	// Key                  *rsa.PrivateKey //保存的公钥和私钥信息
}

// //节点一个小时查询是否在线
// func (this *Node) ticker() {
// 	go func(this *Node) {
// 		//睡眠一小时
// 		time.Sleep(time.Second * this.OverTime)
// 		this.Status = 2
// 		this.Out <- this
// 		go this.timeOut()
// 	}(this)
// }

// //节点发送查询请求5分钟没有回应的，标记为节点下线
// func (this *Node) timeOut() {
// 	fmt.Println("timeout")
// 	time.Sleep(time.Second * this.SelectTime)
// 	switch this.Status {
// 	case 2:
// 		this.Status = 3
// 	case 3:
// 	}
// }
