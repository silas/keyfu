package main

import (
	"labix.org/v2/mgo/bson"
)

func RenameValidation(uid string, srcKey string, dstKey string, content map[string]interface{}) map[string]interface{} {
	srcKeyError := false
	if count, err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(uid), "key": srcKey}).Count(); err == nil && count < 1 {
		srcKeyError = true
	}

	dstKeyError := false
	if count, err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(uid), "key": dstKey}).Count(); err == nil && count > 0 {
		dstKeyError = true
	}

	if srcKeyError && dstKeyError {
		content["error"] = "The source key doesn't exist and the destination key conflicts with an existing keyword."
		content["srckeyerror"] = true
		content["dstkeyerror"] = true
	} else if srcKeyError {
		content["error"] = "The source key doesn't exist."
		content["srckeyerror"] = true
	} else if dstKeyError {
		if srcKey == dstKey {
			content["error"] = "The source and destination keyword are the same."
		} else {
			content["error"] = "The destination key conflicts with an existing keyword."
		}
		content["dstkeyerror"] = true
	}

	return content
}

func RenameHandler(ctx *Context) {
	values := ctx.Request.URL.Query()

	srcKey, value := parseQ(values.Get("q"))
	dstKey, _ := parseQ(value)

	content := map[string]interface{}{
		"title":  "Rename",
		"srckey": srcKey,
		"dstkey": dstKey,
	}

	content = RenameValidation(ctx.Uid, srcKey, dstKey, content)

	ctx.Render("rename.html", content)
}

func RenamePostHandler(ctx *Context) {
	srcKey := ctx.Request.FormValue("srckey")
	dstKey := ctx.Request.FormValue("dstkey")

	content := map[string]interface{}{
		"title":  "Rename",
		"srckey": srcKey,
		"dstkey": dstKey,
	}

	content = RenameValidation(ctx.Uid, srcKey, dstKey, content)
	if _, ok := content["error"]; ok {
		ctx.Render("rename.html", content)
		return
	}

	if err := Config.Keyword.Update(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": srcKey}, bson.M{"$set": bson.M{"key": dstKey}}); err != nil {
		content["error"] = "Failed to rename keyword."
		ctx.Render("rename.html", content)
	} else {
		ctx.SetSuccess("Keyword renamed.")
		ctx.Redirect("/")
	}
}
