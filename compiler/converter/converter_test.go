package converter

import (
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func checkFile(t *testing.T, input, exp string) {
	t.Helper()
	fSet := token.NewFileSet()
	filename := `t.go`
	ps, err := packages.Load(&packages.Config{
		Mode: packages.LoadFiles | packages.LoadSyntax,
		Fset: fSet,
		Overlay: map[string][]byte{
			filename: []byte(input),
		},
	}, filename)
	if err != nil {
		t.Errorf(`Failed to parse input expression: %v`, err)
	}
	if len(ps) != 1 {
		t.Errorf(`Expected there to be one package but there was %d`, len(ps))
	}
	p := ps[0]
	if len(p.Syntax) != 1 {
		t.Errorf(`Expected there to be one file in the package but there was %d`, len(p.Syntax))
	}
	f := p.Syntax[0]
	c := &Converter{
		Info:    p.TypesInfo,
		FileSet: fSet,
	}
	for _, d := range f.Decls {

		// TODO: FINISH
		c.FromNode(d)
	}
}
