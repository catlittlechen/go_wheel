package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"syscall"
)

var fp *os.File

func run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover reason: %s", r)
			fmt.Println("program stack: %s ", debug.Stack())
		}
	}()

	recordPanic("nohup.out")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("usage: nohup utility [arguments]")
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	err := cmd.Start()
	if err != nil {
		fmt.Printf("nohup: %s\n", err)
	}

	fmt.Printf("nohup: Command Success with PID %d\n", cmd.Process.Pid)
	return
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

func main() {
	run()
}
