package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/anmitsu/goful"
	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/infobar"
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progbar"
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)

	// Set a title if in a terminal such as screen and tmux.
	if strings.Contains(os.Getenv("TERM"), "screen") {
		os.Stdout.WriteString("\033kgoful\033\\")
	}

	message.Init()
	infobar.Init()
	progbar.Init()

	const (
		state = "~/.goful/state.json"
		hist  = "~/.goful/history/shell"
	)

	app := goful.New(state)
	config(app)

	cmdline.LoadHistory(hist)
	app.Run()
	app.SaveState(state)
	cmdline.SaveHistory(hist)
}

func config(g *goful.Goful) {
	look.Set("default") // default, midnight, dark or gray

	if runtime.GOOS == "windows" {
		widget.SetBorder('|', '-', '+', '+', '+', '+') // not ambiguous runes for layout collapsing
	} else {
		widget.SetBorder('│', '─', '┌', '┐', '└', '┘') // 0x2502, 0x2500, 0x250c, 0x2510, 0x2514, 0x2518
	}

	message.SetInfoLog("~/.goful/log/info.log")   // "" is not logging
	message.SetErrorLog("~/.goful/log/error.log") // "" is not logging
	message.Sec(5)                                // display second for a message

	// Setup widget keymaps.
	g.ConfigFiler(filerKeymap)
	filer.ConfigFinder(finderKeymap)
	cmdline.Config(cmdlineKeymap)
	cmdline.ConfigCompletion(completionKeymap)
	menu.Config(menuKeymap)

	filer.SetStatView(true, false, true) // size, permission and time

	// Setup open command for C-m (when the enter key is pressed)
	// The macro %f means expanded to a file name, for more see (../../spawn.go)
	if runtime.GOOS == "windows" {
		g.AddKeymap("C-m", func() { g.Spawn("explorer %f %&") })
	} else {
		g.AddKeymap("C-m", func() { g.Spawn("xdg-open %f %&") })
	}

	// Setup a shell and a terminal to execute external commands.
	// The shell is called when execute on background by the macro %&.
	// The terminal is called when the other.
	if runtime.GOOS == "windows" {
		g.ConfigShell(func(cmd string) []string {
			return []string{"cmd", "/c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			return []string{"cmd", "/c", "start", "cmd", "/k", cmd}
		})
	} else {
		g.ConfigShell(func(cmd string) []string {
			return []string{"bash", "-c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			// for not close the terminal when the shell finishes running
			const tail = `;read -p "HIT ENTER KEY"`

			if strings.Contains(os.Getenv("TERM"), "screen") { // such as screen and tmux
				return []string{"tmux", "new-window", "-n", cmd, cmd + tail}
			}
			// To execute bash in gnome-terminal of a new window.
			// Opts -x instead of -- to open in a new tab.
			return []string{"gnome-terminal", "--", "bash", "-c", cmd + tail}
		})
	}

	menu.Add("sort",
		"sort name          ", "n", func() { g.Dir().SortName() },
		"sort name decending", "N", func() { g.Dir().SortNameDec() },
		"sort size          ", "s", func() { g.Dir().SortSize() },
		"sort size decending", "S", func() { g.Dir().SortSizeDec() },
		"sort time          ", "t", func() { g.Dir().SortMtime() },
		"sort time decending", "T", func() { g.Dir().SortMtimeDec() },
		"sort ext           ", "e", func() { g.Dir().SortExt() },
		"sort ext decending ", "E", func() { g.Dir().SortExtDec() },
	)
	g.AddKeymap("s", func() { g.Menu("sort") })

	menu.Add("layout",
		"tile             ", "t", func() { g.Workspace().LayoutTile() },
		"tile-top         ", "T", func() { g.Workspace().LayoutTileTop() },
		"tile-bottom      ", "B", func() { g.Workspace().LayoutTileBottom() },
		"oneline          ", "l", func() { g.Workspace().LayoutOneline() },
		"onecolumn        ", "c", func() { g.Workspace().LayoutOnecolumn() },
		"fullscreen       ", "f", func() { g.Workspace().LayoutFullscreen() },
		"toggle size view ", "S", func() { filer.ToggleSizeView() },
		"toggle perm view ", "P", func() { filer.TogglePermView() },
		"toggle time view ", "M", func() { filer.ToggleTimeView() },
		"size only view   ", "1", func() { filer.SetStatView(true, false, false) },
		"perm only view   ", "2", func() { filer.SetStatView(false, true, false) },
		"time only view   ", "3", func() { filer.SetStatView(false, false, true) },
		"view all state   ", "V", func() { filer.SetStatView(true, true, true) },
		"non view state   ", "0", func() { filer.SetStatView(false, false, false) },
		"set default look ", "d", func() { look.Set("default") },
		"set midnight look", "n", func() { look.Set("midnight") },
		"set dark look    ", "D", func() { look.Set("dark") },
		"set gray look    ", "g", func() { look.Set("gray") },
	)
	g.AddKeymap("l", func() { g.Menu("layout") })

	menu.Add("command",
		"copy         ", "c", func() { g.CopyMode() },
		"move         ", "m", func() { g.MoveMode() },
		"delete       ", "D", func() { g.RemoveMode() },
		"mkdir        ", "k", func() { g.MkdirMode() },
		"newfile      ", "n", func() { g.CreatefileMode() },
		"chmod        ", "h", func() { g.ChmodMode() },
		"rename       ", "r", func() { g.RenameMode() },
		"regexp rename", "R", func() { g.BulkRenameMode() },
		"chdir        ", "d", func() { g.ChdirMode() },
		"glob         ", "g", func() { g.GlobMode() },
		"globdir      ", "G", func() { g.GlobdirMode() },
	)
	g.AddKeymap("M-x", func() { g.Menu("command") })

	g.AddKeymap("v", func() { g.Spawn("less %f") })

	g.MergeExtmap(widget.Extmap{
		"C-m": { // associates by file types with the enter key event
			".dir":   func() { g.Dir().EnterDir() },
			".exec":  func() { g.SpawnMode(" ./" + g.File().Name()) },
			".zip":   func() { g.SpawnMode("unzip %f -d %D") },
			".gz":    func() { g.SpawnMode("tar xvfz %f -C %D") },
			".tgz":   func() { g.SpawnMode("tar xvfz %f -C %D") },
			".bz2":   func() { g.SpawnMode("tar xvfj %f -C %D") },
			".tar":   func() { g.SpawnMode("tar xvf %f -C %D") },
			".rar":   func() { g.SpawnMode("unrar x %f -C %D") },
			".java":  func() { g.SpawnMode("javac %f") },
			".class": func() { g.SpawnMode("java %x") },
			".jar":   func() { g.SpawnMode("java -jar %f") },
			".py":    func() { g.SpawnMode("python %f") },
			".go":    func() { g.SpawnMode("go run %f") },
		},
		"v": {
			".gz":  func() { g.Spawn("tar tvfz %f | less") },
			".tgz": func() { g.Spawn("tar tvfz %f | less") },
			".bz2": func() { g.Spawn("tar tvfj %f | less") },
			".tar": func() { g.Spawn("tar tvf %f | less") },
			".zip": func() { g.Spawn("zipinfo -Ocp932 %f | less") },
			".rar": func() { g.Spawn("unrar l %f | less") },
		},
	})

	menu.Add("editor",
		"vscode        ", "c", func() { g.Spawn("code %f %&") },
		"emacs client  ", "e", func() { g.Spawn("emacsclient -n %f %&") },
		"vim           ", "v", func() { g.Spawn("vim %f") },
	)
	g.AddKeymap("e", func() { g.Menu("editor") })

	menu.Add("image",
		"eog        ", "e", func() { g.Spawn("eog %f %&") },
		"gimp       ", "g", func() { g.Spawn("gimp %m %&") },
	)

	menu.Add("media",
		"mpv     ", "m", func() { g.Spawn("mpv %f") },
		"vlc     ", "v", func() { g.Spawn("vlc %f %&") },
	)

	g.MergeExtmap(widget.Extmap{
		"C-m": { // associates image and media files with the enter key event
			".jpg":  func() { g.Menu("image") },
			".jpeg": func() { g.Menu("image") },
			".gif":  func() { g.Menu("image") },
			".png":  func() { g.Menu("image") },
			".bmp":  func() { g.Menu("image") },
			".avi":  func() { g.Menu("media") },
			".mp4":  func() { g.Menu("media") },
			".mkv":  func() { g.Menu("media") },
			".wmv":  func() { g.Menu("media") },
			".flv":  func() { g.Menu("media") },
			".mp3":  func() { g.Menu("media") },
			".flac": func() { g.Menu("media") },
			".tta":  func() { g.Menu("media") },
		},
	})

	menu.Add("jump",
		"Desktop   ", "t", func() { g.Dir().Chdir("~/Desktop/") },
		"Documents ", "D", func() { g.Dir().Chdir("~/Documents/") },
		"Downloads ", "d", func() { g.Dir().Chdir("~/Downloads/") },
		"Music     ", "m", func() { g.Dir().Chdir("~/Music/") },
		"Pictures  ", "p", func() { g.Dir().Chdir("~/Pictures/") },
		"Videos    ", "v", func() { g.Dir().Chdir("~/Videos/") },
	)
	g.AddKeymap("j", func() { g.Menu("jump") })
}

// Widget keymap functions.

func filerKeymap(g *goful.Goful) widget.Keymap {
	return widget.Keymap{
		"M-C-w":     func() { g.CloseWorkspace() },
		"M-f":       func() { g.MoveWorkspace(1) },
		"M-b":       func() { g.MoveWorkspace(-1) },
		"M-C-m":     func() { g.Workspace().CreateDir() },
		"C-w":       func() { g.Workspace().CloseDir() },
		"C-l":       func() { g.Workspace().ReloadAll() },
		"C-f":       func() { g.Workspace().MoveFocus(1) },
		"C-b":       func() { g.Workspace().MoveFocus(-1) },
		"right":     func() { g.Workspace().MoveFocus(1) },
		"left":      func() { g.Workspace().MoveFocus(-1) },
		"F":         func() { g.Workspace().SwapNextDir() },
		"B":         func() { g.Workspace().SwapPrevDir() },
		"w":         func() { g.Workspace().ChdirNeighbor() },
		"C-h":       func() { g.Dir().Chdir("..") },
		"backspace": func() { g.Dir().Chdir("..") },
		"~":         func() { g.Dir().Chdir("~") },
		"\\":        func() { g.Dir().Chdir("/") },
		"C-n":       func() { g.Dir().MoveCursor(1) },
		"C-p":       func() { g.Dir().MoveCursor(-1) },
		"down":      func() { g.Dir().MoveCursor(1) },
		"up":        func() { g.Dir().MoveCursor(-1) },
		"C-d":       func() { g.Dir().MoveCursor(5) },
		"C-u":       func() { g.Dir().MoveCursor(-5) },
		"C-a":       func() { g.Dir().MoveTop() },
		"C-e":       func() { g.Dir().MoveBottom() },
		"home":      func() { g.Dir().MoveTop() },
		"end":       func() { g.Dir().MoveBottom() },
		"M-n":       func() { g.Dir().Scroll(1) },
		"M-p":       func() { g.Dir().Scroll(-1) },
		"C-v":       func() { g.Dir().PageDown() },
		"M-v":       func() { g.Dir().PageUp() },
		"pgdn":      func() { g.Dir().PageDown() },
		"pgup":      func() { g.Dir().PageUp() },
		"space":     func() { g.Dir().ToggleMark() },
		"*":         func() { g.Dir().ToggleMarkAll() },
		"C-g":       func() { g.Dir().Reset() },
		"f":         func() { g.Dir().Finder() },
		"/":         func() { g.Dir().Finder() },
		"q":         func() { g.QuitMode() },
		"Q":         func() { g.QuitMode() },
		"h":         func() { g.SpawnMode("") },
		"H":         func() { g.ShellMode("") },
		"M-W":       func() { g.ChangeWorkspaceTitle() },
		"n":         func() { g.CreatefileMode() },
		"k":         func() { g.MkdirMode() },
		"c":         func() { g.CopyMode() },
		"m":         func() { g.MoveMode() },
		"r":         func() { g.RenameMode() },
		"R":         func() { g.BulkRenameMode() },
		"D":         func() { g.RemoveMode() },
	}
}

func finderKeymap(w *filer.Finder) widget.Keymap {
	return widget.Keymap{
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-p":       func() { w.MoveHistory(1) },
		"M-n":       func() { w.MoveHistory(-1) },
		"C-g":       func() { w.Exit() },
	}
}

func cmdlineKeymap(w *cmdline.Cmdline) widget.Keymap {
	return widget.Keymap{
		"C-a":       func() { w.MoveTop() },
		"C-e":       func() { w.MoveBottom() },
		"C-f":       func() { w.ForwardChar() },
		"C-b":       func() { w.BackwardChar() },
		"right":     func() { w.ForwardChar() },
		"left":      func() { w.BackwardChar() },
		"M-f":       func() { w.ForwardWord() },
		"M-b":       func() { w.BackwardWord() },
		"C-d":       func() { w.DeleteChar() },
		"delete":    func() { w.DeleteChar() },
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-d":       func() { w.DeleteForwardWord() },
		"M-h":       func() { w.DeleteBackwardWord() },
		"C-k":       func() { w.KillLine() },
		"C-i":       func() { w.StartCompletion() },
		"C-m":       func() { w.Run() },
		"C-g":       func() { w.Exit() },
		"C-n":       func() { w.History.CursorDown() },
		"C-p":       func() { w.History.CursorUp() },
		"down":      func() { w.History.CursorDown() },
		"up":        func() { w.History.CursorUp() },
		"C-v":       func() { w.History.PageDown() },
		"M-v":       func() { w.History.PageUp() },
		"M-<":       func() { w.History.MoveTop() },
		"M->":       func() { w.History.MoveBottom() },
		"pgup":      func() { w.History.MoveTop() },
		"pgdn":      func() { w.History.MoveBottom() },
		"M-n":       func() { w.History.Scroll(1) },
		"M-p":       func() { w.History.Scroll(-1) },
		"C-x":       func() { w.History.Delete() },
	}
}

func completionKeymap(w *cmdline.Completion) widget.Keymap {
	return widget.Keymap{
		"C-n":   func() { w.CursorDown() },
		"C-p":   func() { w.CursorUp() },
		"down":  func() { w.CursorDown() },
		"up":    func() { w.CursorUp() },
		"C-f":   func() { w.CursorToRight() },
		"C-b":   func() { w.CursorToLeft() },
		"right": func() { w.CursorToRight() },
		"left":  func() { w.CursorToLeft() },
		"C-i":   func() { w.CursorToRight() },
		"C-v":   func() { w.PageDown() },
		"M-v":   func() { w.PageUp() },
		"pgdn":  func() { w.PageDown() },
		"pgup":  func() { w.PageUp() },
		"M-<":   func() { w.MoveTop() },
		"M->":   func() { w.MoveBottom() },
		"home":  func() { w.MoveTop() },
		"end":   func() { w.MoveBottom() },
		"M-n":   func() { w.Scroll(1) },
		"M-p":   func() { w.Scroll(-1) },
		"C-m":   func() { w.InsertCompletion() },
		"C-g":   func() { w.Exit() },
	}
}

func menuKeymap(w *menu.Menu) widget.Keymap {
	return widget.Keymap{
		"C-n":  func() { w.MoveCursor(1) },
		"C-p":  func() { w.MoveCursor(-1) },
		"down": func() { w.MoveCursor(1) },
		"up":   func() { w.MoveCursor(-1) },
		"C-v":  func() { w.PageDown() },
		"M-v":  func() { w.PageUp() },
		"M->":  func() { w.MoveBottom() },
		"M-<":  func() { w.MoveTop() },
		"C-m":  func() { w.Exec() },
		"C-g":  func() { w.Exit() },
	}
}
