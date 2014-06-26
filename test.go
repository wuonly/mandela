package main

import (
	"fmt"
)

func main() {
	s := new(Session)
	s.age = "tao"
	s.getAge()

}

var base = sessionBase{}

type sessionBase struct {
	name string
}

type Session struct {
	base
	age string
}

func (this *Session) getAge() {
	fmt.Println("getAge")
}
