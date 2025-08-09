package augmenter

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

var (
	ErrIdAlreadyExists = errors.New(`can not add new identifier via augmenter: Identifier already exist`)
	ErrIdNotForType    = errors.New(`can not add field via augmenter to the non-type found by identifier`)
)

type augAdd struct {
	fileSet   *token.FileSet
	addIds    map[string]token.Pos
	addDecl   []ast.Decl
	addFields map[string]*ast.TypeSpec
}

func (a *augAdd) Modify(fm *fileMod.FileMod) error {
	errs := faults.NewGroup()
	// Check that none of the decls being added already exist.
	for id := range fm.Idents() {
		if pos, has := a.addIds[id.Name]; has {
			if err := errs.Add(faults.From(ErrIdAlreadyExists).
				With(`package path`, id.FileMod.PackagePath()).
				With(`original pos`, id.Start()).
				With(`augmenter pos`, a.fileSet.Position(pos)).
				With(`identifier`, id.Name)); err != nil {
				return err
			}
		}
	}

	// TODO: Check for any import to make sure it doesn't exist.

	for id := range fm.Idents() {
		if ts, has := a.addFields[id.Name]; has {
			if id.TypeSpec == nil {

				return fmt.Errorf(`%w: id=%s.%s`, ErrIdNotForType, id.FileMod.PackagePath(), id.Name)
				// TODO: Not a type so can't match the name
			}
			// TODO: Check the type is a struct or identifier and that it matches the ts

			// TODO: Move all fields into type and check they don't have the same name as any existing types
		}
	}

	// Add all declarations needing to be added.
	if len(a.addDecl) > 0 {
		fm.AddDecls(a.addDecl)
		a.addDecl = []ast.Decl{}
	}
	return nil
}

func (a *augAdd) PackageDone(name, path string) error {
	// TODO: Check if there were any type specs not found to add the fields to.
	return nil
}
