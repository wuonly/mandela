package config

import (
	"github.com/prestonTao/upnp"
)

const (
	C_role_auto   = "auto"   //根据网络环境自适应
	C_role_client = "client" //客户端模式
	C_role_super  = "super"  //超级节点模式
	C_role_root   = "root"   //根节点模式
)

var (
	Init_IsSuperPeer               = false //有公网ip或添加了端口映射则是超级节点
	Init_GlobalUnicastAddress      = ""    //公网地址
	Init_GlobalUnicastAddress_port = 9981  //

	Sys_mapping = new(upnp.Upnp) //端口映射程序

	Init_LocalIP     = ""   //本地ip地址
	Init_LocalPort   = 9981 //本地监听端口
	Init_ExternalIP  = ""   //添加端口映射后的网关公网ip地址
	Init_MappingPort = 9981 //映射到路由器的端口

	// Mode_dev   = false       //是否是开发者模式
	Mode_local = true        //是否是局域网开发模式
	Init_role  = C_role_auto //服务器角色

)

var (
	// IsRoot        bool //是否是第一个节点
	SuperNodeIp   string
	SuperNodePort int
)
