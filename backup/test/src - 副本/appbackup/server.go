package app

import (
	// "./server"
	"fmt"
	"log"
	"net"
)

type RegisterServer struct {
	ready chan bool
}

func (this RegisterServer) LifeCycleEven(even LifeCycleEven) {
	if even.GetEvenType() == Before_start_even {
		this.Start()
		ready := even.GetData().(chan bool)
		ready <- true
		this.ready = ready
	} else if even.GetEvenType() == Stop_even {
		this.Stop()
	}
}
func (this RegisterServer) Start() {
	fmt.Println("server start")
	this.server()
}
func (this RegisterServer) Stop() {
}

func (this RegisterServer) server() {
	listener, err := net.Listen("tcp", ":1990")
	defer listener.Close()
	if err != nil {
		log.Println(err)
	}
	started := true
	for started {
		conn, err := listener.Accept()
		if err != nil {

			log.Println(err.Error())
			started = false
		}
		buf := make([]byte, 3048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
		}

		runes := []rune(string(buf[:n+1]))
		start := 0
		for i, r := range runes {
			//关闭应用命令
			if i == 17 {
				if string(runes[start:i]) == "SHOUTDOWN THE APP" {
					// app.ShoutDown()
					conn.Close()
					started = false
					//停止应用
					log.Println("发送关闭程序信息")
					this.ready <- true
				}
			}
			//#############################
			//判断字符串的结束,返回（0,EOF）
			//有没有更好的方式
			//#############################
			if r == 0 {
				command := string(runes[start:i])
				log.Println(command, i)

				conn.Write([]byte(command))
				conn.Close()

			}
		}

	}
}
