// Package look is attributes setting.
package look

import (
	"github.com/gdamore/tcell/v2"
)

// Set look for a name.
func Set(name string) {
	switch name {
	case "default":
		setDefault()
	case "midnight":
		setMidnight()
	case "black":
		setBlack()
	case "white":
		setWhite()
	default:
		setDefault()
	}
}

// Default is a default look attribute.
func Default() tcell.Style { return defaultAttr }

// MessageInfo is an information message look.
func MessageInfo() tcell.Style { return messageInfo }

// MessageErr is an error message look.
func MessageErr() tcell.Style { return messageErr }

// Prompt is a cmdline and finder prompt look.
func Prompt() tcell.Style { return prompt }

// Cmdline is a cmdline look.
func Cmdline() tcell.Style { return cmdline }

// CmdlineCommand is a highlighted command look in the cmdline.
func CmdlineCommand() tcell.Style { return cmdlineCommand }

// CmdlineMacro is a macro look in the cmdline.
func CmdlineMacro() tcell.Style { return cmdlineMacro }

// CmdlineOption is a option look of the cmdline.
func CmdlineOption() tcell.Style { return cmdlineOption }

// Highlight is a highlight text look in the list box.
func Highlight() tcell.Style { return highlight }

// Title is a title look of the list box.
func Title() tcell.Style { return title }

// Symlink is a symlink file look.
func Symlink() tcell.Style { return symlink }

// SymlinkDir is a symlink directory look.
func SymlinkDir() tcell.Style { return symlinkDir }

// Directory is a directory look.
func Directory() tcell.Style { return directory }

// Executable is an executable file look.
func Executable() tcell.Style { return executable }

// Marked is a marked file look.
func Marked() tcell.Style { return marked }

// Finder is a finder text area look.
func Finder() tcell.Style { return finder }

// Progress is a gauge look of the progress bar.
func Progress() tcell.Style { return progress }

var (
	defaultAttr    tcell.Style
	messageInfo    tcell.Style
	messageErr     tcell.Style
	prompt         tcell.Style
	cmdline        tcell.Style
	cmdlineCommand tcell.Style
	cmdlineMacro   tcell.Style
	cmdlineOption  tcell.Style
	highlight      tcell.Style
	title          tcell.Style
	symlink        tcell.Style
	symlinkDir     tcell.Style
	directory      tcell.Style
	executable     tcell.Style
	marked         tcell.Style
	finder         tcell.Style
	progress       tcell.Style
)

// reference https://jonasjacek.github.io/colors/

func setDefault() {
	d := tcell.StyleDefault
	defaultAttr = d
	messageInfo = d.Foreground(tcell.ColorLime).Bold(true)
	messageErr = d.Foreground(tcell.ColorRed).Bold(true)
	prompt = d.Foreground(tcell.ColorAqua).Bold(true)
	cmdline = d
	cmdlineCommand = d.Foreground(tcell.ColorLime).Bold(true)
	cmdlineMacro = d.Foreground(tcell.ColorFuchsia)
	cmdlineOption = d.Foreground(tcell.ColorYellow)
	highlight = d.Bold(true)
	title = d
	symlink = d.Foreground(tcell.ColorFuchsia)
	symlinkDir = d.Foreground(tcell.ColorFuchsia).Bold(true)
	directory = d.Foreground(tcell.ColorAqua).Bold(true)
	executable = d.Foreground(tcell.ColorRed).Bold(true)
	marked = d.Foreground(tcell.ColorYellow).Bold(true)
	finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
	progress = d.Background(tcell.ColorNavy)
}

func setMidnight() {
	d := tcell.StyleDefault
	bg := tcell.ColorNavy
	defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
	messageInfo = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
	messageErr = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	prompt = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
	cmdline = defaultAttr
	cmdlineCommand = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
	cmdlineMacro = d.Foreground(tcell.ColorFuchsia).Background(bg)
	cmdlineOption = d.Foreground(tcell.ColorYellow).Background(bg)
	highlight = defaultAttr.Bold(true)
	title = defaultAttr
	symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
	symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
	directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
	executable = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
	finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
	progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorAqua)
}

func setBlack() {
	d := tcell.StyleDefault
	bg := tcell.ColorBlack
	defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
	messageInfo = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
	messageErr = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	prompt = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
	cmdline = defaultAttr
	cmdlineCommand = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
	cmdlineMacro = d.Foreground(tcell.ColorFuchsia).Background(bg)
	cmdlineOption = d.Foreground(tcell.ColorYellow).Background(bg)
	highlight = defaultAttr.Bold(true)
	title = defaultAttr
	symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
	symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
	directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
	executable = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
	finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
	progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
}

func setWhite() {
	d := tcell.StyleDefault
	bg := tcell.ColorWhite
	defaultAttr = d.Foreground(tcell.ColorBlack).Background(bg)
	messageInfo = d.Foreground(tcell.ColorGreen).Background(bg).Bold(true)
	messageErr = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	prompt = d.Foreground(tcell.ColorNavy).Background(bg).Bold(true)
	cmdline = defaultAttr
	cmdlineCommand = d.Foreground(tcell.ColorGreen).Background(bg).Bold(true)
	cmdlineMacro = d.Foreground(tcell.ColorFuchsia).Background(bg)
	cmdlineOption = d.Foreground(tcell.ColorOlive).Background(bg)
	highlight = defaultAttr.Bold(true)
	title = defaultAttr
	symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
	symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
	directory = d.Foreground(tcell.ColorNavy).Background(bg).Bold(true)
	executable = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
	marked = d.Foreground(tcell.ColorOlive).Background(bg).Bold(true)
	finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
	progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
}
