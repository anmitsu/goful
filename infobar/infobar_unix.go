// +build !windows

package infobar

import (
	"fmt"
	"os"
	"os/user"
	"syscall"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

func (w *InfoBar) draw(fi os.FileInfo) {
	w.Clear()
	x, y := w.LeftTop()
	x++

	stat := fi.Sys().(*syscall.Stat_t)
	var username, group string
	if u, err := user.LookupId(fmt.Sprintf("%d", stat.Uid)); err != nil {
		username = "unknown"
	} else {
		username = u.Name
	}
	if u, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err != nil {
		group = "unknown"
	} else {
		group = u.Name
	}

	perm := fi.Mode().String()
	size := fi.Size()
	mtime := fi.ModTime().String()
	name := fi.Name()

	info := fmt.Sprintf("%s %s %s %d %d %s %s", perm, username, group, stat.Nlink, size, mtime, name)
	s := runewidth.Truncate(info, w.Width(), "...")
	widget.SetCells(x, y, s, look.Default())
}
