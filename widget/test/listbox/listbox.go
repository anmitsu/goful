package main

import (
	"strings"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
)

func main() {
	widget.Init()
	defer widget.Fini()

	look.Set("default")
	x, y := 0, 0
	width, height := widget.Size()
	lb := widget.NewListBox(x, y, width/2, height/2, "test")
	lb2 := widget.NewListBox(int(float64(width)*0.7), y, int(float64(width)*0.3), height/2/2, "test2")
	lb3 := widget.NewListBox(x, height/2, width/2, height/2, "test3")
	lb4 := widget.NewListBox(width/2, height/2, width/2, height/2, "test4")

	contents := []string{
		"Hello world.",
		strings.Repeat("Hello world. ", 10),
		"こんにちは世界。",
		strings.Repeat("こんにちは世界。", 10),
		"こんにちは○×ﾊﾝｶｸﾓｼﾞの世界。",
		strings.Repeat("こんにちは○×ﾊﾝｶｸﾓｼﾞの世界。", 10),
	}
	lb.AppendString(contents...)
	lb2.AppendString(contents...)
	lb3.AppendString(contents...)
	lb4.AppendString(contents...)

	highlight := [][]string{
		{"Hello world.", "world"},
		{strings.Repeat("Hello world. ", 10), "world"},
		{"こんにちは世界。", "世界"},
		{strings.Repeat("こんにちは世界。 ", 10), "世界"},
		{"こんにちは○×ﾊﾝｶｸﾓｼﾞの世界。", "○×ﾊﾝｶｸﾓｼﾞの"},
		{strings.Repeat("こんにちは○×ﾊﾝｶｸﾓｼﾞの世界。 ", 10), "○×ﾊﾝｶｸﾓｼﾞの"},
	}

	for _, s := range highlight {
		lb.AppendHighlightString(s[0], s[1])
		lb2.AppendHighlightString(s[0], s[1])
		lb3.AppendHighlightString(s[0], s[1])
		lb4.AppendHighlightString(s[0], s[1])
	}

	more := strings.Repeat("ABCあいう○×ﾊﾝｶｸﾓｼﾞ", 10)
	for _, content := range more {
		lb.AppendString(string(content))
		lb2.AppendString(string(content))
		lb3.AppendString(string(content))
		lb4.AppendString(string(content))
	}

	lb.SetCursor(5)
	lb2.SetCursor(9)
	lb3.SetCursor(40)
	lb4.SetCursor(80)

	lb.Draw()
	lb2.Draw()
	lb3.Draw()
	lb4.Draw()
	widget.Show()
	widget.PollEvent()
	widget.PollEvent()
}
