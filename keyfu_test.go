package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func assertLinkKeyword(t *testing.T, q, url, queryURL, body string) {
	c := map[string]string{}

	if url != "" {
		c["url"] = url
	}

	if queryURL != "" {
		c["query_url"] = queryURL
	}

	k, err := NewLinkKeyword(c)
	if assert.Nil(t, err) {
		assert.Equal(t, k.URL, url)
		assert.Equal(t, k.QueryURL, queryURL)

		if req, err := NewRequest(q); assert.Nil(t, err) {
			if res, err := k.Run(req); assert.Nil(t, err) {
				assert.Equal(t, res.Body, body)
			}
		}
	}

}

func TestLinkKeyword(t *testing.T) {
	assertLinkKeyword(t, "key", "url1", "", "url1")
	assertLinkKeyword(t, "key one two", "", "query_url1=%s", "query_url1=one+two")
	assertLinkKeyword(t, "key", "url2", "query_url2", "url2")
	assertLinkKeyword(t, "key one two", "url2", "query_url2=%s", "query_url2=one+two")

	k, err := NewLinkKeyword(map[string]string{})
	assert.NotNil(t, err)
	assert.Nil(t, k)
}

func TestRequest(t *testing.T) {
	q := "example one two"
	req, err := NewRequest(q)
	if assert.Nil(t, err) {
		assert.Equal(t, req.Query, q)
		assert.Equal(t, req.Key, "example")
		assert.Equal(t, req.Value, "one two")
	}
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

	s := Server{}
	err := s.Init("./test/keyfu.conf")
	if assert.Nil(t, err) {
		assert.Equal(t, s.Config.Listen, ":8000")
		assert.Equal(t, s.Config.URL, "http://localhost:8000")
	}

	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("PORT", "1234")

	s = Server{}
	err = s.Init("./test/keyfu.conf")
	if assert.Nil(t, err) {
		assert.Equal(t, s.Config.Listen, "0.0.0.0:1234")
		assert.Equal(t, s.Config.URL, "http://0.0.0.0:1234")
	}

	ts := httptest.NewServer(http.HandlerFunc(s.RunHandler))
	defer ts.Close()

	assertRunLocation(t, ts, "hello+world", "https://encrypted.google.com/search?q=hello+world")
	assertRunLocation(t, ts, "gh", "https://github.com/")
	assertRunLocation(t, ts, "github", "https://github.com/")
	assertRunLocation(t, ts, "gh+code", "https://github.com/search?q=code")
	assertRunLocation(t, ts, "github+code", "https://github.com/search?q=code")
}

func TestOpenSearch(t *testing.T) {
	s := Server{}
	err := s.Init("./test/keyfu.conf")
	if !assert.Nil(t, err) {
		return
	}

	ts := httptest.NewServer(http.HandlerFunc(s.OpenSearchHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/opensearch.xml")
	if assert.Nil(t, err) && assert.NotNil(t, res) {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if assert.Nil(t, err) {
			assert.Equal(t, strings.Index(string(b), "http://0.0.0.0:1234"), 289)
		}
	}
}
