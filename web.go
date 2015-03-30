package mandela

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net"
	"net/http"
	"strconv"
)

func StartWeb() {
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("mandela"))
	m.Use(sessions.Sessions("session", store))

	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
	}))
	m.Use(martini.Static("../../statics"))

	r := martini.NewRouter()

	r.Get("/", Home)
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
func Home(r render.Render, req *http.Request, session sessions.Session) {
	r.Redirect("/index.html")
}
