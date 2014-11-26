package net

import (
	"sync"
)

type MsgHandler func(c Controller, msg GetPacket)

type RouterStore struct {
	lock     *sync.RWMutex
	handlers map[string]MsgHandler
}

func (this *RouterStore) AddRouter(url string, handler MsgHandler) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.handlers[url] = handler
}

func (this *RouterStore) GetHandler(url string) MsgHandler {
	this.lock.Lock()
	defer this.lock.Unlock()
	handler := this.handlers[url]
	return handler
}

func (this *RouterStore) getMapping() map[string]MsgHandler {
	return this.handlers
}

func NewRouter() *RouterStore {
	router := new(RouterStore)
	router.lock = new(sync.RWMutex)
	router.handlers = make(map[string]MsgHandler)
	return router
}
