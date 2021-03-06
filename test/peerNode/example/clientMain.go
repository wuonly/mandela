package main

import (
	"bufio"
	"github.com/prestonTao/mandela/peerNode"
	"os"
)

func main() {
	StartUP()
}

// func StartUp() {
// 	m := peerNode.Manager{}
// 	m.IsRoot = true
// 	m.Run()
// }

func StartUP() {
	m := peerNode.Manager{}
	m.Run()
	running := true
	reader := bufio.NewReader(os.Stdin)

	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch command {
		case "help":

		case "q":
			running = false
		case "info":

		case "oap":
		case "cap":
		case "odp":
		case "cdp":
		case "dump":
		}
	}
}
