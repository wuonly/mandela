package addr_manager

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
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
