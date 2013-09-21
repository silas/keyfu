package main

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/url"
)

func EditHandler(ctx *Context) {
	values := ctx.Request.URL.Query()

	key, value := parseQ(values.Get("q"))

	if k := values.Get("key"); k != "" {
		key = k
	}

	keyword := UserKeyword{}
	Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": key}).One(&keyword)

	kType := keyword.Type

	if kType == "" {
		kType = values.Get("type")
	}

	content := map[string]interface{}{"title": "Edit", "newkey": key}

	if kType == "alias" || kType == "builtin" || kType == "link" {
		content[kType] = "true"
	} else {
		content["link"] = "true"
	}

	if keyword.Body != "" {
		content["body"] = keyword.Body
	} else if body := values.Get("body"); body != "" {
		content["body"] = body
	} else {
		content["body"] = value
	}

	if keyword.Key != "" {
		content["oldkey"] = keyword.Key
	} else {
		content["title"] = "Add"
	}

	ctx.Render("edit.html", content)
}

func EditPostHandler(ctx *Context) {
	oldKey := ctx.Request.FormValue("oldkey")
	newKey := ctx.Request.FormValue("newkey")
	body := ctx.Request.FormValue("body")
	kType := ctx.Request.FormValue("type")

	var errStr string
	var err error

	keyword := UserKeyword{}
	if oldKey != "" {
		err = Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": oldKey}).One(&keyword)
	}

	if err = validateKey(newKey); err != nil {
		errStr = err.Error()
	}

	if err == nil {
		if oldKey != newKey || oldKey == "" {
			if count, _ := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": newKey}).Count(); count > 0 {
				errStr = "Key already exists."
				err = errors.New(errStr)
			}
		}
		if errStr == "" {
			if oldKey == "" {
				err = Config.Keyword.Insert(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": newKey, "type": kType, "body": body})
			} else {
				err = Config.Keyword.Update(bson.M{"_id": keyword.Id}, bson.M{"$set": bson.M{"key": newKey, "type": kType, "body": body}})
			}
		}
	}

	if err != nil {
		content := map[string]interface{}{
			"title":  "Edit",
			"newkey": newKey,
			kType:    true,
			"body":   body,
		}
		if errStr != "" {
			content["error"] = errStr
		}
		ctx.Render("edit.html", content)
	} else {
		ctx.SetSuccess("Keyword updated.")
		ctx.Redirect(fmt.Sprintf("/edit?q=%s", url.QueryEscape(newKey)))
	}
}
