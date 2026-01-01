package augmenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/remapper"
)

func TestAddingType(t *testing.T) {
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
			`type Bar struct{}`),
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

	tempFileSet := artifacts.NewFileSet()
	fm, err := artifacts.Load(tempFileSet, `original/orig.go`, []byte(test.origSrc))
	if err != nil {
		t.Errorf(`failed to load origin file: %v`, err)
		return
	}

	test.errLimit = max(test.errLimit, 1)
	errGroup := faults.NewGroup(test.errLimit)
	a := New(nil, PathRebase(`original`, `base`))

	// Create an augmenter for a package then add the aug file to it
	pkg := fm.Package
	ap := a.addPackage(pkg, errGroup)
	if err := ap.AddFile(nil, `base/aug.go`, []byte(test.augSrc), errGroup); err != nil {
		checkErr(t, `load augment file`, test, err)
		return
	}

	// Perform the augmentation on the file
	con, err := a.Modify(fm, errGroup)
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

	finalFileSet := artifacts.NewFileSet()
	if err := remapper.Remap(fm, finalFileSet, errGroup); err != nil {
		checkErr(t, `remap`, test, err)
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
