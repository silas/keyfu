package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

func SignupHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	email := values.Get("email")

	if ctx.Uid != "" {
		ctx.Redirect("/")
		return
	}

	content := map[string]interface{}{
		"email": email,
		"title": "Sign up",
	}

	ctx.Render("signup.html", content)
}

func SignupPostHandler(ctx *Context) {
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")
	errStr := ""

	user := User{}

	if err := validateEmail(email); err == nil {
		user.Email = email
	} else {
		errStr = err.Error()
		goto END
	}

	if err := user.SetPassword(password); err != nil {
		errStr = err.Error()
		goto END
	}

	if count, err := Config.User.Find(bson.M{"email": email}).Count(); err != nil || count > 0 {
		errStr = "There is an existing account associated with that email."
		goto END
	}

	if err := Config.User.Insert(bson.M{"email": user.Email, "password": user.Password, "created_timestamp": time.Now()}); err != nil {
		errStr = "Failed to create account."
	}

END:

	if errStr != "" {
		content := map[string]interface{}{
			"email": email,
			"title": "Sign up",
			"error": errStr,
		}
		ctx.Render("signup.html", content)
	} else {
		ctx.SetSuccess("Account successfully created.")
		ctx.Redirect("/login?url=%2Fhelp%2Fgetting-started")
	}
}
