package augmenter

import (
	"fmt"
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

func TestAddingType(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`),
		augSrc: lines(
			`package foo`,
			``,
			`// Foo is being added.`,
			`//gozer:add`,
			`type Foo struct{}`),
		expSrc: lines(
			`//line original/orig.go:1`,
			`package foo`,
			``,
			`type Foo struct{}`),
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

func runAugTest(t *testing.T, test augTest) {
	t.Helper()

	fileSet := token.NewFileSet()
	fm, err := file.Load(fileSet, `original/orig.go`, []byte(test.origSrc))
	if err != nil {
		t.Errorf(`failed to load origin file: %v`, err)
		return
	}

	test.errLimit = max(test.errLimit, 1)
	pkgPath := `test/path`
	errGroup := faults.NewGroup(test.errLimit)
	a := New(nil, `base`, pkgPath)
	a.reset()
	if err := a.AddFile(`aug.go`, []byte(test.augSrc), errGroup); err != nil {
		checkErr(t, `load augment file`, test, err)
		return
	}

	if err := a.Modify(fm, errGroup); err != nil {
		checkErr(t, `modify file`, test, err)
		return
	}

	pkg := &mods.Package{Name: `test`, Path: pkgPath}
	if err := a.PackageDone(pkg, errGroup); err != nil {
		checkErr(t, `finish package`, test, err)
		return
	}

	if err := errGroup.Wrap(); err != nil {
		checkErr(t, `accumulated error`, test, err)
		return
	}

	buf := &strings.Builder{}
	if err := fm.Write(buf); err != nil {
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

func checkErr(t *testing.T, prefix string, test augTest, err error) {
	if len(test.expErr) > 0 {
		errStr := fmt.Sprintf(`in %s: %v`, prefix, err)
		if diff := cmp.Diff(strings.Split(errStr, "\n"), strings.Split(test.expErr, "\n")); len(diff) > 0 {
			t.Errorf("resulting error didn't match expected error:\n%s", diff)
		}
		return
	}
	t.Errorf(`failed in %s: %v`, prefix, err)
}
