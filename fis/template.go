package fis

import (
	"html/template"
)

//css placeholder
const CONTENT_FIS_CSS_LINK_TAG = "<!--FIS_CSS_LINK-->"

//javascript placeholder
const CONTENT_FIS_JAVASCRIPT_SCRIPT_TAG = "<!--FIS_JAVASCRIPT_SCRIPT-->"

const EMPTY = ""

func Hello(s string) string {
	return s
}

var Funcs = template.FuncMap{
	"js":        _js,
	"css":       _css,
	"framework": _framework,
	"require":   _require,
	"uri":       _uri,
}

func _js() template.HTML {
	return CONTENT_FIS_JAVASCRIPT_SCRIPT_TAG
}

func _css() template.HTML {
	return CONTENT_FIS_CSS_LINK_TAG
}

func _framework(args ...interface{}) string {
	id := args[0].(string)
	ResourceApi.Framework = ResourceApi.Uri(id)
	return EMPTY
}

func _require(args ...interface{}) string {
	id := args[0].(string)
	async := false
	if len(args) == 2 {
		boolStr := args[1].(string)
		if boolStr == "true" {
			async = true
		}
	}

	ResourceApi.Load(id, async)

	return EMPTY
}

func _uri(args ...interface{}) string {
	id := args[0].(string)
	return ResourceApi.Uri(id)
}
