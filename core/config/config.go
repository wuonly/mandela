package config

import (
	"github.com/prestonTao/mandela/core/utils"
)

const (
	C_Server_name = "mandela" //网络名称

	//服务器角色，只有局域网开发模式才能用
	C_role_auto   = "auto"   //根据网络环境自适应
	C_role_client = "client" //客户端模式
	C_role_super  = "super"  //超级节点模式
	C_role_root   = "root"   //根节点模式
)

var (
	//配置文件存放目录
	Path_configDir = "conf"

	Init_IsGlobalOnlyAddress = false //本地ip是否是公网全球唯一ip
	Init_LocalIP             = ""    //本地ip地址(局域网ip或公网全球唯一ip)
	Init_LocalPort           = 9981  //本地监听端口

	Init_IsMapping = false //是否映射了端口
	// Init_IsSuperPeer    = false //有公网ip或添加了端口映射则是超级节点
	Init_GatewayAddress = ""   //网关地址
	Init_GatewayPort    = 9981 //网关端口

	// Init_ExternalIP  = ""   //添加端口映射后的网关公网ip地址
	// Init_MappingPort = 9981 //映射到路由器的端口

	// Mode_dev   = false       //是否是开发者模式
	Mode_local = true        //是否是局域网开发模式
	Init_role  = C_role_auto //服务器角色

)

var (
	// IsRoot        bool //是否是第一个节点
	SuperNodeIp   string
	SuperNodePort int
)

func init() {
	utils.GlobalInit("console", "", "debug", 1)
	// utils.GlobalInit("file", `{"filename":"/var/log/gd/gd.log"}`, "", 1000)
	// utils.Log.Debug("session handle receive, %d, %v", msg.Code(), msg.Content())
	utils.Log.Debug("test debug")
	utils.Log.Warn("test warn")
	utils.Log.Error("test error")
	AutoRole()
}

/*
	根据网络情况自己确定节点角色
*/
func AutoRole() {
	if Mode_local {
		utils.Log.Debug("局域网模式")
		Init_LocalIP = "127.0.0.1"
	}
	//尝试端口映射
	if !Mode_local {
		portMapping()
	}
	//得到本地ip地址
	if address, ok := utils.GetLocalIntenetIp(); ok {
		Init_LocalIP = address
	} else {
		Init_LocalIP = utils.GetLocalHost()
	}
	//得到本机可用端口
	Init_LocalPort = utils.GetAvailablePortForTCP()
	utils.Log.Debug("本机角色为：%s", Init_role)
	utils.Log.Debug("本机监听地址：%s:%d", Init_LocalIP, Init_LocalPort)
}

/*
	获得本机是否是超级节点
*/
func CheckIsSuperPeer() bool {
	if Mode_local {
		switch Init_role {
		case C_role_client:
			return false
		case C_role_root:
			return true
		case C_role_super:
			return true
		case C_role_auto:
			return false
		}
		return false
	}
	if Init_IsGlobalOnlyAddress {
		return true
	}
	if Init_IsMapping {
		return true
	}
	return false
}
