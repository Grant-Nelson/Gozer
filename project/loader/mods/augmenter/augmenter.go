package augmenter

import (
	"errors"
	"go/build/constraint"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
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

func (a *Augmenter) PackageStart(pkg *mods.Package, errs *faults.Group) error {
	a.reset()
	if err := a.addPackage(pkg.Path, errs); err != nil {
		return err
	}
	return a.Group.PackageStart(pkg, errs)
}

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

func (a *Augmenter) addFile(filename string, src []byte, errGroup *faults.Group) error {
	f, err := file.Load(a.fileSet, filename, src)
	if err != nil {
		return errGroup.Add(err)
	}
	build, err := a.shouldAdd(f, errGroup)
	if err != nil {
		return err
	}
	if !build {
		return nil
	}
	for ds := range f.DeclSpecs() {
		if err := a.addDeclSpec(ds, errGroup); err != nil {
			return err
		}
	}
	return nil
}

func (a *Augmenter) shouldAdd(f *file.File, errGroup *faults.Group) (bool, error) {
	if f.File.Doc == nil || len(f.File.Doc.List) <= 0 {
		return true, nil
	}
	for _, com := range f.File.Doc.List {
		exp, err := constraint.Parse(com.Text)
		if err != nil {
			return false, errGroup.Add(faults.From(ErrParsingBuildConstraints).
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

func (a *Augmenter) addDeclSpec(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	dv, err := readDirectives(ds.Comments(), ds.File.PackagePath(), ds.Start(), errGroup)
	if err != nil {
		return err
	}
	if dv.none {
		return a.addDeclSpecNoDirective(ds, errGroup)
	}

	// TODO: Implement
	return nil
}

func (a *Augmenter) addDeclSpecNoDirective(dv *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
