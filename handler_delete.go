package main

import (
	"labix.org/v2/mgo/bson"
)

func DeleteHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	q := values.Get("q")

	key, _ := parseQ(q)

	content := map[string]interface{}{
		"title": "Delete",
		"key":   key,
	}

	count, err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": key}).Count()
	if err != nil || count < 1 {
		content["error"] = "Keyword not found"
	}

	ctx.Render("delete.html", content)
}

func DeletePostHandler(ctx *Context) {
	key := ctx.Request.FormValue("key")

	content := map[string]interface{}{"title": "Delete", "key": key}

	keyword := UserKeyword{}
	err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": key}).One(&keyword)
	if err == nil {
		err = Config.Keyword.Remove(bson.M{"_id": keyword.Id})
		if err == nil {
			ctx.SetSuccess("Keyword deleted.")
			ctx.Redirect("/")
			return
		} else {
			content["error"] = "Unable to delete keyword"
		}
	} else {
		content["error"] = "Keyword not found"
	}

	ctx.Render("delete.html", content)
}
