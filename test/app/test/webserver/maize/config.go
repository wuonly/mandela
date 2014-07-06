package maize

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	Port          int    = 8080
	WebAppPath    string = Getwd()      //项目路径
	Template_URL  string = "/templates" //模板路径
	Template_PATH string = "/templates" //模板文件夹
	Static_URL    string = "/statics"   //静态文件
	Static_PATH   string = "/statics"   //静态文件夹
	Routers       map[string]interface{}
	Modules       []interface{}
)

type Settings struct {
	Template_URL  string //模板路径
	Template_PATH string //模板文件夹
	Static_URL    string //静态文件
	Static_PATH   string //静态文件夹
	Modules       []interface{}
}

func Getwd() string {
	pwd, e := os.Getwd()
	if e == nil {
		return pwd
	}
	return ""
}

//检查Static_URL参数是否合法
func checkStatic_URL() (url string, boo bool) {
	//包含空格
	if strings.Contains(Static_URL, " ") {
		url = Static_URL
		boo = false
		return
	}
	//以"/"符号开头
	if strings.IndexRune(Static_URL, '/') != 0 {
		url = Static_URL
		boo = false
		return
	}
	//路径中包含多个"/"符号
	if strings.Count(Static_URL, "/") != 1 {
		url = Static_URL
		boo = false
		return
	}
	url = Static_URL
	boo = true
	return
}

func GetStaticPATH() string {
	//是否是根目录
	if !strings.HasPrefix(Static_PATH, "/") {
		return WebAppPath + "\\" + filepath.FromSlash(path.Clean(Static_PATH))
	} else {
		return Static_PATH
	}
	//filepath.FromSlash(path.Clean(Static_URL))

}

func GetStaticURL() string {
	if url, boo := checkStatic_URL(); boo {
		return url
	} else {
		return ""
	}
}
