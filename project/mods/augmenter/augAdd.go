package augmenter

import (
	"errors"
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
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
	fileSet    *token.FileSet
	beingAdded map[string]token.Pos
	addDecl    []ast.Decl
	addFields  map[string]*ast.StructType
	addMethods map[string]*ast.InterfaceType
}

func (a *augAdd) reset(fileSet *token.FileSet) {
	a.fileSet = fileSet
	a.beingAdded = map[string]token.Pos{}
	a.addDecl = []ast.Decl{}
	a.addFields = map[string]*ast.StructType{}
	a.addMethods = map[string]*ast.InterfaceType{}
}

func (a *augAdd) Modify(fm *fileMod.FileMod, errs *faults.Group) error {
	for id := range fm.Idents() {
		if err := a.checkForExistingId(id, errs); err != nil {
			return err
		}
		if err := a.tryToAddFields(id, errs); err != nil {
			return err
		}
		if err := a.tryToAddMethods(id, errs); err != nil {
			return err
		}
	}
	a.addDecls(fm)
	return nil
}

func (a *augAdd) PackageDone(name, path string, errs *faults.Group) error {
	if len(a.addFields) > 0 {
		names := slices.Collect(maps.Keys(a.addFields))
		sort.Strings(names)
		for _, name := range names {
			st := a.addFields[name]
			if err := errs.Add(faults.From(ErrAugAddStructIdDidNotExist).
				With(`package path`, path).
				With(`augmenter pos`, a.fileSet.Position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return err
			}
		}
	}
	if len(a.addMethods) > 0 {
		names := slices.Collect(maps.Keys(a.addMethods))
		sort.Strings(names)
		for _, name := range names {
			st := a.addMethods[name]
			if err := errs.Add(faults.From(ErrAugAddInterfaceIdDidNotExist).
				With(`package path`, path).
				With(`augmenter pos`, a.fileSet.Position(st.Pos())).
				With(`identifier`, name)); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkForExistingId checks that none of the decls being added already exist.
func (a *augAdd) checkForExistingId(id *fileMod.Ident, errs *faults.Group) error {
	pos, has := a.beingAdded[id.Name]
	if !has {
		return nil
	}
	return errs.Add(faults.From(ErrAugAddIdAlreadyExists).
		With(`package path`, id.FileMod.PackagePath()).
		With(`original pos`, id.Start()).
		With(`augmenter pos`, a.fileSet.Position(pos)).
		With(`identifier`, id.Name))
}

// tryToAddFields checks if the given ident is a structure to add fields to.
func (a *augAdd) tryToAddFields(id *fileMod.Ident, errs *faults.Group) error {
	fieldsToAdd, has := a.addFields[id.Name]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return errs.Add(faults.From(ErrAugAddStructIdNotForType).
			With(`package path`, id.FileMod.PackagePath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.fileSet.Position(fieldsToAdd.Pos())).
			With(`identifier`, id.Name))
	}
	structToAddTo, ok := id.TypeSpec.Type.(*ast.StructType)
	if !ok {
		return errs.Add(faults.From(ErrAugAddStructTypeMismatch).
			With(`package path`, id.FileMod.PackagePath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.fileSet.Position(fieldsToAdd.Pos())).
			With(`identifier`, id.Name))
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
			return errs.Add(faults.From(ErrAugAddIdFieldAlreadyExists).
				With(`package path`, id.FileMod.PackagePath()).
				With(`original struct pos`, id.Start()).
				With(`augmenter struct pos`, a.fileSet.Position(fieldsToAdd.Pos())).
				With(`original field pos`, id.Position(name.Pos())).
				With(`augmenter field pos`, a.fileSet.Position(pos)).
				With(`struct identifier`, id.Name).
				With(`field identifier`, name.Name))
		}
	}
	// Add fields to the struct and remove added fields from map so we know it was added.
	structToAddTo.Fields.List = append(structToAddTo.Fields.List, fieldsToAdd.Fields.List...)
	delete(a.addFields, id.Name)
	return nil
}

// tryToAddMethods checks if the given ident is an interface to add methods to.
func (a *augAdd) tryToAddMethods(id *fileMod.Ident, errs *faults.Group) error {
	methodsToAdd, has := a.addMethods[id.Name]
	if !has {
		return nil
	}
	if id.TypeSpec == nil {
		return errs.Add(faults.From(ErrAugAddInterfaceIdNotForType).
			With(`package path`, id.FileMod.PackagePath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.fileSet.Position(methodsToAdd.Pos())).
			With(`identifier`, id.Name))
	}
	interfaceToAddTo, ok := id.TypeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return errs.Add(faults.From(ErrAugAddInterfaceTypeMismatch).
			With(`package path`, id.FileMod.PackagePath()).
			With(`original pos`, id.Start()).
			With(`augmenter pos`, a.fileSet.Position(methodsToAdd.Pos())).
			With(`identifier`, id.Name))
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
			return errs.Add(faults.From(ErrAugAddIdMethodAlreadyExists).
				With(`package path`, id.FileMod.PackagePath()).
				With(`original interface pos`, id.Start()).
				With(`augmenter interface pos`, a.fileSet.Position(methodsToAdd.Pos())).
				With(`original method pos`, id.Position(name.Pos())).
				With(`augmenter method pos`, a.fileSet.Position(pos)).
				With(`interface identifier`, id.Name).
				With(`method identifier`, name.Name))
		}
	}
	// Add methods to the interface and remove added methods from map so we know it was added.
	interfaceToAddTo.Methods.List = append(interfaceToAddTo.Methods.List, methodsToAdd.Methods.List...)
	delete(a.addMethods, id.Name)
	return nil
}

// addDecls dda all declarations needing to be added and clears out the
// list of decls needing to be added so that they're only added once.
func (a *augAdd) addDecls(fm *fileMod.FileMod) {
	if len(a.addDecl) > 0 {
		fm.AddDecls(a.addDecl)
		a.addDecl = []ast.Decl{}
	}
}
