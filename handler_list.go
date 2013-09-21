package main

import (
	"labix.org/v2/mgo/bson"
	"strings"
)

func ListHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	filter := strings.Replace("^"+values.Get("q"), "*", ".*", -1)

	var keywords []UserKeyword

	iter := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": bson.M{"$regex": filter}}).Sort("key").Iter()
	for {
		var result UserKeyword
		if ok := iter.Next(&result); ok {
			keywords = append(keywords, result)
		} else {
			break
		}
	}

	content := map[string]interface{}{
		"title":    "List",
		"keywords": keywords,
	}
	ctx.Render("list.html", content)
}
