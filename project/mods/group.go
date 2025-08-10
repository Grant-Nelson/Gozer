package mods

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type Group []Modifier

func (group Group) Modify(fm *fileMod.FileMod, errGroup *faults.Group) error {
	for _, mod := range group {
		if err := mod.Modify(fm, errGroup); err != nil {
			if errors.Is(err, ErrFileModDone) {
				return nil
			}
			return err
		}
	}
	return nil
}

func (group Group) PackageStart(name, path string, errGroup *faults.Group) error {
	for _, mod := range group {
		if m, ok := mod.(PackageStartExt); ok {
			if err := m.PackageStart(name, path, errGroup); err != nil {
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func (group Group) PackageDone(name, path string, errGroup *faults.Group) error {
	for _, mod := range group {
		if m, ok := mod.(PackageDoneExt); ok {
			if err := m.PackageDone(name, path, errGroup); err != nil {
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
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
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
				return err
			}
		}
	}
	return nil
}
