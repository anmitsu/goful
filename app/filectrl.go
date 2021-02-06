package app

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
	"github.com/anmitsu/goful/util"
	"github.com/anmitsu/goful/widget"
)

func (g *Goful) rename(src, dst string) {
	if _, err := os.Lstat(dst); err != nil {
		if !os.IsNotExist(err) {
			message.Error(err)
			return
		}
	} else {
		message := fmt.Sprintf("Overwrite? %s", dst)
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

func (g *Goful) copy(dst string, src ...string) {
	g.walk(func(dst string, src ...string) {
		walker := g.newWalker(overwriteNo, overwriteNo, copyJob{})
		if err := g.letWalk(walker, dst, src...); err != nil {
			message.Error(err)
		} else {
			message.Infof("Copied to %s form %s", dst, src)
		}
	}, dst, src...)
}

func (g *Goful) move(dst string, src ...string) {
	g.walk(func(dst string, src ...string) {
		walker := g.newWalker(overwriteNo, overwriteNo, moveJob{})
		if err := g.letWalk(walker, dst, src...); err != nil {
			message.Error(err)
		} else {
			message.Infof("Moved to %s form %s", dst, src)
		}
	}, dst, src...)
}

func (g *Goful) walk(walkFn func(dst string, src ...string), dst string, src ...string) {
	srcAbs := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		srcAbs[i], _ = filepath.Abs(src[i])
	}
	dstAbs, _ := filepath.Abs(dst)

	go func() {
		g.task <- 1
	g.ResizeRelative(0, 0, 0, -2)
	if w := g.Next(); w != nil {
		w.ResizeRelative(0, -2, 0, 0)
	}
		defer g.syncCallback(func() {
			g.ResizeRelative(0, 0, 0, 2)
			if w := g.Next(); w != nil {
				w.ResizeRelative(0, 2, 0, 0) // for cmdline and menu
			}
			widget.Show()
			g.Workspace().ReloadAll()
			<-g.task
		})
		walkFn(dstAbs, srcAbs...)
	}()
}

func (g *Goful) letWalk(walker *walker, dst string, src ...string) error {
	size, count := util.CalcSizeCount(src...)
	progress.Start(float64(size))
	progress.StartTaskCount(count)
	var err error
	for _, s := range src {
		if e := walker.walk(s, dst); e != nil {
			err = e
			break
		}
	}
	progress.Finish()
	return err
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
	if dststat, err := os.Stat(dst); err != nil {
		if !os.IsNotExist(err) { // ignore error if not exist dst and create dst
			return err
		}
	} else {
		if dststat.IsDir() { // make to in exist dst directory
			dst = filepath.Join(dst, filepath.Base(src))
		}
	}
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
		if err := w.file2file(src, dst); err != nil {
			return err
		}
	}
	return nil
}

type overWrite int

const (
	overwriteYes overWrite = iota
	overwriteNo
	overwriteYesAll
	overwriteNoAll
	overwriteCancel
)

func (w *walker) confirm(message string) overWrite {
	switch w.dialog(message, "yes", "no", "!", ".") {
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
			w.fileConfirmed = w.confirm(fmt.Sprintf("Overwrite? exists %s", dst))
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
	if _, err := os.Stat(dst); err != nil {
		if os.IsNotExist(err) { // make dst directory if dst not exists
			if err := copyDir(src, dst); err != nil {
				return err
			}
		} else {
			return err
		}
	} else { // dst is already exists
		switch w.dirConfirmed {
		case overwriteNoAll:
			return nil
		case overwriteYesAll:
			break
		default:
			w.dirConfirmed = w.confirm(fmt.Sprintf("Merge? exists %s", dst))
			switch w.dirConfirmed {
			case overwriteNo, overwriteNoAll:
				return nil
			case overwriteCancel:
				return fmt.Errorf("canceled file operation")
			}
		}
	}

	srcdir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcdir.Close()
	for {
		fi, err := srcdir.Readdir(100)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		for _, f := range fi {
			src := filepath.Join(src, f.Name())
			dst := filepath.Join(dst, f.Name())
			if f.IsDir() {
				if err := w.dir2dir(src, dst); err != nil {
					return err
				}
			} else {
				if err := w.file2file(src, dst); err != nil {
					return err
				}
			}
		}
	}

	if err := w.callback.afterVisitDir(src, dst); err != nil {
		return err
	}
	return nil
}

type fileJob interface {
	job(src, dst string) error
	afterVisitDir(src, dst string) error
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

func (job copyJob) afterVisitDir(src, dst string) error {
	if err := copyTimes(src, dst); err != nil {
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

func (job moveJob) afterVisitDir(src, dst string) error {
	if err := copyTimes(src, dst); err != nil {
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

func copyDir(src, dst string) error {
	srcstat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Mkdir(dst, srcstat.Mode()); err != nil {
		return err
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
	if err := copyTimes(src, dst); err != nil {
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

func copyTimes(src, dst string) error {
	srcstat, err := os.Stat(src)
	if err != nil {
		return err
	}
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
	quit := make(chan bool)
	go func() { // drawing progress
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				progress.Draw()
				widget.Show()
			case <-quit:
				return
			}
		}
	}()

	srcstat, err := srcfile.Stat()
	if err != nil {
		return err
	}
	progress.StartTask(srcstat)
	buf := make([]byte, 4096)
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
