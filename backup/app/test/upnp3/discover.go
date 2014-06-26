package upnp

import (
	// "../utils"
	"container/list"
	"log"
	"strconv"
	"strings"
)

//本地ip及端口
//外网ip及端口
type MappingInfo struct {
	GetewayInsideIP    string         //网关内部ip
	GetewayOutsideIP   string         //网关外部ip
	OutsideMappingPort map[string]int //网关外部端口
	InsideMappingPort  map[string]int //映射到本机的端口
}

type Geteway struct {
	GetewayName   string //网关名称
	Host          string //网关ip和端口
	DeviceDescUrl string //网关设备描述路径
	Cache         string //cache
	ST            string
	USN           string
	deviceType    string //设备的urn   "urn:schemas-upnp-org:service:WANIPConnection:1"
	ControlURL    string //设备端口映射请求路径
	ServiceType   string //提供upnp服务的服务类型
}

type DiscoverInfo struct {
	DefaultSearchType []string `[urn:schemas-upnp-org:device:InternetGatewayDevice:1,
	urn:schemas-upnp-org:service:WANIPConnection:1,
	urn:schemas-upnp-org:service:WANPPPConnection:1]`
	MappingPorts *list.List `[1990,1991,1992]`
	Protocols    []string   `[TCP,UDP]` //名称一定要大写
}

type Discover struct {
	MappingInfo            MappingInfo //端口映射信息
	GetwayInfo             Geteway     //网关信息
	DiscoverInfo           DiscoverInfo
	GroupIp                string `"239.255.255.250:1900"`
	TimeOut                int    `3000`
	LocalIPPort            string //本机ip地址
	ExternalIPPort         string //外网ip地址
	TCPLocalMappingPort    string //本地映射的TCP端口
	TCPExternalMappingPort string //外网映射的TCP端口
	UDPLocalMappingPort    string //本地映射的UDP端口
	UDPExternalMappingPort string //外网映射的UDP端口
}

func (d *Discover) init() {
	//构造一个映射端口号队列列表
	list := list.New()
	startPort := 1990
	for i := 0; i < 10; i++ {
		list.PushBack(startPort)
		startPort++
	}
	d.MappingInfo = MappingInfo{OutsideMappingPort: make(map[string]int, 2), InsideMappingPort: make(map[string]int, 2)}
	d.GroupIp = "239.255.255.250:1900"
	d.DiscoverInfo = DiscoverInfo{MappingPorts: list,
		Protocols: []string{"TCP", "UDP"},
		DefaultSearchType: []string{"urn:schemas-upnp-org:service:WANIPConnection:1",
			"urn:schemas-upnp-org:service:WANPPPConnection:1",
			"urn:schemas-upnp-org:device:InternetGatewayDevice:1"}}

}

//发现网关设备
//1.得到本机ip地址
func (d *Discover) SearchGateway() bool {
	d.init()

	searchMessage := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
	log.Println("下一步获得本地ip地址")
	//得到本地IP地址
	d.LocalIPPort = getLocalIntenetIp() + ":"
	log.Println("获得本地ip地址:", d.LocalIPPort)
	c := make(chan string)
	go sendMulticastMsg(d.LocalIPPort, d.GroupIp, searchMessage, c)
	result := <-c
	if result == "" {
		//超时了
		log.Println()
		return false
	}
	d.GetwayInfo = getGetwayInfo(result)

	d.GetwayInfo.ServiceType = "urn:schemas-upnp-org:service:WANIPConnection:1"
	if result != "" {
		log.Println("成功发现网关设备")
		//得到提供端口映射服务的设备类型

		return true
	}
	return false
}

//查看设备描述
//1.得到端口映射请求url
func (d *Discover) DeviceDesc() string {

	requestInfo := RequestInfo{Method: "GET", Host: d.GetwayInfo.Host,
		Url: d.GetwayInfo.DeviceDescUrl, Proto: "HTTP/1.1"}

	headerMap := map[string]string{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"User-Agent": "preston",
		"Host":       d.GetwayInfo.Host,
		"Connection": "keep-alive"}
	msg := Msg{requestInfo: requestInfo, headerMap: headerMap}

	result := dialTCPSendMsg(msg)
	if result != "" {
		log.Println("成功获得网关设备描述")
	}
	d.GetwayInfo.ControlURL = findServiceUrl(result, d.GetwayInfo.ServiceType)
	log.Println("获得网关url  ", d.GetwayInfo.ControlURL)
	return result
}

//查看设备状态
func (d *Discover) SeeDeviceStatusInfo() string {
	log.Println("开始执行查看设备状态")
	requestInfo := RequestInfo{Method: "POST", Host: d.GetwayInfo.Host,
		Url: "http://" + d.GetwayInfo.Host + "/ipc", Proto: "HTTP/1.1"}
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetStatusInfo`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	headerMap := map[string]string{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"SOAPAction":     `"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"`,
		"Content-Type":   "text/xml",
		"Connection":     "Close",
		"Content-Length": ""}
	msg := Msg{requestInfo, headerMap, body}

	result := HttpURLConnect(msg)

	return result
}

//得到已经存在的端口映射列表
func (d *Discover) GetPortMapping() []string {

	return nil
}

//得到网关设备普通的端口列表
func (d *Discover) GetGenericPortMapping() []string {
	return nil
}

//得到网关特殊端口占用列表
func (d *Discover) GetSpecificPortMapping() []string {
	return nil
}

//得到网关外网IP地址
//1.得到外网ip地址
func (d *Discover) GetExternalIPAddress() string {
	requestInfo := RequestInfo{Method: "POST", Host: d.GetwayInfo.Host,
		Url: "http://" + d.GetwayInfo.Host + "/ipc", Proto: "HTTP/1.1"}

	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetExternalIPAddress`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childOne.AddChild(childTwo)
	body.AddChild(childOne)

	headerMap := map[string]string{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"SOAPAction":     `"urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress"`,
		"Content-Type":   "text/xml",
		"Connection":     "Close",
		"Content-Length": ""}
	msg := Msg{requestInfo, headerMap, body}

	result := HttpURLConnect(msg)
	log.Println(result)
	d.MappingInfo.GetewayOutsideIP = "100.3.23.68"

	return result
}

//添加一个端口映射
//1.本机映射端口
//2.外网映射端口
func (d *Discover) AddPortMapping(localPort, remotePort int, protocol string) bool {
	requestInfo := RequestInfo{Method: "POST", Host: d.GetwayInfo.Host,
		Url: "http://" + d.GetwayInfo.Host + "/ipc", Proto: "HTTP/1.1"}

	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:AddPortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}

	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewInternalPort", Content: strconv.Itoa(localPort)}
	childList3 := Node{Name: "NewProtocol", Content: protocol}
	childList4 := Node{Name: "NewEnabled", Content: "1"}
	childList5 := Node{Name: "NewInternalClient", Content: strings.Split(d.LocalIPPort, ":")[0]}
	childList6 := Node{Name: "NewLeaseDuration", Content: "0"}
	childList7 := Node{Name: "NewPortMappingDescription", Content: "test"}
	childList8 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childTwo.AddChild(childList4)
	childTwo.AddChild(childList5)
	childTwo.AddChild(childList6)
	childTwo.AddChild(childList7)
	childTwo.AddChild(childList8)

	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	headerMap := map[string]string{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"SOAPAction":     `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`,
		"Content-Type":   "text/xml",
		"Connection":     "Close",
		"Content-Length": ""}
	msg := Msg{requestInfo, headerMap, body}

	result := HttpURLConnect(msg)
	if result == "" {
		log.Println("添加端口请求返回为空")
		return false
	}
	d.MappingInfo.OutsideMappingPort[protocol] = remotePort
	d.MappingInfo.InsideMappingPort[protocol] = localPort
	return true
}

//删除一个端口映射
func (d *Discover) DeletePortMapping(remotePort int, protocol string) bool {
	requestInfo := RequestInfo{Method: "POST", Host: d.GetwayInfo.Host,
		Url: "http://" + d.GetwayInfo.Host + "/ipc", Proto: "HTTP/1.1"}
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:DeletePortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}

	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewProtocol", Content: protocol}
	childList3 := Node{Name: "NewRemoteHost"}

	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)

	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	headerMap := map[string]string{"Accept": "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2",
		"SOAPAction":     `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`,
		"Content-Type":   "text/xml",
		"Connection":     "Close",
		"Content-Length": ""}
	msg := Msg{requestInfo, headerMap, body}

	result := HttpURLConnect(msg)
	if result != "" {
		return false
	}
	return true
}

var discover *Discover

//添加一个端口映射
//返回是否映射成功

func NewPortMapping() *MappingInfo {
	discover = &Discover{}
	if success := discover.SearchGateway(); success == false {
		//没有发现upnp设备
		return nil
	}

	discover.DeviceDesc()
	// discover.SeeDeviceStatusInfo()
	//添加一个TCP和UDP端口映射并记录下来
	for _, protocol := range discover.DiscoverInfo.Protocols {
	loop:
		//从队列中取出一个
		ele := discover.DiscoverInfo.MappingPorts.Front()
		//从队列中删除
		discover.DiscoverInfo.MappingPorts.Remove(ele)
		port := ele.Value.(int)
		result := discover.AddPortMapping(port, port, protocol)
		if result {
			switch strings.ToUpper(protocol) {
			case "TCP":
				discover.TCPLocalMappingPort = strconv.Itoa(port)
				discover.TCPExternalMappingPort = strconv.Itoa(port)
			case "UDP":
				discover.UDPLocalMappingPort = strconv.Itoa(port)
				discover.UDPExternalMappingPort = strconv.Itoa(port)
			}
			log.Println("成功添加一个端口映射")
		} else {
			log.Println(port, "端口映射失败")
			goto loop
		}
	}
	//得到外网地址
	discover.GetExternalIPAddress()
	return &discover.MappingInfo
}

//删除路由器上的端口映射
func DeletePortMapping() {
	if discover.MappingInfo.GetewayOutsideIP != "" {
		for key, value := range discover.MappingInfo.OutsideMappingPort {
			discover.DeletePortMapping(value, key)
		}
	}
}
