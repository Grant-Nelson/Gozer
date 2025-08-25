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
	"github.com/Grant-Nelson/Gozer/project/loader/mods/augmenter/directives"
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

	for _, d := range f.File.Decls {
		if err := a.readDecl(f, d, errGroup); err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrParsingBuildConstraints = errors.New(`error parsing build constrains for augmentation file`)
	ErrParsingUnexpectedDecl   = errors.New(`unexpected declaration while parsing augmentation file`)
	ErrParsingUnexpectedSpec   = errors.New(`unexpected specification while parsing augmentation file`)
	ErrAugFuncNone             = errors.New(`a function must have a directive`)
	ErrAugSpecNone             = errors.New(`a specification must have a directive`)
	ErrAugGenWithFuncDirective = errors.New(`a general declaration may not have a directive for a function`)
)

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

func (a *Augmenter) readDecl(f *file.File, decl ast.Decl, errGroup *faults.Group) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return a.readFuncDecl(f, d, errGroup)
	case *ast.GenDecl:
		return a.readGenDecl(f, d, errGroup)
	default:
		return errGroup.Add(faults.From(ErrParsingUnexpectedDecl).
			With(`package path`, f.PackagePath()).
			With(`position`, f.FileSet.Position(d.Pos())))
	}
}

func (a *Augmenter) readFuncDecl(f *file.File, fd *ast.FuncDecl, errGroup *faults.Group) error {
	pkgPath := f.PackagePath()
	pos := f.FileSet.Position(fd.Pos())
	dv, err := directives.Read(file.JoinComments(fd.Doc), pkgPath, pos, errGroup)
	if err != nil {
		return err
	}

	switch {
	case dv.Ignore:
		return nil
	case dv.None:
		return errGroup.Add(faults.From(ErrAugFuncNone).
			With(`package path`, pkgPath).
			With(`function name`, fd.Name.Name).
			With(`position`, pos))
	case dv.Add:
		a.add.newDecls = append(a.add.newDecls, fd)
		a.add.beingAdded[fd.Name.Name] = fd.Pos()
		return nil
	case dv.Delete:
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	case dv.Replace:
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	}
	return nil
}

func (a *Augmenter) readGenDecl(f *file.File, gd *ast.GenDecl, errGroup *faults.Group) error {
	pkgPath := f.PackagePath()
	declPos := f.FileSet.Position(gd.Pos())
	declDv, err := directives.Read(file.JoinComments(gd.Doc), pkgPath, declPos, errGroup)
	if err != nil {
		return err
	}

	if len(declDv.ReplaceRecv) > 0 || declDv.ReplaceSig {
		if err := errGroup.Add(faults.From(ErrAugGenWithFuncDirective).
			With(`package path`, f.PackagePath()).
			With(`position`, declPos)); err != nil {
			return err
		}
	}

	if len(declDv.Rename) > 0 && len(gd.Specs) != 1 {
		// TODO: Check that a rename only occurs on single specs
	}

	for _, spec := range gd.Specs {
		specDv, err := a.readSpecDirectives(f, spec, errGroup)
		if err != nil {
			return err
		}

		if declDv.None {
			if specDv.None {
				if err := errGroup.Add(faults.From(ErrAugSpecNone).
					With(`package path`, f.PackagePath()).
					With(`position`, f.FileSet.Position(spec.Pos()))); err != nil {
					return err
				}
			} else if err := a.readSpec(f, gd, spec, errGroup); err != nil {
				return err
			}
		} else {
			if specDv.None {
				if err := a.readSpec(f, gd, spec, errGroup); err != nil {
					return err
				}
			} else {
				// TODO: Implement, Check that the directives don't clash.
				panic(errors.New(`unimplemented`))
			}
		}
	}

	switch {
	case declDv.Add:
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	case declDv.Delete:
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	case declDv.Replace:
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	}
	return nil
}

func (a *Augmenter) readSpecDirectives(f *file.File, spec ast.Spec, errGroup *faults.Group) (*directives.Directives, error) {
	specPos := f.FileSet.Position(spec.Pos())
	var comments []*ast.Comment
	switch s := spec.(type) {
	case *ast.ImportSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	case *ast.TypeSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	case *ast.ValueSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	default:
		return nil, errGroup.Add(faults.From(ErrParsingUnexpectedSpec).
			With(`package path`, f.PackagePath()).
			With(`position`, specPos))
	}
	return directives.Read(comments, f.PackagePath(), specPos, errGroup)
}

func (a *Augmenter) readSpec(f *file.File, gd *ast.GenDecl, s ast.Spec, errGroup *faults.Group) error {
	// TODO: Implement
	panic(errors.New(`unimplemented`))
}
