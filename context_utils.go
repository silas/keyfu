package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func LoginHelper(handler func(ctx *Context)) func(*Context) {
	return func(ctx *Context) {
		if len(ctx.Uid) == 24 {
			handler(ctx)
		} else {
			ctx.Redirect(fmt.Sprintf("/logout?url=%s", url.QueryEscape(ctx.Request.URL.String())))
		}
	}
}

func CsrfHelper(handler func(ctx *Context)) func(*Context) {
	return func(ctx *Context) {
		formCsrf := ctx.Request.FormValue("csrf")
		cookieCsrf := ctx.GetCookie("csrf")
		if cookieCsrf != nil && formCsrf == cookieCsrf.Value {
			handler(ctx)
		} else {
			ctx.Error(http.StatusForbidden)
		}
	}
}

func LoginCsrfHelper(handler func(ctx *Context)) func(*Context) {
	return LoginHelper(CsrfHelper(handler))
}
