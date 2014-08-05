package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertParse(t *testing.T, raw, key, value string, err error) {
	k, v, e := parse(raw)

	assert.Equal(t, k, key)
	assert.Equal(t, v, value)
	assert.Equal(t, e, err)
}

func TestParse(t *testing.T) {
	assertParse(t, "", "", "", errParse)
	assertParse(t, "one", "one", "", nil)
	assertParse(t, " one ", "one", "", nil)
	assertParse(t, "one two", "one", "two", nil)
	assertParse(t, "one two three", "one", "two three", nil)
	assertParse(t, "one  two  three ", "one", "two  three ", nil)
}

func run(t *testing.T, ts *httptest.Server, q string) (url, body string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?q="+q, ts.URL), nil)
	if !assert.Nil(t, err) {
		return
	}

	tr := &http.Transport{}
	res, err := tr.RoundTrip(req)
	if assert.Nil(t, err) && assert.NotNil(t, res) {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if assert.Nil(t, err) {
			return res.Header.Get("Location"), string(b)
		}
	}

	return
}

func assertRunLocation(t *testing.T, ts *httptest.Server, q, value string) {
	location, _ := run(t, ts, q)
	assert.Equal(t, location, value)
}

func assertRunBody(t *testing.T, ts *httptest.Server, q, value string) {
	_, body := run(t, ts, q)
	assert.Equal(t, body, value)
}

func TestServer(t *testing.T) {
	os.Setenv("HOST", "")
	os.Setenv("PORT", "")

	s, err := NewServer("./test/keyfu.conf")
	if assert.Nil(t, err) {
		assert.Equal(t, s.Config.Listen, ":8000")
		assert.Equal(t, s.Config.URL, "http://localhost:8000")
	}

	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("PORT", "1234")

	s, err = NewServer("./test/keyfu.conf")
	if assert.Nil(t, err) {
		assert.Equal(t, s.Config.Listen, "0.0.0.0:1234")
		assert.Equal(t, s.Config.URL, "http://0.0.0.0:1234")
	}

	ts := httptest.NewServer(http.HandlerFunc(s.RunHandler))
	defer ts.Close()

	assertRunLocation(t, ts, "gopher+food", "https://encrypted.google.com/search?q=gopher+food")
	assertRunLocation(t, ts, "empty+test", "https://encrypted.google.com/search?q=empty+test")
	assertRunLocation(t, ts, "amazon", "http://www.amazon.com/")
	assertRunLocation(t, ts, "amazon+test", "http://www.amazon.com/s?url=search-alias%3Daps&field-keywords=test")
	assertRunBody(t, ts, "hello", "Hello, World!")
	assertRunBody(t, ts, "hello+Jane", "Hello, Jane!")
	assertRunBody(t, ts, "timeout", "Error: timeout")
}

func assertBodyContains(t *testing.T, handler http.HandlerFunc, path string, statusCode int, body string) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	res, err := http.Get(ts.URL + path)
	assert.Equal(t, res.StatusCode, statusCode)
	if assert.Nil(t, err) && assert.NotNil(t, res) {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if assert.Nil(t, err) {
			assert.Contains(t, string(b), body)
		}
	}
}

func TestStatic(t *testing.T) {
	s, err := NewServer("./test/keyfu.conf")
	if !assert.Nil(t, err) {
		return
	}

	assertBodyContains(t, s.OpenSearchHandler, "/opensearch.xml", 200, "http://0.0.0.0:1234")
	assertBodyContains(t, s.StaticHandler, "/", 200, "<head>")
	assertBodyContains(t, s.StaticHandler, "/robots.txt", 200, "User-agent")
	assertBodyContains(t, s.StaticHandler, "/bla.html", 404, "not found")
}
