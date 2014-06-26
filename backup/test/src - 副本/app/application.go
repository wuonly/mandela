package app

import (
	"./dao"
	"./upnp"
	"./webserver"
	"log"
	"net"
)

var App *Application

type Application struct {
	Started      bool                   //是否已经启动
	SuperNode    bool                   //是否是超级节点
	LocalWebAddr string                 //本地web服务器地址加端口号: 127.0.0.1:8080
	Modules      map[string]interface{} //各个模块
	Key          string                 //此节点的id
	MappingInfo  *upnp.MappingInfo      //upnp模块
}

func StartUP() {
	//随机产生一个nodeid
	nodeId := RandNodeId(512)
	node := dao.Node{NodeId: nodeId, BeingPinged: true}
	dao.NewNodeStore(node)
	defer func() {
		if r := recover(); r != nil {
			log.Println("发现设备超时")
			//超时了
		}
	}()
	App = &Application{}
	// ready := make(chan int)
	//upnp
	log.Println("11111111111111111")
	mappingInfo := upnp.NewPortMapping()
	if mappingInfo == nil {
		log.Println("没有发现支持upnp协议设备")
	}

	log.Println("22222222222222222")
	App.MappingInfo = mappingInfo
	//webserver

	log.Println("33333333333333333")
	go webserver.StartUP(80)
	//dao
	// dao.NewNodeStore()

	//last setup
	shoutDown := make(chan bool)
	go SocketServer(shoutDown)
	log.Println("444444444444444444444")
	isClose := <-shoutDown

	log.Println("5555555555555555555")
	log.Println(isClose)
	log.Println("关闭应用")
	//upnp
	upnp.DeletePortMapping()
	//webserver
	webserver.StopServer()
}

func ShoutDown() {
	log.Println("关闭应用")
	//upnp
	upnp.DeletePortMapping()
	//webserver
	webserver.StopServer()
}

func SocketServer(c chan bool) {
	listener, err := net.Listen("tcp", ":1990")
	defer listener.Close()
	if err != nil {
		log.Println("--------------1")
		log.Println(err)
	}
	log.Println("--------------11")
	var conn net.Conn
	log.Println("--------------12")

	log.Println("--------------13")
	started := true
	for started {

		conn, err = listener.Accept()
		if err != nil {

			log.Println("--------------2")
			log.Println(err.Error())
			started = false
		}

		buf := make([]byte, 3048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("--------------3")
			log.Println(err.Error())
			//判断错误类型是否是net.OpError
			// if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			// 	log.Println("就是这种")
			// }
		}

		runes := []rune(string(buf[:n+1]))
		start := 0
		for i, r := range runes {
			//关闭应用命令
			if i == 17 {
				if string(runes[start:i]) == "SHOUTDOWN THE APP" {
					// app.ShoutDown()
					conn.Close()
					started = false
					//停止应用
					log.Println("发送关闭程序信息")
					// this.ready <- true
					c <- true
					ShoutDown()
					break
				}
			}
			//#############################
			//判断字符串的结束,返回（0,EOF）
			//有没有更好的方式
			//#############################
			if r == 0 {
				command := string(runes[start:i])
				log.Println(command, i)
				conn.Close()
				break
			}
		}

	}
	defer conn.Close()
}

// func Multicast() {
// 	remotAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
// 	chk(err)
// 	locaAddr, err := net.ResolveUDPAddr("udp", "192.168.1.100:1991")
// 	chk(err)
// 	conn, err := net.ListenUDP("udp", locaAddr)
// 	chk(err)

// 	_, err = conn.WriteToUDP([]byte("M-SEARCH * HTTP/1.1\r\n"+
// 		"HOST: 239.255.255.250:1900\r\n"+
// 		"ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n"+
// 		"MAN: \"ssdp:discover\"\r\n"+
// 		"MX: 3\r\n"+
// 		"\r\n"), remotAddr)
// 	chk(err)
// 	buf := make([]byte, 1024)
// 	_, remoteAddr, err := conn.ReadFromUDP(buf)
// 	chk(err)
// 	fmt.Println(string(buf), "\n  |  ", remoteAddr.IP, "  |  ", remoteAddr.Port)
// }

// func chk(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
