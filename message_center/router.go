package message_center

import (
	"fmt"
	engine "github.com/prestonTao/mandela/net"
	"sync"
	"time"
)

//******************************************
//  注册所有消息
//      1.提供所有消息注册
//      2.保存所有消息编码和回调方法的映射
//      3.通过消息编码获得回调方法
//******************************************

type MsgHandler func(c engine.Controller, packet engine.GetPacket, msg *Message) (bool, string)

var handlersMapping = make(map[int32]MsgHandler)
var router_lock = new(sync.RWMutex)

func addRouter(msgId int32, handler MsgHandler) {
	router_lock.Lock()
	defer router_lock.Unlock()
	handlersMapping[msgId] = handler
}

func getHandler(msgId int32) MsgHandler {
	router_lock.Lock()
	defer router_lock.Unlock()
	handler := handlersMapping[msgId]
	return handler
}

/*
	消息分发程序
*/
func handlerProcess(c engine.Controller, packet engine.GetPacket, msg *Message) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				fmt.Println("网络库：", e.Error())
			}
		}
	}()
	handler := getHandler(msg.ProtoId)
	if handler == nil {
		fmt.Println("消息中心：未注册的消息编号-", packet.MsgID)
		return
	}
	ok, str := handler(c, packet, msg)
	fire(ok, str)
}

//******************************************
//  消息管理中心
//      1.消息超时机制
//******************************************

var timeoutMapping map[int]map[int]string //{ 消息协议id | 时间纳秒数 | 目标peer id字符串 }
var timeoutLock *sync.RWMutex = new(sync.RWMutex)

/*
	添加一个映射
*/
func addTimeoutMapping(protoId, ticker int, id string) {
	timeoutLock.Lock()
	defer timeoutLock.Unlock()
	if timeoutMapping == nil {

	}
}

type Pipe struct {
	c chan string
}

/*
	设置超时
*/
func (this *Pipe) timeout() {

	c := make(chan bool, 1)

	select {
	case <-c:

	case <-time.NewTicker(time.Second * 10).C:
		//超时了
		this.c <- "timeout"
	}
}

/*
	1
*/
func getTimeOutMsg(id string) (c chan bool) {
	time.Now().UnixNano()
	c = make(chan bool, 1)
	return
}

/*
	回调函数执行成功后，给此方法一个事件，看是否执行成功
	用于消除消息的超时
*/
func fire(ok bool, result string) {

}
