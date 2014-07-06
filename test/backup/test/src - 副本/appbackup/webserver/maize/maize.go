package maize

import (
	//"encoding/json"
	"fmt"
	"net/http"
	//"os"
	"reflect"
	//"strings"
	"./session"
	"log"
	"net"
	"runtime"
	"strconv"
)

var app *Application

func Setting(config *Settings) {
	app.config = config
}
func Action(url string, secu_url interface{}, actionMethod interface{}, secu_coll interface{}) {
	methodType := reflect.TypeOf(actionMethod)
	methodNumIn := methodType.NumIn()
	values := make([]reflect.Type, methodNumIn)
	for i := 0; i < methodNumIn; i++ {
		argementType := methodType.In(i)
		fmt.Println("第", i, "个参数类型为：", argementType.Kind().String())

		values[i] = argementType

	}
	methodValue := reflect.ValueOf(actionMethod)
	action := &ActionMethod{methodValue: &methodValue, params: values}
	app.actions[url] = action
	fmt.Println(values)
	//p.params = values
}

type Application struct {
	server          *Server                  //TCP服务
	routerProvider  *RouterProvider          //路由
	config          *Settings                //系统配置
	actions         map[string]*ActionMethod //保存所有的Controller
	filters         map[string]*ActionMethod //过滤器
	attributes      map[string]interface{}   //保存的参数
	controller      *ControllerManager       //Controller的引用
	sessionProvider *session.Manager         //
}

func (a *Application) Init() {
}

func (a *Application) SetAttributes(key string, value interface{}) {
	a.attributes[key] = value
}
func (a *Application) GetAttributes(key string) interface{} {
	return a.attributes[key]
}

func init() {
	app = new(Application)
	app.attributes = make(map[string]interface{})
	app.filters = make(map[string]*ActionMethod)
	app.actions = make(map[string]*ActionMethod)
	manager, e := session.NewManager()
	if e != nil {

	}
	app.sessionProvider = manager
}

func analysisAction() {

}

type Server struct {
	http.Server
	listener net.Listener
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	srv.listener = l
	return srv.Serve(l)
}

func (srv *Server) Stop() {
	e := srv.listener.Close()
	if e != nil {

	}
}

func Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mux := &MaizeMux{}
	fmt.Println("start success!")

	httpserver := http.Server{Addr: ":" + strconv.Itoa(Port), Handler: mux}
	server := Server{Server: httpserver}
	app.server = &server
	e := server.ListenAndServe()
	if e != nil {
		log.Println("web服务器已停止")
	}
	// http.ListenAndServe(":80", mux)
}

func ShoutDown() {
	app.server.Stop()
}
