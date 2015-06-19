package main

import (
	"github.com/prestonTao/mandela"
)

func main() {
	StartUP()
}

func StartUP() {
	mandela.Init_role = mandela.C_role_root
	mandela.Mode_local = true
	mandela.StartCommandWindow()
}
