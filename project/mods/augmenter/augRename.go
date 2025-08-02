package augmenter

import "github.com/Grant-Nelson/Gozer/project/fileMod"

type augRename struct {
}

func (a *augRename) Modify(fm *fileMod.FileMod) error {
	// TODO: Implement
	return nil
}

func (a *augRename) PackageDone(name, path string) error {
	// TODO: Implement
	return nil
}
