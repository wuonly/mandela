package main

import (
	"github.com/astaxie/beego"
	_ "mandela/web/routers"
)

func main() {
	beego.Run()
}
