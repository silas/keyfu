package main

import (
	"net/http"
)

type Handler func(ctx *Context)

type Server struct {
	delete  pathRoute
	get     pathRoute
	head    pathRoute
	options pathRoute
	post    pathRoute
	put     pathRoute
	trace   pathRoute
}

func NewServer() *Server {
	return &Server{
		delete:  pathRoute{},
		get:     pathRoute{},
		head:    pathRoute{},
		options: pathRoute{},
		post:    pathRoute{},
		put:     pathRoute{},
		trace:   pathRoute{},
	}
}

func (s *Server) createHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{w, r, ""}
		if cookie := ctx.GetSecureCookie("uid"); cookie != nil && len(cookie.Value) == 24 {
			ctx.Uid = cookie.Value
		}
		switch r.Method {
		case "DELETE":
			if handler, ok := s.delete[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusMethodNotAllowed)
			}
		case "GET":
			if handler, ok := s.get[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusNotFound)
			}
		case "HEAD":
			if handler, ok := s.head[r.URL.Path]; ok {
				handler(&ctx)
			} else if handler, ok := s.get[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusNotFound)
			}
		case "OPTIONS":
			if handler, ok := s.options[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusMethodNotAllowed)
			}
		case "POST":
			if handler, ok := s.post[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusMethodNotAllowed)
			}
		case "PUT":
			if handler, ok := s.put[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusMethodNotAllowed)
			}
		case "TRACE":
			if handler, ok := s.put[r.URL.Path]; ok {
				handler(&ctx)
			} else {
				ctx.Error(http.StatusMethodNotAllowed)
			}
		default:
			ctx.Error(http.StatusNotImplemented)
		}
	}
}

func (s *Server) Delete(p string, h Handler) {
	s.delete[p] = h
}

func (s *Server) Get(p string, h Handler) {
	s.get[p] = h
}

func (s *Server) Head(p string, h Handler) {
	s.head[p] = h
}

func (s *Server) Options(p string, h Handler) {
	s.options[p] = h
}

func (s *Server) Post(p string, h Handler) {
	s.post[p] = h
}

func (s *Server) Put(p string, h Handler) {
	s.put[p] = h
}

func (s *Server) Trace(p string, h Handler) {
	s.trace[p] = h
}

func (s *Server) Run(addr, path string) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path))))
	http.HandleFunc("/", s.createHandler())
	http.ListenAndServe(addr, nil)
}
