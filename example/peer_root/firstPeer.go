package main

import (
	"github.com/prestonTao/mandela/core"
)

func main() {
	StartUP()
}

func StartUP() {
	core.Init_role = core.C_role_root
	core.Mode_local = true
	core.StartCommandWindow()
}
