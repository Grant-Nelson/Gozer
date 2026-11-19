package coder

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Coder(t *testing.T) {
	c := New(`t.go`, "\t", `// `)
	checkLoc(t, c, `t.go:1:1:0`)
	c.WriteComment(`Entry`, `Point`)
	checkLoc(t, c, `t.go:2:1:15`)
	c.Write(`func `, `main`, `() `)
	checkLoc(t, c, `t.go:2:13:27`)
	c.WriteLn(`{`)
	checkLoc(t, c, `t.go:3:1:29`)
	c.Indent()
	checkLoc(t, c, `t.go:3:1:29`)
	c.WriteF("a := %d\n", 42)
	checkLoc(t, c, `t.go:4:1:38`)
	c.WriteLn(`print(`, `a`, `)`)
	checkLoc(t, c, `t.go:5:1:50`)
	c.Dedent()
	c.Indent()
	c.Indent()
	checkLoc(t, c, `t.go:5:1:50`)
	c.Dedent()
	c.Dedent()
	checkLoc(t, c, `t.go:5:1:50`)
	c.WriteLn(`}`)
	checkLoc(t, c, `t.go:6:1:52`)

	checkString(t, c,
		"// Entry Point\n",
		"func main() {\n",
		"\ta := 42\n",
		"\tprint( a )\n",
		"}\n")
}

func checkLoc(t *testing.T, c *Coder, exp string) {
	got := fmt.Sprintf(`%s:%d:%d:%d`, c.FilePath(), c.Line(), c.Column(), c.Offset())
	if got != exp {
		t.Errorf("Unexpected Coder results:\ngot: %q\nexp: %q",
			got, exp)
	}
}

func checkString(t *testing.T, c *Coder, expLines ...string) {
	got := c.String()
	gotLines := slices.Collect(strings.Lines(got))
	exp := strings.Join(expLines, ``)
	if diff := cmp.Diff(expLines, gotLines); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected Coder results:\ndiff: %s\ngot: %q\nexp: %q",
			diff, got, exp)
	}
}
