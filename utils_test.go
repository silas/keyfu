package main

import (
	"testing"
)

func TestParseQ(t *testing.T) {
	if key, value := parseQ(""); key != "" || value != "" {
		t.Errorf("No key or value")
	}

	if key, value := parseQ(" "); key != "" || value != "" {
		t.Errorf("No key or value (space prefix)")
	}

	if key, value := parseQ(" key"); key != "key" || value != "" {
		t.Errorf("Key and no value (space prefix)")
	}

	if key, value := parseQ(" key value"); key != "key" || value != "value" {
		t.Errorf("Key and no value (space prefix)")
	}

	if key, value := parseQ("key"); key != "key" || value != "" {
		t.Errorf("Key with no value")
	}

	if key, value := parseQ("  key  value1  value2  value3  "); key != "key" || value != "value1  value2  value3  " {
		t.Errorf("Key and value")
	}
}

func TestMd5Sum(t *testing.T) {
	if md5 := md5sum("/dev/null"); md5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("Invalid md5sum")
	}
}

func TestVersionStatic(t *testing.T) {
	if versioned := versionStatic("/dev", "/static/null"); versioned != "/static/null?v=d41d8c" {
		t.Errorf("Invalid versioned file")
	}
}

func TestMinifyHtml(t *testing.T) {
	html := "  <html>\n\t<head>\n\t\t<title> title </title>\n\t</head>  <body>\ttest  </body></html>   "
	validHtml := "<html><head><title> title </title></head><body>\ttest  </body></html>"
	if result := minifyHtml(html); result != validHtml {
		t.Errorf("HTML didn't minify correctly")
	}
}
