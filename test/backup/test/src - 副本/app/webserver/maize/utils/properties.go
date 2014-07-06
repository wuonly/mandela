package utils

import (
	//"bytes"
	//"fmt"
	"os"
	//"unicode/utf8"
)

func NewProperties(path string) properties {
	return properties{make(map[string]string), path}
}

type properties struct {
	keyValues map[string]string
	filePath  string
}

func (p *properties) Get(key string) string {
	return p.keyValues[key]
}
func (p *properties) Set(key, value string) {
	p.keyValues[key] = value
}

func (p *properties) Load() error {
	//#############################
	//这里需要判断读入文件的编码
	//还要处理properties文件大小现在是1024byte
	//#############################
	f, e := os.Open(p.filePath)
	if e != nil {
		return e
	}
	defer f.Close()
	bytes := make([]byte, 1024)
	n, e := f.Read(bytes)
	runes := []rune(string(bytes[:n+1]))
	start := 0
	key := ""
	for i, r := range runes {
		//fmt.Println(i, " rune:", r, "---------string:", string(r))
		//#############################
		//换行码，linux和window系统不一样
		//#############################
		if r == 13 || r == 10 {
			if start != i {
				value := string(runes[start:i])
				p.Set(key, value)
			}
			start = i + 1
		}
		//key和value的分隔符
		//#############################
		//还要处理一行有多个分隔符的情况
		//#############################
		if string(r) == ":" || string(r) == "=" {
			key = string(runes[start:i])
			key = key
			start = i + 1
		}
		//#############################
		//判断字符串的结束,返回（0,EOF）
		//有没有更好的方式
		//#############################
		if r == 0 {
			value := string(runes[start:i])
			p.Set(key, value)
		}

	}
	return nil
}

func (p *properties) write() {

}

func (p *properties) read() {

}

func (p *properties) Write() {

}
