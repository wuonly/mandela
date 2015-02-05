package main

import (
	// "encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println(big.NewInt(1).Bytes())

	bigint, ok := new(big.Int).SetString("0", 16)
	if !ok {
		fmt.Println("no")
	}
	fmt.Println(bigint)

}
