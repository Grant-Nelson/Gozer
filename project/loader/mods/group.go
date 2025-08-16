package mods

import (
	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
)

type Group []Modifier

func (group Group) Modify(f *file.File, errGroup *faults.Group) error {
	for _, mod := range group {
		if err := mod.Modify(f, errGroup); err != nil {
			return err
		}
	}
	return nil
}

func (group Group) PackageStart(pkg *Package, errGroup *faults.Group) error {
	for _, mod := range group {
		if m, ok := mod.(PackageStartExt); ok {
			if err := m.PackageStart(pkg, errGroup); err != nil {
				return err
			}
		}
	}
	return nil
}

func (group Group) PackageDone(pkg *Package, errGroup *faults.Group) error {
	for _, mod := range group {
		if m, ok := mod.(PackageDoneExt); ok {
			if err := m.PackageDone(pkg, errGroup); err != nil {
				return err
			}
		}
	}
	return nil
}

func (group Group) LoadDone(errGroup *faults.Group) error {
	for _, mod := range group {
		if m, ok := mod.(LoadDoneExt); ok {
			if err := m.LoadDone(errGroup); err != nil {
				return err
			}
		}
	}
	return nil
}
