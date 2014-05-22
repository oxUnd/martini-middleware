package fis

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_Renderer(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)
}

func Test_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "res/template",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello", "")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	log.Println(res.Body.String())
}
