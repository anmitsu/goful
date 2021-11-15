package filer

import (
	"regexp"
	"strings"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

// Finder represents textbox for filtering files in the directory.
type Finder struct {
	*widget.TextBox
	dir        *Directory
	names      []string
	startname  string
	historyPos int
}

var finderHistory = make([]string, 0, 100)

var finderKeymap func(*Finder) widget.Keymap

// ConfigFinder sets the finder keymap function.
func ConfigFinder(config func(*Finder) widget.Keymap) {
	finderKeymap = config
}

// NewFinder returns a new finder to position the directory bottom.
func NewFinder(dir *Directory, x, y, width, height int) *Finder {
	names := make([]string, len(dir.List()))
	for i := 0; i < len(dir.List()); i++ {
		names[i] = dir.List()[i].Name()
	}

	finder := &Finder{
		TextBox:    widget.NewTextBox(x, y, width, height),
		dir:        dir,
		names:      names,
		startname:  dir.CurrentContent().Name(),
		historyPos: 0,
	}

	if len(finderHistory) < 1 {
		finderHistory = append(finderHistory, "")
	}
	finder.Edithook = func() { dir.read() }
	return finder
}

// MoveHistory moves histories with specified amounts and sets to textbox.
func (f *Finder) MoveHistory(amount int) {
	f.historyPos += amount
	if f.historyPos > len(finderHistory)-1 {
		f.historyPos = len(finderHistory) - 1
	} else if f.historyPos < 0 {
		f.historyPos = 0
	}
	var text string
	if f.historyPos != 0 {
		text = finderHistory[len(finderHistory)-f.historyPos]
	} else {
		text = finderHistory[0]
	}
	f.SetText(text)
	f.Edithook()
}

func (f *Finder) addHistory() {
	text := f.String()
	if text != "" {
		for i, hist := range finderHistory {
			if hist == text {
				finderHistory = append(finderHistory[:i], finderHistory[i+1:]...)
				break
			}
		}
		if i := len(finderHistory); i != cap(finderHistory) {
			finderHistory = append(finderHistory, text)
		} else {
			finderHistory = append(finderHistory[2:i], text)
		}
	}
}

func (f *Finder) find(callback func(name string)) {
	expr := f.String()
	if expr == strings.ToLower(expr) {
		expr = "(?i)" + expr // case insensitive
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return
	}

	current := ""
	if !f.dir.IsEmpty() {
		current = f.dir.CurrentContent().Name()
	}
	f.dir.ClearList()
	for _, name := range f.names {
		if re.MatchString(name) {
			callback(name)
		}
	}
	if f.dir.IsEmpty() {
		f.dir.AppendList(NewFileStat(f.dir.Path, ".."))
	}
	if current != "" {
		f.dir.SetCursorByName(current)
		f.dir.SetOffsetCenteredCursor()
	}
}

// Draw the finder and show a cursor if focus is true.
func (f *Finder) Draw(focus bool) {
	f.Clear()
	x, y := f.LeftTop()
	s := "Find: " + f.String()
	x = widget.SetCells(x, y, s, look.Finder())
	spacewidth := f.Width() - runewidth.StringWidth(s)
	if spacewidth > 0 {
		widget.SetCells(x, y, strings.Repeat(" ", spacewidth), look.Finder())
	}
	if focus {
		widget.ShowCursor(x, y)
	}
}

// Exit the finder and reload the directory to clear filtering.
func (f *Finder) Exit() {
	f.exitNotRead()
	name := f.startname
	if len(f.dir.List()) > 0 {
		name = f.dir.File().Name()
	}
	f.dir.read()
	f.dir.SetCursorByName(name)
	f.dir.SetOffsetCenteredCursor()
}

func (f *Finder) exitNotRead() {
	f.dir.ResizeRelative(0, 0, 0, 1)
	f.names = nil
	f.dir.finder = nil
	f.addHistory()
	widget.HideCursor()
}
