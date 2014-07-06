package maize

import (
	"net/http"
)

type FilterInterface interface {
	doIn(*http.ResponseWriter, Request) interface{}
	doOut(*http.ResponseWriter, Request) interface{}
}

func Filter(url string, actionMethod interface{}) {

}

type FilterChain struct {
	chain []FilterInterface
}

//func (f *FilterChain) addFilter(filter FilterInterface) *FilterChain {
//	f.chain = append(f.chain, filter)
//	return &f
//}
func (f *FilterChain) doFilter() {
	//for filter := range f.chain{
	//	filter
	//}
}
