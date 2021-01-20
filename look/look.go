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
	// TODO
	// case "midnight":
	// 	setMidnight()
	// case "dark":
	// 	setDark()
	// case "gray":
	// setGray()
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

// Selected is a selected item look of in the list.
func Selected() tcell.Style { return selected }

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
	selected       tcell.Style
	symlink        tcell.Style
	symlinkDir     tcell.Style
	directory      tcell.Style
	executable     tcell.Style
	marked         tcell.Style
	finder         tcell.Style
	progress       tcell.Style
)

func setDefault() {
	defaultAttr = tcell.StyleDefault
	messageInfo = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	messageErr = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	prompt = tcell.StyleDefault.Foreground(tcell.ColorTeal).Bold(true)
	cmdline = tcell.StyleDefault
	cmdlineCommand = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	cmdlineMacro = tcell.StyleDefault.Foreground(tcell.ColorPurple)
	cmdlineOption = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	highlight = tcell.StyleDefault.Bold(true)
	title = tcell.StyleDefault.Foreground(tcell.ColorTeal).Bold(true)
	selected = tcell.StyleDefault.Reverse(true)
	symlink = tcell.StyleDefault.Foreground(tcell.ColorPurple)
	symlinkDir = tcell.StyleDefault.Foreground(tcell.ColorPurple).Bold(true)
	directory = tcell.StyleDefault.Foreground(tcell.ColorTeal).Bold(true)
	executable = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	marked = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	finder = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorTeal)
	progress = tcell.StyleDefault.Background(tcell.ColorBlue)
}

// TODO

// func setDefault() {
// 	defaultAttr = Look{termbox.ColorDefault, termbox.ColorDefault}
// 	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDefault}
// 	messageErr = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorDefault}
// 	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDefault}
// 	cmdline = Look{termbox.ColorDefault, termbox.ColorDefault}
// 	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDefault}
// 	cmdlineMacro = Look{termbox.ColorMagenta, termbox.ColorDefault}
// 	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorDefault}
// 	highlight = Look{termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault}
// 	title = Look{termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault}
// 	selected = Look{termbox.AttrReverse, termbox.ColorDefault}
// 	symlink = Look{termbox.ColorMagenta, termbox.ColorDefault}
// 	symlinkDir = Look{termbox.ColorMagenta | termbox.AttrBold, termbox.ColorDefault}
// 	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDefault}
// 	executable = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorDefault}
// 	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorDefault}
// 	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
// 	progress = Look{termbox.ColorDefault, termbox.ColorBlue}
// }

// func setMidnight() {
// 	defaultAttr = Look{termbox.ColorWhite, termbox.ColorBlue}
// 	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlue}
// 	messageErr = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlue}
// 	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlue}
// 	cmdline = Look{termbox.ColorWhite, termbox.ColorBlue}
// 	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlue}
// 	cmdlineMacro = Look{termbox.ColorMagenta, termbox.ColorBlue}
// 	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorBlue}
// 	highlight = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue}
// 	title = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue}
// 	selected = Look{termbox.AttrReverse, termbox.ColorBlue}
// 	symlink = Look{termbox.ColorMagenta, termbox.ColorBlue}
// 	symlinkDir = Look{termbox.ColorMagenta | termbox.AttrBold, termbox.ColorBlue}
// 	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlue}
// 	executable = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlue}
// 	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorBlue}
// 	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
// 	progress = Look{termbox.ColorWhite, termbox.ColorCyan}
// }

// func setDark() {
// 	defaultAttr = Look{termbox.ColorWhite, termbox.ColorBlack}
// 	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack}
// 	messageErr = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack}
// 	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlack}
// 	cmdline = Look{termbox.ColorWhite, termbox.ColorBlack}
// 	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack}
// 	cmdlineMacro = Look{termbox.ColorMagenta, termbox.ColorBlack}
// 	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorBlack}
// 	highlight = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack}
// 	title = Look{termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack}
// 	selected = Look{termbox.AttrReverse, termbox.ColorBlack}
// 	symlink = Look{termbox.ColorMagenta, termbox.ColorBlack}
// 	symlinkDir = Look{termbox.ColorMagenta | termbox.AttrBold, termbox.ColorBlack}
// 	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorBlack}
// 	executable = Look{termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack}
// 	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorBlack}
// 	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
// 	progress = Look{termbox.ColorWhite, termbox.ColorBlue}
// }

// func setGray() {
// 	defaultAttr = Look{termbox.ColorLightGray, termbox.ColorDarkGray}
// 	messageInfo = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDarkGray}
// 	messageErr = Look{termbox.ColorLightRed | termbox.AttrBold, termbox.ColorDarkGray}
// 	prompt = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDarkGray}
// 	cmdline = Look{termbox.ColorLightGray, termbox.ColorDarkGray}
// 	cmdlineCommand = Look{termbox.ColorGreen | termbox.AttrBold, termbox.ColorDarkGray}
// 	cmdlineMacro = Look{termbox.ColorLightMagenta, termbox.ColorDarkGray}
// 	cmdlineOption = Look{termbox.ColorYellow, termbox.ColorDarkGray}
// 	highlight = Look{termbox.ColorLightGray | termbox.AttrBold, termbox.ColorDarkGray}
// 	title = Look{termbox.ColorLightGray | termbox.AttrBold, termbox.ColorDarkGray}
// 	selected = Look{termbox.AttrReverse, termbox.ColorDarkGray}
// 	symlink = Look{termbox.ColorLightMagenta, termbox.ColorDarkGray}
// 	symlinkDir = Look{termbox.ColorLightMagenta | termbox.AttrBold, termbox.ColorDarkGray}
// 	directory = Look{termbox.ColorCyan | termbox.AttrBold, termbox.ColorDarkGray}
// 	executable = Look{termbox.ColorLightRed | termbox.AttrBold, termbox.ColorDarkGray}
// 	marked = Look{termbox.ColorYellow | termbox.AttrBold, termbox.ColorDarkGray}
// 	finder = Look{termbox.ColorBlack, termbox.ColorCyan}
// 	progress = Look{termbox.ColorLightGray, termbox.ColorBlue}
// }
