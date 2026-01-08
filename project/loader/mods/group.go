package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Group []Modifier

func (group Group) Modify(f *ast.File, errGroup *faults.Group) (bool, error) {
	for _, mod := range group {
		if con, err := mod.Modify(f, errGroup); err != nil || !con {
			return false, err
		}
	}
	return true, nil
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
