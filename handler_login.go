package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/hex"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"strings"
)

func LoginUrl(url_ string) string {
	if url_ == "" ||
		url_[0] != '/' ||
		// Catch login and logout
		strings.Index(url_, "/log") >= 0 ||
		strings.Index(url_, "/run?q=:log") >= 0 {
		return "/"
	}

	return url_
}

func LoginHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	email := values.Get("email")
	url_ := LoginUrl(values.Get("url"))

	if ctx.Uid != "" {
		ctx.Redirect(url_)
		return
	}

	content := map[string]interface{}{
		"url":     url_,
		"email":   email,
		"title":   "Login",
		"persist": "yes",
	}
	ctx.Render("login.html", content)
}

func LoginPostHandler(ctx *Context) {
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")
	url_ := ctx.Request.FormValue("url")
	persist := ctx.Request.FormValue("persist")

	url_, _ = url.QueryUnescape(url_)
	url_ = LoginUrl(url_)

	user := User{}
	err := Config.User.Find(bson.M{"email": email}).One(&user)
	if err == nil && password != "" && user.Password != "" {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
			cookie := http.Cookie{Name: "uid", Value: hex.EncodeToString([]byte(string(user.Id)))}
			if persist == "yes" {
				cookie.MaxAge = int(Config.SessionLength)
			}
			ctx.SetSecureCookie(&cookie)
			ctx.Redirect(url_)
			return
		}
	}

	content := map[string]interface{}{
		"url":   url_,
		"email": email,
		"error": "Invalid email or password.",
		"title": "Login",
	}

	if persist == "yes" {
		content["persist"] = "yes"
	}

	ctx.Render("login.html", content)
}
