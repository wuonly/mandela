package mandela

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

func StartCommandWindow() {
	//命令控制中心发送程序停止命令
	stopChan := make(chan bool, 1)
	//命令行输入的命令和参数
	lineChan := make(chan string, 1)

	reader := bufio.NewReader(os.Stdin)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	StartUpAuto()

	running := true
	for running {
		go ReadLine(reader, lineChan)
		select {
		case dataStr := <-lineChan:
			//执行命令
			CtlCenter(strings.Split(dataStr, " "), stopChan)
		case <-c:
			//Ctrl + c 退出程序
			fmt.Println("Ctrl + c 退出程序")
			running = false
		case <-stopChan:
			//stop 命令退出程序
			fmt.Println("stop 命令退出程序")
			running = false
		}
	}
}

func ReadLine(reader *bufio.Reader, c chan string) {
	data, _, _ := reader.ReadLine()
	c <- string(data)
}

/*
	命令控制中心
*/
func CtlCenter(commands []string, stopChan chan bool) {
	switch commands[0] {
	case "help":

	case "quit":
		stopChan <- false
	case "exit":
		stopChan <- false
	case "info":

	case "send":
		SendMsgAll(commands)
	case "see":
		SelectAllPeer(commands[1])
	}
}

/*
	查询自己保存的逻辑节点
*/
func SelectAllPeer(domain string) {
	switch domain {
	case "all":
		See()
	case "left":
		SeeLeftNode()
	case "right":
		SeeRightNode()
	}
}

/*
	给节点发送消息
*/
var count = 0

func SendMsgAll(commands []string) {
	if len(commands) == 1 {
		SendMsgForAll("hello " + strconv.Itoa(count))
		count += 1
	}
	if len(commands) == 3 {
		SendMsgForOne(commands[1], commands[2])
	}
}
