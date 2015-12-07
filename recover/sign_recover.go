package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var controlChan chan os.Signal
var functionMap map[string][]func()
var signalSlice []os.Signal
var fp *os.File

func init() {
	controlChan = make(chan os.Signal, 1)
	functionMap = make(map[string][]func())
	signalSlice = make([]os.Signal, 0)
}

func register(recoverfunc func(), signSlice ...os.Signal) {
	var signName string
	var ok bool
	for _, sign := range signSlice {
		signName = sign.String()
		if _, ok = functionMap[signName]; !ok {
			functionMap[signName] = make([]func(), 0)
			signalSlice = append(signalSlice, sign)
		}
		functionMap[signName] = append(functionMap[signName], recoverfunc)
	}
}

func startRecover() {
	sign := <-controlChan
	fmt.Printf("Receive Signal %s\n", sign.String())
	recoverfuncSlice := functionMap[sign.String()]
	for _, recoverfunc := range recoverfuncSlice {
		recoverfunc()
	}
}

func run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover reason: %s", r)
			fmt.Println("program stack: %s ", debug.Stack())
		}
	}()

	recordPanic("panic")

	signal.Notify(controlChan, signalSlice...)
	fmt.Println("Start!")

	for {
		startRecover()
	}
}

func recordPanic(filename string) (err error) {
	fp, err = os.Create(filename)
	if err != nil {
		return
	}
	syscall.Dup2(int(fp.Fd()), 1)
	syscall.Dup2(int(fp.Fd()), 2)
	return
}

func recoverfunction() {
	fmt.Println("Don't Stop me!")
}

func main() {
	register(recoverfunction, syscall.SIGINT)
	run()
}
