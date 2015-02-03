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

//需要关闭定时清理超级节点地址列表程序时，向它发送一个信号
var Sys_StopCleanSuperPeerEntry = make(chan bool)

func init() {
	loadSuperPeerEntry()
}

/*
	读取并解析本地的超级节点列表文件
*/
func loadSuperPeerEntry() {
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
	关闭并重启读取并解析本地的超级节点列表文件程序
*/
func reloadSuperPeerEntry() {
	Sys_StopCleanSuperPeerEntry <- true
	loadSuperPeerEntry()
}

/*
	隔时检查地址是否可用
*/
func LoopCheckAddr() {
	//获得一个心跳
	ticker := time.NewTicker(Sys_cleanAddressTicker)
	select {
	case <-Sys_StopCleanSuperPeerEntry: //关闭
		return
	case <-ticker.C:
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

/*
	保存超级节点地址列表到本地配置文件
	@path  保存到本地的磁盘路径
*/
func saveSuperPeerEntry(path string) {

}
