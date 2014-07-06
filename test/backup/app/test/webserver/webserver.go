// sessionManager00010 project main.go
package webserver

import (
	"./books"
	. "./maize"
	//"encoding/json"
	// "fmt"
	//"net/http"
	//"reflect"
)

var ()

func StartUP(port int) {
	SetConfig(port)
	Start()
}

func StopServer() {
	ShoutDown()
}

func init() {

	Controller("/home", One("ROLE_LOGIN"),
		(&books.HomeController{}).Home, All{"ROLE_LOGIN", "ROLE_SESSION"})

	Controller("/struct", One("ROLE_LOGIN"),
		(&books.HomeController{}).Struct, All{"ROLE_LOGIN", "ROLE_SESSION"})

	Controller("/string", One("ROLE_LOGIN"),
		(&books.HomeController{}).String, All{"ROLE_LOGIN", "ROLE_SESSION"})

	Controller("/forward", One("ROLE_LOGIN"),
		(&books.HomeController{}).Forward, All{"ROLE_LOGIN", "ROLE_SESSION"})

	Controller("/redirect", One("ROLE_LOGIN"),
		(&books.HomeController{}).Redirect, All{"ROLE_LOGIN", "ROLE_SESSION"})

	Controller("/actiontest", One("ROLE_LOGIN"), books.ActionTest, One("ROLE_LOGIN"))
	//maize.Setting(settings)
	//maize.NewFilter()
}

func SetConfig(port int) {
	//设置本地监听的端口
	Port = port
	//模板访问路径
	Template_URL = "/templates"
	//模板文件存放文件夹
	Template_PATH = "books/templates"
	//静态文件访问路径
	Static_URL = "/statics"
	//静态文件存放文件夹
	Static_PATH = "books/statics"
	//实体Bean
	Modules = []interface{}{
		&books.User{},
		&books.Role{},
	}
}
