package peerNode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var filePath = "conf/nodeEntry.json"
var serverAddrPort = "127.0.0.1:9981"

type NodeStoreManager struct {
	filePath       string
	superNodeEntry []string
}

func (this *NodeStoreManager) loadPeerEntry() {
	if this.filePath == "" {
		this.filePath = filePath
	}
	fileBytes, err := ioutil.ReadFile(this.filePath)
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileBytes, &this.superNodeEntry); err != nil {
		fmt.Println("文件格式错误")
		return
	}
	// addrPort := strings.Split(this.superNodeEntry[0], ":")

	// this.superNodeIp = addrPort[0]
	// portInt, _ := strconv.Atoi(addrPort[1])
	// this.superNodePort = portInt
}

// func (this *NodeStoreManager) loadHostPeerEntry() bool {
// 	if this.filePath == "" {
// 		this.filePath = filePath
// 	}
// 	fileBytes, err := ioutil.ReadFile(this.filePath)
// 	if err != nil {
// 		return false
// 	}
// 	if err = json.Unmarshal(fileBytes, &this.superNodeEntry); err != nil {
// 		fmt.Println("文件格式错误")
// 		return false
// 	}
// 	// addrPort := strings.Split(this.superNodeEntry[0], ":")

// 	// this.superNodeIp = addrPort[0]
// 	// portInt, _ := strconv.Atoi(addrPort[1])
// 	// this.superNodePort = portInt
// 	return true
// }

// func (this *NodeStoreManager) loadServerPeerEntry() {
// 	//本地没有地址列表文件
// 	//连接服务器获取地址列表
// 	conn, err := net.Dial("tcp4", serverAddrPort)
// 	if err != nil {
// 		fmt.Println("连接超级节点发布服务器失败")
// 		return
// 	}
// 	conn.Write([]byte(this.hostIp + ":" + strconv.Itoa(this.hostPort)))
// 	buf := make([]byte, 1024)
// 	n, _ := conn.Read(buf)
// 	//得到超级节点的ip地址和端口
// 	addrPort := strings.Split(string(buf[:n]), ":")

// 	this.superNodeIp = addrPort[0]
// 	portInt, _ := strconv.Atoi(addrPort[1])
// 	this.superNodePort = portInt
// }
func (this *NodeStoreManager) saveNodeEntry() {
	// ss := make([]string, 0)
	// ss = append(ss, "127.0.0.1:8080")
	// ss = append(ss, "192.168.1.200:9090")
	// bs, _ := json.Marshal(ss)

	// file, _ := os.Create(filePath)
	// file.Write(bs)
	// file.Close()
}

func NewNodeStoreManager() *NodeStoreManager {
	ns := new(NodeStoreManager)
	ns.loadPeerEntry()
	return ns
}
