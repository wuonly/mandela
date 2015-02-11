package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	zaroStr := "0000000000000000000000000000000000000000000000000000000000000000"
	maxNumberStr := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	idStr := "9d1406a76433f43ab752f468e0ca58baf73c9adddfc505a617c5756ea310e44f"
	idInt, _ := new(big.Int).SetString(idStr, 16)
	fmt.Println(len(idStr), len(zaroStr), len(maxNumberStr))

	maxNumberInt, ok := new(big.Int).SetString(maxNumberStr, 16)
	if !ok {
		fmt.Println("失败")
		return
	}

	fmt.Println(maxNumberInt.String(), len(maxNumberInt.String()))
	number_2 := big.NewInt(2)
	halfNumberInt := new(big.Int).Quo(maxNumberInt, number_2)

	fmt.Println(hex.EncodeToString(halfNumberInt.Bytes()))

	halfAndHalfNumberInt := new(big.Int).Quo(halfNumberInt, number_2)
	fmt.Println(hex.EncodeToString(halfAndHalfNumberInt.Bytes()))

	fmt.Println(new(big.Int).Sub(maxNumberInt, idInt))

	// Print(maxNumberInt)

}

func Print(findInt *big.Int) {
	fmt.Println("==================================\r\n")
	bi := ""

	// findInt := new(big.Int).SetBytes([]byte(nodeId))
	lenght := findInt.BitLen()
	for i := 0; i < lenght; i++ {
		tempInt := findInt
		findInt = new(big.Int).Div(tempInt, big.NewInt(2))
		mod := new(big.Int).Mod(tempInt, big.NewInt(2))
		bi = mod.String() + bi
	}
	fmt.Println(bi, "\r\n")
	fmt.Println("==================================\r\n")
}
