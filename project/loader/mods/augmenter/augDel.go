package augmenter

import (
	"errors"
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
)

var (
	ErrAugDelIdentifierNotFunc      = errors.New(`can not delete identifier via augmenter: identifier not function`)
	ErrAugDelIdentifierNotValue     = errors.New(`can not delete identifier via augmenter: identifier not value`)
	ErrAugDelIdentifierNotType      = errors.New(`can not delete identifier via augmenter: identifier not type`)
	ErrAugDelIdentifierNotStruct    = errors.New(`can not delete identifier via augmenter: identifier not struct`)
	ErrAugDelIdentifierNotInterface = errors.New(`can not delete identifier via augmenter: identifier not interface`)
	ErrAugDelIdentifierNotField     = errors.New(`can not delete identifier via augmenter: identifier for field not in struct`)
	ErrAugDelIdentifierNotMethod    = errors.New(`can not delete identifier via augmenter: identifier for method not in interface`)
)

type delHandle func(*file.IdentIteratorValue, *faults.Group) (bool, error)

type augDel struct {
	fileSet    *token.FileSet
	delImport  map[string]bool
	delFunc    map[string]*ast.FuncDecl
	delVar     map[string]*ast.ValueSpec
	delType    map[string]*ast.TypeSpec
	delFields  map[string]*ast.StructType
	delMethods map[string]*ast.InterfaceType
	delHandles []delHandle
}

func (a *augDel) reset() {
	a.delImport = map[string]bool{}
	a.delFunc = map[string]*ast.FuncDecl{}
	a.delVar = map[string]*ast.ValueSpec{}
	a.delType = map[string]*ast.TypeSpec{}
	a.delFields = map[string]*ast.StructType{}
	a.delMethods = map[string]*ast.InterfaceType{}
	if len(a.delHandles) <= 0 {
		a.delHandles = []delHandle{
			a.tryDelFunc,
			a.tryDelVar,
			a.tryDelType,
			a.tryDelFields,
			a.tryDelMethods}
	}
}

func (a *augDel) Modify(f *file.File, errGroup *faults.Group) error {
	for it := range f.Idents() {
		for _, handle := range a.delHandles {
			deleted, err := handle(it, errGroup)
			if err != nil {
				return err
			}
			if deleted {
				break
			}
		}
	}
	return nil
}

func (a *augDel) PackageDone(name, path string, errGroup *faults.Group) error {
	// TODO: Check for any identifiers that weren't found.
	return nil
}

func (a *augDel) tryDelFunc(it *file.IdentIteratorValue, errGroup *faults.Group) (bool, error) {
	d, has := a.delFunc[it.Ident]
	if !has {
		return false, nil
	}
	if it.FuncDecl == nil {
		if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotFunc).
			With(`package path`, it.File.PackagePath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.fileSet.Position(d.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	it.File.File.Decls[it.DeclIndex] = nil
	delete(a.delFunc, it.Ident)
	return true, nil
}

func (a *augDel) tryDelVar(it *file.IdentIteratorValue, errGroup *faults.Group) (bool, error) {
	v, has := a.delVar[it.Ident]
	if !has {
		return false, nil
	}
	if it.ValueSpec == nil {
		if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotValue).
			With(`package path`, it.File.PackagePath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.fileSet.Position(v.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	it.ValueSpec.Names[it.ValueIndex] = nil
	if len(it.ValueSpec.Names) == len(it.ValueSpec.Values) {
		it.ValueSpec.Values[it.ValueIndex] = nil
	}
	delete(a.delVar, it.Ident)
	return true, nil
}

func (a *augDel) tryDelType(it *file.IdentIteratorValue, errGroup *faults.Group) (bool, error) {
	t, has := a.delType[it.Ident]
	if !has {
		return false, nil
	}
	if it.TypeSpec == nil {
		if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotType).
			With(`package path`, it.File.PackagePath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.fileSet.Position(t.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	_, itInter := it.TypeSpec.Type.(*ast.InterfaceType)
	_, tInter := t.Type.(*ast.InterfaceType)
	if itInter != tInter {
		if itInter {
			if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotStruct).
				With(`package path`, it.File.PackagePath()).
				With(`original pos`, it.Start()).
				With(`augmenter pos`, a.fileSet.Position(t.Pos())).
				With(`identifier`, it.Ident)); err != nil {
				return false, err
			}
		} else {
			if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotInterface).
				With(`package path`, it.File.PackagePath()).
				With(`original pos`, it.Start()).
				With(`augmenter pos`, a.fileSet.Position(t.Pos())).
				With(`identifier`, it.Ident)); err != nil {
				return false, err
			}
		}
	}
	it.GenDecl.Specs[it.SpecIndex] = nil
	delete(a.delType, it.Ident)
	return true, nil
}

func (a *augDel) tryDelFields(it *file.IdentIteratorValue, errGroup *faults.Group) (bool, error) {
	fs, has := a.delFields[it.Ident]
	if !has {
		return false, nil
	}
	st, ok := it.TypeSpec.Type.(*ast.StructType)
	if !ok {
		if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotField).
			With(`package path`, it.File.PackagePath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.fileSet.Position(fs.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	// Collect the names of the fields to delete.
	delNames := map[string]token.Pos{}
	for _, field := range fs.Fields.List {
		for _, name := range field.Names {
			delNames[name.Name] = name.Pos()
		}
	}
	// Delete those names from the structure.
	for _, field := range st.Fields.List {
		for i, name := range field.Names {
			if _, del := delNames[name.Name]; del {
				field.Names[i] = nil
				delete(delNames, name.Name)
			}
		}
	}
	// Check that all names were found.

	// TODO: Finish
	delete(a.delFields, it.Ident)
	return true, nil
}

func (a *augDel) tryDelMethods(it *file.IdentIteratorValue, errGroup *faults.Group) (bool, error) {
	ms, has := a.delMethods[it.Ident]
	if !has {
		return false, nil
	}
	st, ok := it.TypeSpec.Type.(*ast.InterfaceType)
	if !ok {
		if err := errGroup.Add(faults.From(ErrAugDelIdentifierNotMethod).
			With(`package path`, it.File.PackagePath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.fileSet.Position(ms.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	// Collect the names of the methods to delete.
	delNames := map[string]token.Pos{}
	for _, method := range ms.Methods.List {
		for _, name := range method.Names {
			delNames[name.Name] = name.Pos()
		}
	}
	// Delete those names from the interface.
	for _, method := range st.Methods.List {
		for i, name := range method.Names {
			if _, del := delNames[name.Name]; del {
				method.Names[i] = nil
				delete(delNames, name.Name)
			}
		}
	}
	// Check that all names were found.
	if len(delNames) > 0 {
		names := slices.Collect(maps.Keys(delNames))
		sort.Strings(names)

		// TODO: Finish

	}
	delete(a.delMethods, it.Ident)
	return true, nil
}
