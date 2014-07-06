package app

import (
	"./webserver"
	"fmt"
)

type RegisterWebServer struct {
}

func (this RegisterWebServer) LifeCycleEven(even LifeCycleEven) {
	if even.GetEvenType() == WebServer_Module {
		this.Start()
		ready := even.GetData().(chan int)
		ready <- 0
	} else if even.GetEvenType() == Before_stop_even {
		this.Stop()
	}
}
func (this RegisterWebServer) Start() {
	fmt.Println("webserver start")
	webserver.StartUP(80)
}

func (this RegisterWebServer) Stop() {
	webserver.StopServer()
}
