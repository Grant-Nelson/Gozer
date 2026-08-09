package converter

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

var wd = sync.OnceValue(func() string {
	dir, err := filepath.Abs(`./`)
	if err != nil {
		panic(fmt.Errorf(`Error getting working directory: %w`, err))
	}
	return dir
})

func workingPath(path string) string {
	return filepath.Join(wd(), path)
}

func parsePackage(inputFiles, extraFiles map[string]string) (*packages.Package, error) {
	patterns := make([]string, 0, len(inputFiles))
	overlay := make(map[string][]byte, len(inputFiles)+len(extraFiles))
	for name, src := range inputFiles {
		patterns = append(patterns, name)
		overlay[workingPath(name)] = []byte(src)
	}
	slices.Sort(patterns)
	for name, src := range extraFiles {
		overlay[workingPath(name)] = []byte(src)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Dir:     wd(),
		Overlay: overlay,
	}, patterns...)
	if err != nil {
		return nil, err
	}

	pkgErrors := []error{}
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			pkgErrors = append(pkgErrors, err)
		}
	})
	if len(pkgErrors) > 0 {
		return nil, errors.Join(pkgErrors...)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf(`Expected exactly one root package but got %d`, len(pkgs))
	}
	return pkgs[0], nil
}

func checkFile(t *testing.T, input, expected string) {
	t.Helper()

	ps, err := parsePackage(map[string]string{`t.go`: input}, nil)
	if err != nil {
		t.Errorf(`Failed to parse input expression: %v`, err)
		return
	}
	if len(ps.Syntax) != 1 {
		t.Errorf(`Expected there to be one file in the package but there was %d`, len(ps.Syntax))
		return
	}
	f := ps.Syntax[0]
	c := &Converter{
		Info:    ps.TypesInfo,
		FileSet: ps.Fset,
	}

	result := []string{}
	for _, d := range f.Decls {

		// TODO: FINISH by making from file and from package
		n := c.FromNode(d)
		result = append(result, n.String())
	}

	exp := slices.Collect(strings.Lines(expected))
	got := slices.Collect(strings.Lines(strings.Join(result, "\n")))
	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		t.Errorf(`Unexpected results:\n%s`, diff)
	}
}

func TestConverter_EmptyFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {}`,
	), lines(
		`func foo {`,
		`  block 0 ()<initial> {}`,
		`}`,
	))
}

func TestConverter_SimpleFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) int {`,
		`	return a + b`,
		`}`,
	), lines(
		`func foo {`,
		`  block 0 (a int, b int)<initial> {`,
		`    return a + b`,
		`  }`,
		`}`,
	))
}

func TestConverter_FuncConstant(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a int) int {`,
		`	const b = 10`,
		`	return a + b`,
		`}`,
	), lines(
		`func foo {`,
		`  block 0 (a int)<initial> {`,
		`    const b (def)untyped int = 10`,
		`    return a + b`,
		`  }`,
		`}`,
	))
}

func TestConverter_ForLoop(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) int {`,
		`	var sum = 0`,
		`	for i := 0; i < b; i++ {`,
		`		sum += a`,
		`	}`,
		`	return sum`,
		`}`,
	), lines(
		`func foo {`,
		`  block 0 (a int, b int)<initial> {`,
		`    var sum (def)int = 0`,
		`    for i := 0; i < b; ++i {`,
		`      sum += a`,
		`    }`,
		`    return sum`,
		`  }`,
		`}`,
	))
}
