package filer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

var statView = fileStatView{true, true, true}

type fileStatView struct {
	size       bool
	permission bool
	time       bool
}

// SetStatView sets file state view.
func SetStatView(size, permission, time bool) { statView = fileStatView{size, permission, time} }

// ToggleSizeView toggles file size view.
func ToggleSizeView() { statView.size = !statView.size }

// TogglePermView toggles file permission view.
func TogglePermView() { statView.permission = !statView.permission }

// ToggleTimeView toggles file time view.
func ToggleTimeView() { statView.time = !statView.time }

// FileStat is file infomation.
type FileStat struct {
	os.FileInfo             // os.Lstat(path)
	stat        os.FileInfo // os.Stat(path)
	path        string      // full path of file
	name        string      // base name of path or ".." as upper directory
	display     string      // display name for draw
	marked      bool        // marked whether
}

// NewFileStat returns FileStat of file name in directory.
func NewFileStat(dir string, name string) *FileStat {
	path := filepath.Join(dir, name)

	lstat, err := os.Lstat(path)
	if err != nil {
		message.Error(err)
		return nil
	}
	stat, err := os.Stat(path)
	if err != nil {
		stat = lstat
	}

	var display string
	if stat.IsDir() {
		display = name
	} else {
		display = utils.RemoveExt(name)
	}

	return &FileStat{
		FileInfo: lstat,
		stat:     stat,
		path:     path,
		name:     name,
		display:  display,
		marked:   false,
	}
}

// Name implements for interface of widget.Enterable.
func (f *FileStat) Name() string {
	return f.name
}

// SetDisplay is setter of display name for draw.
func (f *FileStat) SetDisplay(name string) {
	f.display = name
}

// ResetDisplay resets display name to file name.
func (f *FileStat) ResetDisplay() {
	if f.stat.IsDir() {
		f.display = f.name
	} else {
		f.display = utils.RemoveExt(f.name)
	}
}

// Mark file.
func (f *FileStat) Mark() {
	f.marked = true
}

// Markoff of file.
func (f *FileStat) Markoff() {
	f.marked = false
}

// ToggleMark toggles mark of file.
func (f *FileStat) ToggleMark() {
	f.marked = !f.marked
}

// Path returns the full path of file.
func (f *FileStat) Path() string {
	return f.path
}

// Ext retruns the extension of file.
func (f *FileStat) Ext() string {
	if f.stat.IsDir() {
		return ""
	}
	if ext := filepath.Ext(f.Name()); ext != f.Name() {
		return ext
	}
	return ""
}

// IsLink reports whether symlink.
func (f *FileStat) IsLink() bool {
	return f.Mode()&os.ModeSymlink != 0
}

// IsExec reports whether executable file.
func (f *FileStat) IsExec() bool {
	return f.Mode().Perm()&0111 != 0
}

// IsMarked reports whether marked file.
func (f *FileStat) IsMarked() bool {
	return f.marked
}

func (f *FileStat) suffix() string {
	suffix := ""
	if f.IsLink() {
		if link, err := os.Readlink(f.Path()); err != nil {
			message.Error(err)
		} else {
			suffix += "@ -> " + link
		}
		if f.stat.IsDir() {
			suffix += "/"
		}
	} else if f.IsDir() {
		suffix = "/"
	} else if f.IsExec() {
		suffix = "*"
	}
	return suffix
}

func (f *FileStat) formatFileSize(size int64) string {
	switch {
	case size > 1024*1024*1024:
		return fmt.Sprintf("%6.1fG", float64(size)/(1024*1024*1024))
	case size > 1024*1024:
		return fmt.Sprintf("%6.1fM", float64(size)/(1024*1024))
	case size > 1024:
		return fmt.Sprintf("%6.1fk", float64(size)/1024)
	default:
		return fmt.Sprintf("%7d", size)
	}
}

func (f *FileStat) states() string {
	ret := f.Ext()
	if statView.size {
		if f.stat.IsDir() {
			ret += fmt.Sprintf("%8s", "<DIR>")
		} else {
			ret += " " + f.formatFileSize(f.stat.Size())
		}
	}
	if statView.permission {
		ret += " " + f.stat.Mode().String()
	}
	if statView.time {
		ret += " " + f.stat.ModTime().Format("06-01-02 15:04")
	}
	return ret
}

func (f *FileStat) look() look.Look {
	switch {
	case f.IsMarked():
		return look.Marked()
	case f.IsLink():
		if f.stat.IsDir() {
			return look.SymlinkDir()
		}
		return look.Symlink()
	case f.IsDir():
		return look.Directory()
	case f.IsExec():
		return look.Executable()
	default:
		return look.Default()
	}
}

// Draw implements for interface of widget.Drawer.
func (f *FileStat) Draw(x, y, width int, lk look.Look) {
	lk = lk.And(f.look())
	states := f.states()
	width -= len(states)
	pre := " "
	if f.marked {
		pre = "*"
	}
	s := pre + f.display + f.suffix()
	s = runewidth.Truncate(s, width, "...")
	s = runewidth.FillRight(s, width)
	x = widget.SetCells(x, y, s, lk)
	widget.SetCells(x, y, states, lk)
}
