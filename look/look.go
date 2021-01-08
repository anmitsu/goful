// Package look is attributes setting for termbox-go.
package look

import "github.com/nsf/termbox-go"

// Look is attributes of termbox-go.
type Look struct {
	fg termbox.Attribute
	bg termbox.Attribute
}

// Fg retruns a foreground attribute.
func (l Look) Fg() termbox.Attribute {
	return l.fg
}

// Bg is returns a background attribute.
func (l Look) Bg() termbox.Attribute {
	return l.bg
}

// And combines itself attributes and others.
func (l Look) And(others ...Look) Look {
	fg, bg := l.fg, l.bg
	for _, o := range others {
		fg |= o.fg
		bg |= o.bg
	}
	return Look{fg, bg}
}

// Set look for a name.
func Set(name string) {
	switch name {
	case "default":
		setDefault()
	case "midnight":
		setMidnight()
	default:
		setDefault()
	}
}

// Default is a default look attribute.
func Default() Look { return defaultAttr }

// Blank is a blank look attribute.
func Blank() Look { return Look{0, 0} }

// MessageInfo is an infomation message look.
func MessageInfo() Look { return messageInfo }

// MessageErr is an error message look.
func MessageErr() Look { return messageErr }

// Prompt is a cmdline and finder prompt look.
func Prompt() Look { return prompt }

// Cmdline is a cmdline look.
func Cmdline() Look { return cmdline }

// CmdlineCommand is a highlighted command look in the cmdline.
func CmdlineCommand() Look { return cmdlineCommand }

// CmdlineMacro is a macro look in the cmdline.
func CmdlineMacro() Look { return cmdlineMacro }

// CmdlineOption is a option look of the cmdline.
func CmdlineOption() Look { return cmdlineOption }

// Highlight is a highlight text look in the list box.
func Highlight() Look { return highlight }

// Title is a title look of the list box.
func Title() Look { return title }

// Selected is a selected item look of in the list.
func Selected() Look { return selected }

// Symlink is a symlink file look.
func Symlink() Look { return symlink }

// SymlinkDir is a symlink directory look.
func SymlinkDir() Look { return symlinkDir }

// Directory is a directory look.
func Directory() Look { return directory }

// Executable is an executable file look.
func Executable() Look { return executable }

// Marked is a marked file look.
func Marked() Look { return marked }

// Finder is a finder text area look.
func Finder() Look { return finder }

// Progress is a gauge look of the progress bar.
func Progress() Look { return progress }

var (
	defaultAttr    Look
	messageInfo    Look
	messageErr     Look
	prompt         Look
	cmdline        Look
	cmdlineCommand Look
	cmdlineMacro   Look
	cmdlineOption  Look
	highlight      Look
	title          Look
	selected       Look
	symlink        Look
	symlinkDir     Look
	directory      Look
	executable     Look
	marked         Look
	finder         Look
	progress       Look
)

func setDefault() {
	defaultAttr = Look{termbox.ColorDefault, termbox.ColorDefault}
	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDefault}
	messageErr = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorDefault}
	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDefault}
	cmdline = Look{termbox.ColorDefault, termbox.ColorDefault}
	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDefault}
	cmdlineMacro = Look{termbox.ColorMagenta, termbox.ColorDefault}
	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorDefault}
	highlight = Look{termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault}
	title = Look{termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault}
	selected = Look{termbox.ColorDefault, termbox.AttrReverse}
	symlink = Look{termbox.ColorMagenta, termbox.ColorDefault}
	symlinkDir = Look{termbox.ColorMagenta | termbox.AttrBold, termbox.ColorDefault}
	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDefault}
	executable = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorDefault}
	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorDefault}
	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
	progress = Look{termbox.ColorDefault, termbox.ColorBlue}
}

func setMidnight() {
	defaultAttr = Look{termbox.ColorWhite, termbox.ColorBlue}
	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlue}
	messageErr = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlue}
	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlue}
	cmdline = Look{termbox.ColorWhite, termbox.ColorBlue}
	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlue}
	cmdlineMacro = Look{termbox.ColorMagenta, termbox.ColorBlue}
	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorBlue}
	highlight = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue}
	title = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue}
	selected = Look{termbox.ColorWhite, termbox.AttrReverse}
	symlink = Look{termbox.ColorMagenta, termbox.ColorBlue}
	symlinkDir = Look{termbox.ColorMagenta | termbox.AttrBold, termbox.ColorBlue}
	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlue}
	executable = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlue}
	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorBlue}
	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
	progress = Look{termbox.ColorWhite, termbox.ColorCyan}
}
