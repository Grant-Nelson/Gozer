package mods

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type Group []Modifier

func (group Group) Modify(fm *fileMod.FileMod) error {
	for _, mod := range group {
		if err := mod.Modify(fm); err != nil {
			if errors.Is(err, ErrFileModDone) {
				return nil
			}
			return err
		}
	}
	return nil
}

func (group Group) PackageStart(name, path string) error {
	for _, mod := range group {
		if m, ok := mod.(PackageStartExt); ok {
			if err := m.PackageStart(name, path); err != nil {
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func (group Group) PackageDone(name, path string) error {
	for _, mod := range group {
		if m, ok := mod.(PackageDoneExt); ok {
			if err := m.PackageDone(name, path); err != nil {
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func (group Group) LoadDone() error {
	for _, mod := range group {
		if m, ok := mod.(LoadDoneExt); ok {
			if err := m.LoadDone(); err != nil {
				if errors.Is(err, ErrFileModDone) {
					return nil
				}
				return err
			}
		}
	}
	return nil
}
