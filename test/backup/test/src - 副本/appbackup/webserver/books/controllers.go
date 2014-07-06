package books

import (
	maize "../maize"
	"fmt"
	"net/http"
)

type HomeController struct{}

func (c *HomeController) Home(w http.ResponseWriter, r *maize.Request) interface{} {
	fmt.Println("执行了home方法")
	static := maize.Statics{"/index.html"}
	return static
}

type Utt struct {
	Name string
	Age  int
}

func (c *HomeController) Struct(w http.ResponseWriter, r *maize.Request) interface{} {
	u := Utt{"tao", 23}
	return u
}
func (c *HomeController) String(w http.ResponseWriter, r *maize.Request) interface{} {
	return "it's result String"
}
func (c *HomeController) Forward(w http.ResponseWriter, r *maize.Request) interface{} {
	return maize.Forward{"/home"}
}

func (c *HomeController) Redirect(w http.ResponseWriter, r *maize.Request) interface{} {
	fmt.Println("执行了Redirect方法")
	return maize.Redirect{"/home"}
}

func ActionTest(w http.ResponseWriter, r *maize.Request) interface{} {
	fmt.Println("执行了Action方法")
	u := Utt{"tao", 23}
	return u
}
