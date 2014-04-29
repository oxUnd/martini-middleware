//
//
// res := fis.NewResource(map[string]string{"root": "/home/user/debug/a/config"})
// uri := res.Load("common:static/mod.js")
//  	=> uri = /static/mod.js or http://xxx.com/static/mod.js
//

package fis

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

var Settings = make(map[string]string)

//
// settings["root"] => config_dir
//
func NewResource(settings map[string]string) *Resource {
	Settings = settings
	ret := &Resource{
		"",
		make(map[string]interface{}),
		make(map[string][]string),
		make(map[string]interface{}),
		make(map[string]string),
	}
	return ret
}

//Resource
type Resource struct {
	Framework string
	maps      map[string]interface{}
	staticSet map[string][]string
	asyncSet  map[string]interface{}
	_loaded   map[string]string
}

//reset all data
func (r Resource) Reset() {
	r.maps = nil
	r._loaded = nil
	r.staticSet = nil
	r.Framework = ""
}

func (r Resource) Register(ns string) bool {
	_, ok := r.maps[ns]

	if !ok {
		var path = Settings["root"] + "/" + ns + "-map.json"
		if ns == "__global__" {
			path = Settings["root"] + "/map.json"
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Can't found: " + path)
			return false
		}
		buffer := bytes.NewBuffer(content)
		decoder := json.NewDecoder(buffer) //JSON decoder
		var result map[string]interface{}
		err = decoder.Decode(&result)
		if err != nil {
			log.Println("Can't a JSON file: " + path)
			return false
		}
		r.maps[ns] = result
	}

	return true
}

func (r Resource) getStaticSet(typ string) []string {
	return r.staticSet[typ]
}

func (r Resource) getNamespace(id string) string {
	p := strings.Index(id, ":")
	ret := "__global__"
	if p != -1 {
		ret = id[0:p]
	}
	return ret
}

func (r Resource) getRes(id string) (interface{}, bool) {
	ns := r.getNamespace(id)
	r.Register(ns)
	ret, ok := r.maps[ns]
	if !ok {
		log.Println("Can't load the map of resource: " + id)
		return nil, false
	}
	res := ret.(map[string]interface{})["res"]
	return res.(map[string]interface{})[id], true
}

//get the online url of a js or css file
func (r Resource) Uri(id string) string {
	res, ok := r.getRes(id)
	ret := ""
	if ok {
		//get url return it!
		resT := res.(map[string]interface{})
		if pkg, have := resT["pkg"]; have {
			pkgMap := r.maps[r.getNamespace(id)].(map[string]interface{})["pkg"]
			resMap := pkgMap.(map[string]interface{})
			_resMap := resMap[pkg.(string)].(map[string]interface{})
			ret = _resMap["uri"].(string)
		} else {
			ret = res.(map[string]interface{})["uri"].(string)
		}
	}
	return ret
}

//load static
func (r Resource) Load(id string, async bool) string {
	ret, ok := r._loaded[id]
	if ok {
		return ret
	} else {
		ns := r.getNamespace(id)
		if resMap, ok := r.maps[ns]; ok || r.Register(ns) {
			resMap, _ = r.maps[ns]
			//map.json must have key `res` and `map`, so can't check here.
			ress := resMap.(map[string]interface{})["res"].(map[string]interface{})
			pkgs := resMap.(map[string]interface{})["pkg"].(map[string]interface{})

			if res, ok := ress[id]; ok {
				res := res.(map[string]interface{})
				uri := ""
				if pId, ok := res["pkg"]; ok {
					pkg := pkgs[pId.(string)].(map[string]interface{})
					uri = pkg["uri"].(string)
					hasIds := pkg["has"].([]interface{})
					for _, hasId := range hasIds {
						r._loaded[hasId.(string)] = uri
					}
					//@TODO
					res = pkg
				} else {
					uri = res["uri"].(string)
					r._loaded[id] = uri
				}

				if depIds, ok := res["deps"]; ok {
					for _, depId := range depIds.([]interface{}) {
						r.Load(depId.(string), async)
					}
				}

				ret = uri
				typ := res["type"].(string)
				r.staticSet[typ] = append(r.staticSet[typ], uri)
			}
		}
	}
	return ret
}

func (r Resource) Render(html string) *bytes.Buffer {
	str := ""

	if css, ok := r.staticSet["css"]; ok {
		str += "<link rel=\"stylesheet\" href=\"" + strings.Join(css, "\" /><link rel=\"stylesheet\" href=\"") + "\" />"
	}

	if len(r.Framework) > 0 {
		str += `<script src="` + r.Framework + `"></script>`
	}

	if js, ok := r.staticSet["js"]; ok {
		for _, url := range js {
			if url == r.Framework {
				continue
			}
			str += `<script src="` + url + `"></script>`
		}
	}

	html = strings.Replace(html, "</head>", str+"\n</head>", 1)
	r.Reset()
	return bytes.NewBufferString(html)
}
