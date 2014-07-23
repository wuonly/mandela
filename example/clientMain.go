package main

import (
	"bufio"
	"github.com/prestonTao/mandela"
	"os"
	"strconv"
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
	count := 1

	m := mandela.Manager{}
	m.Run()
	running := true
	reader := bufio.NewReader(os.Stdin)

	// time.Sleep(time.S)
	m.SaveData("tao", "hongfei")

	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch command {
		case "help":

		case "q":
			running = false
		case "info":

		case "send":
			m.SendMsgForAll("hello " + strconv.Itoa(count))
			count += 1
		case "cap":
		case "odp":
		case "cdp":
		case "dump":
		}
	}
}
