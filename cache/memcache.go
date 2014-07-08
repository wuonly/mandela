package cache

import (
	"sync"
)

type Memcache struct {
	lock  *sync.RWMutex
	cache *Cache
}

func (this *Memcache) Add() {

}
