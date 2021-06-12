package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

func main() {
	is_tmux := false
	widget.Init()
	defer widget.Fini()

	if runtime.GOOS == "darwin" {
		is_tmux = strings.Contains(os.Getenv("TERM_PROGRAM"), "tmux")
	} else {
		is_tmux = strings.Contains(os.Getenv("TERM"), "screen")
	}
	// Change a terminal title.
	if is_tmux {
		os.Stdout.WriteString("\033kgoful\033") // for tmux
	} else {
		os.Stdout.WriteString("\033]0;goful\007") // for otherwise
	}

	const state = "~/.goful/state.json"
	const history = "~/.goful/history/shell"

	goful := app.NewGoful(state)
	config(goful, is_tmux)
	_ = cmdline.LoadHistory(history)

	goful.Run()

	_ = goful.SaveState(state)
	_ = cmdline.SaveHistory(history)
}

func config(g *app.Goful, is_tmux bool) {
	look.Set("default") // default, midnight, black, white

	if runewidth.EastAsianWidth {
		// Because layout collapsing for ambiguous runes if LANG=ja_JP.
		widget.SetBorder('|', '-', '+', '+', '+', '+')
	} else {
		// Look good if environment variable RUNEWIDTH_EASTASIAN=0 and
		// ambiguous char setting is half-width for gnome-terminal.
		widget.SetBorder('│', '─', '┌', '┐', '└', '┘') // 0x2502, 0x2500, 0x250c, 0x2510, 0x2514, 0x2518
	}
	g.SetBorderStyle(widget.AllBorder) // AllBorder, ULBorder, NoBorder

	message.SetInfoLog("~/.goful/log/info.log")   // "" is not logging
	message.SetErrorLog("~/.goful/log/error.log") // "" is not logging
	message.Sec(5)                                // display second for a message

	// Setup widget keymaps.
	g.ConfigFiler(filerKeymap)
	filer.ConfigFinder(finderKeymap)
	cmdline.Config(cmdlineKeymap)
	cmdline.ConfigCompletion(completionKeymap)
	menu.Config(menuKeymap)

	filer.SetStatView(true, false, true)  // size, permission and time
	filer.SetTimeFormat("06-01-02 15:04") // ex: "Jan _2 15:04"

	// Setup open command for C-m (when the enter key is pressed)
	// The macro %f means expanded to a file name, for more see (spawn.go)
	opener := "xdg-open %f %&"
	switch runtime.GOOS {
	case "windows":
		opener = "explorer %~f %&"
	case "darwin":
		opener = "open %f %&"
	}
	g.MergeKeymap(widget.Keymap{
		"C-m": func() { g.Spawn(opener) },
		"o":   func() { g.Spawn(opener) },
	})

	// Setup pager by $PAGER
	pager := os.Getenv("PAGER")
	if pager == "" {
		if runtime.GOOS == "windows" {
			pager = "more"
		} else {
			pager = "less"
		}
	}
	if runtime.GOOS == "windows" {
		pager += " %~f"
	} else {
		pager += " %f"
	}
	g.AddKeymap("i", func() { g.Spawn(pager) })

	// Setup a shell and a terminal to execute external commands.
	// The shell is called when execute on background by the macro %&.
	// The terminal is called when the other.
	if runtime.GOOS == "windows" {
		g.ConfigShell(func(cmd string) []string {
			return []string{"cmd", "/c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			return []string{"cmd", "/c", "start", "cmd", "/c", cmd + "& pause"}
		})
	} else {
		g.ConfigShell(func(cmd string) []string {
			return []string{"bash", "-c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			// for not close the terminal when the shell finishes running
			const tail = `;read -p "HIT ENTER KEY"`

			if is_tmux { // such as screen and tmux
				return []string{"tmux", "new-window", "-n", cmd, cmd + tail}
			}
			// To execute bash in gnome-terminal of a new window or tab.
			title := "echo -n '\033]0;" + cmd + "\007';" // for change title
			return []string{"gnome-terminal", "--", "bash", "-c", title + cmd + tail}
		})
	}

	// Setup menus and add to keymap.
	menu.Add("sort",
		"n", "sort name          ", func() { g.Dir().SortName() },
		"N", "sort name decending", func() { g.Dir().SortNameDec() },
		"s", "sort size          ", func() { g.Dir().SortSize() },
		"S", "sort size decending", func() { g.Dir().SortSizeDec() },
		"t", "sort time          ", func() { g.Dir().SortMtime() },
		"T", "sort time decending", func() { g.Dir().SortMtimeDec() },
		"e", "sort ext           ", func() { g.Dir().SortExt() },
		"E", "sort ext decending ", func() { g.Dir().SortExtDec() },
		".", "toggle priority    ", func() { filer.TogglePriority(); g.Workspace().ReloadAll() },
	)
	g.AddKeymap("s", func() { g.Menu("sort") })

	menu.Add("view",
		"s", "stat menu    ", func() { g.Menu("stat") },
		"l", "layout menu  ", func() { g.Menu("layout") },
		"L", "look menu    ", func() { g.Menu("look") },
		".", "toggle show hidden files", func() { filer.ToggleShowHiddens(); g.Workspace().ReloadAll() },
	)
	g.AddKeymap("v", func() { g.Menu("view") })

	menu.Add("layout",
		"t", "tile       ", func() { g.Workspace().LayoutTile() },
		"T", "tile-top   ", func() { g.Workspace().LayoutTileTop() },
		"b", "tile-bottom", func() { g.Workspace().LayoutTileBottom() },
		"r", "one-row    ", func() { g.Workspace().LayoutOnerow() },
		"c", "one-column ", func() { g.Workspace().LayoutOnecolumn() },
		"f", "fullscreen ", func() { g.Workspace().LayoutFullscreen() },
	)

	menu.Add("stat",
		"s", "toggle size  ", func() { filer.ToggleSizeView() },
		"p", "toggle perm  ", func() { filer.TogglePermView() },
		"t", "toggle time  ", func() { filer.ToggleTimeView() },
		"1", "all stat     ", func() { filer.SetStatView(true, true, true) },
		"0", "no stat      ", func() { filer.SetStatView(false, false, false) },
	)

	menu.Add("look",
		"d", "default      ", func() { look.Set("default") },
		"n", "midnight     ", func() { look.Set("midnight") },
		"b", "black        ", func() { look.Set("black") },
		"w", "white        ", func() { look.Set("white") },
		"a", "all border   ", func() { g.SetBorderStyle(widget.AllBorder) },
		"u", "ul border    ", func() { g.SetBorderStyle(widget.ULBorder) },
		"0", "no border    ", func() { g.SetBorderStyle(widget.NoBorder) },
	)

	menu.Add("command",
		"c", "copy         ", func() { g.Copy() },
		"m", "move         ", func() { g.Move() },
		"D", "delete       ", func() { g.Remove() },
		"k", "mkdir        ", func() { g.Mkdir() },
		"n", "newfile      ", func() { g.Touch() },
		"M", "chmod        ", func() { g.Chmod() },
		"r", "rename       ", func() { g.Rename() },
		"R", "bulk rename  ", func() { g.BulkRename() },
		"d", "chdir        ", func() { g.Chdir() },
		"g", "glob         ", func() { g.Glob() },
		"G", "globdir      ", func() { g.Globdir() },
	)
	g.AddKeymap("x", func() { g.Menu("command") })

	if runtime.GOOS == "windows" {
		menu.Add("external-command",
			"c", "copy %~f to %~D2 ", func() { g.Shell("robocopy /e %~f %~D2") },
			"m", "move %~f to %~D2 ", func() { g.Shell("move /-y %~f %~D2") },
			"d", "del /s %~m       ", func() { g.Shell("del /s %~m") },
			"D", "rd /s /q %~m     ", func() { g.Shell("rd /s /q %~m") },
			"k", "make directory   ", func() { g.Shell("mkdir ") },
			"n", "create newfile   ", func() { g.Shell("copy nul ") },
			"r", "move (rename) %f ", func() { g.Shell("move /-y %~f ./") },
			"w", "where . *        ", func() { g.Shell("where . *") },
		)
	} else {
		menu.Add("external-command",
			"c", "copy %m to %D2    ", func() { g.Shell("cp -vai %m %D2") },
			"m", "move %m to %D2    ", func() { g.Shell("mv -vi %m %D2") },
			"D", "remove %m files   ", func() { g.Shell("rm -vR %m") },
			"k", "make directory    ", func() { g.Shell("mkdir -vp ./") },
			"n", "create newfile    ", func() { g.Shell("touch ./") },
			"T", "time copy %f to %m", func() { g.Shell("touch -r %f %m") },
			"M", "change mode %m    ", func() { g.Shell("chmod 644 %m", -3) },
			"r", "move (rename) %f  ", func() { g.Shell("mv -vi %f " + g.File().Name()) },
			"R", "bulk rename %m    ", func() { g.Shell(`rename -v "s///" %m`, -6) },
			"f", "find . -name      ", func() { g.Shell(`find . -name "*"`, -1) },
			"A", "archives menu     ", func() { g.Menu("archive") },
		)
	}
	g.AddKeymap("X", func() { g.Menu("external-command") })

	menu.Add("archive",
		"z", "zip     ", func() { g.Shell(`zip -roD %x.zip %m`, -7) },
		"t", "tar     ", func() { g.Shell(`tar cvf %x.tar %m`, -7) },
		"g", "tar.gz  ", func() { g.Shell(`tar cvfz %x.tgz %m`, -7) },
		"b", "tar.bz2 ", func() { g.Shell(`tar cvfj %x.bz2 %m`, -7) },
		"x", "tar.xz  ", func() { g.Shell(`tar cvfJ %x.txz %m`, -7) },
		"r", "rar     ", func() { g.Shell(`rar u %x.rar %m`, -7) },

		"Z", "extract zip for %m", func() { g.Shell(`for i in %m; do unzip "$i" -d ./; done`, -6) },
		"T", "extract tar for %m", func() { g.Shell(`for i in %m; do tar xvf "$i" -C ./; done`, -6) },
		"G", "extract tgz for %m", func() { g.Shell(`for i in %m; do tar xvfz "$i" -C ./; done`, -6) },
		"B", "extract bz2 for %m", func() { g.Shell(`for i in %m; do tar xvfj "$i" -C ./; done`, -6) },
		"X", "extract txz for %m", func() { g.Shell(`for i in %m; do tar xvfJ "$i" -C ./; done`, -6) },
		"R", "extract rar for %m", func() { g.Shell(`for i in %m; do unrar x "$i" -C ./; done`, -6) },

		"1", "find . *.zip extract", func() { g.Shell(`find . -name "*.zip" -type f -prune -print0 | xargs -n1 -0 unzip -d ./`) },
		"2", "find . *.tar extract", func() { g.Shell(`find . -name "*.tar" -type f -prune -print0 | xargs -n1 -0 tar xvf -C ./`) },
		"3", "find . *.tgz extract", func() { g.Shell(`find . -name "*.tgz" -type f -prune -print0 | xargs -n1 -0 tar xvfz -C ./`) },
		"4", "find . *.bz2 extract", func() { g.Shell(`find . -name "*.bz2" -type f -prune -print0 | xargs -n1 -0 tar xvfj -C ./`) },
		"5", "find . *.txz extract", func() { g.Shell(`find . -name "*.txz" -type f -prune -print0 | xargs -n1 -0 tar xvfJ -C ./`) },
		"6", "find . *.rar extract", func() { g.Shell(`find . -name "*.rar" -type f -prune -print0 | xargs -n1 -0 unrar x -C ./`) },
	)

	menu.Add("bookmark",
		"t", "~/Desktop  ", func() { g.Dir().Chdir("~/Desktop") },
		"c", "~/Documents", func() { g.Dir().Chdir("~/Documents") },
		"d", "~/Downloads", func() { g.Dir().Chdir("~/Downloads") },
		"m", "~/Music    ", func() { g.Dir().Chdir("~/Music") },
		"p", "~/Pictures ", func() { g.Dir().Chdir("~/Pictures") },
		"v", "~/Videos   ", func() { g.Dir().Chdir("~/Videos") },
	)
	if runtime.GOOS == "windows" {
		menu.Add("bookmark",
			"C", "C:/", func() { g.Dir().Chdir("C:/") },
			"D", "D:/", func() { g.Dir().Chdir("D:/") },
			"E", "E:/", func() { g.Dir().Chdir("E:/") },
		)
	} else {
		menu.Add("bookmark",
			"e", "/etc   ", func() { g.Dir().Chdir("/etc") },
			"u", "/usr   ", func() { g.Dir().Chdir("/usr") },
			"x", "/media ", func() { g.Dir().Chdir("/media") },
		)
	}
	g.AddKeymap("b", func() { g.Menu("bookmark") })

	menu.Add("editor",
		"c", "vscode        ", func() { g.Spawn("code %f %&") },
		"e", "emacs client  ", func() { g.Spawn("emacsclient -n %f %&") },
		"v", "vim           ", func() { g.Spawn("vim %f") },
	)
	g.AddKeymap("e", func() { g.Menu("editor") })

	menu.Add("image",
		"x", "default    ", func() { g.Spawn(opener) },
		"e", "eog        ", func() { g.Spawn("eog %f %&") },
		"g", "gimp       ", func() { g.Spawn("gimp %m %&") },
	)

	menu.Add("media",
		"x", "default ", func() { g.Spawn(opener) },
		"m", "mpv     ", func() { g.Spawn("mpv %f") },
		"v", "vlc     ", func() { g.Spawn("vlc %f %&") },
	)

	var associate widget.Keymap
	if runtime.GOOS == "windows" {
		associate = widget.Keymap{
			".dir": func() { g.Dir().EnterDir() },
			".go":  func() { g.Shell("go run %~f") },
			".py":  func() { g.Shell("python %~f") },
			".rb":  func() { g.Shell("ruby %~f") },
			".js":  func() { g.Shell("node %~f") },
		}
	} else {
		associate = widget.Keymap{
			".dir":  func() { g.Dir().EnterDir() },
			".exec": func() { g.Shell(" ./" + g.File().Name()) },

			".zip": func() { g.Shell("unzip %f -d %D") },
			".tar": func() { g.Shell("tar xvf %f -C %D") },
			".gz":  func() { g.Shell("tar xvfz %f -C %D") },
			".tgz": func() { g.Shell("tar xvfz %f -C %D") },
			".bz2": func() { g.Shell("tar xvfj %f -C %D") },
			".xz":  func() { g.Shell("tar xvfJ %f -C %D") },
			".txz": func() { g.Shell("tar xvfJ %f -C %D") },
			".rar": func() { g.Shell("unrar x %f -C %D") },

			".go": func() { g.Shell("go run %f") },
			".py": func() { g.Shell("python %f") },
			".rb": func() { g.Shell("ruby %f") },
			".js": func() { g.Shell("node %f") },

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
		}
	}

	g.MergeExtmap(widget.Extmap{
		"C-m": associate,
		"o":   associate,
	})
}

// Widget keymap functions.

func filerKeymap(g *app.Goful) widget.Keymap {
	return widget.Keymap{
		"M-C-o":     func() { g.CreateWorkspace() },
		"M-C-w":     func() { g.CloseWorkspace() },
		"M-f":       func() { g.MoveWorkspace(1) },
		"M-b":       func() { g.MoveWorkspace(-1) },
		"C-o":       func() { g.Workspace().CreateDir() },
		"C-w":       func() { g.Workspace().CloseDir() },
		"C-l":       func() { g.Workspace().ReloadAll() },
		"C-f":       func() { g.Workspace().MoveFocus(1) },
		"C-b":       func() { g.Workspace().MoveFocus(-1) },
		"right":     func() { g.Workspace().MoveFocus(1) },
		"left":      func() { g.Workspace().MoveFocus(-1) },
		"C-i":       func() { g.Workspace().MoveFocus(1) },
		"l":         func() { g.Workspace().MoveFocus(1) },
		"h":         func() { g.Workspace().MoveFocus(-1) },
		"F":         func() { g.Workspace().SwapNextDir() },
		"B":         func() { g.Workspace().SwapPrevDir() },
		"w":         func() { g.Workspace().ChdirNeighbor() },
		"C-h":       func() { g.Dir().Chdir("..") },
		"backspace": func() { g.Dir().Chdir("..") },
		"u":         func() { g.Dir().Chdir("..") },
		"~":         func() { g.Dir().Chdir("~") },
		"\\":        func() { g.Dir().Chdir("/") },
		"C-n":       func() { g.Dir().MoveCursor(1) },
		"C-p":       func() { g.Dir().MoveCursor(-1) },
		"down":      func() { g.Dir().MoveCursor(1) },
		"up":        func() { g.Dir().MoveCursor(-1) },
		"j":         func() { g.Dir().MoveCursor(1) },
		"k":         func() { g.Dir().MoveCursor(-1) },
		"C-d":       func() { g.Dir().MoveCursor(5) },
		"C-u":       func() { g.Dir().MoveCursor(-5) },
		"C-a":       func() { g.Dir().MoveTop() },
		"C-e":       func() { g.Dir().MoveBottom() },
		"home":      func() { g.Dir().MoveTop() },
		"end":       func() { g.Dir().MoveBottom() },
		"^":         func() { g.Dir().MoveTop() },
		"$":         func() { g.Dir().MoveBottom() },
		"M-n":       func() { g.Dir().Scroll(1) },
		"M-p":       func() { g.Dir().Scroll(-1) },
		"C-v":       func() { g.Dir().PageDown() },
		"M-v":       func() { g.Dir().PageUp() },
		"pgdn":      func() { g.Dir().PageDown() },
		"pgup":      func() { g.Dir().PageUp() },
		" ":         func() { g.Dir().ToggleMark() },
		"C-space":   func() { g.Dir().InvertMark() },
		"C-g":       func() { g.Dir().Reset() },
		"C-[":       func() { g.Dir().Reset() }, // C-[ means ESC
		"f":         func() { g.Dir().Finder() },
		"/":         func() { g.Dir().Finder() },
		"q":         func() { g.Quit() },
		"Q":         func() { g.Quit() },
		":":         func() { g.Shell("") },
		";":         func() { g.ShellSuspend("") },
		"M-W":       func() { g.ChangeWorkspaceTitle() },
		"n":         func() { g.Touch() },
		"K":         func() { g.Mkdir() },
		"c":         func() { g.Copy() },
		"m":         func() { g.Move() },
		"r":         func() { g.Rename() },
		"R":         func() { g.BulkRename() },
		"D":         func() { g.Remove() },
		"d":         func() { g.Chdir() },
		"g":         func() { g.Glob() },
		"G":         func() { g.Globdir() },
	}
}

func finderKeymap(w *filer.Finder) widget.Keymap {
	return widget.Keymap{
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-p":       func() { w.MoveHistory(1) },
		"M-n":       func() { w.MoveHistory(-1) },
		"C-g":       func() { w.Exit() },
		"C-[":       func() { w.Exit() },
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
		"C-[":       func() { w.Exit() },
		"C-n":       func() { w.History.CursorDown() },
		"C-p":       func() { w.History.CursorUp() },
		"down":      func() { w.History.CursorDown() },
		"up":        func() { w.History.CursorUp() },
		"C-v":       func() { w.History.PageDown() },
		"M-v":       func() { w.History.PageUp() },
		"pgdn":      func() { w.History.PageDown() },
		"pgup":      func() { w.History.PageUp() },
		"M-<":       func() { w.History.MoveTop() },
		"M->":       func() { w.History.MoveBottom() },
		"home":      func() { w.History.MoveTop() },
		"end":       func() { w.History.MoveBottom() },
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
		"C-[":   func() { w.Exit() },
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
		"C-[":  func() { w.Exit() },
	}
}
