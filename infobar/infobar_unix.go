// +build !windows

package infobar

import (
	"fmt"
	"os"
	"os/user"
	"syscall"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

func (w *infoWindow) draw(fi os.FileInfo) {
	w.Clear()
	x, y := w.LeftTop()

	var statfs syscall.Statfs_t
	syscall.Statfs(".", &statfs)
	free := statfs.Bavail * uint64(statfs.Bsize)
	all := statfs.Blocks * uint64(statfs.Bsize)
	used := float64(all-free) / float64(all) * 100
	freeSI := utils.FormatSize(int64(free))

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

	info := fmt.Sprintf("%s free (%.1f%s used) %s %s %s %d %d %s %s",
		freeSI, used, "%", perm, username, group, stat.Nlink, size, mtime, name)
	s := runewidth.Truncate(info, w.Width(), "~")
	widget.SetCells(x, y, s, look.Default())
}
