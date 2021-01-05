// Package filer is layout and viewing for files.
package filer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

// Filer is file manger with workspaces to layout directorires to list files
type Filer struct {
	*widget.Window
	keymap     widget.Keymap
	extmap     widget.Extmap
	Workspaces []*Workspace `json:"workspaces"`
	Cursor     int          `json:"cursor"`
}

// New creates a new filer based on specified size and coordinates.
// Creates five workspaces and default path is home directory.
func New(x, y, width, height int) *Filer {
	home, err := os.UserHomeDir()
	if err != nil {
		message.Error(err)
		home = "/"
	}

	workspaces := make([]*Workspace, 5)
	for i := 0; i < 5; i++ {
		title := fmt.Sprintf("%d", i+1)
		ws := NewWorkspace(x, y+1, width, height-1, title)
		ws.Dirs = make([]*Directory, 2)
		for j := 0; j < 2; j++ {
			ws.Dirs[j] = NewDirectory(0, 0, 0, 0)
			ws.Dirs[j].Path = home
			ws.Dirs[j].SetTitle(utils.AbbrPath(home))
		}
		ws.allocate()
		workspaces[i] = ws
	}
	return &Filer{
		Window:     widget.NewWindow(x, y, width, height),
		keymap:     widget.Keymap{},
		extmap:     widget.Extmap{},
		Workspaces: workspaces,
		Cursor:     0,
	}
}

// NewFromJSON creates a new filer form the state json file.
func NewFromJSON(name string, x, y, width, height int) *Filer {
	file, err := os.Open(utils.ExpandPath(name))
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
		Cursor:     0,
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

// SaveJSON saves the filer state to the file.
func (f *Filer) SaveJSON(name string) error {
	jsondata, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(utils.ExpandPath(name))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsondata); err != nil {
		return err
	}
	return nil
}

// CreateWorkspace creates and adds workspace to the end.
func (f *Filer) CreateWorkspace() {
	title := fmt.Sprintf("%d", len(f.Workspaces)+1)
	x, y := f.LeftTop()
	width, height := f.Width(), f.Height()
	ws := NewWorkspace(x, y+1, width, height-1, title)
	ws.CreateDir()
	ws.CreateDir()
	f.Workspaces = append(f.Workspaces, ws)
}

// CloseWorkspace closes workspace on the cursor.
func (f *Filer) CloseWorkspace() {
	if len(f.Workspaces) < 2 {
		return
	}
	i := f.Cursor
	f.Workspaces[i].visible(false)
	f.Workspaces[i] = nil
	f.Workspaces = append(f.Workspaces[:i], f.Workspaces[i+1:]...)
	if f.Cursor > len(f.Workspaces)-1 {
		f.Cursor = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// MoveWorkspace moves to other workspace.
func (f *Filer) MoveWorkspace(amount int) {
	f.Workspace().visible(false)
	f.Cursor += amount
	if f.Cursor >= len(f.Workspaces) {
		f.Cursor = 0
	} else if f.Cursor < 0 {
		f.Cursor = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// Workspace is getter of current workspace.
func (f *Filer) Workspace() *Workspace {
	return f.Workspaces[f.Cursor]
}

// Dir is getter of focused directory on current workspace.
func (f *Filer) Dir() *Directory {
	return f.Workspace().Dir()
}

// File is getter of file on the cursor in focused directory on current workspace.
func (f *Filer) File() *FileStat {
	return f.Dir().File()
}

// AddKeymap adds to keymap for filer.
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

// MergeKeymap merges to keymap for filer.
func (f *Filer) MergeKeymap(m widget.Keymap) {
	for key, callback := range m {
		f.keymap[key] = callback
	}
}

// AddExtmap adds to extmap for filer.
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

// MergeExtmap merges to extmap for filer.
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
		} else if utf8.RuneCountInString(key) == 1 {
			r, _ := utf8.DecodeRuneInString(key)
			finder.InsertChar(r)
			return
		}
	}

	if sub, ok := f.extmap[key]; ok {
		if callback, ok := sub[".link"]; ok && f.File().IsLink() {
			callback()
		} else if callback, ok := sub[".dir"]; ok && f.File().IsDir() {
			callback()
		} else if callback, ok := sub[".exec"]; ok && f.File().IsExec() {
			callback()
		} else if callback, ok := sub[f.File().Ext()]; ok {
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
		if f.Cursor != i {
			x = widget.SetCells(x, y, s, look.Default())
		} else {
			x = widget.SetCells(x, y, s, look.Selected())
		}
	}
	x = widget.SetCells(x, y, " | ", look.Default())

	ws := f.Workspace()
	width := (f.Width() - x) / len(ws.Dirs)
	for i := 0; i < len(ws.Dirs); i++ {
		s := fmt.Sprintf("[%d] %s", i+1, ws.Dirs[i].Title())
		s = runewidth.Truncate(s, width, "...")
		s = runewidth.FillRight(s, width)
		if ws.Cursor != i {
			x = widget.SetCells(x, y, s, look.Default())
		} else {
			x = widget.SetCells(x, y, s, look.Selected())
		}
	}
}

// Draw current workspace.
func (f *Filer) Draw() {
	f.Clear()
	f.drawHeader()
	f.Workspace().Draw()
}

// Resize implements widget.Widget.
func (f *Filer) Resize(x, y, width, height int) {
	f.Window.Resize(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.Resize(x, y+1, width, height-1)
	}
}

// ResizeRelative implements widget.Widget.
func (f *Filer) ResizeRelative(x, y, width, height int) {
	f.Window.ResizeRelative(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.ResizeRelative(x, y, width, height)
	}
}
