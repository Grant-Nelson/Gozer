package augmenter

import (
	"errors"
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

var (
	ErrAugAddIdAlreadyExists        = errors.New(`can not add new identifier via augmenter: identifier already exists`)
	ErrAugAddStructIdNotForType     = errors.New(`can not add fields via augmenter: a non-type found by identifier`)
	ErrAugAddStructTypeMismatch     = errors.New(`can not add fields via augmenter: a non-structure type found by identifier`)
	ErrAugAddIdFieldAlreadyExists   = errors.New(`can not add fields via augmenter: identifier for field already exists`)
	ErrAugAddStructIdDidNotExist    = errors.New(`can not add fields via augmenter: identifier for structure not found`)
	ErrAugAddInterfaceIdNotForType  = errors.New(`can not add methods via augmenter: a non-type found by identifier`)
	ErrAugAddInterfaceTypeMismatch  = errors.New(`can not add methods via augmenter: a non-interface type found by identifier`)
	ErrAugAddIdMethodAlreadyExists  = errors.New(`can not add methods via augmenter: identifier for method already exists`)
	ErrAugAddInterfaceIdDidNotExist = errors.New(`can not add methods via augmenter: identifier for interface not found`)
)

type augAdd struct {
	pkg *artifacts.Package

	// beingAdded is the import paths and identifiers for the decls and specs
	// that are being added, the value is the position value for the node.
	beingAdded map[string]token.Pos

	newImports     []ast.Decl
	newImportSpecs []*ast.ImportSpec
	newGenDecls    []*ast.GenDecl
	newFuncDecls   []*ast.FuncDecl
	newFields      map[string]*ast.StructType
	newMethods     map[string]*ast.InterfaceType
}

func newAdd(pkg *artifacts.Package) *augAdd {
	return &augAdd{
		pkg:            pkg,
		beingAdded:     map[string]token.Pos{},
		newImports:     []ast.Decl{},
		newImportSpecs: []*ast.ImportSpec{},
		newGenDecls:    []*ast.GenDecl{},
		newFuncDecls:   []*ast.FuncDecl{},
		newFields:      map[string]*ast.StructType{},
		newMethods:     map[string]*ast.InterfaceType{},
	}
}

var _ mods.Modifier = (*augAdd)(nil)
var _ mods.LoadDoneExt = (*augAdd)(nil)

func (a *augAdd) Modify(f *artifacts.File, errGroup *faults.Group) (bool, error) {
	for id := range f.Idents() {
		if err := a.checkForExistingId(id, errGroup); err != nil {
			return false, err
		}
		if err := a.tryToAddFields(id, errGroup); err != nil {
			return false, err
		}
		if err := a.tryToAddMethods(id, errGroup); err != nil {
			return false, err
		}
	}
	a.addImports(f)
	a.addDecls(f)
	return true, nil
}

func (a *augAdd) LoadDone(errGroup *faults.Group) error {
	if len(a.newFields) > 0 {
		names := slices.Collect(maps.Keys(a.newFields))
		sort.Strings(names)
		for _, name := range names {
			st := a.newFields[name]
			if err := errGroup.Add(faults.From(ErrAugAddStructIdDidNotExist).
				With(`package path`, a.pkg.Path()).
				With(`augmenter pos`, a.position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return err
			}
		}
	}
	if len(a.newMethods) > 0 {
		names := slices.Collect(maps.Keys(a.newMethods))
		sort.Strings(names)
		for _, name := range names {
			st := a.newMethods[name]
			if err := errGroup.Add(faults.From(ErrAugAddInterfaceIdDidNotExist).
				With(`package path`, a.pkg.Path()).
				With(`augmenter pos`, a.position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *augAdd) position(pos token.Pos) token.Position {
	return a.pkg.TempFileSet().Position(pos)
}

// checkForExistingId checks that none of the decls being added already exist.
func (a *augAdd) checkForExistingId(id *artifacts.IdentIteratorValue, errGroup *faults.Group) error {
	pos, has := a.beingAdded[id.Ident]
	if !has {
		return nil
	}
	return errGroup.Add(faults.From(ErrAugAddIdAlreadyExists).
		With(`package path`, a.pkg.Path()).
		With(`original pos`, id.Start()).
		With(`augmenter pos`, a.position(pos)).
		With(`identifier`, id.Ident))
}

// tryToAddFields checks if the given ident is a structure to add fields to.
func (a *augAdd) tryToAddFields(id *artifacts.IdentIteratorValue, errGroup *faults.Group) error {
	fieldsToAdd, has := a.newFields[id.Ident]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return errGroup.Add(faults.From(ErrAugAddStructIdNotForType).
			With(`package path`, a.pkg.Path()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.position(fieldsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	structToAddTo, ok := id.TypeSpec.Type.(*ast.StructType)
	if !ok {
		return errGroup.Add(faults.From(ErrAugAddStructTypeMismatch).
			With(`package path`, a.pkg.Path()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.position(fieldsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	// Collect the fields being added to check that the field doesn't exist already.
	beingAdded := map[string]token.Pos{}
	for _, field := range fieldsToAdd.Fields.List {
		for _, name := range field.Names {
			beingAdded[name.Name] = name.Pos()
		}
	}
	// Check that the fields being added don't already exist.
	for _, field := range structToAddTo.Fields.List {
		for _, name := range field.Names {
			pos, has := beingAdded[name.Name]
			if !has {
				continue
			}
			return errGroup.Add(faults.From(ErrAugAddIdFieldAlreadyExists).
				With(`package path`, a.pkg.Path()).
				With(`original struct pos`, id.Start()).
				With(`augmenter struct pos`, a.position(fieldsToAdd.Pos())).
				With(`original field pos`, id.Position(name.Pos())).
				With(`augmenter field pos`, a.position(pos)).
				With(`struct identifier`, id.Ident).
				With(`field identifier`, name.Name))
		}
	}
	// Add fields to the struct and remove added fields from map so we know it was added.
	structToAddTo.Fields.List = append(structToAddTo.Fields.List, fieldsToAdd.Fields.List...)
	delete(a.newFields, id.Ident)
	return nil
}

// tryToAddMethods checks if the given ident is an interface to add methods to.
func (a *augAdd) tryToAddMethods(id *artifacts.IdentIteratorValue, errGroup *faults.Group) error {
	methodsToAdd, has := a.newMethods[id.Ident]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return errGroup.Add(faults.From(ErrAugAddInterfaceIdNotForType).
			With(`package path`, a.pkg.Path()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.position(methodsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	interfaceToAddTo, ok := id.TypeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return errGroup.Add(faults.From(ErrAugAddInterfaceTypeMismatch).
			With(`package path`, a.pkg.Path()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.position(methodsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	// Collect the methods being added to check that the method doesn't exist already.
	beingAdded := map[string]token.Pos{}
	for _, method := range methodsToAdd.Methods.List {
		for _, name := range method.Names {
			beingAdded[name.Name] = name.Pos()
		}
	}
	// Check that the methods being added don't already exist.
	for _, method := range interfaceToAddTo.Methods.List {
		for _, name := range method.Names {
			pos, has := beingAdded[name.Name]
			if !has {
				continue
			}
			return errGroup.Add(faults.From(ErrAugAddIdMethodAlreadyExists).
				With(`package path`, a.pkg.Path()).
				With(`original interface pos`, id.Start()).
				With(`augmenter interface pos`, a.position(methodsToAdd.Pos())).
				With(`original method pos`, id.Position(name.Pos())).
				With(`augmenter method pos`, a.position(pos)).
				With(`interface identifier`, id.Ident).
				With(`method identifier`, name.Name))
		}
	}
	// Add methods to the interface and remove added methods from map so we know it was added.
	interfaceToAddTo.Methods.List = append(interfaceToAddTo.Methods.List, methodsToAdd.Methods.List...)
	delete(a.newMethods, id.Ident)
	return nil
}

func findImportInsert(f *artifacts.File) int {
	for i, d := range f.File.Decls {
		switch g := d.(type) {
		case *ast.GenDecl:
			if g.Tok != token.IMPORT {
				return i
			}
		}
	}
	return 0
}

func (a *augAdd) addImports(f *artifacts.File) {
	if len(a.newImports) <= 0 {
		return
	}

	// If there are conflicts in the imports, type checking will catch those conflicts.
	f.File.Imports = append(f.File.Imports, a.newImportSpecs...)

	insert := findImportInsert(f)
	f.File.Decls = slices.Insert(f.File.Decls, insert, a.newImports...)

	a.newImports = []ast.Decl{}
	a.newImportSpecs = []*ast.ImportSpec{}
}

// addDecls dda all declarations needing to be added and clears out the
// list of decls needing to be added so that they're only added once.
func (a *augAdd) addDecls(f *artifacts.File) {
	if len(a.newGenDecls) > 0 {
		for _, d := range a.newGenDecls {
			f.File.Decls = append(f.File.Decls, d)
			f.File.Comments = append(f.File.Comments, d.Doc)
		}
		a.newGenDecls = []*ast.GenDecl{}
	}

	if len(a.newFuncDecls) > 0 {
		for _, d := range a.newGenDecls {
			f.File.Decls = append(f.File.Decls, d)
			f.File.Comments = append(f.File.Comments, d.Doc)
		}
		a.newFuncDecls = []*ast.FuncDecl{}
	}
}
