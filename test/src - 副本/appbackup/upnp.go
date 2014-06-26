package app

import (
	"./upnp"
	"fmt"
)

type RegisterUpnp struct {
}

func (this RegisterUpnp) LifeCycleEven(even LifeCycleEven) {
	if even.GetEvenType() == UPnp_Module {
		this.Start()
		ready := even.GetData().(chan bool)
		ready <- true
	} else if even.GetEvenType() == Stop_even {
		this.Stop()
	}
}
func (this RegisterUpnp) Start() {
	fmt.Println("start upnp")

	discover := upnp.NewPortMapping()
	App.AddModule(discover)
	App.Discover = discover

}
func (this RegisterUpnp) Stop() {
	for key, value := range App.Discover.MappingInfo.OutsideMappingPort {
		App.Discover.DeletePortMapping(value, key)
	}
}
