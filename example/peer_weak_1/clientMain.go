package main

import (
	"bufio"
	m "github.com/prestonTao/mandela"
	"os"
	"strconv"
	"strings"
	// "time"
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
	m.Mode_dev = true

	count := 1

	// m := mandela.Manager{}
	m.StartWeakPeer()
	running := true
	reader := bufio.NewReader(os.Stdin)

	// time.Sleep(time.S)
	m.SaveData("tao", "hongfei")

	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":

		case "q":
			running = false
		case "info":

		case "send":
			if len(commands) == 1 {
				m.SendMsgForAll("hello " + strconv.Itoa(count))
				count += 1
			}
			if len(commands) == 3 {
				m.SendMsgForOne(commands[1], commands[2])
			}

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
		case "odp":
		case "cdp":
		case "dump":
		}
	}
}
