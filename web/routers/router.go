package routers

import (
	"github.com/astaxie/beego"
	"mandela/web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
