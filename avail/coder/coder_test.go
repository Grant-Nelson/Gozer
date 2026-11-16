package coder

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Coder(t *testing.T) {
	c := New(`t.go`, "\t")
	c.Write(`func `, `main`, `() `)
	c.WriteLn(`{`)
	c.Indent()
	c.WriteF("a := %d\n", 42)
	c.WriteLn(`print(`, `a`, `)`)
	c.Dedent()
	c.WriteLn(`}`)

	check(t, c,
		"func main() {\n",
		"\ta := 42\n",
		"\tprint( a )\n",
		"}\n")
}

func check(t *testing.T, c *Coder, expLines ...string) {
	got := c.String()
	gotLines := slices.Collect(strings.Lines(got))
	exp := strings.Join(expLines, ``)
	if diff := cmp.Diff(expLines, gotLines); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ndiff: %s\ngot: %q\nexp: %q",
			diff, got, exp)

	}
}
