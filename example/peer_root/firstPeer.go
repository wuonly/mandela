package main

import (
	"github.com/prestonTao/mandela"
)

func main() {
	StartUP()
}

func StartUP() {
	mandela.IsRoot = true
	mandela.StartCommandWindow()
}
