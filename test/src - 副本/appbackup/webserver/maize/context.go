package maize

import (
	"net/http"
	"reflect"
)

type Context struct {
	w           *http.ResponseWriter
	r           *http.Request
	methodValue *reflect.Value  //方法
	params      []reflect.Value //方法传入参数
	results     []reflect.Value //方法返回值
	resultsObj  interface{}     //controller返回值
}
