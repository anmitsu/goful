# Goful

Goful is a CUI file manager implemented by Go.

## Install

    $ go get github.com/anmitsu/goful/cmd/goful
    ...
    $ goful

![goful](<demo.gif>)

## Usage

Goful operates with the keyboard and key binds are based on such as shell and
emacs. 

Main key binds, for more see [cmd/goful/main.go](cmd/goful/main.go)

| key            | bind |
-----------------|-------
| C-n, down      | Move cursor down |
| C-p, up        | Move cursor up |
| C-d            | Move cursor more down |
| C-u            | Move cursor more up |
| C-a            | Move cursor top |
| C-e            | Move cursor bottom |
| C-f, right     | Move next directory |
| C-b, left      | Move previous directory |
| C-h, backspace | Change to upper directory |
| ~              | Change to home directory |
| \              | Change to root directory |
| space          | Mark file at cursor |
| *              | Toggle mark all files |
| C-l            | Reload files |
| f, /           | Find files |
| h              | Start shell mode |
| n              | Create new file |
| k              | Create new directory |
| c              | Copy file |
| m              | Move file |
| r              | Rename file |
| R              | Rename file by regexp |
| D              | Remove file |
| C-g            | Cancel |
| q, Q           | Quit application |

### Menu

### Command line

## Customize

Goful customizes by edit `cmd/goful/main.go` and rebuild.
