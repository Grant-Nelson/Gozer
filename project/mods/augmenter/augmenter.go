package augmenter

import (
	"errors"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/mods"
)

var (
	ErrParsingBuildConstraints = errors.New(`error parsing build constrains`)
)

type Augmenter struct {
	mods.Group
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd

	build       []string
	basePath    string
	testPkgPath string
	fileSet     *token.FileSet
}

func New(build []string, basePath, testPkgPath string) *Augmenter {
	a := &Augmenter{
		del: &augDel{},
		rep: &augReplace{},
		ren: &augRename{},
		add: &augAdd{},

		build:       build,
		testPkgPath: testPkgPath,
		basePath:    basePath,
	}
	a.Group = mods.Group{a.del, a.rep, a.ren, a.add}
	return a
}

func (a *Augmenter) PackageStart(name, path string, errs *faults.Group) error {
	a.reset()
	if err := a.addPackage(path, errs); err != nil {
		return err
	}
	return a.Group.PackageStart(name, path, errs)
}

const parseMode = parser.AllErrors | parser.ParseComments | parser.SkipObjectResolution

func (a *Augmenter) reset() {
	a.fileSet = token.NewFileSet()
	a.del.reset(a.fileSet)
	a.rep.reset(a.fileSet)
	a.ren.reset(a.fileSet)
	a.add.reset(a.fileSet)
}

func (a *Augmenter) addPackage(path string, errs *faults.Group) error {
	dir := filepath.Join(a.basePath, path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	testDir := path == a.testPkgPath
	for _, entry := range entries {
		if !entry.IsDir() {
			if filepath.Ext(entry.Name()) != `go` {
				continue
			}
			if !testDir || !strings.HasSuffix(entry.Name(), `_test.go`) {
				continue
			}
			if err := a.addFile(entry.Name(), nil, errs); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Augmenter) addFile(filename string, src []byte, errs *faults.Group) error {
	f, err := parser.ParseFile(a.fileSet, filename, src, parseMode)
	if err != nil {
		return errs.Add(err)
	}
	build, err := a.shouldAdd(f, errs)
	if err != nil {
		return err
	}
	if !build {
		return nil
	}
	if err := a.del.AddFile(f, errs); err != nil {
		return err
	}
	if err := a.rep.AddFile(f, errs); err != nil {
		return err
	}
	if err := a.ren.AddFile(f, errs); err != nil {
		return err
	}
	return a.add.AddFile(f, errs)
}

func (a *Augmenter) shouldAdd(f *ast.File, errs *faults.Group) (bool, error) {
	if f.Doc == nil || len(f.Doc.List) <= 0 {
		return true, nil
	}
	for _, com := range f.Doc.List {
		exp, err := constraint.Parse(com.Text)
		if err != nil {
			return false, errs.Add(faults.From(ErrParsingBuildConstraints).
				With(`error`, err).
				With(`position`, a.fileSet.Position(com.Pos())))
		}
		if !exp.Eval(func(tag string) bool {
			return slices.Contains(a.build, tag)
		}) {
			return false, nil
		}

	}
	return true, nil
}
