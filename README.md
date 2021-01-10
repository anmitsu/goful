# Goful

Goful is a CUI file manager implemented by Go.

* Works on cross-platform such as gnome-terminal in Ubuntu and cmd.exe in
  Windows.
* Multi window and workspace to display directory contents.
* A command line to execute external commands such as bash and tmux.
* Provides file search, glob, copy, bulk rename, etc.

## Install

    $ go get github.com/anmitsu/goful/cmd/goful
    ...
    $ goful

![goful](<demo/readme_top.gif>)

## Usage

| key            | function |
-----------------|-------
| C-n, down      | Move cursor down |
| C-p, up        | Move cursor up |
| C-a            | Move cursor top |
| C-e            | Move cursor bottom |
| C-f, right     | Move cursor right |
| C-b, left      | Move cursor left |
| C-h, backspace | Change to upper directory |
| ~              | Change to home directory |
| \              | Change to root directory |
| space          | Mark file on cursor |
| *              | Toggle mark all files |
| C-l            | Reload files |
| C-m            | Open |
| l              | Layout menu |
| s              | Sort menu |
| e              | Editor menu |
| M-x            | Command menu |
| M-f            | Move next workspace |
| M-b            | Move previous workspace |
| f, /           | Find files |
| h              | Start spawn mode |
| H              | Start shell mode |
| n              | Create new file |
| k              | Create new directory |
| c              | Copy file |
| m              | Move file |
| r              | Rename file |
| R              | Bulk rename by regexp |
| D              | Remove file |
| C-g            | Cancel |
| q, Q           | Quit application |

### For more see [cmd/goful/main.go](cmd/goful/main.go)

## Customize

Goful customizes by edit `cmd/goful/main.go` and rebuild.

For example, install your customized binary to `$GOPATH/bin`.

    Copy original main.go to my/goful directory
    $ cd $GOPATH/src/github.com/anmitsu/goful/cmd/goful
    $ mkdir -p my/goful
    $ cp main.go my/goful
    $ cd my/goful
    
    After edited my/goful/main.go
    $ go build
    $ go install

