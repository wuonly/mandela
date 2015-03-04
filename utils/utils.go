package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

/*
	获得域名的hash值
*/
func GetHashForDomain(domain string) string {
	hash := sha256.New()
	hash.Write([]byte(domain))
	md := hash.Sum(nil)
	return FormatIdUtil(new(big.Int).SetBytes(md), 16)
}

/*
	格式化id为十进制或十六进制字符串
*/
func FormatIdUtil(idInt *big.Int, base int) string {
	if idInt.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	return hex.EncodeToString(idInt.Bytes())
}
