package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"strings"
)

func AddHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	filter := strings.Replace("^"+values.Get("q"), "*", ".*", -1)

	var catalog []Catalog

	iter := Config.Catalog.Find(bson.M{"sort": bson.M{"$regex": filter}}).Sort("sort").Iter()
	m := map[string]string{}
	i := 0
	for {
		var result Catalog
		if ok := iter.Next(&result); ok {
			if url, exists := m[result.Root]; exists {
				result.Root = url
			} else {
				if i < 9 {
					i += 1
				} else {
					i = 1
				}
				url := fmt.Sprintf("http://img%d.keyfu.net/%s.png", i, result.Root)
				m[result.Root] = url
				result.Root = url
			}
			catalog = append(catalog, result)
		} else {
			break
		}
	}

	content := map[string]interface{}{
		"title":   "Add",
		"catalog": catalog,
	}
	ctx.Render("add.html", content)
}
