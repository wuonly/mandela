package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	addNode()
}

func addNode() {
	data := make(url.Values)
	data["address"] = []string{"127.0.0.1:9981"}

	resp, _ := http.PostForm("http://127.0.0.1:9981/add", data)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("返回结果：", string(body))
}
