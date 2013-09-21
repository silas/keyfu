package main

import (
	"errors"
	"labix.org/v2/mgo/bson"
	"net/url"
	"strings"
)

type Keyword struct {
	Type string
	Body string
}

func NewKeyword(uid string, key string) (*Keyword, error) {
	if len(key) > 1 && key[0] == ':' {
		if keyword, ok := builtinKeywords[key[1:]]; ok {
			return keyword, nil
		}
	}

	if uid != "" {
		keyword := Keyword{}
		if err := Config.Keyword.Find(bson.M{"uid": bson.ObjectIdHex(uid), "key": key}).One(&keyword); err == nil {
			return &keyword, nil
		}
	}

	return nil, errors.New("No keyword found")
}

func (k *Keyword) Run(uid string, value string) (string, error) {
	switch k.Type {
	case "alias":
		return k.Alias(uid, value)
	case "builtin":
		return k.Builtin(uid, value)
	case "link":
		return k.Link(value)
	}

	return "", errors.New("Keyword not found")
}

func (k *Keyword) Alias(uid string, value string) (string, error) {
	keyword, err := NewKeyword(uid, strings.TrimSpace(k.Body))

	if err != nil {
		return "", errors.New("Aliased keyword not found")
	}

	if keyword.Type == "alias" {
		return "", errors.New("Cannot alias an alias")
	}

	return keyword.Run(uid, value)
}

func (k *Keyword) Builtin(uid string, value string) (string, error) {
	if keyword, ok := builtinKeywords[strings.TrimSpace(k.Body)]; ok {
		return keyword.Run(uid, value)
	}

	return "", errors.New("Unknown builtin keyword")
}

func (k *Keyword) Link(value string) (string, error) {
	links := strings.Fields(k.Body)

	var link string
	switch len(links) {
	case 0:
		return "", nil
	case 1:
		if value != "" || strings.Index(links[0], "%s") >= 0 {
			link = links[0]
		} else {
			return links[0], nil
		}
	default:
		if value == "" {
			link = links[0]
		} else {
			link = links[1]
		}
	}

	return strings.Replace(link, "%s", url.QueryEscape(value), -1), nil
}
