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
	"github.com/Grant-Nelson/Gozer/project/modeler/irc"
	"github.com/Grant-Nelson/Gozer/project/modeler/remodel"
)

func Test_Blocker_Label(t *testing.T) {
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
	exp := lines(`xyz`)
	got := stringForFunc(t, pkg, `doThing`)
	diffString(t, got, exp)
}

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
	fn := pkg.Irc.FindFunc(funcName)
	if fn == nil {
		t.Fatalf(`failed to find function`)
	}
	return fn.String()
}

func blockIrcFunc(t *testing.T, lines ...string) *project.Package {
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
	pkg.Irc = &irc.Package{}
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
				fn := pkg.Irc.NewFunc(fnDecl)
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
	return pkg
}
