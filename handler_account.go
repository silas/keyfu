package main

import (
	"labix.org/v2/mgo/bson"
	"net/url"
)

func AccountHandler(ctx *Context) {
	user := User{}
	if err := Config.User.Find(bson.M{"_id": bson.ObjectIdHex(ctx.Uid)}).One(&user); err == nil {
		ctx.Render("account.html", map[string]interface{}{"title": "Account", "email": user.Email})
	} else {
		ctx.Redirect("/logout?url=" + url.QueryEscape(ctx.Request.URL.String()))
	}
}

func AccountPostHandler(ctx *Context) {
	email := ctx.Request.FormValue("email")
	currentPassword := ctx.Request.FormValue("current-password")
	newPassword := ctx.Request.FormValue("new-password")
	errStr := ""

	user := User{}

	if err := Config.User.Find(bson.M{"_id": bson.ObjectIdHex(ctx.Uid)}).One(&user); err != nil {
		errStr = "Failed to load account."
		goto END
	}

	if email != user.Email {
		if count, err := Config.User.Find(bson.M{"_id": bson.M{"$ne": user.Id}, "email": email}).Count(); err != nil || count > 0 {
			errStr = "There is an existing account associated with that email."
			goto END
		}
	}

	if err := user.CheckPassword(currentPassword); err != nil {
		errStr = err.Error()
		goto END
	}

	if err := user.SetEmail(email); err != nil {
		errStr = err.Error()
		goto END
	}

	if newPassword != "" {
		if err := user.SetPassword(newPassword); err != nil {
			errStr = err.Error()
			goto END
		}
	}

	if err := Config.User.Update(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"email": user.Email, "password": user.Password}}); err != nil {
		errStr = "Failed to update account."
	}

END:

	if errStr != "" {
		content := map[string]interface{}{
			"title": "Account",
			"email": email,
			"error": errStr,
		}

		ctx.Render("account.html", content)
	} else {
		ctx.SetSuccess("Account successfully updated.")
		ctx.Redirect("/account")
	}
}
