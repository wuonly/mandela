package app

import (
	"./upnp"
	"reflect"
	"strings"
	// "crypto/rand"
	// "encoding/hex"
	// "fmt"
	// "bytes"
	"log"
	// "net"
	// "time"
	// "reflect"
	// "strconv"
)

var App *Application

//实例化App
//实例化生命周期
func init() {
	App = &Application{LifeCycle: LifeCycleSupport{}}
	var upnpModule LifeCycleListener
	var serverModule LifeCycleListener
	var webserverModule LifeCycleListener
	//upnp模块
	upnpModule = RegisterUpnp{}
	//server模块
	serverModule = RegisterServer{}
	//webserver模块
	webserverModule = RegisterWebServer{}
	//添加模块
	App.LifeCycle.AddListener(upnpModule)
	App.LifeCycle.AddListener(serverModule)
	App.LifeCycle.AddListener(webserverModule)

	//dao模块
	var daoModule LifeCycleListener
	daoModule = RegisterDao{}
	App.LifeCycle.AddListener(daoModule)
}

type Application struct {
	Started      bool                   //是否已经启动
	SuperNode    bool                   //是否是超级节点
	LocalWebAddr string                 //本地web服务器地址加端口号: 127.0.0.1:8080
	Modules      map[string]interface{} //各个模块
	LifeCycle    LifeCycleSupport       //生命周期
	Key          string                 //此节点的id
	Discover     *upnp.Discover         //upnp模块
}

//开始按照生命周期中的顺序启动应用
func (this *Application) StartUP() {
	ready := make(chan int)
	App.LifeCycle.FireLifeCycleEven(Test, ready)
	if App.Started {
		return
	}

	//1.使用upnp协议映射端口
	App.LifeCycle.FireLifeCycleEven(UPnp_Module, ready)
	//2.启动web服务
	App.LifeCycle.FireLifeCycleEven(WebServer_Module, ready)
	//3.打开浏览器 start http://127.0.0.1:80

	App.LifeCycle.FireLifeCycleEven(Init, ready)

	App.LifeCycle.FireLifeCycleEven(Before_start_even, ready)

	App.LifeCycle.FireLifeCycleEven(Start_even, ready)

	App.LifeCycle.FireLifeCycleEven(After_start_even, ready)

	// App.StartServer()

	listenersLen := len(App.LifeCycle.FindListeners())
	log.Println(listenersLen)
	for i := 0; i < listenersLen; i++ {
		log.Println("第n个模块启动就绪 ", i+1)
		<-ready
	}

	ok := <-ready
	if ok == 1 {
		log.Println("收到关闭程序消息")
		this.ShoutDwon()
	} else {
		log.Println("未收到关闭程序消息")
	}
}

//关闭应用
func (this *Application) ShoutDwon() {
	//关闭web服务器
	App.LifeCycle.FireLifeCycleEven(Before_stop_even, nil)
	//删除upnp端口映射
	App.LifeCycle.FireLifeCycleEven(Stop_even, nil)
	//是否需要从新启动应用
	App.LifeCycle.FireLifeCycleEven(After_stop_even, nil)
}

func (this *Application) AddModule(module interface{}) {
	if this.Modules == nil {
		this.Modules = make(map[string]interface{})
	}
	this.Modules[strings.Split(reflect.TypeOf(module).String(), ".")[0]] = module
}
