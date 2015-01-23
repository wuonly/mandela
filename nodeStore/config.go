package nodeStore

import (
	"time"
)

//节点id长度
var NodeIdLevel uint = 256

//超级节点之间查询的间隔时间
var SpacingInterval time.Duration = time.Second * 30

//存放相邻节点个数(左半边个数或者右半边个数)
var MaxRecentCount int = 2

//id字符串格式为16进制字符串
var IdStrBit int = 16
