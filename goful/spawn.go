package goful

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

// Spawn a process by the shell or the terminal.
func (g *Goful) Spawn(cmd string) {
	cmd, background := g.expandMacro(cmd)
	var args []string
	if background {
		args = g.shell(cmd)
	} else {
		args = g.terminal(cmd)
	}
	execCmd := exec.Command(args[0], args[1:]...)
	message.Info(strings.Join(execCmd.Args, " "))
	if err := spawn(execCmd); err != nil {
		message.Error(err)
	}
}

func spawn(cmd *exec.Cmd) error {
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

// SpawnSuspend spawns a process and suspends termbox.
func (g *Goful) SpawnSuspend(cmd string) {
	cmd, _ = g.expandMacro(cmd)
	args := g.shell(cmd)
	execCmd := exec.Command(args[0], args[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	termbox.Close()
	defer func(cmd string) {
		termbox.Init()
		termbox.SetInputMode(termbox.InputAlt)
		message.Info(cmd)
	}(strings.Join(execCmd.Args, " "))
	execCmd.Run()

	shell := exec.Command(args[0])
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Run()
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
