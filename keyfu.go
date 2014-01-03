package main

import (
	"errors"
	"flag"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/martini"
)

const (
	defaultURL = "https://encrypted.google.com/search?q="
)

type Keyword struct {
	URL      string `toml:"url"`
	QueryURL string `toml:"query_url"`
}

type Server struct {
	Keywords map[string]Keyword `toml:"keyword"`
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

func (s *Server) Run(q string) (string, error) {
	key, value := parseQ(q)
	k, ok := s.Keywords[key]

	if !ok {
		return "", errors.New("not found")
	}

	var u string

	if len(value) > 0 {
		if len(k.QueryURL) == 0 {
			return "", errors.New("no query url for keyword")
		}
		u = k.QueryURL
	} else {
		if len(k.URL) == 0 {
			return "", errors.New("no url for keyword")
		}
		u = k.URL
	}

	return strings.Replace(u, "%s", url.QueryEscape(value), -1), nil
}

func (s *Server) Handle(res http.ResponseWriter, req *http.Request) {
	q := req.FormValue("q")

	u, err := s.Run(q)
	if err != nil {
		u = defaultURL + url.QueryEscape(q)
	}

	http.Redirect(res, req, u, 302)
}

func main() {
	var path = flag.String("c", "keyfu.conf", "KeyFu configuration file")
	flag.Parse()

	s := Server{}

	if _, err := toml.DecodeFile(*path, &s); err != nil {
		panic(err)
	}

	m := martini.Classic()
	m.Get("/run", s.Handle)
	m.Run()
}
