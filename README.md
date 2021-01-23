# Goful

[![Go Report Card](https://goreportcard.com/badge/github.com/anmitsu/goful)](https://goreportcard.com/report/github.com/anmitsu/goful)
[![Go Reference](https://pkg.go.dev/badge/github.com/anmitsu/goful.svg)](https://pkg.go.dev/github.com/anmitsu/goful)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/anmitsu/goful/blob/master/LICENSE)

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

![demo](.github/demo.gif)

### For more [demos](.github/demo.md)

## Usage

| key                  | function |
-----------------------|-----------
| `C-n`, `down`, `j`   | Move cursor down |
| `C-p`, `up`, `k`     | Move cursor up |
| `C-a`, `home`, `^`   | Move cursor top |
| `C-e`, `end`, `$`    | Move cursor bottom |
| `C-f`, `C-i`, `right`, `l`| Move cursor right |
| `C-b`, `left`, `h`   | Move cursor left |
| `C-d`                | More move cursor down |
| `C-u`                | More move cursor up |
| `M-v`, `pgdn`        | Page down |
| `M-p`, `pgup`        | Page up |
| `M-n`                | Scroll down |
| `M-p`                | Scroll up |
| `C-h`, `backspace`   | Change to upper directory |
| `~`                  | Change to home directory |
| `\`                  | Change to root directory |
| `w`                  | Change to neighbor directory |
| `C-o`                | New directory window |
| `C-w`                | Close directory window |
| `M-f`                | Move next workspace |
| `M-b`                | Move previous workspace |
| `M-C-o`              | New workspace |
| `M-C-w`              | Close workspace |
| `space`              | Mark file on cursor |
| `C-space`            | Toggle mark all files |
| `C-l`                | Reload files |
| `C-m`                | Open |
| `s`                  | Sort menu |
| `V`                  | View menu |
| `b`                  | Bookmark menu |
| `e`                  | Editor menu |
| `M-x`                | Command menu |
| `x`                  | External command menu |
| `A`                  | Archive menu |
| `f`, `/`             | Find files |
| `:`                  | Start shell mode |
| `;`                  | Start shell suspend mode |
| `n`                  | Create new file |
| `K`                  | Create new directory |
| `c`                  | Copy file |
| `m`                  | Move file |
| `r`                  | Rename file |
| `R`                  | Bulk rename by regexp |
| `D`                  | Remove file |
| `d`                  | Change directory |
| `g`                  | Glob file |
| `G`                  | Glob directory |
| `C-g`, `esc`         | Cancel |
| `q`, `Q`             | Quit application |

### For more see [cmd/goful/main.go](cmd/goful/main.go)

## Customize

Goful customizes by edit `cmd/goful/main.go` and rebuild.

For example, install your customized binary to `$GOPATH/bin`.

    Clone source code
    $ git clone https://github.com/anmitsu/goful

    Copy original main.go to my/goful directory
    $ cd goful/cmd/goful
    $ mkdir -p my/goful
    $ cp main.go my/goful
    $ cd my/goful
    
    After edited my/goful/main.go
    $ go install

## Contributing

[Contributing Guide](.github/CONTRIBUTING.md)
