package socks5

import (
	"errors"
	"fmt"
	"sync"
)

var connStore = make(map[int64]*Conn, 10000)

var chs_lock sync.RWMutex

//注册Conn
func addConn(id int64, client *Conn) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	connStore[id] = client
}

//移除某个Conn
func removeConn(id int64) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	delete(connStore, id)
}

//获取某个Conn
func getConn(id int64) (client *Conn, err error) {
	chs_lock.RLock()
	defer chs_lock.RUnlock()
	client, ok := connStore[id]
	if !ok {
		err = errors.New(fmt.Sprintf("uid %x 没有对应的 Session", id))
		return nil, err
	}
	return client, nil
}
