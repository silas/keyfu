package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Uid     string
}

func (ctx *Context) GetCookie(name string) *http.Cookie {
	if cookie, err := ctx.Request.Cookie(name); err == nil {
		return cookie
	}
	return nil
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	zone, _ := cookie.Expires.Zone()
	if zone == "" && cookie.MaxAge != 0 {
		cookie.Expires = time.Unix(time.Now().UTC().Unix()+int64(cookie.MaxAge), 0).UTC()
	}
	http.SetCookie(ctx.Writer, cookie)
}

func (ctx *Context) GetMessage() (msgType string, msgText string) {
	if cookie := ctx.GetCookie("msg"); cookie != nil {
		value, _ := url.QueryUnescape(cookie.Value)
		parts := strings.SplitN(value, "|", 2)
		if len(parts) > 1 {
			msgType = parts[0]
			msgText = parts[1]
		}

		ctx.SetCookie(&http.Cookie{Name: "msg", MaxAge: -1})
	}

	return
}

func (ctx *Context) SetMessage(messageType string, messageText string) {
	ctx.SetCookie(&http.Cookie{Name: "msg", Value: url.QueryEscape(messageType + "|" + messageText), MaxAge: 120})
}

func (ctx *Context) SetError(text string) {
	ctx.SetMessage("error", text)
}

func (ctx *Context) SetSuccess(text string) {
	ctx.SetMessage("success", text)
}

func (ctx *Context) Redirect(url_ string) {
	http.Redirect(ctx.Writer, ctx.Request, url_, 302)
}

func (ctx *Context) Error(code int) {
	ctx.Writer.WriteHeader(code)

	switch code {
	case 404:
		ctx.Render("404.html", map[string]interface{}{"title": "Not Found"})
	case 403:
		ctx.Render("403.html", map[string]interface{}{"title": "Forbidden"})
	default:
		ctx.Render("error.html", map[string]interface{}{"title": http.StatusText(code)})
	}
}

func (ctx *Context) RenderString(name string, content map[string]interface{}) string {
	content["cdn"] = Config.CdnUrl
	content["css"] = Config.CssPath
	content["js"] = Config.JsPath
	content["logo"] = Config.LogoPath
	content["opensearch"] = Config.OpenSearchPath

	if ctx.Uid != "" {
		content["uid"] = ctx.Uid
	}

	if cookie := ctx.GetCookie("csrf"); cookie != nil {
		content["csrf"] = cookie.Value
	} else {
		csrf := uuid4()
		content["csrf"] = csrf
		ctx.SetCookie(&http.Cookie{Name: "csrf", Value: csrf, MaxAge: 1800})
	}

	if messageType, messageText := ctx.GetMessage(); messageType != "" {
		content[messageType] = messageText
	}

	html := bytes.NewBufferString("")
	Config.Templates.Lookup(name).Execute(html, content)

	return minifyHtml(string(html.Bytes()))
}

func (ctx *Context) Render(name string, content map[string]interface{}) {
	ctx.WriteString(ctx.RenderString(name, content))
}

func (ctx *Context) WriteString(text string) {
	io.WriteString(ctx.Writer, text)
}
