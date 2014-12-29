package service

import (
	// "code.google.com/p/goprotobuf/proto"
	// "github.com/prestonTao/mandela/message"
	// "github.com/prestonTao/mandela/nodeStore"
	// "fmt"
	// "github.com/prestonTao/mandela/cache"
	engine "github.com/prestonTao/mandela/net"
	// "math/big"
)

type DataStore struct {
}

func (this *DataStore) SaveDataReq(c engine.Controller, msg engine.GetPacket) {
	// memcache := c.GetAttribute("cache").(*cache.Memcache)
	// fmt.Println("收到数据：", string(msg.Date))
	// data := make(map[string]interface{})
	// data["tao"] = "hongfei"
	// memcache.Add(data)
}

func (this *DataStore) SaveDataRsp(c engine.Controller, msg engine.GetPacket) {

}
