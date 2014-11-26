package net

import (
	"sync"
)

//session接口
type Session interface {
	Send(msgID uint32, data *[]byte) error //发送一个消息
	Set(name string, value interface{})    //保存到session中一个键值对
	Get(name string) interface{}           //从session中获得一个键值对
	GetName() string                       //得到这个session的名称
	SetName(name string)                   //设置这个session的名称
}

//实现session接口
type sessionBase struct {
	name      string
	attrbutes map[string]interface{}
}

func (this *sessionBase) Set(name string, value interface{}) {
	this.attrbutes[name] = value
}
func (this *sessionBase) Get(name string) interface{} {
	return this.attrbutes[name]
}
func (this *sessionBase) GetName() string {
	return this.name
}
func (this *sessionBase) SetName(name string) {
	this.name = name
}
func (this *sessionBase) Send(msgID uint32, data *[]byte) (err error) { return }

//session仓库，保存着所有session
type sessionStore struct {
	lock      *sync.RWMutex
	nameStore map[string]Session
}

func (this *sessionStore) addSession(name string, session Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.nameStore[session.GetName()] = session
}

func (this *sessionStore) getSession(name string) (Session, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	s, ok := this.nameStore[name]
	return s, ok
}

func (this *sessionStore) removeSession(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.nameStore, name)
}

func NewSessionStore() *sessionStore {
	sessionStore := new(sessionStore)
	sessionStore.lock = new(sync.RWMutex)
	sessionStore.nameStore = make(map[string]Session)
	return sessionStore
}
