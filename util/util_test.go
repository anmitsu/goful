package util

import (
	"os"
	"path"
	"testing"
)

func TestAbbrPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dirs := []struct {
		orig string
		abbr string
	}{
		{
			path.Join(home, "abc/def"),
			"~/abc/def",
		},
		{
			path.Join("abc", home, "def"),
			path.Join("abc", home, "def"),
		},
	}

	for _, d := range dirs {
		if s := AbbrPath(d.orig); d.abbr != s {
			t.Errorf("AbbrPath(%q)=%q, want %q", d.orig, s, d.abbr)
		}
		if s := ExpandPath(d.abbr); d.orig != s {
			t.Errorf("ExpandPath(%q)=%q, want %q", d.abbr, s, d.orig)
		}
	}
}

func TestShortenPath(t *testing.T) {
	for _, d := range []struct {
		path   string
		width  int
		result string
	}{
		{"/", 0, "/"},
		{"/home/", 2, "/h"},
		{"/home///", 2, "/h"},
		{"/home///", 10, "/home///"},
		{"/home/abc/def/hij", 14, "/h/abc/def/hij"},
		{"/home/abc/def/hij", 12, "/h/a/def/hij"},
		{"/home/abc/def/hij", 10, "/h/a/d/hij"},
		{"/home/abc/def/hij", 9, "/h/a/d/hij"},
		{"/home/abc/def/hij", 1, "/h/a/d/hij"},
		{"/home/あいう/かきく/さしす", 23, "/h/あいう/かきく/さしす"},
		{"/home/あいう/かきく/さしす", 19, "/h/あ/かきく/さしす"},
		{"/home/あいう/かきく/さしす", 15, "/h/あ/か/さしす"},
		{"/home/あいう/かきく/さしす", 14, "/h/あ/か/さしす"},
		{"/home/あいう/かきく/さしす", 1, "/h/あ/か/さしす"},
	} {
		if s := ShortenPath(d.path, d.width); s != d.result {
			t.Errorf("ShortenPath(%q, %d)=%q, want %q", d.path, d.width, s, d.result)
		}
	}
}

func TestRemoveExt(t *testing.T) {
	for _, d := range []struct {
		path   string
		result string
	}{
		{"abc", "abc"},
		{"abc.def", "abc"},
		{".abc", ".abc"},
		{".abc.def", ".abc"},
	} {
		if s := RemoveExt(d.path); s != d.result {
			t.Errorf("PathExt(%q)=%q, want %q", d.path, s, d.result)
		}
	}
}
