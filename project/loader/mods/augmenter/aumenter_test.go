package augmenter

import (
	"fmt"
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/remapper"
)

func Test_Add_Import(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`,
			``,
			`type Foo struct{}`),
		augSrc: lines(
			`package foo`,
			``,
			`//gozer:add`,
			`import "fmt"`),
		expSrc: lines(
			`package foo`,
			``,
			`import "fmt"`,
			``,
			`//line original/orig.go:3:1`,
			`type Foo struct{}`),
	})
}

func Test_Add_Import_Specs(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`,
			``,
			`type Foo struct{}`,
			``,
			`type Bar struct{}`),
		augSrc: lines(
			`package foo`,
			``,
			`import (`,
			`	"fmt" //gozer:add`,
			`	"time" //gozer:ignore`,
			`	gg "log" //gozer:add`,
			`)`),
		expSrc: lines(
			`package foo`,
			``,
			`import (`,
			`	"fmt"`,
			`	gg "log"`,
			`)`,
			``,
			`//line original/orig.go:3:1`,
			`type Foo struct{}`,
			``,
			`type Bar struct{}`),
	})
}

func Test_Add_WholeType(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`,
			``,
			`// Foo already exists.`,
			`type Foo struct{}`),
		augSrc: lines(
			`package foo`,
			``,
			`// Bar is being added.`,
			`//gozer:add`,
			`type Bar struct{}`),
		expSrc: lines(
			`package foo`,
			``,
			`// Foo already exists.`,
			`type Foo struct{}`,
			``,
			`// Bar is being added.`,
			`//`,
			`//line base/aug.go:5:1`,
			`type Bar struct{}`),
	})
}

func Test_Add_Func(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`,
			``,
			`import "fmt"`,
			``,
			`// X marks the spot.`,
			`type X struct{}`,
			``,
			`// Foo does stuff.`,
			`func (x *X) Foo(y int, z string) {`,
			`	fmt.Printf("%d, %s\n", y, z)`,
			`}`),
		augSrc: lines(
			`package foo`,
			``,
			`import "fmt"`,
			``,
			`// Bar is being added.`,
			`//gozer:add`,
			`func (x *X) Bar(y int, z string) {`,
			`	fmt.Printf("%s, %d\n", z, y)`,
			`}`),
		expSrc: lines(
			`package foo`,
			``,
			`import "fmt"`,
			``,
			`// X marks the spot.`,
			`type X struct{}`,
			``,
			`// Foo does stuff.`,
			`func (x *X) Foo(y int, z string) {`,
			`	fmt.Printf("%d, %s\n", y, z)`,
			`}`,
			``,
			`// Bar is being added.`,
			`//`,
			`//line base/aug.go:7:1`,
			`func (x *X) Bar(y int, z string) {`,
			`	fmt.Printf("%s, %d\n", z, y)`,
			`}`),
	})
}

func Test_Add_VarAndConst(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`,
			``,
			`var A int`,
			``,
			`const B = "Hello"`),
		augSrc: lines(
			`package foo`,
			``,
			`//gozer:add`,
			`var (`,
			`	C, D int`,
			`)`,
			``,
			`var X int //gozer:add`,
			``,
			`const (`,
			`	Y int = 42 //gozer:add`,
			`	Z int = 10 //gozer:ignore`,
			`	//gozer:add`,
			`	W string = "World"`,
			`)`),
		expSrc: lines(
			`package foo`,
			``),
	})
}

type augTest struct {
	origSrc  string
	augSrc   string
	expSrc   string
	expErr   string
	errLimit int
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}

func runAugTest(t testing.TB, test augTest) {
	t.Helper()

	tempFileSet := token.NewFileSet()
	remap := remapper.New(tempFileSet, artifacts.DefaultFileParser)
	fm, err := remap.Parse(tempFileSet, `original/orig.go`, []byte(test.origSrc))
	if err != nil {
		t.Errorf(`failed to load origin file: %v`, err)
		return
	}
	f := artifacts.NewFile(tempFileSet, fm)

	test.errLimit = max(test.errLimit, 1)
	errGroup := faults.NewGroup(test.errLimit)
	a := New(nil, PathRebase(`original`, `base`), remap)

	// Create an augmenter for a package then add the aug file to it
	pkg := f.Package
	ap := a.addPackage(pkg, errGroup)
	if err := ap.AddFile(nil, `base/aug.go`, []byte(test.augSrc), errGroup, remap); err != nil {
		checkErr(t, `load augment file`, test, err)
		return
	}

	// Perform the augmentation on the file
	con, err := a.Modify(f, errGroup)
	if err != nil {
		checkErr(t, `modify file`, test, err)
		return
	}
	if !con {
		t.Errorf(`expected Modify to return continue but it did not`)
		return
	}

	if err := a.LoadDone(errGroup); err != nil {
		checkErr(t, `load done`, test, err)
		return
	}

	finalFileSet := token.NewFileSet()
	if err := remap.Remap(fm, finalFileSet, errGroup); err != nil {
		checkErr(t, `remap`, test, err)
		return
	}

	if err := errGroup.Wrap(); err != nil {
		checkErr(t, `accumulated error`, test, err)
		return
	}

	buf := &strings.Builder{}
	if err := remap.Write(fm, buf); err != nil {
		t.Errorf(`failed to write result: %v`, err)
		return
	}

	got := buf.String()
	if diff := cmp.Diff(strings.Split(got, "\n"), strings.Split(test.expSrc, "\n")); len(diff) > 0 {
		t.Logf("Got:\n%s\n", got)
		t.Errorf("resulting source didn't match expected:\n%s", diff)
		return
	}
}

func checkErr(t testing.TB, prefix string, test augTest, err error) {
	t.Helper()
	if len(test.expErr) > 0 {
		errStr := fmt.Sprintf(`in %s: %v`, prefix, err)
		if diff := cmp.Diff(strings.Split(errStr, "\n"), strings.Split(test.expErr, "\n")); len(diff) > 0 {
			t.Errorf("resulting error didn't match expected error:\n%s", diff)
		}
		return
	}
	t.Errorf(`failed in %s: %v`, prefix, err)
}
