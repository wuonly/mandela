package maize

import (
	//"../books"
	//"encoding/json"
	"fmt"
	//"html/template"
	"net/http"
	//"reflect"
	//"strings"
	"errors"
)

func init() {
	staticControllers = make([]ControllerFunc, 1)
	synamicControllers = make([]ControllerFunc, 1)
	staticControllerNumber = map[string]int{"": 0}
	dynamicControllerNumber = map[string]int{"": 0}
}

type ControllerFunc func(w http.ResponseWriter, r *Request) interface{}

var staticControllers []ControllerFunc
var synamicControllers []ControllerFunc

var staticControllerNumber map[string]int
var dynamicControllerNumber map[string]int

func Controller(url string, secu_url interface{}, actionMethod ControllerFunc, secu_coll interface{}) {
	fmt.Println("这是添加第几个controller：", len(staticControllerNumber))
	staticControllerNumber[url] = len(staticControllerNumber)
	staticControllers = append(staticControllers, actionMethod)
}

func GetController(url string) (int, error) {
	if number := staticControllerNumber[url]; number == 0 {
		return 0, errors.New("error")
	} else {
		return number, nil
	}
}

type Parser struct {
	w http.ResponseWriter
	r *Request
}

func (p *Parser) buildParam() {
	p.r.ParseForm()
}

func (p *Parser) RunController(collNumber int) interface{} {
	p.buildParam()
	collerFunc := staticControllers[collNumber]
	result := collerFunc(p.w, p.r)
	return result
}
