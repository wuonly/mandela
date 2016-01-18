package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/prestonTao/mandela/core"
	"net"
	"net/http"
	"strconv"
	"time"
)

func StartWeb() {
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("mandela"))
	m.Use(sessions.Sessions("session", store))

	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
	}))
	m.Use(martini.Static("../../web/statics"))

	r := martini.NewRouter()

	r.Get("/", Home_handler)               //首页
	r.Get("/getdomain", GetDomain_handler) //获得本机的域名

	m.Action(r.Handle)

	webPort := 80
	for i := 0; i < 1000; i++ {
		_, err := net.ListenPacket("udp", ":"+strconv.Itoa(webPort))
		if err != nil {
			webPort = webPort + 1
		} else {
			break
		}
	}
	m.RunOnAddr(":" + strconv.Itoa(webPort))
	m.Run()
}

/*
	首页
*/
func Home_handler(r render.Render, req *http.Request, session sessions.Session) {
	r.Redirect("/index.html")
}

/*
	获得本机域名
*/
func GetDomain_handler() map[string]interface{} {
	retmap := make(map[string]interface{})
	if len(core.Init_IdInfo.Id) == 0 {
		retmap["ret"] = -1
	} else {
		retmap["ret"] = 0
		retmap["domain"] = core.Init_IdInfo.Domain
	}
	return retmap
}

/*
	创建域名
*/
func CreateDomain_handler(params martini.Params) map[string]interface{} {
	domain := params["domain"]
	name := params["name"]
	email := params["email"]
	core.CreateAccount(name, email, domain)

	retmap := make(map[string]interface{})
	retmap["ret"] = 0
	return retmap
}

/*
	发送一个消息
*/
func SendMsg_handler(params martini.Params) map[string]interface{} {
	tid := params["tid"]
	msg := params["msg"]
	core.SendMsgForOne_opt(tid, msg)

	retmap := make(map[string]interface{})
	retmap["ret"] = 0
	return retmap
}

/*
	获得消息
*/
var webMsgChan = make(chan string, 1000)

func GetMessage_handler() map[string]interface{} {
	retmap := make(map[string]interface{})
	select {
	case <-time.NewTicker(time.Minute).C:
		retmap["ret"] = -1
	case msg := <-webMsgChan:
		retmap["ret"] = 0
		retmap["msg"] = msg
	}
	return retmap
}
