package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertRequestParse(t *testing.T, q, key, value string) {
	r, err := NewRequest(q)

	if assert.Nil(t, err) {
		assert.Equal(t, r.Key, key)
		assert.Equal(t, r.Value, value)
	}
}

func assertNotRequestParse(t *testing.T, q string) {
	r, err := NewRequest(q)

	if assert.NotNil(t, err) {
		assert.Equal(t, err, errParse)
		assert.Nil(t, r)
	}
}

func TestRequestParse(t *testing.T) {
	assertNotRequestParse(t, "")
	assertRequestParse(t, "one", "one", "")
	assertRequestParse(t, " one ", "one", "")
	assertRequestParse(t, "one two", "one", "two")
	assertRequestParse(t, "one two three", "one", "two three")
	assertRequestParse(t, "one  two  three ", "one", "two  three ")
}
