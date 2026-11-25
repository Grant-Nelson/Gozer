package ssaConverter

import (
	"slices"
	"strings"
	"testing"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/ssa"
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func checkFile(t *testing.T, input, expected string) {
	t.Helper()

	ps, err := astTools.ParsePackage(map[string]string{`t.go`: input}, nil)
	if err != nil {
		t.Errorf(`Failed to parse package: %v`, err)
		return
	}

	// TODO: Make sure this mode is reasonable.
	mode := ssa.SanityCheckFunctions |
		ssa.GlobalDebug |
		ssa.NaiveForm

	prog := ssa.NewProgram(ps.Fset, mode)
	sp := prog.CreatePackage(ps.Types, ps.Syntax, ps.TypesInfo, true)
	sp.Build()

	p, err := ConvertPackage(sp, nil)
	if err != nil {
		t.Errorf(`Failed to convert package: %v`, err)
		return
	}

	result := p.String()
	result = strings.ReplaceAll(result, `command-line-arguments.`, `$.`)
	exp := slices.Collect(strings.Lines(expected))
	got := slices.Collect(strings.Lines(result))
	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		t.Errorf("Unexpected results:\n%s", diff)
	}
}

func TestConverter_EmptyFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {}`,
		`    }`,
		`}`,
	))
}
