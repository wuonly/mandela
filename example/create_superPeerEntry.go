package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	create()
}

func create() {
	superPeerEntry := make(map[string]string)
	superPeerEntry["mandela.io:9981"] = ""
	superPeerEntry["192.168.6.30:9981"] = ""

	fileBytes, _ := json.Marshal(superPeerEntry)

	file, _ := os.Create("addrEntry.json")
	file.Write(fileBytes)
	file.Close()
	fmt.Println("done !")
}
