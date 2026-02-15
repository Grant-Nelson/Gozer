package augmenter

import (
	"errors"
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
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

type delHandle func(*astTools.IdentIteratorValue) (bool, error)

type augDel struct {
	pkg        *project.Package
	errGroup   *faults.Group
	delImport  map[string]bool
	delFunc    map[string]*ast.FuncDecl
	delVar     map[string]*ast.ValueSpec
	delType    map[string]*ast.TypeSpec
	delFields  map[string]*ast.StructType
	delMethods map[string]*ast.InterfaceType
	delHandles []delHandle
}

func newDel(pkg *project.Package, errGroup *faults.Group) *augDel {
	a := &augDel{
		pkg:        pkg,
		errGroup:   errGroup,
		delImport:  map[string]bool{},
		delFunc:    map[string]*ast.FuncDecl{},
		delVar:     map[string]*ast.ValueSpec{},
		delType:    map[string]*ast.TypeSpec{},
		delFields:  map[string]*ast.StructType{},
		delMethods: map[string]*ast.InterfaceType{},
	}
	a.delHandles = []delHandle{
		a.tryDelFunc,
		a.tryDelVar,
		a.tryDelType,
		a.tryDelFields,
		a.tryDelMethods,
	}
	return a
}

var (
	_ mods.Modifier         = (*augDel)(nil)
	_ mods.ModifyAstFileExt = (*augDel)(nil)
)

func (a *augDel) ModifyAstFile(f *ast.File) (bool, error) {
	for it := range astTools.Idents(a.pkg.Ast.Fset, f) {
		for _, handle := range a.delHandles {
			deleted, err := handle(it)
			if err != nil {
				return false, err
			}
			if deleted {
				break
			}
		}
	}
	return true, nil
}

func (a *augDel) PackageDone() (bool, error) {
	// TODO: Check for any identifiers that weren't found.
	return true, nil
}

func (a *augDel) tryDelFunc(it *astTools.IdentIteratorValue) (bool, error) {
	d, has := a.delFunc[it.Ident]
	if !has {
		return false, nil
	}
	if it.FuncDecl == nil {
		if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotFunc).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.pkg.Position(d.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	it.File.Decls[it.DeclIndex] = nil
	delete(a.delFunc, it.Ident)
	return true, nil
}

func (a *augDel) tryDelVar(it *astTools.IdentIteratorValue) (bool, error) {
	v, has := a.delVar[it.Ident]
	if !has {
		return false, nil
	}
	if it.ValueSpec == nil {
		if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotValue).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.pkg.Position(v.Pos())).
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

func (a *augDel) tryDelType(it *astTools.IdentIteratorValue) (bool, error) {
	t, has := a.delType[it.Ident]
	if !has {
		return false, nil
	}
	if it.TypeSpec == nil {
		if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotType).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.pkg.Position(t.Pos())).
			With(`identifier`, it.Ident)); err != nil {
			return false, err
		}
	}
	_, itInter := it.TypeSpec.Type.(*ast.InterfaceType)
	_, tInter := t.Type.(*ast.InterfaceType)
	if itInter != tInter {
		if itInter {
			if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotStruct).
				With(`package path`, a.pkg.PkgPath()).
				With(`original pos`, it.Start()).
				With(`augmenter pos`, a.pkg.Position(t.Pos())).
				With(`identifier`, it.Ident)); err != nil {
				return false, err
			}
		} else {
			if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotInterface).
				With(`package path`, a.pkg.PkgPath()).
				With(`original pos`, it.Start()).
				With(`augmenter pos`, a.pkg.Position(t.Pos())).
				With(`identifier`, it.Ident)); err != nil {
				return false, err
			}
		}
	}
	it.GenDecl.Specs[it.SpecIndex] = nil
	delete(a.delType, it.Ident)
	return true, nil
}

func (a *augDel) tryDelFields(it *astTools.IdentIteratorValue) (bool, error) {
	fs, has := a.delFields[it.Ident]
	if !has {
		return false, nil
	}
	st, ok := it.TypeSpec.Type.(*ast.StructType)
	if !ok {
		if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotField).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.pkg.Position(fs.Pos())).
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

func (a *augDel) tryDelMethods(it *astTools.IdentIteratorValue) (bool, error) {
	ms, has := a.delMethods[it.Ident]
	if !has {
		return false, nil
	}
	st, ok := it.TypeSpec.Type.(*ast.InterfaceType)
	if !ok {
		if err := a.errGroup.Add(faults.From(ErrAugDelIdentifierNotMethod).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, it.Start()).
			With(`augmenter pos`, a.pkg.Position(ms.Pos())).
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
