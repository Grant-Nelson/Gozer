package mods

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type Group []Modifier

func (group Group) Modify(f *artifacts.File, errGroup *faults.Group) error {
	for _, mod := range group {
		if err := mod.Modify(f, errGroup); err != nil {
			return err
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
