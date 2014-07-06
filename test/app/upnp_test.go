package app

import (
	"./upnp"
	"log"
	"testing"
	"time"
)

func TestUpnp(t *testing.T) {
	u := upnp.Upnp{}
	if !u.AddPortMapping(1990, 1990, "TCP") {
		log.Println("映射失败")
	}
	if !u.AddPortMapping(1991, 1991, "UDP") {
		log.Println("映射失败")
	}
	if !u.AddPortMapping(1992, 1992, "TCP") {
		log.Println("映射失败")
	}
	time.Sleep(time.Second * 6)
	u.Reclaim()
}
