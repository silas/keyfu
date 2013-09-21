// Copyright (c) 2009 Michael Hoisie
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//
// Source: https://github.com/hoisie/web.go

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getCookieSig(key string, val []byte, timestamp string) string {
	hm := hmac.New(sha1.New, []byte(key))

	hm.Write(val)
	hm.Write([]byte(timestamp))

	hex := fmt.Sprintf("%02x", hm.Sum(nil))
	return hex
}

func (ctx *Context) GetSecureCookie(name string) *http.Cookie {
	cookie := ctx.GetCookie(name)

	if cookie == nil {
		return cookie
	}

	parts := strings.SplitN(cookie.Value, "|", 3)

	if len(parts) < 3 {
		return nil
	}

	val, timestamp, sig := parts[0], parts[1], parts[2]

	if getCookieSig(Config.SessionCookieSecret, []byte(val), timestamp) != sig {
		return nil
	}

	ts, _ := strconv.ParseInt(timestamp, 10, 64)

	if (time.Now().UTC().Unix() - 31*86400) > ts {
		return nil
	}

	buf := bytes.NewBufferString(val)
	encoder := base64.NewDecoder(base64.StdEncoding, buf)

	res, _ := ioutil.ReadAll(encoder)
	cookie.Value = string(res)

	return cookie
}

func (ctx *Context) SetSecureCookie(cookie *http.Cookie) {
	//base64 encode the val
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(cookie.Value))
	encoder.Close()
	vs := buf.String()
	vb := buf.Bytes()
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	sig := getCookieSig(Config.SessionCookieSecret, vb, timestamp)
	cookie.Value = strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.SetCookie(cookie)
}
