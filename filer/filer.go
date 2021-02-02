// Package filer draws directories and files and handles inputs.
package filer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/util"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

// Filer is a file manager with workspaces to layout directorires to list files.
type Filer struct {
	*widget.Window
	keymap     widget.Keymap
	extmap     widget.Extmap
	Workspaces []*Workspace `json:"workspaces"`
	Current    int          `json:"current"`
}

// New creates a new filer based on specified size and coordinates.
// Creates five workspaces and default path is home directory.
func New(x, y, width, height int) *Filer {
	home, err := os.UserHomeDir()
	if err != nil {
		message.Error(err)
		home = "/"
	}

	workspaces := make([]*Workspace, 3)
	for i := 0; i < 3; i++ {
		title := fmt.Sprintf("%d", i+1)
		ws := NewWorkspace(x, y+1, width, height-1, title)
		ws.Dirs = make([]*Directory, 2)
		for j := 0; j < 2; j++ {
			ws.Dirs[j] = NewDirectory(0, 0, 0, 0)
			ws.Dirs[j].Path = home
			ws.Dirs[j].SetTitle(util.AbbrPath(home))
		}
		ws.allocate()
		workspaces[i] = ws
	}
	return &Filer{
		Window:     widget.NewWindow(x, y, width, height),
		keymap:     widget.Keymap{},
		extmap:     widget.Extmap{},
		Workspaces: workspaces,
		Current:    0,
	}
}

// NewFromState creates a new filer form the state json file.
func NewFromState(path string, x, y, width, height int) *Filer {
	file, err := os.Open(util.ExpandPath(path))
	if err != nil {
		return New(x, y, width, height)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return New(x, y, width, height)
	}

	filer := &Filer{
		Window:     widget.NewWindow(x, y, width, height),
		keymap:     widget.Keymap{},
		extmap:     widget.Extmap{},
		Workspaces: []*Workspace{},
		Current:    0,
	}
	if err := json.Unmarshal(data, filer); err != nil {
		return New(x, y, width, height)
	}
	if len(filer.Workspaces) < 1 {
		return New(x, y, width, height)
	}
	for _, ws := range filer.Workspaces {
		if len(ws.Dirs) < 1 {
			return New(x, y, width, height)
		}
		ws.init4json(x, y+1, width, height-1)
		for _, dir := range ws.Dirs {
			dir.init4json()
		}
		ws.allocate()
	}
	return filer
}

// SaveState saves the filer state to the file.
func (f *Filer) SaveState(path string) error {
	jsondata, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(util.ExpandPath(path))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsondata); err != nil {
		return err
	}
	return nil
}

// CreateWorkspace creates and adds a workspace to the end.
func (f *Filer) CreateWorkspace() {
	title := fmt.Sprintf("%d", len(f.Workspaces)+1)
	x, y := f.LeftTop()
	width, height := f.Width(), f.Height()
	ws := NewWorkspace(x, y+1, width, height-1, title)
	ws.CreateDir()
	ws.CreateDir()
	f.Workspaces = append(f.Workspaces, ws)
}

// CloseWorkspace closes a workspace on the current.
func (f *Filer) CloseWorkspace() {
	if len(f.Workspaces) < 2 {
		return
	}
	i := f.Current
	f.Workspaces[i].visible(false)
	f.Workspaces[i] = nil
	f.Workspaces = append(f.Workspaces[:i], f.Workspaces[i+1:]...)
	if f.Current > len(f.Workspaces)-1 {
		f.Current = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// MoveWorkspace moves to the other workspace.
func (f *Filer) MoveWorkspace(amount int) {
	f.Workspace().visible(false)
	f.Current += amount
	if f.Current >= len(f.Workspaces) {
		f.Current = 0
	} else if f.Current < 0 {
		f.Current = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// Workspace returns the current workspace.
func (f *Filer) Workspace() *Workspace {
	return f.Workspaces[f.Current]
}

// Dir returns the focused directory on the current workspace.
func (f *Filer) Dir() *Directory {
	return f.Workspace().Dir()
}

// File returns the cursor file in the focused directory on the current workspace.
func (f *Filer) File() *FileStat {
	return f.Dir().File()
}

// AddKeymap adds to the filer keymap.
func (f *Filer) AddKeymap(keys ...interface{}) {
	if len(keys)%2 != 0 {
		panic("items must be a multiple of 2")
	}

	for i := 0; i < len(keys); i += 2 {
		key := keys[i].(string)
		callback := keys[i+1].(func())
		f.keymap[key] = callback
	}
}

// MergeKeymap merges to the filer keymap.
func (f *Filer) MergeKeymap(m widget.Keymap) {
	for key, callback := range m {
		f.keymap[key] = callback
	}
}

// AddExtmap adds to the filer extmap.
func (f *Filer) AddExtmap(a ...interface{}) {
	if len(a)%3 != 0 {
		panic("items must be a multiple of 3")
	}

	for i := 0; i < len(a); i += 3 {
		key := a[i].(string)
		ext := a[i+1].(string)
		callback := a[i+2].(func())
		f.extmap[key][ext] = callback
	}
}

// MergeExtmap merges to the filer extmap.
func (f *Filer) MergeExtmap(m widget.Extmap) {
	for key, submap := range m {
		if _, found := f.extmap[key]; !found {
			f.extmap[key] = map[string]func(){}
		}
		for ext, callback := range submap {
			f.extmap[key][ext] = callback
		}
	}
}

// Input for key events.
func (f *Filer) Input(key string) {
	if finder := f.Dir().finder; finder != nil {
		if callback, ok := finderKeymap(finder)[key]; ok {
			callback()
			return
		} else if utf8.RuneCountInString(key) == 1 && key != " " {
			r, _ := utf8.DecodeRuneInString(key)
			finder.InsertChar(r)
			return
		}
	}

	if ext, ok := f.extmap[key]; ok {
		if callback, ok := ext[".dir"]; ok && (f.File().IsDir() || f.File().stat.IsDir()) {
			callback()
		} else if callback, ok := ext[".exec"]; ok && f.File().IsExec() {
			callback()
		} else if callback, ok := ext[f.File().Ext()]; ok {
			callback()
		} else if callback, ok := f.keymap[key]; ok {
			callback()
		}
	} else if callback, ok := f.keymap[key]; ok {
		callback()
	}
}

func (f *Filer) drawHeader() {
	x, y := f.LeftTop()
	for i, ws := range f.Workspaces {
		s := fmt.Sprintf(" %s ", ws.Title)
		if f.Current != i {
			x = widget.SetCells(x, y, s, look.Default())
		} else {
			x = widget.SetCells(x, y, s, look.Default().Reverse(true))
		}
	}
	x = widget.SetCells(x, y, " | ", look.Default())

	ws := f.Workspace()
	width := (f.Width() - x) / len(ws.Dirs)
	for i := 0; i < len(ws.Dirs); i++ {
		style := look.Default()
		if ws.Focus == i {
			style = style.Reverse(true)
		}
		s := fmt.Sprintf("[%d] ", i+1)
		x = widget.SetCells(x, y, s, style)
		w := width - len(s)
		s = util.ShortenPath(ws.Dirs[i].Title(), w)
		s = runewidth.Truncate(s, w, "~")
		s = runewidth.FillRight(s, w)
		x = widget.SetCells(x, y, s, style)
	}
}

// Draw the current workspace.
func (f *Filer) Draw() {
	f.Clear()
	f.drawHeader()
	f.Workspace().Draw()
}

// Resize all workspaces.
func (f *Filer) Resize(x, y, width, height int) {
	f.Window.Resize(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.Resize(x, y+1, width, height-1)
	}
}

// ResizeRelative resize relative to current sizes.
func (f *Filer) ResizeRelative(x, y, width, height int) {
	f.Window.ResizeRelative(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.ResizeRelative(x, y, width, height)
	}
}
