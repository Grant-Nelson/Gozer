package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
)

type Group []Modifier

func (group Group) ModName() string { return `Modifier Group` }

func (group Group) ModifyFile(f *ast.File, errGroup *faults.Group) (bool, error) {
	for _, mod := range group {
		if m, ok := mod.(ModifyFileExt); ok {
			if con, err := m.ModifyFile(f, errGroup); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}

func (group Group) PackageStart(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	for _, mod := range group {
		if m, ok := mod.(PackageStartExt); ok {
			if con, err := m.PackageStart(pkg, errGroup); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}

func (group Group) PackageDone(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	for _, mod := range group {
		if m, ok := mod.(PackageDoneExt); ok {
			if con, err := m.PackageDone(pkg, errGroup); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}
