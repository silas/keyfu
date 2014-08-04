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

	"github.com/BurntSushi/toml"
	"github.com/robertkrimen/otto"
)

var (
	defaultURL = "https://encrypted.google.com/search?q="
	minTimeout = time.Duration(10 * time.Millisecond)

	errHalt    = errors.New("halt")
	errTimeout = errors.New("timeout")
	errParse   = errors.New("parse error")
)

type Config struct {
	Path    []string
	URL     string
	Listen  string
	Timeout time.Duration
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
	Config    Config
	StartTime time.Time
	vm        *otto.Otto
}

// NewServer creates and sets up a new server.
func NewServer(path string) (*Server, error) {
	var err error

	s := Server{}

	s.vm = otto.New()

	if b, err := Asset("src/runtime.js"); err == nil {
		s.vm.Run(b)
	} else {
		return nil, err
	}

	s.StartTime = time.Now()

	if _, err = toml.DecodeFile(path, &s.Config); err != nil {
		return nil, err
	}

	if s.Config.Timeout < minTimeout {
		s.Config.Timeout = minTimeout
	}

	for _, path := range strings.Split(os.Getenv("KEYFU_PATH"), ":") {
		s.Config.Path = append(s.Config.Path, path)
	}

	for i, path := range s.Config.Path {
		if s.Config.Path[i], err = filepath.Abs(path); err != nil {
			return nil, err
		}
	}

	if s.Config.Listen == "" {
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")

		if port == "" {
			port = "8000"
		}

		s.Config.Listen = net.JoinHostPort(host, port)
	}

	if s.Config.URL == "" {
		host, port, err := net.SplitHostPort(s.Config.Listen)

		if err != nil {
			return nil, err
		}

		if host == "" {
			host = "localhost"
		}

		s.Config.URL = fmt.Sprintf("http://%s:%s", host, port)
	}

	return &s, nil
}

// StopRun reports errors or redirects to default URL.
func (s *Server) StopRun(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
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

	for _, dirPath := range s.Config.Path {
		path := filepath.Join(dirPath, fileName)

		stat, err := os.Stat(path)

		if err == nil && stat.Mode().IsRegular() {
			filePath = path
			break
		}
	}

	if filePath == "" {
		s.StopRun(w, r, nil)
	}

	code, err := ioutil.ReadFile(filePath)

	if err != nil {
		s.StopRun(w, r, errTimeout)
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
func (s *Server) Run() {
	http.HandleFunc("/run", s.RunHandler)
	http.HandleFunc("/opensearch.xml", s.OpenSearchHandler)
	http.HandleFunc("/", s.StaticHandler)

	log.Fatal(http.ListenAndServe(s.Config.Listen, nil))
}

func main() {
	var path = flag.String("c", "keyfu.conf", "KeyFu configuration file")
	flag.Parse()

	s, err := NewServer(*path)

	if err != nil {
		log.Fatal(err.Error())
	}

	s.Run()
}
