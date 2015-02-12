package main

import (
	"fmt"
	"os"
	"unsafe"
)

func main() {
	example3()
}

func signalHandle() {
	for {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
		sig := <-ch
		Utils.LogInfo("Signal received: %v", sig)
		switch sig {
		default:
			Utils.LogInfo("get sig=%v\n", sig)
		case syscall.SIGHUP:
			Utils.LogInfo("get sighup\n") //Utils.LogInfo是我自己封装的输出信息函数
		case syscall.SIGINT:
			os.Exit(1)
		case syscall.SIGUSR1:
			Utils.LogInfo("usr1\n")
		case syscall.SIGUSR2:
			Utils.LogInfo("usr2\n")

		}
	}
}

func example2() {
	b := make([]byte, 100)
	f := os.Stdin
	w := os.Stdout
	defer f.Close()
	defer w.Close()
	for {
		w.WriteString("input:")
		c, _ := f.Read(b)
		bb := b[:c-2]
		str := *(*string)(unsafe.Pointer(&bb))
		fmt.Println(str)
		if str == "exit" {
			break
		}
	}
}

func example() {
	maps := make(map[string]string, 1)
	maps["tao"] = "tao"
	maps["taopopoo"] = "taopopoo"

	oldMaps := maps
	oldMaps["nimei"] = "nimei"
	fmt.Println(oldMaps)
}
