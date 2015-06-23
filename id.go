/*
	加载本地配置文件中的idinfo
		1.读取并解析本地idinfo配置文件。
*/
package mandela

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prestonTao/mandela/nodeStore"
	"github.com/prestonTao/mandela/utils"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
)

const (
	Str_zaro          = "0000000000000000000000000000000000000000000000000000000000000000" //字符串0
	Str_maxNumber     = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //256位的最大数十六进制表示id
	Str_halfNumber    = "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //最大id的二分之一
	Str_quarterNumber = "3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //最大id的四分之一
)

var (
	Path_Id = filepath.Join(Path_configDir, "idinfo.json")

	Init_IdInfo nodeStore.IdInfo

	Number_max     *big.Int //最大id数
	Number_half    *big.Int //最大id的二分之一
	Number_quarter *big.Int //最大id的四分之一
)

func init() {
	var ok bool
	Number_max, ok = new(big.Int).SetString(Str_maxNumber, 16)
	if !ok {
		panic("id string format error")
	}
	Number_half, ok = new(big.Int).SetString(Str_halfNumber, 16)
	if !ok {
		panic("id string format error")
	}
	Number_quarter, ok = new(big.Int).SetString(Str_quarterNumber, 16)
	if !ok {
		panic("id string format error")
	}
}

/*
	加载本地的idInfo
*/
func loadIdInfo() {
	data, err := ioutil.ReadFile(Path_Id)
	//本地没有idinfo文件
	if err != nil {
		fmt.Println("读取idinfo.json文件出错")
		return
	}
	err = json.Unmarshal(data, &Init_IdInfo)
	if err != nil {
		fmt.Println("解析idinfo.json文件错误")
		return
	}
}

/*
	保存idinfo到本地文件
*/
func saveIdInfo(path string) {
	fileBytes, _ := json.Marshal(Init_IdInfo)
	file, _ := os.Create(path)
	file.Write(fileBytes)
	file.Close()
}

/*
	连接超级节点，得到一个id
	@ addr   超级节点ip地址
*/
func GetId(idInfo nodeStore.IdInfo) (newIdInfo *nodeStore.IdInfo, err error) {
	// idInfo = nodeStore.IdInfo{
	// 	Id:     Str_zaro,
	// 	Name:   "nimei",
	// 	Email:  "qqqqq@qq.com",
	// 	Domain: "djfkafjkls",
	// }

	conn, err := net.Dial("tcp", getSuperAddrOne())
	if err != nil {
		err = errors.New("连接超级节点失败")
		return
	}

	/*
		向对方发送自己的名称
	*/
	lenght := int32(len(idInfo.Build()))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, lenght)
	buf.Write(idInfo.Build())
	conn.Write(buf.Bytes())

	/*
		对方服务器创建好id后，发送给自己
	*/
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	nameLenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, nameLenght)
	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	//得到对方生成的名称
	newIdInfo = new(nodeStore.IdInfo)
	json.Unmarshal(nameByte[:n], newIdInfo)
	conn.Close()
	return
}

/*
	得到保存数据的逻辑节点
	@idStr  id十六进制字符串
*/
func GetLogicIds(idStr string) (logicIds []string, ok bool) {
	ok = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			ok = false
		}
	}()
	logicIds = make([]string, 0)
	var idInt *big.Int
	idInt, ok = new(big.Int).SetString(idStr, nodeStore.IdStrBit)
	if !ok {
		return
	}
	//先获取5个逻辑id
	//1
	oppositeId := new(big.Int).Not(idInt)
	//2
	logicIds = append(logicIds, utils.FormatIdUtil(oppositeId))
	id_2 := new(big.Int).Add(oppositeId, Number_quarter)
	if id_2.Cmp(Number_max) == 1 {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(id_2, Number_max)))
	} else {
		logicIds = append(logicIds, utils.FormatIdUtil(id_2))
	}
	//3
	id_3 := new(big.Int).Add(oppositeId, Number_half)
	if id_3.Cmp(Number_max) == 1 {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(id_3, Number_max)))
	} else {
		logicIds = append(logicIds, utils.FormatIdUtil(id_3))
	}
	//4
	if oppositeId.Cmp(Number_quarter) == -1 {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(Number_max, new(big.Int).Sub(Number_quarter, oppositeId))))
	} else {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(oppositeId, Number_quarter)))
	}
	//5
	if oppositeId.Cmp(Number_half) == -1 {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(Number_half, new(big.Int).Sub(Number_half, oppositeId))))
	} else {
		logicIds = append(logicIds, utils.FormatIdUtil(new(big.Int).Sub(oppositeId, Number_half)))
	}
	return
}
