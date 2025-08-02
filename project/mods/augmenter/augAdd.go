package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type augAdd struct {
	addIds     map[string]bool
	addImports []*ast.ImportSpec
	addDecl    []ast.Decl
	addFields  []*ast.TypeSpec
}

func (a *augAdd) Modify(fm *fileMod.FileMod) error {
	// TODO: Implement
	return nil
}

func (a *augAdd) PackageDone(name, path string) error {
	// TODO: Implement
	return nil
}
