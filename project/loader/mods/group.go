package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/project"
)

type (
	Group    []ModFactory
	ModGroup []Modifier
)

func (group Group) StartPackage(pkg *project.Package) (bool, Modifier, error) {
	mg := make(ModGroup, 0, len(group))
	for _, factory := range group {
		con, mod, err := factory.StartPackage(pkg)
		if err != nil || !con {
			return false, nil, err
		}
		if mod != nil {
			mg = append(mg, mod)
		}
	}
	if len(mg) <= 0 {
		return true, nil, nil
	}
	return true, mg, nil
}

func (mg ModGroup) PackageDone() (bool, error) {
	for _, mod := range mg {
		if con, err := mod.PackageDone(); err != nil || !con {
			return false, err
		}
	}
	return true, nil
}

func (mg ModGroup) ModifyAstFile(f *ast.File) (bool, error) {
	for _, mod := range mg {
		if m, ok := mod.(ModifyAstFileExt); ok {
			if con, err := m.ModifyAstFile(f); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}
