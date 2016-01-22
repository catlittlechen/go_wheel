package main

import (
	"fmt"
)

func main() {
	str := "你好，人类"
	fmt.Println(len([]rune(str)))
	fmt.Println(len(str))
}
