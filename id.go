package mandela

import (
	"encoding/json"
	"github.com/prestonTao/mandela/nodeStore"
	"io/ioutil"
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
