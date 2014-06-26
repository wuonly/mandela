package maize

import (
	"reflect"
	"strings"
)

type RouterInfo struct {
	actionName string //
	methodName string //funName
	url        string //url
}

type RouterProvider struct {
	actionMap  map[string]interface{}
	routerInfo []*RouterInfo
}

func (rp *RouterProvider) AddRouter(url, actionName_Method string, action interface{}) {
	if rp.actionMap == nil {
		rp.actionMap = make(map[string]interface{})
	}

	if rp.routerInfo == nil {
		rp.routerInfo = make([]*RouterInfo, 0)
	}
	//#############################################
	//去掉前后空格，对各种情况的字符判断
	//#############################################
	list := strings.Split(actionName_Method, ".")
	actionName := list[0]
	methodName := list[1]
	router := &RouterInfo{actionName: actionName, methodName: methodName, url: url}
	rp.routerInfo = append(rp.routerInfo, router)
	rp.actionMap[actionName] = action

}
func (rp *RouterProvider) GetMethod(url string) *reflect.Value {
	//#############################################
	//判断routerInfo是否为空
	//#############################################
	routerInfos := rp.routerInfo
	for _, router := range routerInfos {
		if router.url == url {
			action := rp.actionMap[router.actionName]
			method := reflect.ValueOf(action).MethodByName(router.methodName)
			return &method
		}
	}
	return nil
}
