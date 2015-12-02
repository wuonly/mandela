package addr_manager

import (
	"encoding/json"
	"errors"
	"github.com/prestonTao/mandela/core/config"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (

	//超级节点地址列表文件地址
	Path_SuperPeerAddress = filepath.Join(config.Path_configDir, "nodeEntry.json")

	//超级节点地址最大数量
	Sys_config_entryCount = 1000
	//本地保存的超级节点地址列表
	Sys_superNodeEntry = make(map[string]string, Sys_config_entryCount)
	//清理本地保存的超级节点地址间隔时间
	Sys_cleanAddressTicker = time.Minute * 1
	//需要关闭定时清理超级节点地址列表程序时，向它发送一个信号
	Sys_StopCleanSuperPeerEntry = make(chan bool)

	startLoadChan     = make(chan bool, 1) //当本机没有可用的超级节点地址，这里会收到一个信号
	AvailableAddrChan = make(chan bool, 1) //当本机有可用的超级节点地址，这里会收到一个信号

)

/*
	启动本地服务
*/
func init() {
	go smartLoadAddr()
	startLoadChan <- true
}

/*
	根据信号加载超级节点地址列表
*/
func smartLoadAddr() {
	for {
		<-startLoadChan
		LoadAddrForAll()
	}
}

/*
	从所有渠道加载超级节点地址列表
*/
func LoadAddrForAll() {
	//加载本地文件
	//官网获取
	//私网获取
	//局域网组播获取
	LoadByMulticast()
}

/*
	添加一个地址
*/
func AddSuperPeerAddr(addr string) {
	Sys_superNodeEntry[addr] = ""
}

/*
	随机得到一个可用的超级节点地址
	这个地址不能是自己的地址
	@return  addr  随机获得的地址
*/
func GetSuperAddrOne() (string, error) {
	addr, port := config.GetHost()
	myaddr := addr + strconv.Itoa(port)
	rand.Seed(int64(time.Now().Nanosecond()))
	for len(Sys_superNodeEntry) != 0 {
		if len(Sys_superNodeEntry) == 1 {
			if _, ok := Sys_superNodeEntry[myaddr]; ok {
				return "", errors.New("超级节点地址只有自己")
			}
		}
		// 随机取[0-1000)
		r := rand.Intn(len(Sys_superNodeEntry))
		count := 0
		for key, _ := range Sys_superNodeEntry {
			if count == r {
				if key == myaddr {
					break
				}
				if CheckOnline(key) {
					return key, nil
				} else {
					delete(Sys_superNodeEntry, key)
					break
				}
			}
			count = count + 1
		}
	}
	return "", errors.New("没有可用的超级节点地址")
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
