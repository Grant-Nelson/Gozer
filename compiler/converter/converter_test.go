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
		t.Errorf(`Failed to parse package: %v`, err)
		return
	}

	p, err := ConvertPackage(ps, nil)
	if err != nil {
		t.Errorf(`Failed to convert package: %v`, err)
		return

	}

	result := p.String()

	exp := slices.Collect(strings.Lines(expected))
	got := slices.Collect(strings.Lines(result))
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
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 ()<initial> {}`,
		`    }`,
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
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 (a int, b int)<initial> {`,
		`        return a + b`,
		`      }`,
		`    }`,
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
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 (a int)<initial> {`,
		`        const b (def)untyped int = 10`,
		`        return a + b`,
		`      }`,
		`    }`,
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
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 (a int, b int)<initial> {`,
		`        var sum (def)int = 0`,
		`        for (i := 0; i < b; ++i)`,
		`          sum += a`,
		`        return sum`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_ForRange(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) int {`,
		`	var sum = 0`,
		`	for i := range b {`,
		`		_ = i`,
		`		sum += a`,
		`	}`,
		`	return sum`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 (a int, b int)<initial> {`,
		`        var sum (def)int = 0`,
		`        for (i := range b) {`,
		`          _ = i`,
		`          sum += a`,
		`        }`,
		`        return sum`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_IfElse(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) int {`,
		`	if a > 4 {`,
		`		return 4`,
		`	} else if b > 4 {`,
		`		return -4`,
		`	} else {`,
		`		return a + b`,
		`	}`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func foo {`,
		`      block 0 (a int, b int)<initial> {`,
		`        if (a > 4)`,
		`          return 4`,
		`        else if (b > 4)`,
		`          return -4`,
		`        else`,
		`          return a + b`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_SimpleStructAndMethod(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`type Person struct{`,
		`	Name string`,
		`	Age  int`,
		`}`,
		``,
		`func (p *Person) String() string {`,
		`	return p.Name`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Types:`,
		`    Person struct{Name string; Age int}`,
		`  Funcs:`,
		`    func String {`,
		`      block 0 ()<initial> {`,
		`        return p.Name`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_Globals(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`const start = 0`,
		`var current = start`,
		`const step = 1`,
		``,
		`var update = func() {`,
		`	current += step`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Consts: {`,
		`    const start (def)untyped int = 0`,
		`    const step (def)untyped int = 1`,
		`  }`,
		`  Vars: {`,
		`    var current (def)int = start`,
		`    var update (def)func() = funcLit: func unnamed {`,
		`      block 0 ()<initial> {`,
		`        current += step`,
		`      }`,
		`    }`,
		`  }`,
		`}`,
	))
}
