package goful

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

// match shell separators, macros, options and spaces.
var re = regexp.MustCompile(`([;|>&])|(%(?:[&mMfFxX]|[dD]2?))|([[:space:]]-[[:word:]-=]+)|[[:space:]]`)

// SpawnMode starts the spawn mode.
func (g *Goful) SpawnMode(cmd string) cmdline.Mode {
	commands, err := utils.SearchCommands()
	if err != nil {
		message.Error(err)
	}
	mode := &spawnMode{g, commands}
	c := cmdline.New(mode, g)
	c.SetText(cmd)
	g.next = c
	return mode
}

type spawnMode struct {
	*Goful
	commands map[string]bool
}

func (m *spawnMode) String() string          { return "shell" }
func (m *spawnMode) Prompt() string          { return "$ " }
func (m *spawnMode) Result() string          { return "" }
func (m *spawnMode) Init(c *cmdline.Cmdline) {}

func (m *spawnMode) Draw(c *cmdline.Cmdline) {
	c.Clear()
	x, y := c.LeftTop()
	x++
	x = widget.SetCells(x, y, m.Prompt(), look.Prompt())
	termbox.SetCursor(x+c.Cursor(), y)
	m.drawCommand(x, y, c.String())
}

func (m *spawnMode) drawCommand(x, y int, cmd string) {
	start := 0
	// match is index [start, end, sep_start, sep_end, macro_start, macro_end, opt_start, opt_end]
	for _, match := range re.FindAllStringSubmatchIndex(cmd, -1) {
		s := cmd[start:match[0]]
		if _, ok := m.commands[s]; ok { // as command
			x = widget.SetCells(x, y, s, look.CmdlineCommand())
		} else {
			x = widget.SetCells(x, y, s, look.Cmdline())
		}
		start = match[0]
		s = cmd[start:match[1]]
		if match[2] != -1 { // as shell separator ;|>&
			x = widget.SetCells(x, y, s, look.Cmdline())
		} else if match[4] != -1 { // as macro %& %m %M %f %F %x %X %d2 %D %d2 %D2
			x = widget.SetCells(x, y, s, look.CmdlineMacro())
		} else if match[6] != -1 { // as option -a --bcd-efg
			x = widget.SetCells(x, y, s, look.CmdlineOption())
		} else {
			x = widget.SetCells(x, y, s, look.Cmdline())
		}
		start = match[1]
	}
	// draw the rest
	s := cmd[start:]
	if _, ok := m.commands[s]; ok { // as command
		x = widget.SetCells(x, y, s, look.CmdlineCommand())
	} else {
		x = widget.SetCells(x, y, s, look.Cmdline())
	}
}

func (m *spawnMode) Run(c *cmdline.Cmdline) {
	m.commands = nil
	m.Spawn(c.String())
	c.Exit()
}

// ShellMode starts the shell mode.
func (g *Goful) ShellMode(cmd string) cmdline.Mode {
	commands, err := utils.SearchCommands()
	if err != nil {
		message.Error(err)
	}
	mode := &shellMode{&spawnMode{g, commands}}
	c := cmdline.New(mode, g)
	c.SetText(cmd)
	g.next = c
	return mode
}

type shellMode struct {
	*spawnMode
}

func (m *shellMode) Prompt() string { return "Shell: " }
func (m *shellMode) Draw(c *cmdline.Cmdline) {
	c.Clear()
	x, y := c.LeftTop()
	x++
	x = widget.SetCells(x, y, m.Prompt(), look.Prompt())
	termbox.SetCursor(x+c.Cursor(), y)
	m.drawCommand(x, y, c.String())
}

func (m *shellMode) Run(c *cmdline.Cmdline) {
	m.commands = nil
	m.SpawnShell(c.String())
	c.Exit()
}

// DialogMode starts dialog.
func (g *Goful) DialogMode(message string, options ...string) cmdline.Mode {
	mode := &dialogMode{message, options, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type dialogMode struct {
	message string
	options []string
	result  string
}

func (m *dialogMode) String() string { return "dialog" }
func (m *dialogMode) Prompt() string {
	return fmt.Sprintf("%s [%s]: ", m.message, strings.Join(m.options, "/"))
}
func (m *dialogMode) Result() string          { return m.result }
func (m *dialogMode) Init(c *cmdline.Cmdline) {}
func (m *dialogMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *dialogMode) Run(c *cmdline.Cmdline) {
	for _, opt := range m.options {
		if c.String() == opt {
			m.result = opt
			c.Exit()
			return
		}
	}
	c.SetText("")
}

// QuitMode starts quit cmdline mode.
func (g *Goful) QuitMode() cmdline.Mode {
	mode := &quitMode{g}
	g.next = cmdline.New(mode, g)
	return mode
}

type quitMode struct {
	*Goful
}

func (m quitMode) String() string          { return "quit" }
func (m quitMode) Prompt() string          { return "Quit? [yes/no]: " }
func (m quitMode) Result() string          { return "" }
func (m quitMode) Init(c *cmdline.Cmdline) {}
func (m quitMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m quitMode) Run(c *cmdline.Cmdline) {
	switch c.String() {
	case "yes":
		c.Exit()
		m.quit()
	case "no":
		c.Exit()
	default:
		c.SetText("")
	}
}

// CopyMode starts copy.
func (g *Goful) CopyMode() cmdline.Mode {
	mode := &copyMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type copyMode struct {
	*Goful
	src string
}

func (m *copyMode) String() string { return "copy" }
func (m *copyMode) Prompt() string {
	if m.Dir().IsMark() {
		return fmt.Sprintf("Copy %d mark files to: ", m.Dir().MarkCount())
	} else if m.src != "" {
		return fmt.Sprintf("Copy from %s to: ", m.src)
	} else {
		return "Copy from: "
	}
}
func (m *copyMode) Result() string { return "" }
func (m *copyMode) Init(c *cmdline.Cmdline) {
	if m.Dir().IsMark() {
		c.SetText(m.Workspace().NextDir().Path)
	} else {
		c.SetText(m.File().Name())
	}
}
func (m *copyMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *copyMode) Run(c *cmdline.Cmdline) {
	if m.Dir().IsMark() {
		dst := c.String()
		src := m.Dir().MarkfilePaths()
		m.copy(dst, src...)
		c.Exit()
	} else if m.src != "" {
		dst := c.String()
		m.copy(dst, m.src)
		c.Exit()
	} else {
		m.src = c.String()
		c.SetText(m.Workspace().NextDir().Path)
	}
}

// MoveMode starts move.
func (g *Goful) MoveMode() cmdline.Mode {
	mode := &moveMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type moveMode struct {
	*Goful
	src string
}

func (m *moveMode) String() string { return "move" }
func (m *moveMode) Prompt() string {
	if m.Dir().IsMark() {
		return fmt.Sprintf("Move %d mark files to: ", m.Dir().MarkCount())
	} else if m.src != "" {
		return fmt.Sprintf("Move from %s to: ", m.src)
	} else {
		return "Move from: "
	}
}
func (m *moveMode) Result() string { return "" }
func (m *moveMode) Init(c *cmdline.Cmdline) {
	if m.Dir().IsMark() {
		c.SetText(m.Workspace().NextDir().Path)
	} else {
		c.SetText(m.File().Name())
	}
}
func (m *moveMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *moveMode) Run(c *cmdline.Cmdline) {
	if m.Dir().IsMark() {
		dst := c.String()
		src := m.Dir().MarkfilePaths()
		m.move(dst, src...)
		c.Exit()
	} else if m.src != "" {
		dst := c.String()
		m.move(dst, m.src)
		c.Exit()
	} else {
		m.src = c.String()
		c.SetText(m.Workspace().NextDir().Path)
	}
}

// RenameMode starts rename.
func (g *Goful) RenameMode() cmdline.Mode {
	mode := &renameMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type renameMode struct {
	*Goful
	src string
}

func (m *renameMode) String() string { return "rename" }
func (m *renameMode) Prompt() string {
	return fmt.Sprintf("Rename: %s -> ", m.src)
}
func (m *renameMode) Result() string { return "rename" }
func (m *renameMode) Init(c *cmdline.Cmdline) {
	m.src = m.File().Name()
	c.SetText(m.src)
	c.MoveCursor(-len(filepath.Ext(m.src)))
}
func (m *renameMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *renameMode) Run(c *cmdline.Cmdline) {
	dst := c.String()
	if dst == "" {
		return
	}
	m.rename(m.src, dst)
	m.Workspace().ReloadAll()
	c.Exit()
}

// BulkRenameMode starts regexp rename.
func (g *Goful) BulkRenameMode() cmdline.Mode {
	mode := &bulkRenameMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type bulkRenameMode struct {
	*Goful
	src string
}

func (m *bulkRenameMode) String() string { return "renameregexp" }
func (m *bulkRenameMode) Prompt() string {
	return "Rename by regexp: %s/"
}
func (m *bulkRenameMode) Result() string          { return "" }
func (m *bulkRenameMode) Init(c *cmdline.Cmdline) {}
func (m *bulkRenameMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *bulkRenameMode) Run(c *cmdline.Cmdline) {
	var pattern, repl string
	patterns := strings.Split(c.String(), "/")
	if len(patterns) > 1 {
		pattern = patterns[0]
		repl = patterns[1]
	} else {
		message.Errorf("Input must be like `regexp/replaced'")
		return
	}
	c.Exit()
	m.renameRegexp(pattern, repl, m.Dir().Markfiles()...)
}

// RemoveMode starts remove.
func (g *Goful) RemoveMode() cmdline.Mode {
	mode := &removeMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type removeMode struct {
	*Goful
	src string
}

func (m *removeMode) String() string { return "remove" }
func (m *removeMode) Result() string { return "" }
func (m *removeMode) Prompt() string {
	if m.Dir().IsMark() {
		return fmt.Sprintf("Remove %d mark files? [yes/no]: ", m.Dir().MarkCount())
	} else if m.src != "" {
		return fmt.Sprintf("Remove? %s [yes/no]: ", m.src)
	} else {
		return "Remove: "
	}
}
func (m *removeMode) Init(c *cmdline.Cmdline) {
	if !m.Dir().IsMark() {
		c.SetText(m.File().Name())
	}
}
func (m *removeMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *removeMode) Run(c *cmdline.Cmdline) {
	if marked := m.Dir().IsMark(); marked || m.src != "" {
		switch c.String() {
		case "yes":
			if marked {
				m.remove(m.Dir().MarkfilePaths()...)
			} else {
				m.remove(m.src)
			}
			c.Exit()
		case "no":
			c.Exit()
		default:
			c.SetText("")
		}
	} else {
		m.src = c.String()
		c.SetText("")
	}
}

// MkdirMode starts make directory.
func (g *Goful) MkdirMode() cmdline.Mode {
	mode := &mkdirMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type mkdirMode struct {
	*Goful
	path string
}

func (m *mkdirMode) String() string { return "mkdir" }
func (m *mkdirMode) Result() string { return "" }
func (m *mkdirMode) Prompt() string {
	if m.path != "" {
		return "Mode (default 0755): "
	}
	return "Make directory: "
}
func (m *mkdirMode) Init(c *cmdline.Cmdline) {}
func (m *mkdirMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *mkdirMode) Run(c *cmdline.Cmdline) {
	if m.path != "" {
		mode := c.String()
		if mode != "" {
			if mode, err := strconv.ParseUint(mode, 8, 32); err != nil {
				message.Error(err)
			} else if err := os.MkdirAll(m.path, os.FileMode(mode)); err != nil {
				message.Error(err)
			}
		} else {
			if err := os.MkdirAll(m.path, 0755); err != nil {
				message.Error(err)
			}
		}
		message.Info("Made directory " + m.path)
		c.Exit()
		m.Workspace().ReloadAll()
	} else {
		m.path = c.String()
		c.SetText("")
	}
}

// CreatefileMode starts file creation.
func (g *Goful) CreatefileMode() cmdline.Mode {
	mode := &createFileMode{g, ""}
	g.next = cmdline.New(mode, g)
	return mode
}

type createFileMode struct {
	*Goful
	path string
}

func (m *createFileMode) String() string { return "createfile" }
func (m *createFileMode) Result() string { return "" }
func (m *createFileMode) Prompt() string {
	if m.path != "" {
		return "Mode (default 0664): "
	}
	return "New file: "
}
func (m *createFileMode) Init(c *cmdline.Cmdline) {}
func (m *createFileMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *createFileMode) Run(c *cmdline.Cmdline) {
	if m.path != "" {
		mode := c.String()
		if mode != "" {
			if mode, err := strconv.ParseUint(mode, 8, 32); err != nil {
				message.Error(err)
			} else {
				m.touch(m.path, os.FileMode(mode))
			}
		} else {
			m.touch(m.path, 0644)
		}
		c.Exit()
		m.Workspace().ReloadAll()
	} else {
		m.path = c.String()
		c.SetText("")
	}
}

// ChmodMode starts change mode.
func (g *Goful) ChmodMode() cmdline.Mode {
	mode := &chmodMode{g, nil}
	g.next = cmdline.New(mode, g)
	return mode
}

type chmodMode struct {
	*Goful
	fi os.FileInfo
}

func (m *chmodMode) String() string { return "chmod" }
func (m *chmodMode) Result() string { return "" }
func (m *chmodMode) Prompt() string {
	if m.Dir().IsMark() {
		return fmt.Sprintf("Chmod %d mark files to: ", m.Dir().MarkCount())
	} else if m.fi != nil {
		return fmt.Sprintf("Chmod %s: %o to ", m.fi.Name(), m.fi.Mode())
	}
	return "Chmod: "
}
func (m *chmodMode) Init(c *cmdline.Cmdline) {
	if !m.Dir().IsMark() {
		c.SetText(m.File().Name())
	}
}
func (m *chmodMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *chmodMode) Run(c *cmdline.Cmdline) {
	if m.Dir().IsMark() || m.fi != nil {
		mode, err := strconv.ParseUint(c.String(), 8, 32)
		if err != nil {
			message.Error(err)
			c.Exit()
			return
		}
		if m.fi != nil {
			m.chmod(os.FileMode(mode), m.fi.Name())
		} else {
			files := m.Dir().MarkfilePaths()
			m.chmod(os.FileMode(mode), files...)
		}
		c.Exit()
		m.Workspace().ReloadAll()
	} else {
		file := c.String()
		lstat, err := os.Lstat(file)
		if err != nil {
			message.Error(err)
			c.Exit()
			return
		}
		m.fi = lstat
		c.SetText("")
	}
}

// ChangeWorkspaceTitle starts changing workspace title.
func (g *Goful) ChangeWorkspaceTitle() cmdline.Mode {
	mode := &changeWorkspaceTitle{g}
	g.next = cmdline.New(mode, g)
	return mode
}

type changeWorkspaceTitle struct {
	*Goful
}

func (m *changeWorkspaceTitle) String() string          { return "changeworkspacetitle" }
func (m *changeWorkspaceTitle) Result() string          { return "" }
func (m *changeWorkspaceTitle) Prompt() string          { return "Change workspace title: " }
func (m *changeWorkspaceTitle) Init(c *cmdline.Cmdline) {}
func (m *changeWorkspaceTitle) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *changeWorkspaceTitle) Run(c *cmdline.Cmdline) {
	title := c.String()
	if title != "" {
		m.Workspace().SetTitle(title)
	}
	c.Exit()
}

// ChdirMode starts change directory.
func (g *Goful) ChdirMode() cmdline.Mode {
	mode := &chdirMode{g}
	g.next = cmdline.New(mode, g)
	return mode
}

type chdirMode struct {
	*Goful
}

func (m *chdirMode) String() string          { return "chdir" }
func (m *chdirMode) Result() string          { return "" }
func (m *chdirMode) Prompt() string          { return "Chdir to: " }
func (m *chdirMode) Init(c *cmdline.Cmdline) {}
func (m *chdirMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *chdirMode) Run(c *cmdline.Cmdline) {
	if path := c.String(); path != "" {
		m.Dir().Chdir(path)
		c.Exit()
	}
}

// GlobMode starts glob.
func (g *Goful) GlobMode() cmdline.Mode {
	mode := &globMode{g}
	g.next = cmdline.New(mode, g)
	return mode
}

type globMode struct {
	*Goful
}

func (m *globMode) String() string { return "glob" }
func (m *globMode) Result() string { return "" }
func (m *globMode) Prompt() string {
	return "Glob pattern: "
}
func (m *globMode) Init(c *cmdline.Cmdline) {}
func (m *globMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *globMode) Run(c *cmdline.Cmdline) {
	if pattern := c.String(); pattern != "" {
		m.Dir().Glob(pattern)
		c.Exit()
	}
}

// GlobdirMode starts globdir.
func (g *Goful) GlobdirMode() cmdline.Mode {
	mode := &globdirMode{g}
	g.next = cmdline.New(mode, g)
	return mode
}

type globdirMode struct {
	*Goful
}

func (m *globdirMode) String() string { return "globdir" }
func (m *globdirMode) Result() string { return "" }
func (m *globdirMode) Prompt() string {
	return "Globdir pattern: "
}
func (m *globdirMode) Init(c *cmdline.Cmdline) {}
func (m *globdirMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }
func (m *globdirMode) Run(c *cmdline.Cmdline) {
	if pattern := c.String(); pattern != "" {
		m.Dir().Globdir(pattern)
		c.Exit()
	}
}
