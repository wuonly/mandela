package main

import (
	"fmt"
	"github.com/prestonTao/mandela"
	"os"
)

var zaro = "0000000000000000000000000000000000000000000000000000000000000000"

func main() {
	fmt.Println(zaro, "\n", len(zaro))

	idInfo, err := mandela.NewIdInfo("prestonTao", "taopopoo@126.com", "mandela", zaro)
	if err != nil {
		fmt.Println(err)
		panic("create id error")
	}
	fmt.Println(idInfo)
	private := idInfo.Build()
	file, _ := os.Create("private")
	file.Write(private)
	file.Close()

}