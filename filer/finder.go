package filer

import (
	"regexp"
	"strings"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Finder represents textbox for filtering files in the directory.
type Finder struct {
	*widget.TextBox
	dir        *Directory
	names      []string
	startname  string
	historyPos int
}

var finderHistroy = make([]string, 0, 100)

var finderKeymap func(*Finder) widget.Keymap

// ConfigFinder sets the finder keymap function.
func ConfigFinder(config func(*Finder) widget.Keymap) {
	finderKeymap = config
}

// NewFinder returns a new finder to position the directory bottom.
func NewFinder(dir *Directory, x, y, width, height int) *Finder {
	names := make([]string, len(dir.List())-1)
	for i := 0; i < len(dir.List())-1; i++ {
		names[i] = dir.List()[i+1].Name() // not contain ".."
	}

	finder := &Finder{
		TextBox:   widget.NewTextBox(x, y, width, height),
		dir:       dir,
		names:     names,
		startname: dir.CurrentContent().Name(),
	}

	if len(finderHistroy) < 1 {
		finderHistroy = append(finderHistroy, "")
	}
	finder.Edithook = func() { dir.read() }
	return finder
}

// MoveHistory moves histories with specified amounts and sets to textbox.
func (f *Finder) MoveHistory(amount int) {
	f.historyPos += amount
	if f.historyPos > len(finderHistroy)-1 {
		f.historyPos = len(finderHistroy) - 1
	} else if f.historyPos < 0 {
		f.historyPos = 0
	}
	var text string
	if f.historyPos != 0 {
		text = finderHistroy[len(finderHistroy)-f.historyPos]
	} else {
		text = finderHistroy[0]
	}
	f.SetText(text)
	f.Edithook()
}

func (f *Finder) addHistory() {
	text := f.String()
	if text != "" {
		for i, hist := range finderHistroy {
			if hist == text {
				finderHistroy = append(finderHistroy[:i], finderHistroy[i+1:]...)
				break
			}
		}
		if i := len(finderHistroy); i != cap(finderHistroy) {
			finderHistroy = append(finderHistroy, text)
		} else {
			finderHistroy = append(finderHistroy[2:i], text)
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
		message.Error(err)
		return
	}

	current := f.dir.CurrentContent().Name()
	f.dir.ClearList()
	f.dir.AppendList(NewFileStat(f.dir.Path, ".."))
	for _, name := range f.names {
		if re.MatchString(name) {
			callback(name)
		}
	}
	f.dir.SetCursorByName(current)
	f.dir.SetOffsetCenteredCursor()
}

// Draw the finder.
func (f *Finder) Draw() {
	f.Clear()
	x, y := f.LeftTop()
	const prompt = "Finder: "
	x = widget.SetCells(x, y, prompt, look.Prompt())
	x = widget.SetCells(x, y, f.String(), look.Finder())
	spacewidth := f.Width() - len(prompt) - runewidth.StringWidth(f.String())
	if spacewidth > 0 {
		widget.SetCells(x, y, strings.Repeat(" ", spacewidth), look.Finder())
	}
	termbox.SetCursor(x, y)
}

// Exit the finder and reload the directory to clear filtering.
func (f *Finder) Exit() {
	f.exitNotRead()
	name := f.startname
	if len(f.dir.List()) > 1 {
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
	termbox.HideCursor()
}
