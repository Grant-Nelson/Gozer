package ssaConverter

import (
	"slices"
	"strings"
	"testing"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func checkFile(t *testing.T, input, expected string) {
	t.Helper()

	ps, err := astTools.ParsePackage(map[string]string{`t.go`: input}, nil)
	if err != nil {
		t.Errorf(`Failed to parse package: %v`, err)
		return
	}

	mode := ssa.SanityCheckFunctions
	prog, sps := ssautil.Packages([]*packages.Package{ps}, mode)
	prog.Build()
	sp := sps[0]

	p, err := ConvertPackage(ps, sp, nil)
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
		// TODO: Finish
		`}`,
	))
}

func TestConverter_CondFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(cond bool, a, b int) int {`,
		`	x := a`,
		`	if cond {`,
		`		x = b`,
		`	}`,
		`	return x + 1`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		// TODO: Finish
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
		// TODO: Finish
		`}`,
	))
}

func TestConverter_MathFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func foo(x, y float64) float64 {`,
		`	return x*x + y*y`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		// TODO: Finish
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
		// TODO: Finish
		`}`,
	))
}

func TestConverter_OpFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`func scale(a, b int) int {`,
		`	return a * b`,
		`}`,
		``,
		`func foo(a, b int) int {`,
		`	x := 3*scale(a, 2) + b`,
		`	return x`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
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
		// TODO: Finish
		`}`,
	))
}

func TestConverter_RecursiveLocalFunc(t *testing.T) {
	checkFile(t, lines(
		`package t`,
		``,
		`import "errors"`,
		``,
		`func main() {`,
		`	var fn func(string, int) (string, error)`,
		``,
		`	fn = func(phrase string, count int) (string, error) {`,
		`		if count <= 0 {`,
		`			return "", errors.New("invalid count")`,
		`		}`,
		`		if count == 1 {`,
		`			return phrase, nil`,
		`		}`,
		`		result, err := fn(phrase, count-1)`,
		`		if err != nil {`,
		`			return "", err`,
		`		}`,
		`		return phrase + " " + result, nil`,
		`	}`,
		``,
		`	println(fn("buffalo", 5))`,
		`}`,
	), lines(
		`package{`,
		`  Name: t`,
		// TODO: Finish
		`}`,
	))
}
