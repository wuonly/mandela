package net

import (
	"time"
)

type TimeOut struct {
	isTimeOutChan chan bool
	duration      time.Duration
	f             func()
}

func (this *TimeOut) Do(duration time.Duration) bool {
	this.duration = duration
	go this.run()

	select {
	case <-this.isTimeOutChan:
		close(this.isTimeOutChan)
		return false
	case <-time.After(this.duration):
		return true
	}
}

func (this *TimeOut) run() {
	this.f()
	this.isTimeOutChan <- false

}

func NewTimeOut(f func()) *TimeOut {
	to := TimeOut{
		isTimeOutChan: make(chan bool),
		f:             f,
	}
	return &to
}
