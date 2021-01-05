package cmdline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
)

// Completion is the list box based on cmdline text.
type Completion struct {
	*widget.ListBox
	cmdline widget.Widget
}

var completionKeymap func(*Completion) widget.Keymap

// ConfigCompletion sets keymap function for completion.
func ConfigCompletion(config func(*Completion) widget.Keymap) {
	completionKeymap = config
}

// NewCompletion creates a new completion.
func NewCompletion(x, y, width, height int, cmdline *Cmdline) *Completion {
	comp := &Completion{
		ListBox: widget.NewListBox(x, y, width, height, "Completion"),
		cmdline: cmdline,
	}

	parser := parseCmdline(cmdline)
	candidates := []string{}
	if cmdline.mode.String() == "shell" && parser.cmdname == "" {
		candidates = append(parser.compCommands(), parser.compFiles()...)
	} else {
		candidates = parser.compFiles()
	}
	for _, v := range candidates {
		comp.AppendHighlightString(v, parser.current)
	}
	comp.ColumnAdjustContentsWidth()
	return comp
}

// Next implements widget.Widget.
func (c *Completion) Next() widget.Widget { return nil }

// Disconnect implements widget.Widget.
func (c *Completion) Disconnect() {}

// InsertCompletion inserts selected completion to cmdline and exits completion.
func (c *Completion) InsertCompletion() {
	start := c.cmdline.(*Cmdline).TextBeforeCursor()
	compname := c.CurrentContent().Name()

	for i := len(compname); i >= 0; i-- {
		if strings.HasSuffix(start, compname[:i]) {
			c.cmdline.(*Cmdline).InsertString(compname[i:])
			break
		}
	}
	c.cmdline.Disconnect()
}

// Input implements widget.Widget.
func (c *Completion) Input(key string) {
	if cb, ok := completionKeymap(c)[key]; ok {
		cb()
	} else {
		c.cmdline.Disconnect()
		c.cmdline.Input(key)
	}
}

// Exit completion mode in cmdline.
func (c *Completion) Exit() { c.cmdline.Disconnect() }

type parser struct {
	cmdname string
	current string
	preword string
}

func parseCmdline(c *Cmdline) *parser {
	text := c.TextBeforeCursor()
	words, _ := shlex.Split(text, true)

	switch i := len(words); i {
	case 0:
		return &parser{"", "", ""}
	case 1:
		if isSep(text[len(text)-1]) {
			return &parser{words[0], "", ""}
		}
		return &parser{"", words[0], ""}
	default:
		if isSep(text[len(text)-1]) {
			return &parser{words[0], "", words[i-1]}
		}
		return &parser{words[0], words[i-1], words[i-2]}
	}
}

func isSep(b byte) bool {
	return b == ' ' || b == ';' || b == '|' || b == '>' || b == '&'
}

func (p *parser) compFiles() (candidates []string) {
	candidates = make([]string, 0, 100)
	dirname, file := filepath.Split(p.current)
	if dirname == "" {
		dirname = "."
	}
	dir, err := os.Open(dirname)
	if err != nil {
		return candidates
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return candidates
	}

	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, file) {
			if f.IsDir() {
				name += "/"
			}
			candidates = append(candidates, name)
		}
	}
	sort.Strings(candidates)
	return candidates
}

func (p *parser) compCommands() (candidates []string) {
	commands, _ := utils.SearchCommands()
	candidates = make([]string, 0, len(commands))
	for name := range commands {
		if strings.HasPrefix(name, p.current) {
			candidates = append(candidates, name)
		}
	}
	sort.Strings(candidates)
	return candidates
}
