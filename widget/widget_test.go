package widget

import (
	"bytes"
	"testing"
)

func TestInsertBytes(t *testing.T) {
	for _, d := range []struct {
		s      []byte
		data   []byte
		offset int
		result []byte
	}{
		{
			[]byte("Hello world"),
			[]byte("my "),
			6,
			[]byte("Hello my world"),
		},
		{
			[]byte("こんにちは世界"),
			[]byte("私の"),
			15,
			[]byte("こんにちは私の世界"),
		},
		{
			[]byte("こんにちは△□の世界"),
			[]byte("○✕"),
			15,
			[]byte("こんにちは○✕△□の世界"),
		},
	} {
		s := InsertBytes(d.s, d.data, d.offset)
		if !bytes.Equal(s, d.result) {
			t.Errorf("InsertBytes(%q, %q, %q)=%q, want %q", d.s, d.data, d.offset, s, d.result)
		}
	}
}

func TestDeleteBytes(t *testing.T) {
	for _, d := range []struct {
		s      []byte
		offset int
		length int
		result []byte
	}{
		{
			[]byte("Hello my world"),
			6,
			3,
			[]byte("Hello world"),
		},
		{
			[]byte("こんにちは私の世界"),
			15,
			6,
			[]byte("こんにちは世界"),
		},
		{
			[]byte("こんにちは○✕△□の世界"),
			15,
			15,
			[]byte("こんにちは世界"),
		},
	} {
		s := DeleteBytes(d.s, d.offset, d.length)
		if !bytes.Equal(s, d.result) {
			t.Errorf("DeleteBytes(%q, %q, %q)=%q, want %q", d.s, d.offset, d.length, s, d.result)
		}
	}
}
