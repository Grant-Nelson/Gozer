package mods

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

// ErrFileModDone can be returned from a modifier to skip running any following
// modifiers and finishes loading a file.
var ErrFileModDone = errors.New(`file modification is done`)

// Modifier performs a set changes to the given file.
type Modifier interface {
	Modify(fm *fileMod.FileMod) error
	Finished() error
}

func Modify(fm *fileMod.FileMod, mods ...Modifier) error {
	for _, mod := range mods {
		if err := mod.Modify(fm); err != nil {
			return err
		}
	}
	return nil
}

func Finished(mods ...Modifier) error {
	for _, mod := range mods {
		if err := mod.Finished(); err != nil {
			return err
		}
	}
	return nil
}
