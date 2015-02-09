package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	//超级节点地址列表文件地址
	Path_SuperPeerAddress = "node_entry.json"
)

//保存文件配置
var config map[string]string

//超级节点地址最大数量
var Sys_config_entryCount = 10000

//本地保存的超级节点地址列表
var Sys_superNodeEntry = make(map[string]string, Sys_config_entryCount)

//清理本地保存的超级节点地址间隔时间
var Sys_cleanAddressTicker = time.Minute * 1

/*
	解析config.json文件
	解析node_entry.json文件
*/
func init() {
	/*
		解析config.json文件
	*/
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("read config.json file error: " + err.Error())
	}
	if err = json.Unmarshal(data, &config); err != nil {
		panic("marshal config.json file error: " + err.Error())
	}
	/*
		解析node_entry.json文件
	*/
	loadSuperPeerEntry()
	// LoopCheckAddr()
	go func() {
		//获得一个心跳
		for range time.NewTicker(Sys_cleanAddressTicker).C {
			LoopCheckAddr()
		}
	}()
}

func main() {
	server := new(Server)
	go server.start()

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":
			fmt.Println("       ------------------------------")
			fmt.Println("       stop    关闭目录服务器")
			fmt.Println("       ------------------------------")
		case "stop":
			running = false
		case "info":
		default:
			fmt.Println("--server: not find " + commands[0] + " command")
		}
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
	var tempSuperPeerEntry map[string]string
	if err = json.Unmarshal(fileBytes, &tempSuperPeerEntry); err != nil {
		return
	}
	for key, _ := range tempSuperPeerEntry {
		addSuperPeerAddr(key)
	}
}

/*
	定时检查地址是否可用
*/
func LoopCheckAddr() {
	/*
		先获得一个拷贝
	*/
	oldSuperPeerEntry := make(map[string]string)
	for key, value := range Sys_superNodeEntry {
		if key == "mandela.io:9981" {
			continue
		}
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

/*
	检查一个地址的计算机是否在线
	@return idOnline    是否在线
*/
func CheckOnline(addr string) (isOnline bool) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

type Server struct{}

func (this *Server) start() {
	http.HandleFunc("/", this.entry)
	http.HandleFunc("/add", this.addNode)
	fmt.Println("listen to ", config["listen_ip"], ":", config["listen_port"])
	fmt.Println("directory server startup...")
	err := http.ListenAndServe(config["listen_ip"]+":"+config["listen_port"], nil)
	if err != nil {
		panic("server error: " + err.Error())
	}
}

/*
	返回超级节点地址列表
*/
func (this *Server) entry(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(Sys_superNodeEntry)
	io.WriteString(w, string(str))
}

/*
	添加超级节点地址
*/
func (this *Server) addNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	address := r.FormValue("address")
	addSuperPeerAddr(address)
}
