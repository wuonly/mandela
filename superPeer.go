/*
	获取超级节点地址方式：
		1.本地配置文件方式获取
		2.官方目录服务器获取

	工作流程：
		1.判断配置文件夹是否存在，不存在则创建空文件夹。
		2.读取本地超级节点地址文件，添加配置中的地址。
		3.添加官方地址。
		4.启动心跳检查本地地址是否可用。
*/
package mandela

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
// Local_
)

var (
	//官方节点地址
	Path_SuperPeerdomain = "mandela.io:9981"
	//官方目录服务器地址
	Path_DirectotyServerAddr = []string{"mandela.io:19981"}
)

var (
	//配置文件存放目录
	Path_configDir = "conf"
	//超级节点地址列表文件地址
	Path_SuperPeerAddress = filepath.Join(Path_configDir, "nodeEntry.json")

	//超级节点地址最大数量
	Sys_config_entryCount = 1000
	//本地保存的超级节点地址列表
	Sys_superNodeEntry = make(map[string]string, Sys_config_entryCount)
	//清理本地保存的超级节点地址间隔时间
	Sys_cleanAddressTicker = time.Minute * 1
	//需要关闭定时清理超级节点地址列表程序时，向它发送一个信号
	Sys_StopCleanSuperPeerEntry = make(chan bool)
)

func InitSuperPeer() {
	if Mode_local {
		Path_SuperPeerdomain = Init_LocalIP + ":9981"
		Path_DirectotyServerAddr = []string{Init_LocalIP + ":19981"}
	}
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
	CheckAddr()
	go func() {
		//获得一个心跳
		for range time.NewTicker(Sys_cleanAddressTicker).C {
			CheckAddr()
		}
	}()
}

/*
	从目录服务器获取超级节点地址
	@ ds   目录服务器地址列表
*/
func pullSuperPeerAddrForDS(ds []string) {
	for _, addrOne := range ds {
		resp, _ := http.Get("http://" + addrOne)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("返回结果：", string(body))
		parseSuperPeerEntry(body)
	}
}

/*
	读取并解析本地的超级节点列表文件
*/
func loadSuperPeerEntry() {
	fileBytes, err := ioutil.ReadFile(Path_SuperPeerAddress)
	if err != nil {
		return
	}
	parseSuperPeerEntry(fileBytes)
}

/*
	解析超级节点地址列表
*/
func parseSuperPeerEntry(fileBytes []byte) {
	var tempSuperPeerEntry map[string]string
	if err := json.Unmarshal(fileBytes, &tempSuperPeerEntry); err != nil {
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
	检查地址是否可用
*/
func CheckAddr() {
	/*
		先获得一个拷贝
	*/
	oldSuperPeerEntry := make(map[string]string)
	for key, value := range Sys_superNodeEntry {
		oldSuperPeerEntry[key] = value
	}
	/*
		一个地址一个地址的判断是否可用
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
	Sys_superNodeEntry[addr] = ""
}

/*
	随机得到一个超级节点地址
	@return  addr  随机获得的地址
*/
func getSuperAddrOne() (addr string) {
	timens := int64(time.Now().Nanosecond())
	rand.Seed(timens)
	fmt.Println(len(Sys_superNodeEntry))
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
