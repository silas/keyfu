package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/martini"
)

const (
	defaultURL     = "https://encrypted.google.com/search?q="
	defaultTimeout = 5 * time.Second
)

var (
	errEnvKilled  = errors.New("keyfu: env killed")
	errNoQueryUrl = errors.New("keyfu: no query url")
	errNoUrl      = errors.New("keyfu: no url")
	errNotFound   = errors.New("keyfu: not found")
	errUnknownEnv = errors.New("keyfu: unknown environment")
)

type Env struct {
	Path    string   `toml:"path"`
	Timeout duration `toml:"timeout"`
}

type Keyword struct {
	Env      string `toml:"env"`
	URL      string `toml:"url"`
	QueryURL string `toml:"query_url"`
}

type Server struct {
	Keywords map[string]Keyword `toml:"keyword"`
	Envs     map[string]Env     `toml:"env"`
}

type duration struct {
	Duration time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
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

func (e *Env) Run(k *Keyword, v string) (string, error) {
	cmd := exec.Command(e.Path, v)

	cmd.Env = []string{}

	if k.URL != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("URL=%s", k.URL))
	}

	if k.QueryURL != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("QUERY_URL=%s", k.QueryURL))
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(e.Timeout.Duration * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			cmd.Wait()
			return "", err
		}
		<-done
		return "", errEnvKilled
	case err := <-done:
		if err != nil {
			return "", err
		}
	}

	text, err := out.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

func (k *Keyword) Run(v string) (string, error) {
	var u string

	if len(v) > 0 {
		if len(k.QueryURL) == 0 {
			return "", errNoQueryUrl
		}
		u = k.QueryURL
	} else {
		if len(k.URL) == 0 {
			return "", errNoUrl
		}
		u = k.URL
	}

	return strings.Replace(u, "%s", url.QueryEscape(v), -1), nil
}

func (s *Server) Init(path string) error {
	if _, err := toml.DecodeFile(path, &s); err != nil {
		return err
	}

	for _, env := range s.Envs {
		if env.Timeout.Duration == 0 {
			env.Timeout.Duration = defaultTimeout
		}
	}

	return nil
}

func (s *Server) Run(q string) (string, error) {
	key, value := parseQ(q)

	k, ok := s.Keywords[key]
	if !ok {
		return "", errNotFound
	}

	if k.Env != "" {
		e, ok := s.Envs[k.Env]
		if !ok {
			return "", errUnknownEnv
		}

		return e.Run(&k, value)
	}

	return k.Run(value)
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
	if err := s.Init(*path); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	m := martini.Classic()
	m.Get("/run", s.Handle)
	m.Run()
}
