package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	sample1()
}

func sample1() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/nimei", nimei)
	http.HandleFunc("/src.js", srcjs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
	fmt.Println("webServer startup...")
}

func hello(w http.ResponseWriter, r *http.Request) {
	requestStr := resolveRequest(r)
	log.Println(requestStr)

	connWeb, _ := net.Dial("tcp", "127.0.0.1:80")
	connWeb.Write([]byte(requestStr))
	buf := make([]byte, 1024)
	n, _ := connWeb.Read(buf)
	io.WriteString(w, string(buf[:n]))
}

func resolveRequest(r *http.Request) (result string) {
	result = r.Method + " " + r.URL.Path + " " + r.Proto + "\r\n"
	result += "Host: " + r.Host + "\r\n"
	for key, value := range r.Header {
		result += key + ": " + value[0] + "\r\n\r\n"
	}

	return
}

func nimei(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.URL.Path
	scheme := r.URL.Scheme
	log.Println(path, "     ", scheme)
	io.WriteString(w, "nimei")
}

func srcjs(w http.ResponseWriter, r *http.Request) {
	log.Println("srcjs")
	io.WriteString(w, "alert(1);")
}

//----------------------------------------------

type Server struct {
	connProxy net.Conn
	connWeb   net.Conn
}

func (this *Server) readProxy() {
	for {
		buf := make([]byte, 1024)
		n, _ := this.connProxy.Read(buf)
		log.Println(string(buf[:n]))
		this.connWeb.Write(buf[:n])

		buf = make([]byte, 1024)
		n, _ = this.connWeb.Read(buf)
		this.connProxy.Write(buf[:n])
	}
}

func startUP() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		connProxy, _ := ln.Accept()

		connWeb, _ := net.Dial("tcp", "127.0.0.1:80")

		s := Server{connProxy: connProxy, connWeb: connWeb}
		go s.readProxy()
	}
}
