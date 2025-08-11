package augmenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
	"github.com/google/go-cmp/cmp"
)

func TestAddingType(t *testing.T) {
	runAugTest(t, augTest{
		origSrc: lines(
			`package foo`),
		augSrc: lines(
			`package foo`,
			``,
			`//gozer:add`,
			`type Foo struct{}`),
		expSrc: lines(
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

	fm := fileMod.New(`original`)
	if err := fm.AddFile(`orig.go`, []byte(test.origSrc)); err != nil {
		t.Errorf(`failed to load origin file: %v`, err)
		return
	}

	test.errLimit = max(test.errLimit, 1)
	pkgPath := `test/path`
	errs := faults.NewGroup(test.errLimit)
	a := New(nil, `base`, pkgPath)
	a.reset()
	if err := a.addFile(`aug.go`, []byte(test.augSrc), errs); err != nil {
		checkErr(t, `load augment file`, test, err)
		return
	}

	if err := a.Modify(fm, errs); err != nil {
		checkErr(t, `modify file`, test, err)
		return
	}

	if err := a.PackageDone(`test`, pkgPath, errs); err != nil {
		checkErr(t, `finish package`, test, err)
		return
	}

	if err := errs.Wrap(); err != nil {
		checkErr(t, `accumulated error`, test, err)
		return
	}

	buf := &strings.Builder{}
	if err := fm.Write(buf); err != nil {
		t.Errorf(`failed to write result: %v`, err)
		return
	}

	if diff := cmp.Diff(strings.Split(buf.String(), "\n"), strings.Split(test.expSrc, "\n")); len(diff) > 0 {
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
