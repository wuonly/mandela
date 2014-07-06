package app

import (
// "./dao"
// "fmt"
// "math/big"
)

type RegisterDao struct {
}

func (this RegisterDao) LifeCycleEven(even LifeCycleEven) {
	if even.GetEvenType() == Store_Module {
		// App.Started = true
		this.Start()
		ready := even.GetData().(chan bool)
		ready <- true
	} else if even.GetEvenType() == Stop_even {
		this.Stop()
	}
}

// var store *NodeStore

func (this RegisterDao) Start() {
	// dao.NewNodeStore()
}

func (this RegisterDao) Stop() {

}

// type NodeStoreManager struct {
// }
// func (this *NodeStoreManager)
