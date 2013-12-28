package main

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/martini"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

type Keyword struct {
	Url      string `toml:"url"`
	QueryUrl string `toml:"query_url"`
}

type Config struct {
	Keywords map[string]Keyword `toml:"keyword"`
}

var config Config

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
	k, ok := config.Keywords[key]

	if !ok {
		return "", errors.New("not found")
	}

	var u string

	if len(value) > 0 {
		if len(k.QueryUrl) == 0 {
			return "", errors.New("no query url for keyword")
		}
		u = k.QueryUrl
	} else {
		if len(k.Url) == 0 {
			return "", errors.New("no url for keyword")
		}
		u = k.Url
	}

	return strings.Replace(u, "%s", url.QueryEscape(value), -1), nil
}

func main() {
	var configPath = flag.String("c", "keyfu.conf", "KeyFu configuration file")
	flag.Parse()

	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		panic(err)
	}

	m := martini.Classic()

	m.Get("/run", func(res http.ResponseWriter, req *http.Request) {
		q := req.FormValue("q")
		if v, err := Run(q); err == nil {
			http.Redirect(res, req, v, 302)
		} else {
			http.Redirect(res, req, "https://encrypted.google.com/search?q="+url.QueryEscape(q), 302)
		}
	})

	m.Run()
}
