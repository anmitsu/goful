package utils

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
