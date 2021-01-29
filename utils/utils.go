package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mattn/go-runewidth"
)

// ExpandPath expands path beginning of ~  to the home directory.
// Not use for file operations because unexpected behavior when exist a file beginning of ~
func ExpandPath(name string) string {
	if name == "" {
		return ""
	}
	if name[:1] == "~" {
		home, _ := os.UserHomeDir()
		return strings.Replace(name, "~", home, 1)
	}
	return name
}

// AbbrPath abbreviates path beginning of home directory to ~.
func AbbrPath(name string) string {
	home, _ := os.UserHomeDir()
	lenhome := len(home)
	if len(name) >= lenhome && name[:lenhome] == home {
		return "~" + name[lenhome:]
	}
	return name
}

// ShortenPath returns a shorten path to be shorter than width.
func ShortenPath(path string, width int) string {
	if width < runewidth.StringWidth(path) {
		root := filepath.VolumeName(path)
		names := strings.Split(path, string(filepath.Separator))
		for i, name := range names[:len(names)-1] {
			if name == root {
				if name == "" {
					names[i] = string(filepath.Separator)
				}
				continue
			}
			for _, r := range name {
				names[i] = string(r)
				break
			}
			path = filepath.Join(names...)
			if runewidth.StringWidth(path) <= width {
				break
			}
		}
	}
	return path
}

// RemoveExt removes extension from the name.
func RemoveExt(name string) string {
	if ext := filepath.Ext(name); ext != name {
		return name[:len(name)-len(ext)]
	}
	return name
}

// SplitWithSep splits string with separator.
func SplitWithSep(s, sep string) []string {
	n := strings.Count(s, sep)*2 + 1
	a := make([]string, n)
	c := sep[0]
	start := 0
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
			a[na] = s[start:i]
			na++
			a[na] = s[i : i+len(sep)]
			na++

			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	a[na] = s[start:]
	return a[0 : na+1]
}

// Quote encloses string double quotes and escapes by backslash if this string contains double quotes.
func Quote(s string) string {
	return fmt.Sprintf(`"%s"`, strings.Replace(s, `"`, `\"`, -1))
}

// FormatSize returns formated to SI prefix unit.
func FormatSize(n int64) string {
	const (
		Tb = 1024 * 1024 * 1024 * 1024
		Gb = 1024 * 1024 * 1024
		Mb = 1024 * 1024
		kb = 1024
	)
	if n > Tb {
		return fmt.Sprintf("%.1fT", float64(n)/Tb)
	} else if n > Gb {
		return fmt.Sprintf("%.1fG", float64(n)/Gb)
	} else if n > Mb {
		return fmt.Sprintf("%.1fM", float64(n)/Mb)
	} else if n > kb {
		return fmt.Sprintf("%.1fk", float64(n)/kb)
	} else {
		return fmt.Sprintf("%d", n)
	}
}

func searchPath(results map[string]bool, path string) (map[string]bool, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return results, err
	}
	for {
		names, err := dir.Readdirnames(100)
		if err == io.EOF {
			break
		} else if err != nil {
			return results, err
		}
		for _, name := range names {
			results[name] = true
		}
	}
	return results, nil
}

// SearchCommands returns map with key is command name in $PATH and by bash compgen -c.
func SearchCommands() (map[string]bool, error) {
	results := map[string]bool{}
	for _, path := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		if results, err := searchPath(results, path); err != nil {
			if !os.IsNotExist(err) {
				return results, err
			}
		}
	}

	if runtime.GOOS == "windows" {
		for key := range results {
			if filepath.Ext(key) == ".exe" {
				results[RemoveExt(key)] = true
			}
		}
		return results, nil
	}
	cmd := exec.Command("/bin/bash", "-c", "compgen -c")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	r := bufio.NewReader(out)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		results[string(line)] = true
	}
	return results, nil
}
