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
	pkg      *project.Package
	errGroup *faults.Group

	// beingAdded is the import paths and identifiers for the decls and specs
	// that are being added, the value is the position value for the node.
	beingAdded map[string]token.Pos

	newImports       []ast.Decl
	newImportSpecs   []*ast.ImportSpec
	newDecls         []ast.Decl
	newDeclsComments [][]*ast.CommentGroup
	newFields        map[string]*ast.StructType
	newMethods       map[string]*ast.InterfaceType
}

func newAdd(pkg *project.Package, errGroup *faults.Group) *augAdd {
	return &augAdd{
		pkg:              pkg,
		errGroup:         errGroup,
		beingAdded:       map[string]token.Pos{},
		newImports:       []ast.Decl{},
		newImportSpecs:   []*ast.ImportSpec{},
		newDecls:         []ast.Decl{},
		newDeclsComments: [][]*ast.CommentGroup{},
		newFields:        map[string]*ast.StructType{},
		newMethods:       map[string]*ast.InterfaceType{},
	}
}

var (
	_ mods.Modifier         = (*augAdd)(nil)
	_ mods.ModifyAstFileExt = (*augAdd)(nil)
)

func (a *augAdd) ModifyAstFile(f *ast.File) (bool, error) {
	for id := range astTools.Idents(a.pkg.Ast.Fset, f) {
		if err := a.checkForExistingId(id); err != nil {
			return false, err
		}
		if err := a.tryToAddFields(id); err != nil {
			return false, err
		}
		if err := a.tryToAddMethods(id); err != nil {
			return false, err
		}
	}
	a.addImports(f)
	a.addDecls(f)
	return true, nil
}

func (a *augAdd) PackageDone() (bool, error) {
	if len(a.newFields) > 0 {
		names := slices.Collect(maps.Keys(a.newFields))
		sort.Strings(names)
		for _, name := range names {
			st := a.newFields[name]
			if err := a.errGroup.Add(faults.From(ErrAugAddStructIdDidNotExist).
				With(`package path`, a.pkg.PkgPath()).
				With(`augmenter pos`, a.pkg.Position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return false, err
			}
		}
	}
	if len(a.newMethods) > 0 {
		names := slices.Collect(maps.Keys(a.newMethods))
		sort.Strings(names)
		for _, name := range names {
			st := a.newMethods[name]
			if err := a.errGroup.Add(faults.From(ErrAugAddInterfaceIdDidNotExist).
				With(`package path`, a.pkg.PkgPath()).
				With(`augmenter pos`, a.pkg.Position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

// checkForExistingId checks that none of the decls being added already exist.
func (a *augAdd) checkForExistingId(id *astTools.IdentIteratorValue) error {
	pos, has := a.beingAdded[id.Ident]
	if !has {
		return nil
	}
	return a.errGroup.Add(faults.From(ErrAugAddIdAlreadyExists).
		With(`package path`, a.pkg.PkgPath()).
		With(`original pos`, id.Start()).
		With(`augmenter pos`, a.pkg.Position(pos)).
		With(`identifier`, id.Ident))
}

// tryToAddFields checks if the given ident is a structure to add fields to.
func (a *augAdd) tryToAddFields(id *astTools.IdentIteratorValue) error {
	fieldsToAdd, has := a.newFields[id.Ident]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return a.errGroup.Add(faults.From(ErrAugAddStructIdNotForType).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.pkg.Position(fieldsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	structToAddTo, ok := id.TypeSpec.Type.(*ast.StructType)
	if !ok {
		return a.errGroup.Add(faults.From(ErrAugAddStructTypeMismatch).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.pkg.Position(fieldsToAdd.Pos())).
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
			return a.errGroup.Add(faults.From(ErrAugAddIdFieldAlreadyExists).
				With(`package path`, a.pkg.PkgPath()).
				With(`original struct pos`, id.Start()).
				With(`augmenter struct pos`, a.pkg.Position(fieldsToAdd.Pos())).
				With(`original field pos`, id.Position(name.Pos())).
				With(`augmenter field pos`, a.pkg.Position(pos)).
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
func (a *augAdd) tryToAddMethods(id *astTools.IdentIteratorValue) error {
	methodsToAdd, has := a.newMethods[id.Ident]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return a.errGroup.Add(faults.From(ErrAugAddInterfaceIdNotForType).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.pkg.Position(methodsToAdd.Pos())).
			With(`identifier`, id.Ident))
	}
	interfaceToAddTo, ok := id.TypeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return a.errGroup.Add(faults.From(ErrAugAddInterfaceTypeMismatch).
			With(`package path`, a.pkg.PkgPath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.pkg.Position(methodsToAdd.Pos())).
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
			return a.errGroup.Add(faults.From(ErrAugAddIdMethodAlreadyExists).
				With(`package path`, a.pkg.PkgPath()).
				With(`original interface pos`, id.Start()).
				With(`augmenter interface pos`, a.pkg.Position(methodsToAdd.Pos())).
				With(`original method pos`, id.Position(name.Pos())).
				With(`augmenter method pos`, a.pkg.Position(pos)).
				With(`interface identifier`, id.Ident).
				With(`method identifier`, name.Name))
		}
	}
	// Add methods to the interface and remove added methods from map so we know it was added.
	interfaceToAddTo.Methods.List = append(interfaceToAddTo.Methods.List, methodsToAdd.Methods.List...)
	delete(a.newMethods, id.Ident)
	return nil
}

func findImportInsert(f *ast.File) int {
	for i, d := range f.Decls {
		switch g := d.(type) {
		case *ast.GenDecl:
			if g.Tok != token.IMPORT {
				return i
			}
		}
	}
	return 0
}

func (a *augAdd) addImports(f *ast.File) {
	if len(a.newImports) <= 0 {
		return
	}

	// If there are conflicts in the imports, type checking will catch those conflicts.
	f.Imports = append(f.Imports, a.newImportSpecs...)

	insert := findImportInsert(f)
	f.Decls = slices.Insert(f.Decls, insert, a.newImports...)

	a.newImports = []ast.Decl{}
	a.newImportSpecs = []*ast.ImportSpec{}
}

// addDecls adds all declarations needing to be added and
// clears out the list of decls needing to be added so that
// they are only added once.
func (a *augAdd) addDecls(f *ast.File) {
	if len(a.newDecls) > 0 {
		for i, d := range a.newDecls {
			f.Decls = append(f.Decls, d)
			f.Comments = append(f.Comments, a.newDeclsComments[i]...)
		}
		a.newDecls = []ast.Decl{}
	}
}
