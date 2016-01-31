package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var filename = flag.String("f", "", "文件名")
var command = flag.String("c", "", "命令")

func cmd() {
	out, err := exec.Command(*command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

func main() {
	flag.Parse()
	var newInfo os.FileInfo
	var oldInfo os.FileInfo
	var err error
	cmd()
	for {
		time.Sleep(1e9)
		newInfo, err = os.Stat(*filename)
		if err != nil {
			continue
		}
		if oldInfo != nil && !newInfo.ModTime().Equal(oldInfo.ModTime()) {
			cmd()
		}
		oldInfo = newInfo
	}
}
