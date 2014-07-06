package upnp

import (
	// "fmt"
	// "net"
	// "reflect"
	"encoding/xml"
	// "log"
	"strings"
)

func getGetwayInfo(resultMsg string) Geteway {
	geteway := Geteway{}

	lines := strings.Split(resultMsg, "\r\n")
	for _, line := range lines {

		//按照第一个冒号分为两个字符串
		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			geteway.ST = nameValues[1]
		case "CACHE-CONTROL":
			geteway.Cache = nameValues[1]
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			geteway.Host = urls[0]
			geteway.DeviceDescUrl = "/" + urls[1]
		case "SERVER":
			geteway.GetewayName = nameValues[1]
		default:
		}
	}
	return geteway
}

func findServiceUrl(resultStr, serviceType string) string {

	inputReader := strings.NewReader(resultStr)

	// 从文件读取，如可以如下：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`
	// var eventSubURL string //`eventSubURL`
	// var SCPDURL string     //`SCPDURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		// 处理元素结束（标签）
		case xml.EndElement:
			// log.Println("结束标记：", token.Name.Local)
		// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			//得到url后其他标记就不处理了
			content := string([]byte(token))

			//找到提供端口映射的服务
			if content == serviceType {
				ISUpnpServer = true
				continue
			}
			//urn:upnp-org:serviceId:WANIPConnection
			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":

					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
					// eventSubURL = content
				case "SCPDURL":
					// SCPDURL = content
				}
			}
		default:
			// ...
		}
	}
	return controlURL
}

// func main() {
// 	msg := "HTTP/1.1 200 OK\r\n" +
// 		"CACHE-CONTROL: max-age=100\r\n" +
// 		"DATE: Thu, 16 Jan 2014 10:51:47 GMT\r\n" +
// 		"EXT:\r\n" +
// 		"LOCATION: http://192.168.1.1:1900/igd.xml\r\n" +
// 		"SERVER: Wireless N Router WR845N, UPnP/1.0\r\n" +
// 		"ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n" +
// 		"USN: uuid:upnp-InternetGatewayDevice-192168115678900001::urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n\r\n"
// 	strs := getIPAndUrl(msg)
// 	for _, str := range strs {
// 		fmt.Println("-------------------------------------")
// 		temps := strings.SplitAfterN(str, ":", 2)
// 		for _, temp := range temps {
// 			fmt.Println(temp)
// 		}
// 	}
// }
func findExternalIPAddress(result string) {
	// inputReader := strings.NewReader(result)
	// decoder := xml.NewDecoder(inputReader)
	// for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
	// 	switch token := t.(type) {
	// 	// 处理元素开始（标签）
	// 	case xml.StartElement:
	// 		if ISUpnpServer {
	// 			name := token.Name.Local
	// 			lastLabel = name
	// 		}

	// 	// 处理元素结束（标签）
	// 	case xml.EndElement:
	// 		// log.Println("结束标记：", token.Name.Local)
	// 	// 处理字符数据（这里就是元素的文本）
	// 	case xml.CharData:
	// 		//得到url后其他标记就不处理了
	// 		content := string([]byte(token))

	// 		//找到提供端口映射的服务
	// 		if content == serviceType {
	// 			ISUpnpServer = true
	// 			continue
	// 		}
	// 		//urn:upnp-org:serviceId:WANIPConnection
	// 		if ISUpnpServer {
	// 			switch lastLabel {
	// 			case "controlURL":

	// 				controlURL = content
	// 				IScontrolURL = true
	// 			case "eventSubURL":
	// 				// eventSubURL = content
	// 			case "SCPDURL":
	// 				// SCPDURL = content
	// 			}
	// 		}
	// 	default:
	// 		// ...
	// 	}
	// }

}
