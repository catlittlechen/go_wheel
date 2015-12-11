package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func run() {
	var line string
	var tip string
	count := 50
	for i := 0; i <= count; i++ {
		time.Sleep(100 * time.Millisecond)
		line = strings.Repeat("=", i) + ">" + strings.Repeat(" ", count-i)
		tip = ""
		if i != count {
			switch i % 4 {
			case 0:
				tip = "-"
			case 1:
				tip = "\\"
			case 2:
				tip = "|"
			case 3:
				tip = "/"
			}
		}
		fmt.Printf("\r%s%d%%[%s]", tip, i*100/count, line)
		os.Stdout.Sync()
	}
}

func main() {
	run()
}
