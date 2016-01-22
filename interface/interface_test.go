package util

import (
	"encoding/json"
	"testing"
)

var (
	rsps = `
	{"c":0,
	 "d":{
		"id":60007,
		"username":"catchen",
		"platform":1,
		"type":0,
		"time":"2015-03-03 09:41:47",
		"status":true,
		"nick":"catlittlechen",
	  }
	}
	`
)

type userEx struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Platform int    `json:"platform"`
	Type     int    `json:"type"`
	Time     int    `json:"time"`
	Status   bool   `json:"status"`
	Nick     string `json:"nick"`
}

type response struct {
	C int         `json:"c"`
	M string      `json:"m"`
	D interface{} `json:"d"`
}

func testInterfaceToStruct(t testing.T) {
	rsp := response{}
	err := json.Unmarshal([]byte(rsps), &rsp)
	if err != nil {
		t.Fatal(err)
	}
	u := new(userEx)
	err = InterfaceToSimpleStruct("jspn", rsp.D, u)
	if err != nil {
		t.Fatal("!!fail!!")
		t.Fatal(err)
	}
	t.Log("ok")
}
