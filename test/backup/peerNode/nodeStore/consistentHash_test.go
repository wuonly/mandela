package nodeStore

import (
	"fmt"
	"math/big"
	"testing"
	// "time"
)

func TestCHash(t *testing.T) {
	// distanceTest()
}

func distanceTest() {
	root, _ := new(big.Int).SetString("8", 10)
	num9, _ := new(big.Int).SetString("9", 10)
	num10, _ := new(big.Int).SetString("10", 10)
	num12, _ := new(big.Int).SetString("12", 10)
	num0, _ := new(big.Int).SetString("0", 10)
	num16, _ := new(big.Int).SetString("16", 10)
	fmt.Println(new(big.Int).Xor(root, num9))
	fmt.Println(new(big.Int).Xor(root, num10))
	fmt.Println(new(big.Int).Xor(root, num12))
	fmt.Println(new(big.Int).Xor(root, num0))
	fmt.Println(new(big.Int).Xor(root, num16))

}
