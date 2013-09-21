package main

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"regexp"
)

func AutocompleteHandler(ctx *Context) {
	values := ctx.Request.URL.Query()
	key, _ := parseQ(values.Get("q"))

	enc := json.NewEncoder(ctx.Writer)

	var v = make([]interface{}, 0, 2)
	v = append(v, key)

	if ctx.Uid == "" || key == "" {
		v = append(v, []string{})
	} else {
		var results []struct{ Key string }
		var keys = make([]string, 0)
		if err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(ctx.Uid), "key": bson.RegEx{"^" + regexp.QuoteMeta(key), ""}}).Select(bson.M{"key": 1}).Limit(5).All(&results); err == nil {
			for _, r := range results {
				keys = append(keys, r.Key)
			}
			v = append(v, keys)
		} else {
			v = append(v, []string{})
		}
	}

	if err := enc.Encode(&v); err != nil {
		// handle error
	}
}
