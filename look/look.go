package look

import "github.com/nsf/termbox-go"

// Look is attributes of termbox.
type Look struct {
	fg termbox.Attribute
	bg termbox.Attribute
}

// Fg is getter of foreground attribute.
func (l Look) Fg() termbox.Attribute {
	return l.fg
}

// Bg is getter of background attribute.
func (l Look) Bg() termbox.Attribute {
	return l.bg
}

// And combines this and others.
func (l Look) And(others ...Look) Look {
	fg, bg := l.fg, l.bg
	for _, o := range others {
		fg |= o.fg
		bg |= o.bg
	}
	return Look{fg, bg}
}

// Set look for name.
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

// Default is default look attribute.
func Default() Look { return defaultAttr }

// Blank is blank look attribute.
func Blank() Look { return Look{0, 0} }

// MessageInfo is look for infomation message.
func MessageInfo() Look { return messageInfo }

// MessageErr is look for error message.
func MessageErr() Look { return messageErr }

// Prompt is look for prompt of cmdline and finder.
func Prompt() Look { return prompt }

// Cmdline is look for cmdline text area.
func Cmdline() Look { return cmdline }

// CmdlineCommand is look for highlighted commands of cmdline text.
func CmdlineCommand() Look { return cmdlineCommand }

// CmdlineMacro is look for macro of cmdline text.
func CmdlineMacro() Look { return cmdlineMacro }

// CmdlineOption is look for options of cmdline text.
func CmdlineOption() Look { return cmdlineOption }

// Highlight is look for highlight text of listbox.
func Highlight() Look { return highlight }

// Title is look for title of listbox.
func Title() Look { return title }

// Selected is look for selected item of container.
func Selected() Look { return selected }

// Symlink is look for symlink file.
func Symlink() Look { return symlink }

// SymlinkDir is look for symlink directory.
func SymlinkDir() Look { return symlinkDir }

// Directory is look for directory.
func Directory() Look { return directory }

// Executable is look for executable file.
func Executable() Look { return executable }

// Marked is look for marked file.
func Marked() Look { return marked }

// Finder is look for finder of directory.
func Finder() Look { return finder }

// Progress is look for gauge of progress bar.
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
