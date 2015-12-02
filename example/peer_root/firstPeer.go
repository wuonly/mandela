package main

import (
	m "github.com/prestonTao/mandela/core"
	"github.com/prestonTao/mandela/core/config"
)

func main() {
	StartUP()
}

func StartUP() {
	config.Init_role = config.C_role_super
	config.Mode_local = true
	m.StartCommandWindow()
}
