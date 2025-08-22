package augmenter

import (
	"errors"
	"go/ast"
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
		if err := a.readDeclSpec(ds, errGroup); err != nil {
			return err
		}
	}
	return nil
}

var ErrParsingBuildConstraints = errors.New(`error parsing build constrains`)

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

func (a *Augmenter) readDeclSpec(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	dv, err := readDirectives(ds.Comments(), ds.File.PackagePath(), ds.Start(), errGroup)
	switch {
	case err != nil:
		return err
	case dv.none:
		return a.readDeclSpecNoDirective(ds, errGroup)
	case dv.add:
		return a.readDeclSpecAdd(ds, errGroup)
	case dv.delete:
		return a.readDeclSpecDelete(ds, errGroup)
	case dv.ignore:
		return a.readDeclSpecIgnore(ds, errGroup)
	default:
		return a.readDeclSpecReplace(ds, dv, errGroup)
	}
}

var (
	ErrAugFuncNoDirectives     = errors.New(`function is missing directive`)
	ErrAugImportNoDirectives   = errors.New(`import is missing directive`)
	ErrAugVarConstNoDirectives = errors.New(`var or const is missing directive`)
)

func (a *Augmenter) readDeclSpecNoDirective(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	if ds.FuncDecl != nil {
		return errGroup.Add(faults.From(ErrAugFuncNoDirectives).
			With(`package path`, ds.File.PackagePath()).
			With(`position`, ds.Start()).
			With(`ident`, ds.FuncDecl.Name.Name))
	}
	if ds.ImportSpec != nil {
		return errGroup.Add(faults.From(ErrAugImportNoDirectives).
			With(`package path`, ds.File.PackagePath()).
			With(`position`, ds.Start()).
			With(`ident`, ds.ImportSpec.Path.Value))
	}
	if ds.ValueSpec != nil {
		return errGroup.Add(faults.From(ErrAugVarConstNoDirectives).
			With(`package path`, ds.File.PackagePath()).
			With(`position`, ds.Start()))
	}
	// Check that there are directives on the fields or methods, error otherwise.
	switch t := ds.TypeSpec.Type.(type) {
	case *ast.StructType:
		return a.readStructNoDirective(ds, t, errGroup)
	case *ast.InterfaceType:
		return a.readInterfaceNoDirective(ds, t, errGroup)
	}
	return nil
}

func (a *Augmenter) readStructNoDirective(ds *file.DeclSpecIteratorValue, st *ast.StructType, errGroup *faults.Group) error {
	/*
		var (
			addFields     map[*ast.Field]struct{}
			deleteFields  map[*ast.Field]struct{}
			replaceFields map[*ast.Field]struct{}
		)
		for _, f := range st.Fields.List {
			dv, err := readDirectives(file.JoinComments(f.Doc, f.Comment), ds.File.PackagePath(), ds.Position(), errGroup)
			if err != nil {
				return err
			}
			// TODO: Implement
		}
	*/
	return nil
}

func (a *Augmenter) readInterfaceNoDirective(ds *file.DeclSpecIteratorValue, it *ast.InterfaceType, errGroup *faults.Group) error {
	// for _, m := range it.Methods.List {
	// TODO: Implement
	// }
	return nil
}

func (a *Augmenter) readDeclSpecIgnore(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	// Check that no directives had been added onto fields and methods if type decl.
	if ds.TypeSpec == nil {
		return nil
	}
	/*
		switch t := ds.TypeSpec.Type.(type) {
		case *ast.StructType:
			for _, f := range t.Fields.List {
				// TODO: Implement
			}
		case *ast.InterfaceType:
			for _, m := range t.Methods.List {
				// TODO: Implement
			}
		}
	*/
	return nil
}

func (a *Augmenter) readDeclSpecAdd(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *Augmenter) readDeclSpecDelete(ds *file.DeclSpecIteratorValue, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *Augmenter) readDeclSpecReplace(ds *file.DeclSpecIteratorValue, dv *directives, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
