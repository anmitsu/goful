package goful

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExpandMacro(t *testing.T) {
	g := New("")
	g.Workspace().ReloadAll() // in home directory

	home, _ := os.UserHomeDir()
	macros := []struct {
		in  string
		out string
	}{
		{`%f`, `".."`},
		{`%F`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%x`, `"."`},
		{`%X`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%m`, `".."`},
		{`%M`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%d`, fmt.Sprintf(`"%s"`, filepath.Base(home))},
		{`%D`, fmt.Sprintf(`"%s"`, home)},
		{`%d2`, fmt.Sprintf(`"%s"`, filepath.Base(home))},
		{`%D2`, fmt.Sprintf(`"%s"`, home)},
		{`%~f`, ".."},
		{`%~F`, filepath.Dir(home)},
		{`%~x`, "."},
		{`%~X`, filepath.Dir(home)},
		{`%~m`, ".."},
		{`%~M`, filepath.Dir(home)},
		{`%~d`, filepath.Base(home)},
		{`%~D`, home},
		{`%~d2`, filepath.Base(home)},
		{`%~D2`, home},
		{`%%%f`, `%%".."`},
		{`%%%~f`, `%%..`},
		{`%~~f`, `%~~f`},
		{`\%f%f`, `%f".."`},
		{`\%~f%~f`, `%~f..`},
		{`%\f%f`, `%f".."`},
		{`%\~f%~f`, `%~f..`},
		{"%AA%ff", `%AA".."f`},
		{"%~A~A%~ff", `%~A~A..f`},
	}

	for _, macro := range macros {
		ret, _ := g.expandMacro(macro.in)
		if ret != macro.out {
			t.Errorf("%s -> %s result %s\n", macro.in, macro.out, ret)
		}
	}
}
