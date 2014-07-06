package app

import ()

const (
	UPnp_Module      string = "upnp_module"
	WebServer_Module string = "webserver_module"
	Store_Module     string = "store_module"

	Test              string = "test"
	Init              string = "init"
	Before_start_even string = "before_start"
	Start_even        string = "start"
	After_start_even  string = "after_start"

	Before_stop_even string = "before_stop"
	Stop_even        string = "stop"
	After_stop_even  string = "after_stop"
)

type LifeCycleEven struct {
	Data      interface{}
	EvenType  string
	LifeCycle LifeCycle
}

func (e *LifeCycleEven) GetData() interface{} {
	return e.Data
}
func (e *LifeCycleEven) GetLifeCycle() LifeCycle {
	return e.LifeCycle
}
func (e *LifeCycleEven) GetEvenType() string {
	return e.EvenType
}

type LifeCycleListener interface {
	LifeCycleEven(even LifeCycleEven)
}

type LifeCycle interface {
	AddListener(listener LifeCycleListener)
	RemoveListener(listener LifeCycleListener)
	FindListeners() []LifeCycleListener
	Start()
	Stop()
}
type LifeCycleSupport struct {
	lifeCycle LifeCycle
	listeners []LifeCycleListener
}

func (support *LifeCycleSupport) AddListener(listener LifeCycleListener) {
	//#####################################
	//考虑support.listeners变量的线程安全性
	//#####################################
	support.listeners = append(support.listeners, listener)
}
func (support *LifeCycleSupport) RemoveListener(listener LifeCycleListener) {
	//#####################################
	//考虑support.listeners变量的线程安全性
	//#####################################
	for i, listenerOne := range support.listeners {
		if listenerOne == listener {
			support.listeners = append(support.listeners[:i], support.listeners[i+1:]...)
			break
		}
	}
}
func (support *LifeCycleSupport) FindListeners() []LifeCycleListener {
	return support.listeners
}
func (support *LifeCycleSupport) FireLifeCycleEven(evenType string, data interface{}) {
	//#####################################
	//考虑support.listeners变量的线程安全性
	//#####################################
	event := LifeCycleEven{data, evenType, support.lifeCycle}
	// fmt.Println(reflect.TypeOf(event))
	listenerClone := make([]LifeCycleListener, len(support.listeners))
	copy(listenerClone, support.listeners[:])
	for _, listener := range listenerClone {
		go listener.LifeCycleEven(event)
	}

}
