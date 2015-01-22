package mandela

import (
	"encoding/json"
	"io/ioutil"
)

const (
	Path_SuperPeerAddress = "conf/nodeEntry.json"
)

var Sys_superNodeEntry []string = []string{}

func init() {
	fileBytes, err := ioutil.ReadFile(Path_SuperPeerAddress)
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileBytes, &Sys_superNodeEntry); err != nil {
		return
	}
}
