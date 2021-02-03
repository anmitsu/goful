package filer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/util"
	"github.com/anmitsu/goful/widget"
)

// Directory is a list box to store file stats.
type Directory struct {
	*widget.ListBox
	reader  reader
	history map[string]string // key: path, value: file name on cursor
	finder  *Finder
	Path    string   `json:"path"`
	Sort    sortType `json:"sort_kind"`
}

// NewDirectory creates a new directory based on specified size and coordinates.
func NewDirectory(x, y, width, height int) *Directory {
	path, _ := filepath.Abs(".")
	listbox := widget.NewListBox(x, y, width, height, path)
	listbox.SetBorderStyle(borderStyle)
	return &Directory{
		ListBox: listbox,
		reader:  defaultReader("."),
		history: map[string]string{},
		Path:    path,
		Sort:    sortName,
	}
}

// default border style
var borderStyle widget.BorderStyle = widget.ULBorder

// SetBorderStyle sets a directory default border style.
func SetBorderStyle(style widget.BorderStyle) {
	borderStyle = style
}

type sortType string

const (
	sortName     sortType = "Name[^]"
	sortNameRev           = "Name[$]"
	sortSize              = "Size[^]"
	sortSizeRev           = "Size[$]"
	sortMtime             = "Time[^]"
	sortMtimeRev          = "Time[$]"
	sortExt               = "Ext[^]"
	sortExtRev            = "Ext[$]"
)

var priorityDir = true

// TogglePriority toggles the priority for sorting files.
// The directory is prioritized in sorting if this is true.
func TogglePriority() {
	priorityDir = !priorityDir
}

var showHiddens = true

// ToggleShowHiddens toggles the showing of hidden files.
func ToggleShowHiddens() {
	showHiddens = !showHiddens
}

type reader interface {
	Read(callback func(name string))
	String() string
}

type defaultReader string

func (s defaultReader) String() string { return "" }
func (s defaultReader) Read(callback func(string)) {
	fd, err := os.Open(string(s))
	if err != nil {
		message.Error(err)
		return
	}
	defer fd.Close()

	for {
		names, err := fd.Readdirnames(100)
		for _, name := range names {
			if !showHiddens && strings.HasPrefix(name, ".") {
				continue
			}
			callback(name)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			message.Error(err)
			return
		}
	}
}

type globPattern string

func (s globPattern) String() string {
	return fmt.Sprintf("Glob:(%s)", string(s))
}

func (s globPattern) Read(callback func(name string)) {
	matches, err := filepath.Glob(string(s))
	if err != nil {
		message.Error(err)
		return
	}
	for _, name := range matches {
		if !showHiddens && strings.HasPrefix(name, ".") {
			continue
		}
		callback(name)
	}
}

type globDirPattern string

func (s globDirPattern) String() string {
	return fmt.Sprintf("Globdir:(%s)", string(s))
}

func (s globDirPattern) Read(callback func(string)) {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "." {
			return nil
		}
		if ok, _ := filepath.Match(string(s), info.Name()); ok {
			if !showHiddens {
				if strings.HasPrefix(path, ".") || strings.HasPrefix(info.Name(), ".") {
					return nil
				}
			}
			callback(path)
		}
		return nil
	})
}

func (d *Directory) init4json() {
	d.ListBox = widget.NewListBox(0, 0, 0, 0, "")
	d.history = map[string]string{}
	d.SetTitle(util.AbbrPath(d.Path))
	d.SetColumn(1)
	d.reader = defaultReader(".")
}

// Resize the window and the finder.
func (d *Directory) Resize(x, y, width, height int) {
	d.ListBox.Resize(x, y, width, height)
	if d.finder != nil {
		d.finder.Resize(x, y+d.Height()-1, d.Width(), 1)
		d.ResizeRelative(0, 0, 0, -1)
	}
}

// Finder starts a finder in the directory for filtering files.
func (d *Directory) Finder() {
	x, y := d.LeftTop()
	d.finder = NewFinder(d, x, y+d.Height()-1, d.Width(), 1)
	d.ResizeRelative(0, 0, 0, -1)
}

// EnterDir changes the directory to a path on the cursor.
func (d *Directory) EnterDir() {
	d.Chdir(d.File().Name())
}

// Reset marking or reader.
func (d *Directory) Reset() {
	if d.IsMark() {
		d.MarkClear()
	} else if _, ok := d.reader.(defaultReader); !ok {
		d.reader = defaultReader(".")
		d.read()
	}
}

// Chdir changes the current directory and reads a new path by the default reader.
// Sets the cursor to the history name or to the previous directory name if parent destinats.
func (d *Directory) Chdir(path string) {
	path = util.ExpandPath(path)
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(filepath.Join(d.Path, path))
	}
	olddir := filepath.Base(d.Path)
	parent := filepath.Dir(d.Path)

	if d.finder != nil {
		d.finder.exitNotRead()
	}

	if err := os.Chdir(path); err != nil {
		message.Error(err)
		return
	}
	if !d.IsEmpty() {
		d.history[d.Path] = d.File().Name()
	}
	d.SetTitle(util.AbbrPath(path))
	d.Path = path
	d.reader = defaultReader(".")
	d.read()

	if name, ok := d.history[d.Path]; ok {
		d.SetCursorByName(name)
		d.SetOffsetCenteredCursor()
	} else if path == parent {
		d.SetCursorByName(olddir)
		d.SetOffsetCenteredCursor()
	} else {
		d.SetCursor(0)
	}
}

// Glob sets a reader to matching pattern in the current directory.
func (d *Directory) Glob(pattern string) {
	d.reader = globPattern(pattern)
	d.read()
}

// Globdir sets a reader to matching pattern in the directory includeing sub directories.
func (d *Directory) Globdir(pattern string) {
	d.reader = globDirPattern(pattern)
	d.read()
}

func (d *Directory) read() {
	marked := make(map[string]bool, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			marked[e.(*FileStat).Path()] = true
		}
	}

	callback := func(name string) {
		if fs := NewFileStat(d.Path, name); fs != nil {
			d.AppendList(fs)
		}
	}
	if d.finder != nil {
		d.finder.find(callback)
	} else {
		d.ClearList()
		d.reader.Read(callback)
	}
	if d.IsEmpty() {
		d.AppendList(NewFileStat(d.Path, ".."))
	}
	sort.Sort(d)

	for _, e := range d.List() {
		if _, ok := marked[e.(*FileStat).Path()]; ok {
			e.(*FileStat).Mark()
		}
	}
}

func (d *Directory) reload() {
	if err := os.Chdir(d.Path); err != nil {
		message.Error(err)
		home, _ := os.UserHomeDir()
		d.Chdir(home)
		return
	}
	d.read()
}

// File returns a file on the cursor.
func (d *Directory) File() *FileStat {
	return d.CurrentContent().(*FileStat)
}

// Base returns the directory name.
func (d *Directory) Base() string { return filepath.Base(d.Path) }

func (d *Directory) sortBy(typ sortType) {
	d.Sort = typ
	name := d.File().Name()
	sort.Sort(d)
	d.SetCursorByName(name)
}

// SortName sorts files in ascending order by the file name.
func (d *Directory) SortName() { d.sortBy(sortName) }

// SortNameDec sorts files in descending order by the file name.
func (d *Directory) SortNameDec() { d.sortBy(sortNameRev) }

// SortMtime sorts files in ascending order by the modified time.
func (d *Directory) SortMtime() { d.sortBy(sortMtime) }

// SortMtimeDec sorts files in descending order by the modified time.
func (d *Directory) SortMtimeDec() { d.sortBy(sortMtimeRev) }

// SortSize sorts files in ascending order by the file size.
func (d *Directory) SortSize() { d.sortBy(sortSize) }

// SortSizeDec sorts files in descending order by the file size.
func (d *Directory) SortSizeDec() { d.sortBy(sortSizeRev) }

// SortExt sorts files in ascending order by the file extension.
func (d *Directory) SortExt() { d.sortBy(sortExt) }

// SortExtDec sorts files in descending order by the file extension.
func (d *Directory) SortExtDec() { d.sortBy(sortExtRev) }

// Less compares based on Sort.
func (d *Directory) Less(i, j int) bool {
	if priorityDir {
		id := d.List()[i].(*FileStat).stat.IsDir()
		jd := d.List()[j].(*FileStat).stat.IsDir()
		if !(id && jd) && (id || jd) {
			return id
		}
	}
	switch d.Sort {
	case sortName:
		return d.List()[i].Name() < d.List()[j].Name()
	case sortNameRev:
		return d.List()[i].Name() > d.List()[j].Name()
	case sortMtime:
		return d.lessMtime(i, j)
	case sortMtimeRev:
		return d.lessMtime(j, i)
	case sortSize:
		return d.lessSize(i, j)
	case sortSizeRev:
		return d.lessSize(j, i)
	case sortExt:
		return d.lessExt(i, j)
	case sortExtRev:
		return d.lessExt(j, i)
	}
	return d.List()[i].Name() < d.List()[j].Name()
}

func (d *Directory) lessMtime(i, j int) bool {
	f1 := d.List()[i].(*FileStat)
	f2 := d.List()[j].(*FileStat)
	t1 := f1.ModTime().Unix()
	t2 := f2.ModTime().Unix()
	if t1 != t2 {
		return t1 < t2
	}
	return f1.Name() < f2.Name()
}

func (d *Directory) lessSize(i, j int) bool {
	f1 := d.List()[i].(*FileStat)
	f2 := d.List()[j].(*FileStat)
	s1 := f1.Size()
	s2 := f2.Size()
	if s1 != s2 {
		return s1 < s2
	}
	return f1.Name() < f2.Name()
}

func (d *Directory) lessExt(i, j int) bool {
	f1 := d.List()[i].(*FileStat)
	f2 := d.List()[j].(*FileStat)
	e1 := f1.Ext()
	e2 := f2.Ext()
	if e1 != e2 {
		return e1 < e2
	}
	return f1.Name() < f2.Name()
}

// IsMark reports whether even one file marked.
func (d *Directory) IsMark() bool {
	return d.MarkCount() != 0
}

// ToggleMark toggles the file mark on the cursor.
func (d *Directory) ToggleMark() {
	fs := d.CurrentContent().(*FileStat)
	if fs.Name() == ".." {
		d.MoveCursor(1)
	} else {
		fs.ToggleMark()
		d.MoveCursor(1)
	}
}

// InvertMark toggles all file marks.
func (d *Directory) InvertMark() {
	for _, e := range d.List() {
		if e.(*FileStat).Name() != ".." {
			e.(*FileStat).ToggleMark()
		}
	}
}

// MarkClear clears all file marks.
func (d *Directory) MarkClear() {
	for _, e := range d.List() {
		e.(*FileStat).Markoff()
	}
}

// MarkCount returns a number of marked files.
func (d *Directory) MarkCount() int {
	c := 0
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			c++
		}
	}
	return c
}

// Markfiles returns marked file lists.
func (d *Directory) Markfiles() []*FileStat {
	if d.MarkCount() < 1 {
		return []*FileStat{d.File()}
	}
	markfiles := make([]*FileStat, 0, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			markfiles = append(markfiles, e.(*FileStat))
		}
	}
	return markfiles
}

// MarkfileNames returns marked file names.
func (d *Directory) MarkfileNames() []string {
	if d.MarkCount() < 1 {
		return []string{d.File().Name()}
	}
	markfiles := make([]string, 0, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			markfiles = append(markfiles, e.(*FileStat).Name())
		}
	}
	return markfiles
}

// MarkfilePaths returns marked file paths.
func (d *Directory) MarkfilePaths() []string {
	if d.MarkCount() < 1 {
		return []string{d.File().Path()}
	}
	markfiles := make([]string, 0, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			markfiles = append(markfiles, e.(*FileStat).Path())
		}
	}
	return markfiles
}

// MarkfileQuotedNames returns quoted file names for marked.
func (d *Directory) MarkfileQuotedNames() []string {
	if d.MarkCount() < 1 {
		return []string{util.Quote(d.File().Name())}
	}
	markfiles := make([]string, 0, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			markfiles = append(markfiles, util.Quote(e.(*FileStat).Name()))
		}
	}
	return markfiles
}

// MarkfileQuotedPaths returns quoted file paths for marked.
func (d *Directory) MarkfileQuotedPaths() []string {
	if d.MarkCount() < 1 {
		return []string{util.Quote(d.File().Path())}
	}
	markfiles := make([]string, 0, d.MarkCount())
	for _, e := range d.List() {
		if e.(*FileStat).IsMarked() {
			markfiles = append(markfiles, util.Quote(e.(*FileStat).Path()))
		}
	}
	return markfiles
}

func (d *Directory) drawFooter() {
	count := len(d.List()) - 1
	s := fmt.Sprintf("[%d/%d] %s(%d) %s %s",
		d.MarkCount(), count, d.ScrollRate(), d.Cursor(), d.Sort, d.reader.String())
	x, y := d.LeftBottom()
	widget.SetCells(x, y, s, look.Default())
}

func (d *Directory) drawFiles(focus bool) {
	height := d.Height() - 2
	row := 1
	shift := 0
	width := d.Width() - 1
	if d.BorderStyle() == widget.AllBorder {
		shift++
		width--
	}
	for i := d.Offset(); i < d.Upper(); i++ {
		if row > height {
			break
		}
		x, y := d.LeftTop()
		y += row
		x += shift
		if focus && i == d.Cursor() {
			d.List()[i].Draw(x, y, width, true)
		} else {
			d.List()[i].Draw(x, y, width, false)
		}
		row++
	}
}

func (d *Directory) draw(focus bool) {
	d.AdjustCursor()
	d.AdjustOffset()
	d.Border()
	d.drawFiles(focus)
	d.drawFooter()
	if d.finder != nil {
		d.finder.Draw()
	}
}
