package goful

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

// Spawn a process by an external terminal.
func (g *Goful) Spawn(cmd string) {
	cmd, background := g.expandMacro(cmd)
	newSpawn(cmd, background, g.shell).spawn()
}

const (
	macroPrefix             = '%'
	macroEscape             = '\\' // \ is an escape sequence for the macro prefix %
	macroFile               = 'f'  // %f is expanded a file name on the cursor
	macroFilePath           = 'F'  // %F is expanded a file path on the cursor
	macroFileWithoutExt     = 'x'  // %x is expanded a file name excluded the extention on the cursor
	macroFileWithoutExtPath = 'X'  // %x is expanded a file path excluded the extention on the cursor
	macroMarkfile           = 'm'  // %m is expanded mark file names joined by spaces
	macroMarkfilePath       = 'M'  // %M is expanded mark file paths joined by spaces
	macroDir                = 'd'  // %d is expanded a directory name on the cursor
	macroDirPath            = 'D'  // %D is expanded a directory path on the cursor
	macroNextDir            = '2'  // %D2 is expanded the neighbor directory path
	macroRunBackground      = '&'  // %& is a flag runned in background
)

func (g *Goful) expandMacro(cmd string) (result string, background bool) {
	data := []byte(cmd)
	ret := make([]byte, len(data))
	copy(ret, data)

	background = false
	escape := false
	prefix := false
	offset := 0
	for i, b := range data {
		if escape { // skip the escape sequence
			ret = widget.DeleteBytes(ret, offset-1, 1)
			escape = false
			continue
		}

		if prefix {
			prefix = false
			src := ""
			macrolen := 2
			switch b {
			case macroFile:
				src = utils.Quote(g.File().Name())
			case macroFilePath:
				src = utils.Quote(g.File().Path())
			case macroFileWithoutExt:
				src = utils.Quote(utils.RemoveExt(g.File().Name()))
			case macroFileWithoutExtPath:
				src = utils.Quote(utils.RemoveExt(g.File().Path()))
			case macroMarkfile:
				src = strings.Join(g.Dir().MarkfileQuotedNames(), " ")
			case macroMarkfilePath:
				src = strings.Join(g.Dir().MarkfileQuotedPaths(), " ")
			case macroDir:
				if i != len(data)-1 && data[i+1] == macroNextDir {
					src = g.Workspace().NextDir().Base()
					macrolen = 3
				} else {
					src = g.Dir().Base()
				}
				src = utils.Quote(src)
			case macroDirPath:
				if i != len(data)-1 && data[i+1] == macroNextDir {
					src = g.Workspace().NextDir().Path
					macrolen = 3
				} else {
					src = g.Dir().Path
				}
				src = utils.Quote(src)
			case macroRunBackground:
				background = true
			default:
				goto other
			}
			ret = widget.DeleteBytes(ret, offset-1, macrolen)
			ret = widget.InsertBytes(ret, []byte(src), offset-1)
			offset += len(src) - macrolen
			offset++
			continue
		}
	other:
		switch b {
		case macroPrefix:
			prefix = true
		case macroEscape:
			escape = true
		}
		offset++
	}
	return string(ret), background
}

// type Spawner interface {
// 	Spawn() error
// 	Shell() string
// 	ShellCmd()
// 	Title()
// }

type shell struct {
	name       string
	opts       []string
	screen     string
	screenopts []string
}

type spawn struct {
	cmd        string
	background bool
	shell      shell
}

func newSpawn(cmd string, background bool, s shell) *spawn {
	return &spawn{cmd, background, s}
}

func (s *spawn) spawn() {
	screen := false
	if strings.Contains(os.Getenv("TERM"), "screen") {
		screen = true
	}

	var err error
	var cmd *exec.Cmd
	if screen && !s.background {
		opts := append(s.shell.screenopts, s.cmd, s.cmd+`; read -p "HIT ENTER KEY"`)
		cmd = exec.Command(s.shell.screen, opts...)
		err = s.runNewscreen(cmd)
	} else {
		opts := append(s.shell.opts, s.cmd)
		cmd = exec.Command(s.shell.name, opts...)
		if s.background {
			err = s.runBackground(cmd)
		} else {
			err = s.runForeground(cmd)
		}
	}
	message.Info(strings.Join(cmd.Args, " "))
	if err != nil {
		message.Error(err)
	}
}

func (s *spawn) runNewscreen(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (s *spawn) runBackground(cmd *exec.Cmd) error {
	var bufout, buferr bytes.Buffer
	cmd.Stdout = &bufout
	cmd.Stderr = &buferr
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		cmd.Wait()
		if bufout.Len() > 0 {
			message.Info(bufout.String())
		}
		if buferr.Len() > 0 {
			message.Errorf(buferr.String())
		}
	}()
	return nil
}

func (s *spawn) runForeground(cmd *exec.Cmd) error {
	termbox.Close()
	defer func() {
		termbox.Init()
		termbox.SetInputMode(termbox.InputAlt)
	}()

	c := exec.Command("clear")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	var buf string
	fmt.Println("\nHIT ENTER KEY")
	fmt.Scanln(&buf)
	return nil
}
