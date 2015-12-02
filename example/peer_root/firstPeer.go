package main

import (
	"github.com/prestonTao/mandela/core"
	"github.com/prestonTao/mandela/core/config"
)

func main() {
	StartUP()
}

func StartUP() {
	config.Init_role = config.C_role_root
	config.Mode_local = true
	core.StartCommandWindow()
}
