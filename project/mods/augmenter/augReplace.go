package augmenter

import "github.com/Grant-Nelson/Gozer/project/fileMod"

type augReplace struct {
}

func (a *augReplace) Modify(fm *fileMod.FileMod) error {
	// TODO: Implement
	return nil
}

func (a *augReplace) PackageDone(name, path string) error {
	// TODO: Implement
	return nil
}
