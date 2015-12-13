package net

import (
	"fmt"
	"net"
)

//实例化
var engine *Engine

/*
	启动一个消息引擎
*/
func InitEngine(name string) {
	engine = NewEngine(name)
}

/*
	注册一个普通消息
*/
func RegisterMsg(msgId int32, handler MsgHandler) {
	if msgId <= 100 {
		fmt.Println("该消息不能注册，消息编号0-100被系统占用。")
		return
	}
	AddRouter(msgId, handler)
}

func Listen(listener *net.TCPListener) {
	engine.run()
	engine.net.Listen(listener)
}

/*
	添加一个连接，给这个连接取一个名字，连接名字可以在自定义权限验证方法里面修改
	@powerful      是否是强连接
	@return  name  对方的名称
*/
func AddClientConn(ip string, port int32, powerful bool) (name string) {
	engine.run()
	session, err := engine.net.AddClientConn(ip, engine.name, port, powerful)
	if err != nil {
		fmt.Println("连接服务器失败")
		return ""
	}
	name = session.GetName()
	return
}

//给一个session绑定另一个名称
func LinkName(name string, session Session) {

}

//添加一个拦截器，所有消息到达业务方法之前都要经过拦截器处理
func AddInterceptor(itpr Interceptor) {
	engine.interceptor.addInterceptor(itpr)
}

//得到控制器
func GetController() Controller {
	return engine.controller
}

//设置自定义权限验证
func SetAuth(auth Auth) {
	if auth == nil {
		return
	}
	defaultAuth = auth
}

//设置关闭连接回调方法
func SetCloseCallback(call CloseCallback) {
	engine.net.closecallback = call
}
