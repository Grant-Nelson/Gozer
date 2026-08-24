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
	result = strings.ReplaceAll(result, `command-line-arguments.`, `$.`)
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
		`        return (ref var a int)+(ref var b int)`,
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
		`        const b untyped int = 10`,
		`        return (ref var a int)+(ref const b untyped int)`,
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
		`        var sum int = 0`,
		`        for ((ref var i int):=(0); (ref var i int)<(ref var b int); ++(ref var i int))`,
		`          (ref var sum int)+=(ref var a int)`,
		`        return ref var sum int`,
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
		`        var sum int = 0`,
		`        for (ref var i int := range ref var b int) {`,
		`          _=(ref var i int)`,
		`          (ref var sum int)+=(ref var a int)`,
		`        }`,
		`        return ref var sum int`,
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
		`        if ((ref var a int)>(4))`,
		`          return 4`,
		`        else if ((ref var b int)>(4))`,
		`          return -(4)`,
		`        else`,
		`          return (ref var a int)+(ref var b int)`,
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
		`    type $.Person struct{Name string; Age int}`,
		`  Funcs:`,
		`    func (*$.Person).String {`,
		`      block 0 ()<initial> {`,
		`        return (ref var p *$.Person).Name`,
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
		`    const $.start untyped int = 0`,
		`    const $.step untyped int = 1`,
		`  }`,
		`  Vars: {`,
		`    var $.current int = ref const $.start untyped int`,
		`    var $.update func() = funcLit {`,
		`      block 0 ()<initial> {`,
		`        (ref var $.current int)+=(ref const $.step untyped int)`,
		`      }`,
		`    }`,
		`  }`,
		`}`,
	))
}

func TestConverter_AtomicFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`//gozer:atomic`,
		`func foo() {}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo {`,
		`      atomic`,
		`      block 0 ()<initial> {}`,
		`    }`,
		`}`,
	))
}
