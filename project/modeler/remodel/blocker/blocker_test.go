package blocker

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/ir"
	"github.com/Grant-Nelson/Gozer/project/modeler/remodel"
)

func Test_Blocker_Label_ForwardJump(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`	if i > 10 {`,
		`		goto Finished`,
		`	}`,
		`	i += 10`,
		`Finished:`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    if i > 10 {`,
		`      goto(block 1)`,
		`    }`,
		`    i+=10`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <Label Finished> {`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_BackwardJump(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`Loop:`,
		`   i++`,
		`	if i < 10 {`,
		`		goto Loop`,
		`	}`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <Label Loop> {`,
		`    i++`,
		`    if i < 10 {`,
		`      goto(block 1)`,
		`    }`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_JumpToIf(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`Loop:`,
		`	if i < 10 {`,
		`       i++`,
		`		goto Loop`,
		`	}`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <Label Loop> {`,
		`    if i < 10 {`,
		`      i++`,
		`      goto(block 1)`,
		`    }`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_BreakFor(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`ForLoop:`,
		`   for j := 0; j < 10; j++ {`,
		`       i++`,
		`		if i > 10 {`,
		`			break ForLoop`,
		`       }`,
		`	}`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <For-loop Init for ForLoop> {`,
		`    j:=0`,
		`    goto(block 2)`,
		`  }`,
		`  block 2 <For-loop Body for ForLoop> {`,
		`    if !(j < 10) {`,
		`      goto(block 3)`,
		`    }`,
		`    i++`,
		`    if i > 10 {`,
		`      goto(block 3)`,
		`    }`,
		`    j++`,
		`    goto(block 2)`,
		`  }`,
		`  block 3 <After For-loop for ForLoop> {`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_JumpForwardMiddleAndBack(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`	if i < 3 {`,
		`		goto ForLoop`,
		`	}`,
		`	i *= 2`,
		`ForLoop:`,
		`	for j := 0; j < 5; j++ {`,
		`		if i > 10 {`,
		`			break ForLoop`,
		`		}`,
		`		i++`,
		`	}`,
		`	if i < 3 {`,
		`		goto ForLoop`,
		`	}`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    if i < 3 {`,
		`      goto(block 1)`,
		`    }`,
		`    i*=2`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <For-loop Init for ForLoop> {`,
		`    j:=0`,
		`    goto(block 2)`,
		`  }`,
		`  block 2 <For-loop Body for ForLoop> {`,
		`    if !(j < 5) {`,
		`      goto(block 3)`,
		`    }`,
		`    if i > 10 {`,
		`      goto(block 3)`,
		`    }`,
		`    i++`,
		`    j++`,
		`    goto(block 2)`,
		`  }`,
		`  block 3 <After For-loop for ForLoop> {`,
		`    if i < 3 {`,
		`      goto(block 1)`,
		`    }`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_WhileLoop(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(i int) int {`,
		`Loop:`,
		`	for {`,
		`		if i < 10 {`,
		`       	i++`,
		`			continue Loop`,
		`		}`,
		`		break Loop`,
		`	}`,
		`	return i`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`    goto(block 1)`,
		`  }`,
		`  block 1 <For-loop Body for Loop> {`,
		`    if i < 10 {`,
		`      i++`,
		`      goto(block 1)`,
		`    }`,
		`    goto(block 2)`,
		`  }`,
		`  block 2 <After For-loop for Loop> {`,
		`    return i`,
		`  }`,
		`}`))
}

func Test_Blocker_Label_NestedFor(t *testing.T) {
	pkg := blockIrcFunc(t,
		`package main`,
		``,
		`func doThing(k int) int {`,
		`OuterLoop:`,
		`	for i := 0; i < 10; i++ {`,
		`InnerLoop:`,
		`		for j := 3; j >=0; j-- {`,
		`       	k++`,
		`			if k < 3 {`,
		`				continue InnerLoop`,
		`			}`,
		`			if k > 10 {`,
		`				break OuterLoop`,
		`			}`,
		`			continue OuterLoop`,
		`		}`,
		`	}`,
		`	return k`,
		`}`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, lines(
		`func doThing {`,
		`  block 0 <initial> {`,
		`  }`,
		`}`))
}

// TODO: Check that continue hits post expression of the for-loop.
// TODO: Check nested for-loops jump.
// TODO: Check adding an unused label.
// TODO: Check break, continue, and fallthrough with and without labels
// TODO: Check returns being added where implicit returns occur.
// TODO: Check for-ranges.

func diffString(t *testing.T, got, exp string) {
	gotLines := slices.Collect(strings.Lines(got))
	expLines := slices.Collect(strings.Lines(exp))
	if diff := cmp.Diff(gotLines, expLines); len(diff) > 0 {
		t.Errorf("unexpected result (-want, +got):\n%s", diff)
	}
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func stringForFunc(t *testing.T, pkg *project.Package, funcName string) string {
	fn := pkg.Ir.FindFunc(funcName)
	if fn == nil {
		t.Fatalf(`failed to find function`)
	}
	return fn.String()
}

func blockIrcFunc(t *testing.T, lines ...string) *project.Package {
	t.Helper()

	fileName := `blockTestFunc.go`
	dirPath := `c:\`
	fileSrc := strings.Join(lines, "\n")
	fileSet := token.NewFileSet()
	packageCfg := &packages.Config{
		Dir:  dirPath,
		Mode: packages.LoadAllSyntax,
		Fset: fileSet,
		Overlay: map[string][]byte{
			dirPath + fileName: []byte(fileSrc),
		},
	}
	roots, err := packages.Load(packageCfg, fileName)
	if err != nil {
		t.Fatalf(`failed to parse function for blocker test: %v`, err)
	}

	proj := project.New(fileSet, roots)
	errGroup := faults.NewErrGroup(10)
	proj.CollectErrors(errGroup)
	if err := errGroup.AnyOrNil(); err != nil {
		t.Fatalf(`errors in project: %v`, err)
	}

	if len(proj.Roots) > 1 {
		t.Fatalf(`expected only one root package but got %d`, len(proj.Roots))
	}
	pkg := proj.Roots[0]
	pkg.Ir = &ir.Package{}
	blocker := New(&Config{
		ErrGroup: errGroup,
	})
	_, rm, err := blocker.StartPackage(pkg)
	if err != nil {
		t.Fatalf(`error starting blocker: %v`, err)
	}

	brm := rm.(remodel.RemodelFuncExt)
	for _, file := range pkg.Ast.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			if fnDecl, ok := n.(*ast.FuncDecl); ok {
				fn := pkg.Ir.NewFunc(fnDecl)
				if _, err := brm.RemodelFunc(fn); err != nil {
					t.Errorf(`error updating function: %v`, err)
				}
			}
			return true
		})
	}
	if _, err := rm.PackageDone(); err != nil {
		t.Errorf(`error ending blocker: %v`, err)
	}

	if err := errGroup.AnyOrNil(); err != nil {
		t.Errorf(`errors blocker: %v`, err)
	}
	return pkg
}
