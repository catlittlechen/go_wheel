package main

import (
	"fmt"
	"sync"
)

//Object 对于空对象的话，一定是指向同一个内存空间的，所以这里必须是有参数的。
type Object struct {
	num int
}

var once sync.Once

var instance *Object

func SingleNew() *Object {
	once.Do(func() {
		instance = new(Object)
	})
	return instance
}

func NormalNew() *Object {
	obj := new(Object)
	return obj
}

func signleTest() {
	fmt.Println("In the Signle Test")
	first := SingleNew()
	second := SingleNew()
	if first == second {
		fmt.Println("there are same object!")
	} else {
		fmt.Println("there are different objects!")
	}
	return
}

func normalTest() {
	fmt.Println("In the Normal Test")
	first := NormalNew()
	second := NormalNew()
	if first == second {
		fmt.Println("there are same object!")
	} else {
		fmt.Println("there are different objects!")
	}
	return
}

func main() {
	signleTest()
	normalTest()
}
