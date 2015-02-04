package mandela

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/prestonTao/mandela/nodeStore"
	"io"
	"io/ioutil"
	"net"
	"os"
)

const (
	Path_Id  = "conf/id.json"
	Str_zaro = "0000000000000000000000000000000000000000000000000000000000000000"
)

//节点是否是新节点，
//新节点需要连接超级节点，然后超级节点给她生成id
var Init_HaveId = true

var Init_IdInfo nodeStore.IdInfo

func init() {
	loadIdInfo()
}

/*
	加载本地的idInfo
*/
func loadIdInfo() {
	data, err := ioutil.ReadFile(Path_Id)
	//本地没有idinfo文件
	if err != nil {
		Init_HaveId = true
		return
	}
	err = json.Unmarshal(data, Init_IdInfo)
	if err != nil {
		Init_HaveId = true
		return
	}
	Init_HaveId = false
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
func GetId(addr string) (idInfo *nodeStore.IdInfo, err error) {
	idInfo = &nodeStore.IdInfo{
		Id:       Str_zaro,
		UserName: "nimei",
		Email:    "qqqqq@qq.com",
		Local:    "djfkafjkls",
	}

	conn, err := net.Dial("tcp", addr)
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
	idInfo = new(nodeStore.IdInfo)
	json.Unmarshal(nameByte[:n], idInfo)
	return
}
