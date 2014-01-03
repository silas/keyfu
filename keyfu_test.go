package main

import (
	"testing"

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
