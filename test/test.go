package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	x := big.NewInt(5897878)
	fmt.Println(hex.EncodeToString(x.Bytes()))

	b, ok := new(big.Int).SetString("013755a17ac00e6d483b0be97a91dcfefdfbb6cf1177cc731bf39301730a6e35", 16)
	if !ok {
		fmt.Println("解析错误", ok)
	}
	fmt.Println(b)
}
