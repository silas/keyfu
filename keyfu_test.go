package main

import (
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

func runProgramKeyword(t *testing.T, q, name, timeout string) (*Response, error) {
	c := map[string]string{
		"name":    name,
		"timeout": timeout,
	}

	k, err := NewProgramKeyword(c)
	if assert.Nil(t, err) {
		if req, err := NewRequest(q); assert.Nil(t, err) {
			return k.Run(req)
		}
	}

	return nil, err
}

func TestProgramKeyword(t *testing.T) {
	res, err := runProgramKeyword(t, "ok", "./test/redirect.sh", "100ms")
	if assert.Nil(t, err) {
		assert.Equal(t, res.Body, "http://www.keyfu.com/")
	}

	res, err = runProgramKeyword(t, "ok", "./test/render.sh", "100ms")
	if assert.Nil(t, err) {
		assert.Equal(t, res.Body, "hello world\n")
	}

	_, err = runProgramKeyword(t, "ok", "./test/timeout.sh", "50ms")
	assert.Equal(t, err, errProgramTimeout)

	_, err = runProgramKeyword(t, "ok", "./test/exit.sh", "100ms")
	assert.NotNil(t, err)
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
