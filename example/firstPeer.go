package main

import (
	"bufio"
	"github.com/prestonTao/mandela"
	"os"
	"strings"
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
	m := mandela.Manager{}
	m.IsRoot = true

	//---------------------------------------
	//  手动设置端口
	//---------------------------------------
	m.HostPort = 9990

	m.Run()
	running := true
	reader := bufio.NewReader(os.Stdin)

	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":

		case "q":
			running = false
		case "info":

		case "send":
		case "see":
			if len(commands) == 1 {
				m.See()
			}
			if len(commands) == 2 {
				if commands[1] == "left" {
					m.SeeLeftNode()
				}
				if commands[1] == "right" {
					m.SeeRightNode()
				}
			}
		case "cap":
		case "odp":
		case "cdp":
		case "dump":
		}
	}
}
