package main

import (
	m "github.com/prestonTao/mandela"
)

func main() {
	StartUP()
}

func StartUP() {
	m.Init_role = m.C_role_super
	m.Mode_local = true
	m.StartCommandWindow()

}