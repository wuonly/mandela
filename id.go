package mandela

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

//Id信息
type IdInfo struct {
	Id           string `json:"id"`           //id
	CreateTime   string `json:"createtime"`   //创建时间
	UserName     string `json:"username"`     //用户名
	Email        string `json:"email"`        //email
	Local        string `json:"local"`        //地址
	SuperNodeId  string `json:"supernodeid"`  //创建者节点id
	SuperNodeKey string `json:"supernodekey"` //创建者公钥
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

//userName      用户名，最大长度100
//email         email，最大长度100
//local         地址，最大长度100
//superNodeId   超级节点id，最大长度
//superNodeKey  超级节点密钥
//rerutn idInfo
//return err
func NewIdInfo(userName, email, local, superNodeId, superNodeKey string) (idInfo *IdInfo, err error) {
	if len(userName) > 100 {
		err = errors.New("userName 长度不能超过100个字符")
		return
	}
	if len(email) > 100 {
		err = errors.New("email 长度不能超过100个字符")
		return
	}
	if len(local) > 100 {
		err = errors.New("local 长度不能超过100个字符")
		return
	}
	if len(superNodeId) != 64 {
		err = errors.New("superNodeId 参数长度不正确")
		return
	}
	if len(superNodeKey) > 100 {
		err = errors.New("superNodeKey 长度不能超过100个字符")
		return
	}

	hash := sha256.New()
	hash.Write([]byte(userName + "#" + email + "#" + local + "#" + superNodeId + "#" + superNodeKey))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	idInfo = &IdInfo{
		Id:           mdStr,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05.999999999"),
		UserName:     userName,
		Email:        email,
		Local:        local,
		SuperNodeId:  superNodeId,
		SuperNodeKey: superNodeKey,
	}
	return
}
