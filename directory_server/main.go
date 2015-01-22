package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// "time"
)

//保存文件配置
var config map[string]string

//保存超级节点ip地址和端口
var nodeEntry []string

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
	data, err = ioutil.ReadFile("node_entry.json")
	if err != nil {
		panic("read node_entry.json file error: " + err.Error())
	}
	if err = json.Unmarshal(data, &nodeEntry); err != nil {
		panic("marshal node_entry.json file error: " + err.Error())
	}
}

func main() {
	server := new(Server)
	server.Run()

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":
			fmt.Println(`       ------------------------------
       stop    关闭目录服务器
       ------------------------------`)
		case "stop":
			running = false
		case "info":
		default:
			fmt.Println("--server: not find " + commands[0] + " command")
		}
	}
}

type Server struct{}

func (this *Server) Run() {
	go this.start()
	go this.hold()
}

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
func (this *Server) hold() {
	// for range time.NewTicker(time.Hour).C {
	// 	for _, node := range nodeEntry {

	// 	}
	// }
}

/*
	返回超级节点地址列表
*/
func (this *Server) entry(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(nodeEntry)
	io.WriteString(w, string(str))
}

/*
	添加超级节点地址
*/
func (this *Server) addNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	address := r.FormValue("address")
	nodeEntry = append(nodeEntry, address)
}
