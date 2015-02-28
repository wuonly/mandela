package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	getAddrs()
	// addNode()
}

func addNode() {
	data := make(url.Values)
	data["address"] = []string{"127.0.0.1:9981"}

	resp, _ := http.PostForm("http://127.0.0.1:19981/add", data)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("返回结果：", string(body))
}

/*
	得到超级节点地址列表
*/
func getAddrs() {
	resp, _ := http.Get("http://mandela.io:19981")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("返回结果：", string(body))
}
