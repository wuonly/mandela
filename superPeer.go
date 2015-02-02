package mandela

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	//超级节点地址列表文件地址
	Path_SuperPeerAddress = "conf/nodeEntry.json"
)

//本地保存的超级节点地址列表
var Sys_superNodeEntry []string = []string{}

//清理本地保存的超级节点地址间隔时间
var Sys_cleanAddressTicker = time.Minute * 1

/*
	读取并解析本地的超级节点列表文件
*/
func init() {
	fileBytes, err := ioutil.ReadFile(Path_SuperPeerAddress)
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileBytes, &Sys_superNodeEntry); err != nil {
		return
	}
	go LoopCheckAddr()
}

/*
	隔时检查地址是否可用
*/
func LoopCheckAddr() {
	//获得一个心跳
	for range time.NewTicker(Sys_cleanAddressTicker).C {
		isChange := false
		tempAddrEntry := []string{}
		for _, addrOne := range Sys_superNodeEntry {
			if CheckOnline(addrOne) {
				tempAddrEntry = append(tempAddrEntry, addrOne)
			} else {
				isChange = true
			}
		}
		if isChange {
			Sys_superNodeEntry = tempAddrEntry
		}
	}
}
