package main

import (
	"net/http"
	"net/url"
)

func LogoutHandler(ctx *Context) {
	values := ctx.Request.URL.Query()

	ctx.SetSecureCookie(&http.Cookie{Name: "uid", Value: "", MaxAge: -1})

	if url_ := values.Get("url"); url_ != "" {
		ctx.Redirect("/login?url=" + url.QueryEscape(url_))
	} else {
		ctx.Redirect("/")
	}
}
