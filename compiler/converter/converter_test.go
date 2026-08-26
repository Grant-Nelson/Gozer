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
		`    func $.foo () {`,
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
		`    func $.foo (a int, b int) int {`,
		`      block 0 (a int, b int)<initial> {`,
		`        return (ref var a int) + (ref var b int)`,
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
		`    func $.foo (a int) int {`,
		`      block 0 (a int)<initial> {`,
		`        const b untyped int = 10`,
		`        return (ref var a int) + (ref const b untyped int)`,
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
		`    func $.foo (a int, b int) int {`,
		`      block 0 (a int, b int)<initial> {`,
		`        var sum int = 0`,
		`        for ((ref var i int) := 0; (ref var i int) < (ref var b int); (ref var i int)++) {`,
		`          (ref var sum int) += (ref var a int)`,
		`        }`,
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
		`    func $.foo (a int, b int) int {`,
		`      block 0 (a int, b int)<initial> {`,
		`        var sum int = 0`,
		`        for ((ref var i int) := range (ref var b int)) {`,
		`          _ = (ref var i int)`,
		`          (ref var sum int) += (ref var a int)`,
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
		`    func $.foo (a int, b int) int {`,
		`      block 0 (a int, b int)<initial> {`,
		`        if ((ref var a int) > 4) {`,
		`          return 4`,
		`        } else if ((ref var b int) > 4) {`,
		`          return -4`,
		`        } else {`,
		`          return (ref var a int) + (ref var b int)`,
		`        }`,
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
		`    func (*$.Person).String () string {`,
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
		`  Consts:`,
		`    const $.start untyped int = 0`,
		`    const $.step untyped int = 1`,
		`  Vars:`,
		`    var $.current int = ref const $.start untyped int`,
		`    var $.update func() = funcLit {`,
		`      block 0 ()<initial> {`,
		`        (ref var $.current int) += (ref const $.step untyped int)`,
		`      }`,
		`    }`,
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
		`    func $.foo () <atomic> {`,
		`      block 0 ()<initial> {}`,
		`    }`,
		`}`,
	))
}

func TestConverter_BuildInFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`  println("Hello World")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        builtin println("Hello World")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_FmtImport(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`import "fmt"`,
		``,
		`func foo() {`,
		`	fmt.Println("Hello World")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Imports:`,
		`    import package fmt`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        (import package fmt).Println("Hello World")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_Ellipse(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(vs ...int) int {`,
		`	sum := 0`,
		`   for _, v := range vs {`,
		`		sum += v`,
		`	}`,
		`	return sum`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo (vs ...int) int {`,
		`      block 0 (vs []int)<initial> {`,
		`        (ref var sum int) := 0`,
		`        for ((ref var _ int), (ref var v int) := range (ref var vs []int)) {`,
		`          (ref var sum int) += (ref var v int)`,
		`        }`,
		`        return ref var sum int`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_Params(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) bool {`,
		`	return !(a*(b+1) > 10)`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo (a int, b int) bool {`,
		`      block 0 (a int, b int)<initial> {`,
		`        return !(((ref var a int) * ((ref var b int) + 1)) > 10)`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_SliceLit(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(a, b int) []int {`,
		`	v := []int { a, a+1, b }`,
		`   return append([]int{}, v...)`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo (a int, b int) []int {`,
		`      block 0 (a int, b int)<initial> {`,
		`        (ref var v []int) := ([]int {`,
		`          ref var a int`,
		`          (ref var a int) + 1`,
		`          ref var b int`,
		`        })`,
		`        return builtin append([]int {}, ref var v []int...)`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_SliceIndex(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func last(v []int) int {`,
		`	if count := len(v); count > 0 {`,
		`		return v[count-1]`,
		`	}`,
		`   return 0`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.last (v []int) int {`,
		`      block 0 (v []int)<initial> {`,
		`        if ((ref var count int) := (builtin len(ref var v []int)); (ref var count int) > 0) {`,
		`          return (ref var v []int)[(ref var count int) - 1]`,
		`        }`,
		`        return 0`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_MapIndex(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func has(m map[string]int, k string) bool {`,
		`	_, ok := m[k]`,
		`   return ok`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.has (m map[string]int, k string) bool {`,
		`      block 0 (m map[string]int, k string)<initial> {`,
		`        ref var _ int, ref var ok bool := (ref var m map[string]int)[ref var k string]`,
		`        return ref var ok bool`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_StringIndex(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func last(v string) byte {`,
		`	if count := len(v); count > 0 {`,
		`		return v[count-1]`,
		`	}`,
		`   return 0`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.last (v string) byte {`,
		`      block 0 (v string)<initial> {`,
		`        if ((ref var count int) := (builtin len(ref var v string)); (ref var count int) > 0) {`,
		`          return (ref var v string)[(ref var count int) - 1]`,
		`        }`,
		`        return 0`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_GenericSliceIndex(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func last[V any, S ~[]V](s S) V {`,
		`	if count := len(s); count > 0 {`,
		`		return s[count-1]`,
		`	}`,
		`	var zero V`,
		`   return zero`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.last [V any, S ~[]V](s S) V {`,
		`      block 0 (s S)<initial> {`,
		`        if ((ref var count int) := (builtin len(ref var s S)); (ref var count int) > 0) {`,
		`          return (ref var s S)[(ref var count int) - 1]`,
		`        }`,
		`        var zero V`,
		`        return ref var zero V`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_GenericMapIndex(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func has[K comparable, V any, M ~map[K]V](m M, k K) bool {`,
		`	_, ok := m[k]`,
		`   return ok`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.has [K comparable, V any, M ~map[K]V](m M, k K) bool {`,
		`      block 0 (m M, k K)<initial> {`,
		`        ref var _ V, ref var ok bool := (ref var m M)[ref var k K]`,
		`        return ref var ok bool`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_GotoLabel(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	print("start")`,
		`	goto End`,
		`	print("unreachable")`,
		`	End:`,
		`   print("done")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        builtin print("start")`,
		`        goto End`,
		`        builtin print("unreachable")`,
		`        End:`,
		`        builtin print("done")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_Branch(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	for i := 0; i < 10; i++ {`,
		`		if i % 2 == 0 {`,
		`			println("even")`,
		`			continue`,
		`		}`,
		`		println("odd")`,
		`		break`,
		`	}`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        for ((ref var i int) := 0; (ref var i int) < 10; (ref var i int)++) {`,
		`          if (((ref var i int) % 2) == 0) {`,
		`            builtin println("even")`,
		`            continue`,
		`          }`,
		`          builtin println("odd")`,
		`          break`,
		`        }`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_LabelBranch(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	Loop: for i := 0; i < 10; i++ {`,
		`		if i % 2 == 0 {`,
		`			println("even")`,
		`			continue Loop`,
		`		}`,
		`		println("odd")`,
		`		break Loop`,
		`	}`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        Loop:`,
		`        for ((ref var i int) := 0; (ref var i int) < 10; (ref var i int)++) {`,
		`          if (((ref var i int) % 2) == 0) {`,
		`            builtin println("even")`,
		`            continue Loop`,
		`          }`,
		`          builtin println("odd")`,
		`          break Loop`,
		`        }`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_SimpleCall(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	bar()`,
		`}`,
		``,
		`func bar() {`,
		`	println("Boop")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        $.bar()`,
		`      }`,
		`    }`,
		`    func $.bar () {`,
		`      block 0 ()<initial> {`,
		`        builtin println("Boop")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_GoFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	go bar()`,
		`}`,
		``,
		`func bar() {`,
		`	println("Boop")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        go $.bar()`,
		`      }`,
		`    }`,
		`    func $.bar () {`,
		`      block 0 ()<initial> {`,
		`        builtin println("Boop")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_GenericCall(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	bar[int](42)`,
		`	bar("Hello")`,
		`}`,
		``,
		`func bar[T any](t T) {`,
		`	print(">", t, "<\n")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        $.bar[int](42)`,
		`        $.bar[string]("Hello")`,
		`      }`,
		`    }`,
		`    func $.bar [T any](t T) {`,
		`      block 0 (t T)<initial> {`,
		`        builtin print(">", ref var t T, "<\n")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_MultiGenericCall(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo() {`,
		`	m := map[string]int{`,
		`		"banana": 6,`,
		`		"cat":    3,`,
		`	}`,
		`	bar[string, int](m)`,
		`}`,
		``,
		`func bar[K comparable, V any, M ~map[K]V](t M) {`,
		`	print(">", t, "<\n")`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Funcs:`,
		`    func $.foo () {`,
		`      block 0 ()<initial> {`,
		`        (ref var m map[string]int) := (map[string]int {`,
		`          "banana": 6`,
		`          "cat": 3`,
		`        })`,
		`        $.bar[string, int, map[string]int](ref var m map[string]int)`,
		`      }`,
		`    }`,
		`    func $.bar [K comparable, V any, M ~map[K]V](t M) {`,
		`      block 0 (t M)<initial> {`,
		`        builtin print(">", ref var t M, "<\n")`,
		`      }`,
		`    }`,
		`}`,
	))
}

func TestConverter_Struct(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`type Foo struct { name string }`,
		``,
		`func (f *Foo) String() string { return f.name }`,
		``,
		`func main() {`,
		`	f1 := &Foo{name: "Bob"}`,
		`	println(f1.String())`,
		``,
		`	f2 := &Foo{"Bob"}`,
		`   println(f1 == f2)`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		`  Types:`,
		`    type $.Foo struct{name string}`,
		`  Funcs:`,
		`    func (*$.Foo).String () string {`,
		`      block 0 ()<initial> {`,
		`        return (ref var f *$.Foo).name`,
		`      }`,
		`    }`,
		`    func $.main () {`,
		`      block 0 ()<initial> {`,
		`        (ref var f1 *$.Foo) := (&($.Foo {`,
		`          ref field name string: "Bob"`,
		`        }))`,
		`        builtin println((ref var f1 *$.Foo).String())`,
		`        (ref var f2 *$.Foo) := (&($.Foo {`,
		`          "Bob"`,
		`        }))`,
		`        builtin println((ref var f1 *$.Foo) == (ref var f2 *$.Foo))`,
		`      }`,
		`    }`,
		`}`,
	))
}
