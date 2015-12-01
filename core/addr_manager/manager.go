package addr_manager

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
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

/*
	启动本地服务
*/
func init() {
	LoadAddrForAll()
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
	随机得到一个超级节点地址
	@return  addr  随机获得的地址
*/
func GetSuperAddrOne() (addr string, err error) {
	timens := int64(time.Now().Nanosecond())
	rand.Seed(timens)
	// 随机取[0-1000)
	r := rand.Intn(len(Sys_superNodeEntry))
	count := 0
	for key, _ := range Sys_superNodeEntry {
		addr = key
		if count == r {
			return key, nil
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
