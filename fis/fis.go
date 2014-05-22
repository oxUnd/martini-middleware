// Package fis is a middleware for Martini that provides easy JSON serialization and HTML template rendering and support FIS
//
// package main
//
// import (
//   "github.com/go-martini/martini"
//   "github.com/xiangshouding/martini-middleware/fis"
// )
//
// func main() {
//   m := martini.Classic()
//   m.Use(martini.Static("public"))
//   m.Use(fis.Renderer(fis.Options{
//     Directory:  "template",
//     Extensions: []string{".tpl"},
//   }))
//   m.Get("/", func(r fis.Render) {
//     r.HTML(200, "page/index", "")
//   })
//   m.Run()
// }
package fis

import (
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
)

//common function

var ResourceApi *Resource

//inject martini, it dep for github.com/xiangshouding/martini-middleware/fis
func Renderer(options ...Options) martini.Handler {
	opt := prepareOptions(options)
	s := map[string]string{
		"root": opt.Directory + "/config",
	}

	opt.Funcs = append(opt.Funcs, Funcs)
	cs := prepareCharset(opt.Charset)
	ResourceApi = NewResource((map[string]string)(s))
	t := compile(opt)

	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		var tc *template.Template
		if martini.Env == martini.Dev {
			ResourceApi = NewResource((map[string]string)(s))
			// recompile for easy development
			tc = compile(opt)
		} else {
			// use a clone of the initial template
			tc, _ = t.Clone()
		}

		c.MapTo(&renderer{res, req, tc, opt, cs}, (*Render)(nil))
	}
}
