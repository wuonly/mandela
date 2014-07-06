package messageEngine

import (
	"errors"
	"fmt"
	"sync"
)

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
func (this *sessionBase) Send(msgID uint32, data *[]byte) {}
func (this *sessionBase) Close()                          {}

type Session interface {
	Send(msgID uint32, data *[]byte)
	Close()
	Set(name string, value interface{})
	Get(name string) interface{}
}

type sessionProvider struct {
	lock *sync.RWMutex
	// store     map[int64]Session
	nameStore map[string]Session
}

func addSession(name string, session Session) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	sessionStore.nameStore[session.(*sessionBase).name] = session
	// sessionStore.store[sessionId] = session
}

func getSession(name string) (Session, bool) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	s, ok := sessionStore.nameStore[name]
	return s, ok
}

func removeSession(name string) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	delete(sessionStore.nameStore, name)
}

var sessionStore = new(sessionProvider)

func init() {
	sessionStore.lock = new(sync.RWMutex)
	sessionStore.nameStore = make(map[string]Session, 10000)
}

//---------------------------------
//---------------------------------
//---------------------------------
//---------------------------------
//---------------------------------
//---------------------------------

type Conn interface {
	Send(msgID uint32, data *[]byte)
	Close()
}

//根据sessionId保存连接
var connStore = make(map[uint64]Conn, 10000)

var chs_lock sync.RWMutex

//注册Conn
func addConn(id uint64, client Conn) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	connStore[id] = client
	fmt.Println("添加一个conn", client, id, connStore)
}

//移除某个Conn
func removeConn(id uint64) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	delete(connStore, id)
}

//获取某个Conn
func getConn(id uint64) (client Conn, err error) {
	chs_lock.RLock()
	defer chs_lock.RUnlock()
	client, ok := connStore[id]
	if !ok {
		err = errors.New(fmt.Sprintf("uid %x 没有对应的 Session", id))
		return nil, err
	}
	return client, nil
}
