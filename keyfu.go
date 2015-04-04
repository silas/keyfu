package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/robertkrimen/otto"
)

const (
	defaultPort = "8000"
	defaultURL  = "https://encrypted.google.com/search?q="
)

var (
	defaultTimeout = time.Duration(10 * time.Millisecond)

	errHalt    = errors.New("halt")
	errTimeout = errors.New("timeout")
	errParse   = errors.New("parse error")
)

type Config struct {
	Path    string
	URL     string
	Listen  string
	Timeout time.Duration
}

func (c *Config) setup() error {
	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}

	if c.Listen == "" {
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")

		if port == "" {
			port = defaultPort
		}

		c.Listen = net.JoinHostPort(host, port)
	}

	if c.URL == "" {
		host, port, err := net.SplitHostPort(c.Listen)

		if err != nil {
			return err
		}

		if host == "" {
			host = "127.0.0.1"
		}

		c.URL = fmt.Sprintf("http://%s:%s", host, port)
	}

	c.Path += ":" + os.Getenv("KEYFU_PATH")
	c.Path += ":" + filepath.Join(os.Getenv("HOME"), ".keyfu")

	return nil
}

// parse takes a string in the form "keyword [query]" and returns the parsed
// key/value pair or an error.
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

// Server holds the application state.
type Server struct {
	Config    *Config
	StartTime time.Time
	vm        *otto.Otto
	path      []string
}

// NewServer creates and sets up a new server.
func NewServer(c Config) (*Server, error) {
	if err := c.setup(); err != nil {
		return nil, err
	}

	s := Server{Config: &c}
	s.vm = otto.New()

	if b, err := Asset("lib/runtime.js"); err == nil {
		s.vm.Run(b)
	} else {
		return nil, err
	}

	for _, path := range strings.Split(s.Config.Path, ":") {
		if path == "" {
			continue
		}
		if absPath, err := filepath.Abs(path); err == nil {
			for _, p := range s.path {
				if p == absPath {
					continue
				}
			}
			s.path = append(s.path, absPath)
		}
	}

	s.StartTime = time.Now()

	for _, dir := range s.path {
		paths, err := filepath.Glob(filepath.Join(dir, "lib", "*.js"))

		if err != nil {
			return nil, err
		}

		for _, path := range paths {
			if b, err := ioutil.ReadFile(path); err == nil {
				s.vm.Run(b)
			} else {
				return nil, err
			}
		}
	}

	return &s, nil
}

// StopRun reports errors or redirects to default URL.
func (s *Server) StopRun(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil || err.Error() == "skip" {
		http.Redirect(w, r, defaultURL+url.QueryEscape(r.FormValue("q")), 302)
	} else {
		io.WriteString(w, fmt.Sprintf("Error: %s", err.Error()))
	}
}

// RunHandler executes the keyword code and sends a response.
func (s *Server) RunHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")

	key, value, err := parse(q)

	if err != nil || key == "" {
		s.StopRun(w, r, err)
		return
	}

	filePath := ""
	fileName := fmt.Sprintf("%s.js", key)

	for _, dirPath := range s.path {
		path := filepath.Join(dirPath, fileName)

		stat, err := os.Stat(path)

		if err == nil && stat.Mode().IsRegular() {
			filePath = path
			break
		}
	}

	if filePath == "" {
		s.StopRun(w, r, nil)
		return
	}

	code, err := ioutil.ReadFile(filePath)

	if err != nil {
		s.StopRun(w, r, err)
		return
	}

	defer func() {
		if caught := recover(); caught != nil {
			if caught == errHalt {
				s.StopRun(w, r, errTimeout)
				return
			}
			if err, ok := caught.(error); ok {
				s.StopRun(w, r, err)
			}
			panic(caught)
		}
	}()

	vm := s.vm.Copy()
	vm.Interrupt = make(chan func(), 1)

	go func() {
		time.Sleep(s.Config.Timeout)
		vm.Interrupt <- func() {
			panic(errHalt)
		}
	}()

	vm.Set("query", value)

	if _, err := vm.Run(code); err != nil {
		s.StopRun(w, r, err)
		return
	}

	if location, err := vm.Get("location"); err == nil && location.IsString() {
		if location, err := location.ToString(); err == nil {
			http.Redirect(w, r, location, 302)
			return
		}
	}

	if body, err := vm.Get("body"); err == nil && body.IsString() {
		if body, err := body.ToString(); err == nil {
			io.WriteString(w, body)
			return
		}
	}

	s.StopRun(w, r, nil)
}

func (s *Server) OpenSearchHandler(w http.ResponseWriter, r *http.Request) {
	if b, err := Asset("static/opensearch.xml"); err == nil {
		w.Write(bytes.Replace(b, []byte("http://www.keyfu.com"), []byte(s.Config.URL), 1))
	} else {
		http.NotFound(w, r)
	}
}

// StaticHandler serves embeded static content.
func (s *Server) StaticHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	if p[len(p)-1] == '/' {
		p = p + "index.html"
	}

	b, err := Asset("static" + p)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, name := path.Split(p)

	http.ServeContent(w, r, name, s.StartTime, bytes.NewReader(b))
}

// Init reads configuration and sets up server state.
// Run starts HTTP server.
func (s *Server) Run() error {
	http.HandleFunc("/run", s.RunHandler)
	http.HandleFunc("/opensearch.xml", s.OpenSearchHandler)
	http.HandleFunc("/", s.StaticHandler)

	return http.ListenAndServe(s.Config.Listen, nil)
}

func main() {
	var listen = flag.String("listen", ":"+defaultPort, "listen address")
	var path = flag.String("path", "", "keyfu path")
	var timeout = flag.Duration("timeout", defaultTimeout, "run timeout")
	var url = flag.String("url", "", "serve URL")

	flag.Parse()

	c := Config{
		Listen:  *listen,
		Path:    *path,
		Timeout: *timeout,
		URL:     *url,
	}

	s, err := NewServer(c)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(s.Run())
}
