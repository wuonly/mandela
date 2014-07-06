package nodeStore

import (
	"math/big"
)

//从大到小排序
type IdDESC []*big.Int

func (this IdDESC) Len() int {
	return len(this)
}

func (this IdDESC) Less(i, j int) bool {
	qu := new(big.Int).Sub(this[i], this[j])

	quInt := qu.Cmp(big.NewInt(0))
	//从大到小排序
	return quInt > 0

	// return this[i].NodeId < this[j].Val // 按值排序
	//return ms[i].Key < ms[j].Key // 按键排序
}

func (this IdDESC) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

//从小到大排序
type IdASC []*big.Int

func (this IdASC) Len() int {
	return len(this)
}

func (this IdASC) Less(i, j int) bool {
	qu := new(big.Int).Sub(this[i], this[j])

	quInt := qu.Cmp(big.NewInt(0))
	//从大到小排序
	return quInt < 0

	// return this[i].NodeId < this[j].Val // 按值排序
	//return ms[i].Key < ms[j].Key // 按键排序
}

func (this IdASC) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
