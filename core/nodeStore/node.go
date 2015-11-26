package nodeStore

import (
	// "crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

//保存节点的id
//ip地址
//不同协议的端口
type Node struct {
	// NodeId               *big.Int  //节点id的10进制字符串
	IdInfo               IdInfo    //节点id信息，id字符串以16进制显示
	IsSuper              bool      //是不是超级节点，超级节点有外网ip地址，可以为其他节点提供代理服务
	Addr                 string    //外网ip地址
	TcpPort              int32     //TCP端口
	UdpPort              int32     //UDP端口
	LastContactTimestamp time.Time //最后检查的时间戳
	// NodeIdShould         *big.Int  //影子id
	// Status               int       //节点状态，1：在线，2：正在查询中，3：下线
	// Out                  chan *Node //需要查询是否在线的节点
	// OverTime             time.Duration `1 * 60 * 60` //超时时间，单位为秒
	// SelectTime           time.Duration `5 * 60`      //查询时间，单位为秒
	// Key                  *rsa.PrivateKey //保存的公钥和私钥信息
}

//Id信息
type IdInfo struct {
	Id          string `json:"id"`          //id
	CreateTime  string `json:"createtime"`  //创建时间
	Domain      string `json:"domain"`      //域名
	Name        string `json:"local"`       //姓名
	Email       string `json:"email"`       //email
	SuperNodeId string `json:"supernodeid"` //创建者节点id
	// SuperNodeKey string `json:"supernodekey"` //创建者公钥
}

func (this *IdInfo) GetId() string {
	return this.Id
}

func (this *IdInfo) GetBigIntId() *big.Int {
	bigInt, _ := new(big.Int).SetString(this.Id, IdStrBit)
	return bigInt
}

/*
	解析一个idInfo
*/
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
	if len(idInfo.Name) > 100 {
		fmt.Println("userName 长度不能超过100个字符")
		return false
	}
	if len(idInfo.Email) > 100 {
		fmt.Println("email 长度不能超过100个字符")
		return false
	}
	if len(idInfo.Domain) > 100 {
		fmt.Println("local 长度不能超过100个字符")
		return false
	}
	if len(idInfo.SuperNodeId) != 64 {
		fmt.Println("superNodeId 参数长度不正确")
		return false
	}
	return true
}

/*
	得到id
*/
func ParseId(idInfoStr string) (id string) {
	idInfo := IdInfo{}
	idInfo.Parse([]byte(idInfoStr))
	return idInfo.Id
}

//userName      用户名，最大长度100
//email         email，最大长度100
//local         地址，最大长度100
//superNodeId   超级节点id，最大长度
//superNodeKey  超级节点密钥
//rerutn idInfo
//return err
func NewIdInfo(name, email, domain, superNodeId string) (idInfo IdInfo) {
	createTime := time.Now().Format("2006-01-02 15:04:05.999999999")
	hash := sha256.New()
	hash.Write([]byte(name + "#" + email + "#" + domain + "#" + superNodeId + "#" + createTime))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	idInfo = IdInfo{
		Id:          mdStr,
		CreateTime:  createTime,
		Name:        name,
		Email:       email,
		Domain:      domain,
		SuperNodeId: superNodeId,
		// SuperNodeKey: superNodeKey,
	}
	return
}
