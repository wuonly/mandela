package app

import (
	"fmt"
	"testing"
)

//通过一个域名和用户名得到节点的id
func TestGetHashKey(t *testing.T) {
	str := GetHashKey("百度", "nihao")
	fmt.Println("百度字符串：", str)
}
