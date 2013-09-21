package main

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func parseQ(value string) (string, string) {
	if value == "" {
		return "", ""
	}

	begin := -1
	end := len(value)
	for i, rune := range value {
		if !unicode.IsSpace(rune) {
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

// Author: Russ Cox
// Source: http://groups.google.com/group/golang-nuts/msg/5ebbdd72e2d40c09
func uuid4() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func md5sum(path string) string {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return ""
	}
	h := md5.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func versionStatic(root string, path string) string {
	// assumes path: /static/bla.jpg
	if sum := md5sum(root + path[7:]); sum != "" {
		return path + "?v=" + sum[0:6]
	}
	return path
}

// this probably isn't unicode safe
func minifyHtml(html string) string {
	if html == "" {
		return html
	}
	src := []byte(html)
	l := len(src)
	dst := make([]byte, l)

	srcLast := 0
	dstLast := 0
	dropping := src[0] == ' ' || src[0] == '\n' || src[0] == '\t'

	for i := 0; i < l; i++ {
		rune := src[i]
		if rune == '>' || (rune == '<' && !dropping) {
			dstEnd := dstLast + i + 1 - srcLast
			copy(dst[dstLast:dstEnd], src[srcLast:i+1])
			dstLast, srcLast = dstEnd, i+1
			dropping = true
		} else if rune == '<' {
			srcLast = i
			dropping = false
		} else if rune != ' ' && rune != '\n' && rune != '\t' {
			dropping = false
		}
	}

	return string(dst[:dstLast])
}
