package main

import (
	"fmt"
	"net/url"
)

func RunHandler(ctx *Context) {
	if ctx.Uid == "" {
		ctx.Redirect(fmt.Sprintf("/login?url=%s", url.QueryEscape(ctx.Request.URL.String())))
		return
	}

	values := ctx.Request.URL.Query()
	key, value := parseQ(values.Get("q"))

	if key == "" {
		ctx.Redirect("/")
		return
	}

	keyword, err := NewKeyword(ctx.Uid, key)
	if err != nil {
		keyword, err = NewKeyword(ctx.Uid, ":")
		if err != nil {
			keyword = builtinKeywords["com.google"]
		}
		if key != ":" {
			value = values.Get("q")
		}
	}

	var data string
	data, err = keyword.Run(ctx.Uid, value)
	if err == nil {
		ctx.Redirect(data)
	} else {
		ctx.Render("run.html", map[string]interface{}{"title": "Error", "bodytitle": "", "runtimeerror": err.Error(), "q": values.Get("q")})
	}
}
