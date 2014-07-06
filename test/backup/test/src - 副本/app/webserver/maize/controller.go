package maize

import (
	//"../books"
	//"encoding/json"
	"fmt"
	//"html/template"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type ControllerManager interface {
	RunController()
	Redirect()
}
type Statics struct {
	Url string
}
type Template struct {
	Url    string
	Locals interface{}
}
type Redirect struct {
	Url string
}
type Forward struct {
	Url string
}

//=====================================
type ActionMethod struct {
	methodValue *reflect.Value //方法
	params      []reflect.Type //方法传入参数
}

type MaizeMux struct {
	ControllerParser *Parser //Controller解析器
}

func (m *MaizeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer m.deferRecover()
	/**
	 * 建立session
	 */
	sessionStore := app.sessionProvider.SessionStart(w, r) //得到SessionStore
forward:
	//判断是不是访问静态文件
	fmt.Println("本次请求路径:", r.URL.Path)
	if strings.HasPrefix(r.URL.Path, Static_URL) {
		fmt.Println("本次请求访问的是静态文件")

		RunStaticServer(w, r)
		return
	}
	/**
	 * static过滤器
	 * restful过滤器
	 */
	/**
	 * static路由器
	 * restful路由器
	 * 得到Controller
	 */
	controllerNumber, e := GetController(r.URL.Path)
	if e != nil {
		//未找到页面
		http.NotFound(w, r)
		return
	}
	/**
	 * 调用staticController
	 * 调用dynamicController
	 * 1.分析controller参数
	 * 2.注入参数
	 * 3.调用方法并返回值
	 * 4.分析返回值
	 */
	request := &Request{r, &sessionStore}
	m.ControllerParser = &Parser{w: w, r: request}
	result := m.ControllerParser.RunController(controllerNumber)
	/**
	 * 分析controller返回
	 */
	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
	}
	//返回静态文件
	structType := reflect.TypeOf(Statics{})
	if resultType.PkgPath() == structType.PkgPath() && resultType.Name() == structType.Name() {
		fmt.Println("返回的是静态文件")
		StaticServer(w, r, result.(Statics).Url)
		return
	}
	//返回模板
	structType = reflect.TypeOf(Template{})
	if resultType.PkgPath() == structType.PkgPath() && resultType.Name() == structType.Name() {
		fmt.Println("返回的是模板")
		return
	}
	//服务器重定向
	structType = reflect.TypeOf(Forward{})
	if resultType.PkgPath() == structType.PkgPath() && resultType.Name() == structType.Name() {
		fmt.Println("返回的是服务器重定向")
		url := result.(Forward).Url
		r.URL.Path = url
		goto forward
		return
	}
	//客户端重定向
	structType = reflect.TypeOf(Redirect{})
	if resultType.PkgPath() == structType.PkgPath() && resultType.Name() == structType.Name() {
		fmt.Println("返回的是客户端重定向")
		http.Redirect(w, r, result.(Redirect).Url, http.StatusFound)
		return
	}
	fmt.Println("返回的是其他类型")
	resultJson, _ := json.Marshal(result)
	fmt.Fprintf(w, string(resultJson))

}

func (m *MaizeMux) ResultJson() {}
func (m *MaizeMux) deferRecover() {
	if r := recover(); r != nil {

	}
}

type ControllerCommandInterface interface {
	DoFilter()
	DoController()
	Redirect()
	ResultJSONOrString()
}

type ControllerCommand struct {
	//filer
	//controller
}

func (c *ControllerCommand) DoFilter() {

}
func (c *ControllerCommand) DoController() {

}
func (c *ControllerCommand) Redirect() {

}
func (c *ControllerCommand) ResultJSONOrString() {

}
