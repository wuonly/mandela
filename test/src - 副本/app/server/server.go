package server

import (
	"log"
	"net"
)

type AppManager interface {
	StartUP()
	ShoutDwon()
}

func StartServer(app *AppManager) {
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
					app.ShoutDwon()
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

func StopServer(app *AppManager) {

}

type Server struct {
	in  chan<- []byte
	out <-chan []byte
}

func newServer(app *AppManager) {
	s := &Server{}
	go s.read()
	go s.write()
	go s.coordinate()
}

func (s *Server) coordinate() {
	for {

	}
}
func (s *Server) write() {

}
func (s *Server) read() {

}
