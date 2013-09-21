package main

import (
	"testing"
)

func TestKeywordBuiltin(t *testing.T) {
	k := Keyword{"builtin", "com.google"}

	if url, err := k.Builtin("", ""); url != "http://www.google.com/" || err != nil {
		t.Errorf("Builtin direct link failed")
	}

	if url, err := k.Builtin("", "test"); url != "http://www.google.com/search?q=test" || err != nil {
		t.Errorf("Builtin query link failed")
	}
}

func TestKeywordLink(t *testing.T) {
	k := Keyword{"link", "http://www.google.com/\nhttp://www.google.com/search?q=%s"}

	if url, err := k.Link(""); url != "http://www.google.com/" || err != nil {
		t.Errorf("Link direct")
	}
	if url, err := k.Link("test"); url != "http://www.google.com/search?q=test" || err != nil {
		t.Errorf("Link query")
	}

	k.Body = "http://www.google.com/"
	if url, err := k.Link(""); url != "http://www.google.com/" || err != nil {
		t.Errorf("Link direct (1 direct link)")
	}
	if url, err := k.Link("test"); url != "http://www.google.com/" || err != nil {
		t.Errorf("Link query (1 direct link)")
	}

	k.Body = "http://www.google.com/search?q=%s"
	if url, err := k.Link(""); url != "http://www.google.com/search?q=" || err != nil {
		t.Errorf("Link direct (1 query link)")
	}
	if url, err := k.Link("test"); url != "http://www.google.com/search?q=test" || err != nil {
		t.Errorf("Link query (1 query link)")
	}
}
