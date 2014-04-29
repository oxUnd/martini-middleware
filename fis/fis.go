package fis

import (
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
)

//common function

var ResourceApi *Resource

//inject martini, it dep for martini-contrib/render
func Renderer(options ...Options) martini.Handler {
	opt := prepareOptions(options)
	s := map[string] string {
		"root": opt.Directory + "/config",
	}

	ResourceApi = NewResource((map[string]string)(s))

	opt.Funcs = append(opt.Funcs, Funcs)
	cs := prepareCharset(opt.Charset)
	t := compile(opt)

	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		var tc *template.Template
		if martini.Env == martini.Dev {
			// recompile for easy development
			tc = compile(opt)
		} else {
			// use a clone of the initial template
			tc, _ = t.Clone()
		}
		c.MapTo(&renderer{res, req, tc, opt, cs}, (*Render)(nil))
	}
}


