package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	processDir = "/proc/"
)

// Stat 记录程序的io信息 对应linux下/proc/pid/io
type Stat struct {
	pid                 int
	cmd                 string
	rchar               int64
	wchar               int64
	syscr               int64
	syscw               int64
	readBytes           int64
	writeBytes          int64
	cancelledWriteBytes int64
}

var (
	pid2stat map[int]Stat
)

func init() {
	pid2stat = make(map[int]Stat)
}

func getDataFromFile(name string) ([]byte, error) {
	var (
		fp   *os.File
		err  error
		data []byte
	)
	if fp, err = os.Open(name); err != nil {
		fmt.Printf("open '%s' fail %s\n", name, err)
		return nil, err
	}
	if data, err = ioutil.ReadAll(fp); err != nil {
		fmt.Printf("read %s fail %s\n", name, err)
		return nil, err
	}
	return data, nil
}

func getDiff(stat *Stat, s string) string {
	array := strings.Split(s, ":")
	array[0] = strings.TrimSpace(array[0])
	number, _ := strconv.ParseInt(strings.TrimSpace(array[1]), 10, 64)
	var diff int64
	switch array[0] {
	case "rchar":
		diff = number - stat.rchar
		stat.rchar = number
	case "wchar":
		diff = number - stat.wchar
		stat.wchar = number
	case "syscr":
		diff = number - stat.syscr
		stat.syscr = number
	case "syscw":
		diff = number - stat.syscw
		stat.syscw = number
	case "read_bytes":
		diff = number - stat.readBytes
		stat.readBytes = number
	case "write_bytes":
		diff = number - stat.writeBytes
		stat.writeBytes = number
	case "cancelled_write_bytes":
		diff = number - stat.cancelledWriteBytes
		stat.cancelledWriteBytes = number
	}
	return strconv.FormatInt(diff, 10)
}

func getStat(name string) {
	var (
		err      error
		pid      int
		fileName string
		cmdline  string
		data     []byte
		stat     Stat
		ok       bool
		array    []string
	)

	if pid, err = strconv.Atoi(name); err != nil {
		return
	}

	fileName = processDir + name + "/cmdline"
	if data, err = getDataFromFile(fileName); err != nil {
		return
	}

	cmdline = string(data)

	if stat, ok = pid2stat[pid]; !ok || stat.cmd != cmdline {
		ok = false
		stat = Stat{
			pid: pid,
			cmd: cmdline,
		}
	}

	fileName = processDir + name + "/io"
	if data, err = getDataFromFile(fileName); err != nil {
		return
	}

	array = strings.Split(string(data), "\n")
	line := strconv.Itoa(stat.pid) + "\t"
	for _, s := range array {
		if len(s) == 0 {
			continue
		}
		line += getDiff(&stat, s) + "\t"
	}
	line += stat.cmd + "\t"

	if ok {
		fmt.Println(line)
	}

	pid2stat[pid] = stat
}

var t = flag.Int("t", 1, "间隔时间")
var p = flag.String("p", "", "程序的PID")

func main() {
	flag.Parse()
	if *p == "" {
		flag.PrintDefaults()
		return
	}
	fmt.Println("pid\trchar\twchar\tsyscr\tsyscw\trbytes\twbytes\tcwbytes\tcmd\t")
	tick := time.Tick(time.Second * time.Duration(*t))
	for range tick {
		getStat(*p)
	}
}
