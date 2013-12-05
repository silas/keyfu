package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"net/url"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer("templates"))

	m.Get("/run", func(res http.ResponseWriter, req *http.Request) {
		q := req.FormValue("q")
		if v, err := Run(q); err == nil {
			http.Redirect(res, req, v, 302)
		} else {
			http.Redirect(res, req, "https://www.google.com/search?q="+url.QueryEscape(q), 302)
		}
	})

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Run()
}
