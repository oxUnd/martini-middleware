package fis

import (
	"testing"
//	"fmt"
	"path/filepath"
//	"runtime"
)


func Test_Register(t *testing.T) {
	__dir, _ := filepath.Abs("./")
	Settings["root"] =  __dir + "/res/config"
	res := &Resource{
		make(map[string]interface {}),
		make(map[string][]string),
		make(map[string]interface {}),
		make(map[string]string),
	}

	ret := res.Register("test")
	if (!ret) {
		t.Fail();
	} else {
		expect(t, res.maps["test"].(map[string]interface{})["test"].(string), "test")
	}
}

func Test_getRes(t *testing.T) {
	__dir, _ := filepath.Abs("./")
	Settings["root"] = __dir + "/res/config"
	res := &Resource{
		make(map[string]interface {}),
		make(map[string][]string),
		make(map[string]interface {}),
		make(map[string]string),
	}

	ret, ok := res.getRes("common:static/mod.js")
	expect(t, ok, true)

	r := ret.(map[string]interface {})

	expect(t,  r["uri"].(string), "/static/mod.js")
	expect(t, r["type"].(string), "js")
}


func Test_Load(t *testing.T) {
	__dir, _ := filepath.Abs("./")
	Settings["root"] = __dir + "/res/config"
	res := &Resource{
		make(map[string]interface {}),
		make(map[string][]string),
		make(map[string]interface {}),
		make(map[string]string),
	}

	uri := res.Load("common:static/mod.js", false);
	expect(t, uri, "/static/pkg.js")
	expect(t, res.staticSet["js"][0], "/static/pkg.js")
}
