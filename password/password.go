package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
)

var website = flag.String("w", "", "网站")
var username = flag.String("u", "", "用户名")
var salt = flag.String("s", "", "参数")

func main() {
	flag.Parse()

	keyword := "%s&&%s&&1232123&"
	originPW := fmt.Sprintf(keyword, *website, *username, *salt)
	fmt.Println(originPW)
	h := md5.New()
	io.WriteString(h, originPW)
	fmt.Printf("%x\n", h.Sum(nil))
	return
}
