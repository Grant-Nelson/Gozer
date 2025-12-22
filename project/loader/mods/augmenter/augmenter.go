package augmenter

import (
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type Augmenter struct {
	packages map[string]*augPackage
	build    []string
	fileSet  *artifacts.FileSet
}

func New(build []string, fileSet *artifacts.FileSet) *Augmenter {
	a := &Augmenter{
		packages: map[string]*augPackage{},
		build:    build,
		fileSet:  fileSet,
	}
	return a
}

func (a *Augmenter) Modify(f *artifacts.File, errGroup *faults.Group) error {
	key := f.PackageKey()
	ap, exists := a.packages[key]

	if !exists {
		ap, err := a.addPackage(pkg.Path, errGroup)
		if err != nil {
			return err
		}
		a.packages[key] = ap
	}

	if err := ap.Modify(f, errGroup); err != nil {
		return err
	}
	return nil
}

func (a *Augmenter) LoadDone(errGroup *faults.Group) error {
	keys := slices.Collect(maps.Keys(a.packages))
	sort.Strings(keys)
	for _, key := range keys {
		a.packages[key].LoadDone(errGroup)
	}
	return nil
}

func (a *Augmenter) addPackage(path string, errGroup *faults.Group) (p *augPackage, err error) {
	defer faults.Recover(&err)
	ap := newPackage(a.build, path, testPkgPath, a.fileSet)
	ar := &augReader{augPackage: ap, errGroup: errGroup}
	ar.addPackage(path)
	return ap, nil
}
