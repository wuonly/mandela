package main

import (
	m "github.com/prestonTao/mandela"
)

func main() {
	StartUP()
}

func StartUP() {
	m.IsRoot = false
	m.StartCommandWindow()

}
