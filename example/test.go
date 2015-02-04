package main

import (
	"fmt"
)

func main() {
	example()
}

func example() {
	maps := make(map[string]string, 1)
	maps["tao"] = "tao"
	maps["taopopoo"] = "taopopoo"

	oldMaps := maps
	oldMaps["nimei"] = "nimei"
	fmt.Println(oldMaps)
}
