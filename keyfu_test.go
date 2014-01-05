package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertParseQ(t *testing.T, v, r1, r2 string) {
	p1, p2 := parseQ(v)

	assert.Equal(t, p1, r1)
	assert.Equal(t, p2, r2)
}

func TestParseQ(t *testing.T) {
	assertParseQ(t, "", "", "")
	assertParseQ(t, "one", "one", "")
	assertParseQ(t, " one ", "one", "")
	assertParseQ(t, "one two", "one", "two")
	assertParseQ(t, "one two three", "one", "two three")
	assertParseQ(t, "one  two  three ", "one", "two  three ")
}

func server() *Server {
	s := Server{}

	s.Keywords = map[string]Keyword{
		"one":  Keyword{"", "one", ""},
		"two":  Keyword{"", "", "two=%s"},
		"both": Keyword{"", "url", "query_url=%s"},
	}

	return &s
}

func assertRun(t *testing.T, s *Server, q, url string) {
	u, err := s.Run(q)
	if assert.Nil(t, err) {
		assert.Equal(t, u, url)
	}
}

func assertNotRun(t *testing.T, s *Server, q string) {
	_, err := s.Run(q)
	assert.NotNil(t, err)
}

func TestEnvRun(t *testing.T) {
	env := Env{"echo", duration{5 * time.Second}}
	keyword := Keyword{"", "url", "query_url"}

	url, err := env.Run(&keyword, "one")
	if assert.Nil(t, err) {
		assert.Equal(t, url, "one")
	}
}

func TestServerRun(t *testing.T) {
	s := server()

	// not found
	assertNotRun(t, s, "")
	assertNotRun(t, s, "notfound")

	// only url
	assertRun(t, s, "one", "one")
	assertNotRun(t, s, "one a")

	// only query_url
	assertRun(t, s, "two 2", "two=2")
	assertNotRun(t, s, "two")

	// url and query_url
	assertRun(t, s, "both a", "query_url=a")
	assertRun(t, s, "both a", "query_url=a")
	assertRun(t, s, "both a b", "query_url=a+b")
}
