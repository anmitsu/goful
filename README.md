# Goful

[![Go Report Card](https://goreportcard.com/badge/github.com/anmitsu/goful)](https://goreportcard.com/report/github.com/anmitsu/goful)
[![Go Reference](https://pkg.go.dev/badge/github.com/anmitsu/goful.svg)](https://pkg.go.dev/github.com/anmitsu/goful)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/anmitsu/goful/blob/master/LICENSE)

Goful is a CUI file manager written in Go.

* Works on cross-platform such as gnome-terminal and cmd.exe
* Displays multiple windows and workspaces
* A command line to execute using such as bash and tmux
* Provides filtering search, async copy, glob, bulk rename, etc.

![demo](.github/demo.gif)

## Install

### Pre-build releases

See [releases](https://github.com/anmitsu/goful/releases).

### Go version >= 1.16

    $ go install github.com/anmitsu/goful@latest
    ...
    $ goful

### Go version < 1.16

    $ go get github.com/anmitsu/goful
    ...
    $ goful

## Usage

### [Tutorial Demos](.github/demo.md)

key                  | function
---------------------|-----------
`C-n` `down` `j`     | Move cursor down
`C-p` `up` `k`       | Move cursor up
`C-a` `home` `^`     | Move cursor top
`C-e` `end` `$`      | Move cursor bottom
`C-f` `C-i` `right` `l`| Move cursor right
`C-b` `left` `h`     | Move cursor left
`C-d`                | More move cursor down
`C-u`                | More move cursor up
`C-v` `pgdn`         | Page down
`M-v` `pgup`         | Page up
`M-n`                | Scroll down
`M-p`                | Scroll up
`C-h` `backspace` `u`| Change to upper directory
`~`                  | Change to home directory
`\`                  | Change to root directory
`w`                  | Change to neighbor directory
`C-o`                | Create directory window
`C-w`                | Close directory window
`M-f`                | Move next workspace
`M-b`                | Move previous workspace
`M-C-o`              | Create workspace
`M-C-w`              | Close workspace
`space`              | Toggle mark
`C-space`            | Invert mark
`C-l`                | Reload
`C-m` `o`            | Open
`i`                  | Open by pager
`s`                  | Sort
`v`                  | View
`b`                  | Bookmark
`e`                  | Editor
`x`                  | Command
`X`                  | External command
`f` `/`              | Find
`:`                  | Shell
`;`                  | Shell suspend
`n`                  | Make file
`K`                  | Make directory
`c`                  | Copy
`m`                  | Move
`r`                  | Rename
`R`                  | Bulk rename by regexp
`D`                  | Remove
`d`                  | Change directory
`g`                  | Glob
`G`                  | Glob recursive
`C-g` `C-[`          | Cancel
`q` `Q`              | Quit

For more see [main.go](main.go)

## Customize

Goful don't have a config file, instead you can customize by edit `main.go`.

Examples of customizing:

* Change and add keybindings
* Change terminal and shell
* Change file opener (editor, pager and more)
* Adding bookmarks
* Setting colors and looks

Recommend remain original `main.go` and copy to own `main.go` for example:

Go to source directory

    $ cd $GOPATH/src/github.com/anmitsu/goful

Copy original `main.go` to `my/goful` directory

    $ mkdir -p my/goful
    $ cp main.go my/goful
    $ cd my/goful

Install after edit `my/goful/main.go`

    $ go install

## Contributing

[Contributing Guide](.github/CONTRIBUTING.md)
