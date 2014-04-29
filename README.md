Martini Middleware
==================


### fis

#### get 

```bash
$ go get github.com/xiangshouding/martini-middleware/fis
```

#### use

```go
package main

import (
	"github.com/go-martini/martini"
	"github.com/xiangshouding/martini-middleware/fis"
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("src/app/public"))
	m.Use(fis.Renderer(fis.Options{
		Directory: "src/app/template",
		Extensions: []string {".tpl"},
	}))

	m.Get("/", func (r fis.Render) {
		r.HTML(200, "index", "")
	})
	m.Run()
}
```

