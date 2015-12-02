package config

import (
	"github.com/prestonTao/mandela/core/utils"
	"github.com/prestonTao/upnp"
)

var sys_mapping = new(upnp.Upnp) //端口映射程序

/*
	判断自己是否有公网ip地址
	若支持upnp协议，则添加一个端口映射
*/
func portMapping() {
	// utils.Log.Debug("监听一个本地地址：%s:%d", Init_LocalIP, Init_LocalPort)

	// fmt.Println("监听一个本地地址：", Init_LocalIP, ":", Init_LocalPort)
	//本地地址是全球唯一公网地址
	if utils.IsOnlyIp(Init_LocalIP) {
		Init_IsGlobalOnlyAddress = true
		utils.Log.Debug("本机ip是公网全球唯一地址")
		return
	}
	//获得网关公网地址
	err := sys_mapping.ExternalIPAddr()
	if err != nil {
		// fmt.Println(err.Error())
		utils.Log.Warn("网关不支持端口映射", err)
		return
	}
	utils.Log.Debug("正在尝试端口映射")
	for i := 0; i < 1000; i++ {
		if err := sys_mapping.AddPortMapping(Init_LocalPort, Init_GatewayPort, "TCP"); err == nil {
			Init_IsMapping = true
			utils.Log.Debug("映射到公网地址：%s:%d", Init_GatewayAddress, Init_GatewayPort)
			return
		}
		Init_GatewayPort = Init_GatewayPort + 1
	}
	utils.Log.Warn("端口映射失败")
	// fmt.Println("端口映射失败")
}

/*
	关闭服务器时回收端口
*/
func Reclaim() {
	sys_mapping.Reclaim()
}
