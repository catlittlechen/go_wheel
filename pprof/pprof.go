package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"
)

var proFile *os.File
var mutex *sync.Mutex

func init() {
	mutex = new(sync.Mutex)
}

func getURLValues(r *http.Request) (url.Values, error) {
	urlMap, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	keyword, err := url.ParseQuery(urlMap.RawQuery)
	if err != nil {
		return nil, err
	}

	return keyword, nil
}

func work(action string) {
	switch action {
	case "start_cpuprofile":
		pprof.StartCPUProfile(proFile)
	case "stop_cpuprofile":
		pprof.StopCPUProfile()
	case "write_heap":
		pprof.WriteHeapProfile(proFile)
	case "lookup_goroutine":
		profile := pprof.Lookup("goroutine")
		profile.WriteTo(proFile, 2)
	case "lookup_heap":
		profile := pprof.Lookup("heap")
		profile.WriteTo(proFile, 2)
	case "lookup_threadcreate":
		profile := pprof.Lookup("threadcreate")
		profile.WriteTo(proFile, 2)
	case "lookup_block":
		profile := pprof.Lookup("block")
		profile.WriteTo(proFile, 2)
	case "start_trace":
		trace.Start(proFile)
	case "stop_trace":
		trace.Stop()
	}
	return
}

func pprofInfo(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("request actioln %+v\n", r)

	keyword, err := getURLValues(r)
	if err != nil {
		fmt.Printf("url Parse Error[%s]\n", err)
		return
	}

	action := keyword.Get("action")
	if strings.HasPrefix(action, "stop") {
		if proFile != nil {
			fmt.Println("Profile work is not Started")
			return
		}
	} else {
		if proFile != nil {
			fmt.Println("Profile is recording")
			return
		}
		filename := keyword.Get("filename")
		if proFile, err = os.Create(filename); err != nil {
			proFile = nil
			fmt.Printf("create file fail. Error[%s]\n", err)
			return
		}
	}

	if !strings.HasPrefix(action, "start") {
		defer func() {
			proFile.Close()
			proFile = nil
		}()
	}

	work(action)
	return
}

func run() {
	http.HandleFunc("/pprof", pprofInfo)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Printf("http ListenAndServe fail Error[%s]\n", err)
	}
	return
}

func doSomething() {
	for {

	}
}

func main() {
	go run()
	doSomething()
}
