package addr_manager

import (
	"encoding/json"
	"errors"
	"github.com/prestonTao/mandela/core/config"
	"math/rand"
	"net"
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

	startLoadChan   = make(chan bool, 1)     //当本机没有可用的超级节点地址，这里会收到一个信号
	subscribesChans = make([]chan string, 0) //当本机有可用的超级节点地址，这里会收到一个信
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
	@bool    contain    是否包含自己的地址
	@return  addr       随机获得的地址
*/
func GetSuperAddrOne(contain bool) (string, int, error) {
	addr, port := config.GetHost()
	myaddr := addr + strconv.Itoa(port)
	rand.Seed(int64(time.Now().Nanosecond()))
	for len(Sys_superNodeEntry) != 0 {
		if !contain && len(Sys_superNodeEntry) == 1 {
			if _, ok := Sys_superNodeEntry[myaddr]; ok {
				return "", 0, errors.New("超级节点地址只有自己")
			}
		}
		// 随机取[0-1000)
		r := rand.Intn(len(Sys_superNodeEntry))
		count := 0
		for key, _ := range Sys_superNodeEntry {
			if count == r {
				if !contain && key == myaddr {
					break
				}
				if CheckOnline(key) {
					host, portStr, _ := net.SplitHostPort(key)
					port, err := strconv.Atoi(portStr)
					if err != nil {
						return "", 0, errors.New("IP地址解析失败")
					}
					return host, port, nil
				} else {
					delete(Sys_superNodeEntry, key)
					break
				}
			}
			count = count + 1
		}
	}
	return "", 0, errors.New("没有可用的超级节点地址")
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

/*
	添加一个消息订阅
	当服务器有可用的ip地址时，广播给每一个订阅
*/
func AddSubscribe(c chan string) {
	subscribesChans = append(subscribesChans, c)
}

/*
	给所有订阅者广播消息
*/
func BroadcastSubscribe(addr string) {
	for _, one := range subscribesChans {
		select {
		case one <- addr:
		default:
		}
	}
}
