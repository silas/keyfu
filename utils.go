package main

import (
	"strings"
	"unicode"
)

func parseQ(value string) (string, string) {
	if value == "" {
		return "", ""
	}

	begin := -1
	end := len(value)
	for i, r := range value {
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
		return "", ""
	}

	return value[begin:end], strings.TrimLeftFunc(value[end:], unicode.IsSpace)
}
