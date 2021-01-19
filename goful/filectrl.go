package goful

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progress"
	"github.com/anmitsu/goful/widget"
)

func (g *Goful) rename(src, dst string) {
	if _, err := os.Lstat(dst); err != nil {
		if !os.IsNotExist(err) {
			message.Error(err)
			return
		}
	} else {
		message := fmt.Sprintf("Overwrite? (%s -> %s)", src, dst)
		switch g.dialog(message, "yes", "no") {
		case "yes":
			break
		default:
			return
		}
	}
	if err := os.Rename(src, dst); err != nil {
		message.Error(err)
	} else {
		message.Infof("Renamed %s -> %s", src, dst)
	}
}

func (g *Goful) bulkRename(pattern, repl string, files ...*filer.FileStat) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		message.Error(err)
		return
	}
	count := 0
	newnames := make([]string, len(files))
	for i, file := range files {
		newname := re.ReplaceAllString(file.Name(), repl)
		if newname == file.Name() {
			newnames[i] = ""
		} else {
			newnames[i] = newname
			file.SetDisplay(file.Name() + " -> " + newname)
			count++
		}
	}
	if count == 0 {
		message.Errorf("No matches found for %s", pattern)
		return
	}

	switch g.dialog(fmt.Sprintf("Rename(%d)? origin -> result", count), "yes", "no") {
	case "yes":
		renames := make([]string, 0, count)
		for i, file := range files {
			if newnames[i] == "" {
				continue
			}
			if err := os.Rename(file.Name(), newnames[i]); err != nil {
				message.Error(err)
				// error handling confirm
			} else {
				renames = append(renames, file.Name())
			}
			file.ResetDisplay()
		}
		message.Infof(`Renamed(%d) "%s" to "%s" for %s`, count, pattern, repl, renames)
		g.Workspace().ReloadAll()
	default:
		for _, file := range files {
			file.ResetDisplay()
		}
	}
}

func (g *Goful) chmod(mode os.FileMode, names ...string) {
	for _, name := range names {
		if err := os.Chmod(name, mode); err != nil {
			message.Error(err)
			return
		}
	}
	message.Infof("Changed mode %s", names)
}

func (g *Goful) touch(name string, mode os.FileMode) {
	file, err := os.OpenFile(name, os.O_CREATE, mode)
	if err != nil {
		message.Error(err)
		return
	}
	if err := file.Close(); err != nil {
		message.Error(err)
	}
	message.Infof("Touched file %s", name)
}

func (g *Goful) remove(files ...string) {
	filesAbs := make([]string, len(files))
	for i := 0; i < len(files); i++ {
		filesAbs[i], _ = filepath.Abs(files[i])
	}
	go func() {
		defer g.syncCallback(func() { g.Workspace().ReloadAll() })

		if err := removeFiles(filesAbs...); err != nil {
			message.Error(err)
		} else {
			message.Infof("Removed %s", files)
		}
	}()
}

type overWrite int

const (
	overwriteYes overWrite = iota
	overwriteNo
	overwriteYesAll
	overwriteNoAll
	overwriteCancel
)

func (g *Goful) confirm(message string) overWrite {
	switch g.dialog(message, "yes", "no", "!", ".") {
	case "yes":
		return overwriteYes
	case "no":
		return overwriteNo
	case "!":
		return overwriteYesAll
	case ".":
		return overwriteNoAll
	default:
		return overwriteCancel
	}
}

func (g *Goful) copy(dst string, src ...string) {
	srcAbs := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		srcAbs[i], _ = filepath.Abs(src[i])
	}
	dstAbs, _ := filepath.Abs(dst)

	walker := g.newWalker(overwriteNo, overwriteNo, copyJob{})
	go func() {
		g.task <- 1
		defer g.syncCallback(func() {
			g.Workspace().ReloadAll()
			<-g.task
		})
		g.goWalk(walker, dstAbs, srcAbs...)
		message.Infof("Copied to %s form %s", dst, src)
	}()
}

func (g *Goful) move(dst string, src ...string) {
	srcAbs := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		srcAbs[i], _ = filepath.Abs(src[i])
	}
	dstAbs, _ := filepath.Abs(dst)

	walker := g.newWalker(overwriteNo, overwriteNo, moveJob{})
	go func() {
		g.task <- 1
		defer g.syncCallback(func() {
			g.Workspace().ReloadAll()
			<-g.task
		})
		g.goWalk(walker, dstAbs, srcAbs...)
		message.Infof("Moved to %s form %s", dst, src)
	}()
}

func (g *Goful) goWalk(walker *walker, dst string, src ...string) {
	g.ResizeRelative(0, 0, 0, -2)

	size, count := calcSizeCount(src...)
	progress.Start(size)
	progress.StartTaskCount(count)
	for _, s := range src {
		if err := walker.walk(s, dst); err != nil {
			message.Error(err)
		}
	}
	progress.Finish()

	g.ResizeRelative(0, 0, 0, 2)
	widget.Flush()
}

func calcSizeCount(src ...string) (float64, int) {
	size := int64(0)
	count := 0
	for _, s := range src {
		filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil
			}
			size += info.Size()
			count++
			return nil
		})
	}
	return float64(size), count
}

type walker struct {
	*Goful
	fileConfirmed overWrite
	dirConfirmed  overWrite
	callback      fileJob
}

func (g *Goful) newWalker(fileConfirmed, dirConfirmed overWrite, f fileJob) *walker {
	return &walker{g, fileConfirmed, dirConfirmed, f}
}

func (w *walker) walk(src, dst string) error {
	srcstat, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if srcstat.IsDir() {
		if strings.HasPrefix(dst, src) {
			return fmt.Errorf("Cannot copy/move directory %s into itself %s", src, dst)
		}
		if err := w.dir2dir(src, dst); err != nil {
			return err
		}
	} else {
		dst = filepath.Join(dst, filepath.Base(src))
		if err := w.file2file(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func (w *walker) file2file(src, dst string) error {
	if _, err := os.Lstat(dst); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		switch w.fileConfirmed {
		case overwriteNoAll:
			return nil
		case overwriteYesAll:
			break
		default:
			w.fileConfirmed = w.confirm(fmt.Sprintf("Overwrite? %s -> %s", filepath.Base(src), dst))
			switch w.fileConfirmed {
			case overwriteNo, overwriteNoAll:
				return nil
			case overwriteCancel:
				return fmt.Errorf("canceled file operation")
			}
		}
	}

	if err := w.callback.job(src, dst); err != nil {
		return err
	}
	return nil
}

func (w *walker) dir2dir(src, dst string) error {
	srcdir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcdir.Close()

	dirstat, err := srcdir.Stat()
	if err != nil {
		return err
	}

	dstdir := filepath.Join(dst, filepath.Base(src))
	if _, err := os.Lstat(dstdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dstdir, dirstat.Mode()); err != nil {
				return err
			}
		} else {
			return err
		}
	} else { // dstdir is already exists
		switch w.dirConfirmed {
		case overwriteNoAll:
			return nil
		case overwriteYesAll:
			break
		default:
			w.dirConfirmed = w.confirm(fmt.Sprintf("Merge directory? %s -> %s", filepath.Base(src), dstdir))
			switch w.dirConfirmed {
			case overwriteNo, overwriteNoAll:
				return nil
			case overwriteCancel:
				return fmt.Errorf("canceled file operation")
			}
		}
	}

	for {
		fi, err := srcdir.Readdir(100)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		for _, f := range fi {
			if f.IsDir() {
				srcdir := filepath.Join(src, f.Name())
				if err := w.dir2dir(srcdir, dstdir); err != nil {
					return err
				}
			} else {
				srcfile := filepath.Join(src, f.Name())
				dstfile := filepath.Join(dstdir, f.Name())
				if err := w.file2file(srcfile, dstfile); err != nil {
					return err
				}
			}
		}
	}

	if err := w.callback.afterVisitDir(src, dirstat, dstdir); err != nil {
		return err
	}
	return nil
}

type fileJob interface {
	job(src, dst string) error
	afterVisitDir(src string, srcstat os.FileInfo, dst string) error
}

type (
	copyJob struct{}
	moveJob struct{}
)

func (job copyJob) job(src, dst string) error {
	if err := copyFile(src, dst); err != nil {
		return err
	}
	return nil
}

func (job copyJob) afterVisitDir(src string, srcstat os.FileInfo, dst string) error {
	if err := copyTimes(srcstat, dst); err != nil {
		return err
	}
	return nil
}

func (job moveJob) job(src, dst string) error {
	if err := moveFile(src, dst); err != nil {
		return err
	}
	return nil
}

func (job moveJob) afterVisitDir(src string, srcstat os.FileInfo, dst string) error {
	if err := copyTimes(srcstat, dst); err != nil {
		return err
	}
	if err := removeEmptyDir(src); err != nil {
		return err
	}
	return nil
}

func removeFiles(files ...string) error {
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error { // not make directories in this function
	// copy symlink
	if lstat, err := os.Lstat(src); err != nil {
		return err
	} else if lstat.Mode()&os.ModeSymlink != 0 {
		if err := copySymlink(src, dst); err != nil {
			return err
		}
		return nil
	}

	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	srcstat, err := srcfile.Stat()
	if err != nil {
		return err
	}
	dstfile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcstat.Mode().Perm())
	if err != nil {
		return err
	}
	defer dstfile.Close()

	if err := letCopy(srcfile, dstfile); err != nil {
		return err
	}
	if err := copyTimes(srcstat, dst); err != nil {
		return err
	}
	return nil
}

func copySymlink(src, dst string) error {
	linksrc, err := os.Readlink(src) // not eval link path
	if err != nil {
		return err
	}
	if err := os.Symlink(linksrc, dst); err != nil {
		return err
	}
	return nil
}

func copyTimes(srcstat os.FileInfo, dst string) error {
	mtime := srcstat.ModTime()
	atime := srcstat.ModTime() // atime := time.Unix(srcstat.Sys().(*syscall.Stat_t).Atim.Unix())
	if err := os.Chtimes(dst, atime, mtime); err != nil {
		return err
	}
	return nil
}

func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		if err.(*os.LinkError).Err == syscall.EXDEV { // cross-device link
			if err := copyFileAfterRemove(src, dst); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func copyFileAfterRemove(src, dst string) error {
	if err := copyFile(src, dst); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}
	return nil
}

func removeEmptyDir(src string) error {
	srcdir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcdir.Close()

	remain, err := srcdir.Readdir(1)
	if err != nil && err != io.EOF {
		return err
	}
	if len(remain) < 1 {
		if err := os.Remove(src); err != nil {
			return err
		}
	}
	return nil
}

func letCopy(srcfile, dstfile *os.File) error {
	srcstat, err := srcfile.Stat()
	if err != nil {
		return err
	}

	quit := make(chan bool)
	go func() { // drawing progress
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				progress.Draw()
				widget.Flush()
			case <-quit:
				return
			}
		}
	}()

	progress.StartTask(srcstat)
	buf := make([]byte, 1024*32)
	for {
		n, err := srcfile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := dstfile.Write(buf[:n]); err != nil {
			return err
		}
		progress.Update(float64(n))
	}
	close(quit)
	progress.FinishTask()
	return nil
}
