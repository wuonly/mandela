package mandela

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	Path_SuperPeerdomain = "mandela.io"
)

var (
	//配置文件存放目录
	Path_configDir = "conf"
	//超级节点地址列表文件地址
	Path_SuperPeerAddress = filepath.Join("conf", "nodeEntry.json")

	//超级节点地址最大数量
	Sys_config_entryCount = 1000
	//本地保存的超级节点地址列表
	Sys_superNodeEntry = make(map[string]string, Sys_config_entryCount)
	//清理本地保存的超级节点地址间隔时间
	Sys_cleanAddressTicker = time.Minute * 1
	//需要关闭定时清理超级节点地址列表程序时，向它发送一个信号
	Sys_StopCleanSuperPeerEntry = make(chan bool)
)

func init() {
	//判断文件夹是否存在
	if _, err := os.Stat(Path_configDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(Path_configDir, 0755)
		}
		panic(err.Error())
	}
}

/*
	开始加载超级节点地址
*/
func startLoadSuperPeer() {
	Sys_superNodeEntry[Path_SuperPeerdomain] = ""
	loadSuperPeerEntry()
	LoopCheckAddr()
	go func() {
		//获得一个心跳
		for range time.NewTicker(Sys_cleanAddressTicker).C {
			LoopCheckAddr()
		}
	}()
}

/*
	读取并解析本地的超级节点列表文件
*/
func loadSuperPeerEntry() {
	fileBytes, err := ioutil.ReadFile(Path_SuperPeerAddress)
	if err != nil {
		return
	}
	var tempSuperPeerEntry map[string]string
	if err = json.Unmarshal(fileBytes, &tempSuperPeerEntry); err != nil {
		return
	}
	for key, _ := range tempSuperPeerEntry {
		addSuperPeerAddr(key)
	}
	addSuperPeerAddr(Path_SuperPeerdomain)
}

/*
	关闭并重启读取并解析本地的超级节点列表文件程序
*/
// func reloadSuperPeerEntry() {
// 	Sys_StopCleanSuperPeerEntry <- true
// 	loadSuperPeerEntry()
// }

/*
	定时检查地址是否可用
*/
func LoopCheckAddr() {
	/*
		先获得一个拷贝
	*/
	oldSuperPeerEntry := make(map[string]string)
	for key, value := range Sys_superNodeEntry {
		oldSuperPeerEntry[key] = value
	}
	/*
		一个地址一个地址判断是否可用
	*/
	for key, _ := range oldSuperPeerEntry {
		if CheckOnline(key) {
			addSuperPeerAddr(key)
		} else {
			delete(Sys_superNodeEntry, key)
		}
	}
}

/*
	添加一个地址
*/
func addSuperPeerAddr(addr string) {
	if Mode_dev && addr == "mandela.io:9981" {
		return
	}
	Sys_superNodeEntry[addr] = ""
}

/*
	随机得到一个超级节点地址
	@return  addr  随机获得的地址
*/
func getSuperAddrOne() (addr string) {
	timens := int64(time.Now().Nanosecond())
	rand.Seed(timens)
	// 随机取[0-1000)
	r := rand.Intn(len(Sys_superNodeEntry))
	count := 0
	for key, _ := range Sys_superNodeEntry {
		addr = key
		if count == r {
			return key
		}
		count = count + 1
	}
	return
}

/*
	保存超级节点地址列表到本地配置文件
	@path  保存到本地的磁盘路径
*/
func saveSuperPeerEntry(path string) {
	fileBytes, _ := json.Marshal(Sys_superNodeEntry)
	file, _ := os.Create(path)
	file.Write(fileBytes)
	file.Close()
}
