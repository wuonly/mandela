package upnp

import (
	// "fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	// "reflect"
	"strconv"
	"strings"
	"time"
)

//发送组播消息，要带上端口，格式如："239.255.255.250:1900"
func sendMulticastMsg(localAddr, groupIP, msg string, c chan string) {
	var conn *net.UDPConn
	defer func() {
		if r := recover(); r != nil {
			log.Println("发现设备超时")
			//超时了
		}
	}()
	go func(conn *net.UDPConn) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("发现设备连接未超时")
				//没超时
			}
		}()
		//超时时间为3秒
		time.Sleep(time.Second * 3)
		c <- ""
		conn.Close()
	}(conn)
	remotAddr, err := net.ResolveUDPAddr("udp", groupIP)
	if err != nil {
		log.Println("组播：组播地址格式不正确")
	}
	locaAddr, err := net.ResolveUDPAddr("udp", localAddr)

	if err != nil {
		log.Println("组播：本地ip地址格式不正确")
	}
	conn, err = net.ListenUDP("udp", locaAddr)
	defer conn.Close()
	if err != nil {
		log.Println("组播：监听udp出错")
	}
	_, err = conn.WriteToUDP([]byte(msg), remotAddr)
	if err != nil {
		log.Println("组播：发送msg到组播地址出错")
	}
	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		log.Println("组播：从组播地址接搜消息出错")
	}

	result := string(buf)
	c <- result
}

func dialTCPSendMsg(msg Msg) string {
	log.Println(msg.requestInfo.Host)
	conn, err := net.Dial("tcp", msg.requestInfo.Host)
	if err != nil {
		log.Println("TCP：连接网关设备出错")
	}
	_, err = conn.Write([]byte(msg.BuildString()))
	if err != nil {
		log.Println("TCP：向网关设备发送消息出错")
	}
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Println("TCP：从网关设备接收消息出错")
	}
	buf = make([]byte, 3048)
	_, err = conn.Read(buf)
	if err != nil {
		log.Println("TCP：从网关设备接收消息出错")
	}
	return string(buf)
}

func HttpURLConnect(msg Msg) string {

	body := msg.body.BuildXML()
	client := &http.Client{}
	// 第三个参数设置body部分
	reqest, _ := http.NewRequest(msg.requestInfo.Method, msg.requestInfo.Url, strings.NewReader(body))

	reqest.Proto = msg.requestInfo.Proto
	reqest.Host = msg.requestInfo.Host

	for key, value := range msg.headerMap {
		reqest.Header.Set(key, value)
	}

	reqest.Header.Set("Content-Length", strconv.Itoa(len([]byte(body))))

	response, _ := client.Do(reqest)

	resultBody, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	log.Println("http请求返回", response.StatusCode)
	if response.StatusCode == 200 {
		// log.Println(response.Header)
		return string(resultBody)
	}
	return ""
}

//获取本机能联网的ip地址
func getLocalIntenetIp() string {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}

// //获取本机能联网的ip地址
// func getLocalIntenetIp() string {
// 	/*
// 	  获得所有本机地址
// 	  判断能联网的ip地址
// 	*/
// 	// c := make(chan bool)
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("发现设备超时111111111111111111")
// 			//超时了
// 		}
// 	}()
// 	var conn net.Conn
// 	defer conn.Close()
// 	go func(conn net.Conn) {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				log.Println("发现设备连接未超时1111111111111")
// 				//没超时
// 			}
// 		}()
// 		//超时时间为3秒
// 		time.Sleep(time.Second * 10)
// 		conn.Close()
// 	}(conn)
// 	conn, err := net.Dial("udp", "google.com:80")
// 	if err != nil {
// 		log.Println("预计超时错误")
// 		log.Println(err.Error())
// 	}

// 	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
// 	return ip
// }
