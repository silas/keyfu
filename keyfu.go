package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
)

const (
	Redirect = iota
	Render
)

var (
	defaultTimeout = 5 * time.Second
	defaultURL     = "https://encrypted.google.com/search?q="
	types          = map[string]int{"redirect": Redirect, "render": Render}
	static         = map[string]func() []byte{}

	errEnvTimeout      = errors.New("keyfu: env timed out")
	errNoUrl           = errors.New("keyfu: no url")
	errNoQueryUrl      = errors.New("keyfu: no query url")
	errParse           = errors.New("keyfu: parse error")
	errLinkConfig      = errors.New("keyfu: invalid link keyword")
	errProgramName     = errors.New("keyfu: program name invalid")
	errProgramResponse = errors.New("keyfu: program response invalid")
)

type Config struct {
	Listen   string                       `toml:"listen"`
	Keywords map[string]map[string]string `toml:"keyword"`
}

type Response struct {
	Type int
	Body string
}

type Keyword interface {
	Run(r *Request) (*Response, error)
}

func parse(v string) (string, string, error) {
	if v == "" {
		return "", "", errParse
	}

	begin := -1
	end := len(v)
	for i, r := range v {
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
		return "", "", errParse
	}

	return v[begin:end], strings.TrimLeftFunc(v[end:], unicode.IsSpace), nil
}

type LinkKeyword struct {
	URL      string
	QueryURL string
}

func NewLinkKeyword(c map[string]string) (*LinkKeyword, error) {
	k := LinkKeyword{c["url"], c["query_url"]}
	if k.URL == "" && k.QueryURL == "" {
		return nil, errLinkConfig
	}
	return &k, nil
}

func (k LinkKeyword) Run(req *Request) (*Response, error) {
	var u string

	if len(req.Value) > 0 {
		if len(k.QueryURL) == 0 {
			return nil, errNoQueryUrl
		}
		u = k.QueryURL
	} else {
		if len(k.URL) == 0 {
			return nil, errNoUrl
		}
		u = k.URL
	}

	return &Response{Redirect, strings.Replace(u, "%s", url.QueryEscape(req.Value), -1)}, nil
}

type ProgramKeyword struct {
	Name    string
	Timeout time.Duration
}

func NewProgramKeyword(c map[string]string) (*ProgramKeyword, error) {
	k := ProgramKeyword{}

	if name, ok := c["name"]; ok && len(name) > 0 {
		k.Name = name
	} else {
		return nil, errProgramName
	}

	if timeout, ok := c["timeout"]; ok {
		if t, err := time.ParseDuration(timeout); err != nil {
			return nil, err
		} else {
			k.Timeout = t
		}
	} else {
		k.Timeout = defaultTimeout
	}

	return &k, nil
}

func (r ProgramKeyword) Run(req *Request) (*Response, error) {
	cmd := exec.Command(r.Name, req.Query, req.Key, req.Value)
	cmd.Env = []string{}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(r.Timeout):
		if err := cmd.Process.Kill(); err != nil {
			cmd.Wait()
			return nil, err
		}
		<-done
		return nil, errEnvTimeout
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	c, b, err := parse(out.String())
	if err != nil {
		return nil, err
	}

	var res Response

	var ok bool
	if res.Type, ok = types[c]; !ok {
		return nil, errProgramResponse
	}

	res.Body = b

	return &res, nil
}

type Request struct {
	Query string
	Key   string
	Value string
}

func NewRequest(q string) (*Request, error) {
	r := new(Request)

	if err := r.Parse(q); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Request) Parse(q string) (err error) {
	r.Query = q
	r.Key, r.Value, err = parse(q)
	return
}

type Server struct {
	Config    Config
	Keywords  map[string]Keyword
	StartTime time.Time
}

func (s *Server) RunError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		log.Print(err)
	}
	q := r.FormValue("q")
	http.Redirect(w, r, defaultURL+url.QueryEscape(q), 302)
}

func (s *Server) RunHandler(w http.ResponseWriter, r *http.Request) {
	req, err := NewRequest(r.FormValue("q"))
	if err != nil {
		s.RunError(w, r, err)
		return
	}

	k, ok := s.Keywords[req.Key]
	if !ok {
		s.RunError(w, r, nil)
		return
	}

	res, err := k.Run(req)
	if err != nil {
		s.RunError(w, r, err)
		return
	}

	if res.Type == Redirect {
		http.Redirect(w, r, res.Body, 302)
	} else {
		io.WriteString(w, res.Body)
	}
}

func (s *Server) StaticHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p[len(p)-1] == '/' {
		p = p + "index.html"
	}
	b, found := static["/static"+p]
	if !found {
		http.NotFound(w, r)
		return
	}
	_, name := path.Split(p)
	http.ServeContent(w, r, name, s.StartTime, bytes.NewReader(b()))
}

func (s *Server) Load() error {
	var err error
	var keyword Keyword

	s.Keywords = map[string]Keyword{}

	for k, v := range s.Config.Keywords {
		switch v["type"] {
		case "program":
			keyword, err = NewProgramKeyword(v)
		case "":
			keyword, err = NewLinkKeyword(v)
		default:
			log.Printf("unknown type: keyword %s (%s)", k, v["type"])
			continue
		}

		if err == nil {
			s.Keywords[k] = keyword
		}
	}

	return nil
}

func (s *Server) Init(path string) error {
	if _, err := toml.DecodeFile(path, &s.Config); err != nil {
		return err
	}

	s.StartTime = time.Now()

	if s.Config.Listen == "" {
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")

		if port == "" {
			port = "8000"
		}

		s.Config.Listen = host + ":" + port
	}

	if err := s.Load(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	http.HandleFunc("/run", s.RunHandler)

	if len(static) > 0 {
		http.HandleFunc("/", s.StaticHandler)
	} else {
		http.Handle("/", http.FileServer(http.Dir("./static/")))
	}

	log.Fatal(http.ListenAndServe(s.Config.Listen, nil))
}

func main() {
	var path = flag.String("c", "keyfu.conf", "KeyFu configuration file")
	flag.Parse()

	s := Server{}

	if err := s.Init(*path); err != nil {
		log.Fatal(err.Error())
	}

	s.Run()
}
