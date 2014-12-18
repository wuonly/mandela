package net

import (
	"sync"
)

type MsgHandler func(c Controller, msg interface{})

type RouterStore struct {
	lock     *sync.RWMutex
	handlers map[int]MsgHandler
}

func (this *RouterStore) AddRouter(msgId int, handler MsgHandler) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.handlers[msgId] = handler
}

func (this *RouterStore) GetHandler(msgId int) MsgHandler {
	this.lock.Lock()
	defer this.lock.Unlock()

	handler := this.handlers[msgId]
	return handler
}

// func (this *RouterStore) getMapping() map[int]MsgHandler {
// 	return this.handlers
// }

func NewRouter() *RouterStore {
	router := new(RouterStore)
	router.lock = new(sync.RWMutex)
	router.handlers = make(map[int]MsgHandler)
	return router
}
