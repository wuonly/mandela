package app

import (
	"../peerNode"
	// "./net"
	server "../messageEngine"
	"../upnp"
	"log"
)

type Option struct {
	upnp     *upnp.Upnp
	peerNode *peerNode.NodeStore
	// socketServer *net.SocketServer
}

//启动
func (this *Option) StartUP() {
	this.startMsgEngine()
	this.startUpnp()
	this.startPeerNode()
}

func (this *Option) startMsgEngine() {
	// server.ConfigPath = "config.ini"
	server.IP = getLocalIntenetIp()
	server.PORT = 9090
	man := server.ServerManager{}

	man.Run()
}

func (this *Option) startUpnp() {
	//启动upnp模块
	this.upnp = new(upnp.Upnp)

	if !this.upnp.AddPortMapping(server.PORT, server.PORT, "TCP") && !this.upnp.Active {
		log.Println("设备不支持upnp协议")
	} else {
		for i := uint16(server.PORT); ; i++ {
			if this.upnp.AddPortMapping(server.PORT, int(i), "TCP") {
				log.Println("映射成功")
				//映射成功
				continue
			}
		}
	}
}

func (this *Option) startPeerNode() {
	// mappings := this.upnp.GetAllMapping()
	// mappingPort := mappings["TCP"][0][1]
	//启动dao模块
	this.peerNode = peerNode.NewNodeStore("", "")
}

//创建一个域名
func (this *Option) CreateDomain(name string) {

}

func (this *Option) Shoutdown() {
	//关闭upnp模块
	this.upnp.Reclaim()
}
