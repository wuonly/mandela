package mandela

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"time"
)

const (
	Path_Id  = "conf/id.json"
	Str_zaro = "0000000000000000000000000000000000000000000000000000000000000000"
)

//节点是否是新节点，
//新节点需要连接超级节点，然后超级节点给她生成id
var Init_HaveId = true

var Init_IdInfo IdInfo

func init() {
	data, err := ioutil.ReadFile(Path_Id)
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

//Id信息
type IdInfo struct {
	Id          string `json:"id"`          //id
	CreateTime  string `json:"createtime"`  //创建时间
	UserName    string `json:"username"`    //用户名
	Email       string `json:"email"`       //email
	Local       string `json:"local"`       //地址
	SuperNodeId string `json:"supernodeid"` //创建者节点id
	// SuperNodeKey string `json:"supernodekey"` //创建者公钥
}

func (this *IdInfo) Parse(code []byte) (err error) {
	err = json.Unmarshal(code, this)
	return
}

//将此节点id详细信息构建为标准code
func (this *IdInfo) Build() []byte {
	str, _ := json.Marshal(this)
	return str
}

/*
	检查idInfo是否合法
	@return   true:合法;false:不合法;
*/
func CheckIdInfo(idInfo IdInfo) bool {
	if len(idInfo.UserName) > 100 {
		err = errors.New("userName 长度不能超过100个字符")
		return false
	}
	if len(idInfo.Email) > 100 {
		err = errors.New("email 长度不能超过100个字符")
		return false
	}
	if len(idInfo.Local) > 100 {
		err = errors.New("local 长度不能超过100个字符")
		return false
	}
	if len(idInfo.SuperNodeId) != 64 {
		err = errors.New("superNodeId 参数长度不正确")
		return false
	}
	return true
}

//userName      用户名，最大长度100
//email         email，最大长度100
//local         地址，最大长度100
//superNodeId   超级节点id，最大长度
//superNodeKey  超级节点密钥
//rerutn idInfo
//return err
func NewIdInfo(userName, email, local, superNodeId string) (idInfo *IdInfo, err error) {

	// if len(superNodeKey) > 100 {
	// 	err = errors.New("superNodeKey 长度不能超过100个字符")
	// 	return
	// }

	hash := sha256.New()
	hash.Write([]byte(userName + "#" + email + "#" + local + "#" + superNodeId))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	idInfo = &IdInfo{
		Id:          mdStr,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05.999999999"),
		UserName:    userName,
		Email:       email,
		Local:       local,
		SuperNodeId: superNodeId,
		// SuperNodeKey: superNodeKey,
	}
	return
}

/*
	连接超级节点，得到一个id
	@ addr   超级节点ip地址
*/
func GetId(addr string) (idInfo *IdInfo, err error) {

	idInfo, err = NewIdInfo("", "", "", zaro)
	if err != nil {
		fmt.Println(err)
		err = errors.New("生成id错误")
		return
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		err = errors.New("连接超级节点错误")
		return
	}

	//第一次连接，向对方发送自己的名称
	lenght := int32(len(name))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, lenght)
	buf.Write([]byte(name))
	conn.Write(buf.Bytes())

	//对方服务器验证成功后发送给自己的名称
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	nameLenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, nameLenght)
	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	//得到对方名称
	remoteName = string(nameByte[:n])

}
