package maize

import (
	//"fmt"
	"net/http"
	//"path"
	"path/filepath"
)

func RunStaticServer(w http.ResponseWriter, r *http.Request) {
	fileUrl := GetStaticURL()
	//##############################
	//这里的截串需要优化
	//##############################
	runes := []rune(r.URL.Path)[len(fileUrl):]
	//ss := make([]string, 0)
	ss := ""
	for _, s := range runes {
		//ss = append(ss, string(s))
		ss += string(s)
	}
	staticServer(w, r, filepath.FromSlash(ss))
}

func StaticServer(w http.ResponseWriter, r *http.Request, url string) {
	staticServer(w, r, url)
}
func staticServer(w http.ResponseWriter, r *http.Request, url string) {
	filePath := GetStaticPATH()
	http.ServeFile(w, r, filePath+url)
}
