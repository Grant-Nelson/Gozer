package augmenter

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

var ErrIdAlreadyExists = errors.New(`can not add new identifier via augmenter: Identifier already exist`)

type augAdd struct {
	addIds    map[string]bool
	addDecl   []ast.Decl
	addFields map[string]*ast.TypeSpec
}

func (a *augAdd) Modify(fm *fileMod.FileMod) error {
	for it := range fm.IdentIter() {
		if a.addIds[it.Name] {
			return fmt.Errorf(`%w: id=%q`, ErrIdAlreadyExists, it.Name)
		}
	}

	// TODO: Check for any import to make sure it doesn't exist.
	// TODO: Find the type spec to add the field to.

	fm.AddDecls(a.addDecl)
	a.addDecl = []ast.Decl{}
	return nil
}

func (a *augAdd) checkName(name string) error {

	return nil
}

func (a *augAdd) PackageDone(name, path string) error {
	// TODO: Check if there were any type specs not found to add the fields to.
	return nil
}
