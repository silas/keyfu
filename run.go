package main

import (
	"errors"
	"strings"
	"net/url"
	"unicode"
)

var keywords = map[string]string{
	"a":    "http://www.amazon.com/\nhttp://www.amazon.com/s?url=search-alias%3Daps&field-keywords=%s",
	"ab":   "http://www.amazon.com/b?node=283155\nhttp://www.amazon.com/s?url=search-alias%3Dstripbooks&field-keywords=%s",
	"av":   "https://www.amazon.com/gp/video/library\nhttp://www.amazon.com/s/?url=search-alias%3Dinstant-video&field-keywords=%s",
	"digg":   "http://digg.com/reader",
	"gc":   "https://www.google.com/calendar/render\nhttps://www.google.com/calendar/render?q=%s",
	"gh":   "https://github.com/\nhttps://github.com/search?q=%s",
	"gi":   "https://www.google.com/imghp\nhttps://www.google.com/search?site=imghp&tbm=isch&q=%s",
	"gl":   "https://www.google.com/search?q=%s&btnI=I'm+Feeling+Lucky",
	"gm":   "https://mail.google.com/mail/\nhttps://mail.google.com/mail/#search/%s",
	"gd":   "https://drive.google.com/#starred\nhttps://drive.google.com/#search/%s",
	"maps":  "https://maps.google.com/\nhttps://maps.google.com/maps?q=%s",
	"name": "https://www.name.com/\nhttps://www.name.com/name?domain=%s",
	"w": "https://en.wikipedia.org/wiki/Main_Page\nhttps://en.wikipedia.org/wiki/Special:Search?search=%s",
}

func parseQ(value string) (string, string) {
	if value == "" {
		return "", ""
	}

	begin := -1
	end := len(value)
	for i, r := range value {
		if !unicode.IsSpace(r) {
			if begin == -1 {
				begin = i
			}
		} else if begin != -1 {
			end = i
			break
		}
	}

	if begin == -1 {
		return "", ""
	}

	return value[begin:end], strings.TrimLeftFunc(value[end:], unicode.IsSpace)
}

func Run(q string) (string, error) {
	key, value := parseQ(q)
	body, ok := keywords[key]

	if !ok {
		return "", errors.New("not found")
	}

	links := strings.Fields(body)

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
