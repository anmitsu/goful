// +build windows

package info

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/util"
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

	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	var free, all int64
	_, _, _ = c.Call(uintptr(
		unsafe.Pointer(syscall.StringToUTF16Ptr("."))),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&all)),
		uintptr(unsafe.Pointer(nil)))
	used := float64(all-free) / float64(all) * 100
	freeSI := util.FormatSize(free)

	info := fmt.Sprintf("%s free %.1f%% used %s %d %s %s", freeSI, used, perm, size, mtime, name)
	s := runewidth.Truncate(info, w.Width(), "~")
	widget.SetCells(x, y, s, look.Default())
}
