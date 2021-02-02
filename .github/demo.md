# Goful Demos

## Copy Progress

File/directory copy (default `c`) and move (default `m`).

For example, copy to a neighbor directory from selected and mark files:

![demo_copy](demo_copy.gif)

During copy processing draws progress percent, gauge, bps and estimated time of
arrival.

## Bulk Rename

Bulk renaming (default `R`) for mark files:

![demo_bulk](demo_bulk.gif)

Rename by the regexp pattern.  Input to the cmdline like the vim substituting
style (regexp/replaced).  Display and confirm matched files and replaced names
before rename.

## Glob

Glob is matched by wild card pattern in the current directory (default `g` and
recursive `G`).

![demo_glob](demo_glob.gif)

## Layout

Directory windows position are allocated by layouts of tile, tile-top,
tile-bottom, one-row, one-column and fullscreen.

View menu (default `v`), run layout menu and select layout:

![demo_layout](demo_layout.gif)

## Execute Terminal and Shell

Shell mode (default `:` and suspended `;`) runs a terminal and execute shell
such as bash and tmux.  The cmdline completion (file names and commands in
$PATH) is available (default `C-i` that means `tab`).

For example, spawns commands by bash in a gnome-terminal new tab:

![demo_shell](demo_shell.gif)

The terminal immediately doesn't close when command finished because check
outputs.

If goful is running in tmux, it creates a new window and executes the command.

## Expand Macro

macro        | expanded string
-------------|------------------
`%f` `%F`   | File name/path on cursor
`%x` `%X`   | File name/path with extension excluded on cursor
`%m` `%M`   | Marked file names/paths joined by spaces
`%d` `%D`   | Directory name/path on cursor
`%d2` `%D2` | Neighbor directory name/path
`%~f` ...   | Expand by non quote
`%&`        | Flag to run command in background

The macro is useful if do not want to specify a file name when run the shell.

Macros starts with `%` are expanded surrounded by quote, and those starts with
`%~` are expanded by non quote.  The `%~` mainly uses to for cmd.exe.

Use `%&` when background execute the shell such as GUI apps launching.

![demo_macro](demo_macro.gif)

<!-- demo size 120x35 -->