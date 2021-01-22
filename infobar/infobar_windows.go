// +build windows

package infobar

import (
	"fmt"
	"os"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

func (w *infoWindow) draw(fi os.FileInfo) {
	w.Clear()
	x, y := w.LeftTop()

	perm := fi.Mode().String()
	size := fi.Size()
	mtime := fi.ModTime().String()
	name := fi.Name()

	info := fmt.Sprintf("%s %d %s %s", perm, size, mtime, name)
	s := runewidth.Truncate(info, w.Width(), "...")
	widget.SetCells(x, y, s, look.Default())
}
