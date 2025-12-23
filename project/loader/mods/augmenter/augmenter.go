package augmenter

import (
	"maps"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type Augmenter struct {
	packages map[string]*augPackage
	build    []string
	pathConv PathConverter
}

// PathConverter takes the given import path for a package and returns
// the paths to where the augmentation files are found.
// Return false to skip this path.
type PathConverter func(string) (bool, string)

var _ mods.Modifier = (*Augmenter)(nil)
var _ mods.LoadDoneExt = (*Augmenter)(nil)

func New(build []string, pathConv PathConverter) *Augmenter {
	return &Augmenter{
		packages: map[string]*augPackage{},
		build:    build,
		pathConv: pathConv,
	}
}

func (a *Augmenter) Modify(f *artifacts.File, errGroup *faults.Group) (con bool, err error) {
	defer faults.Recover(&err)
	key := f.PackageKey()
	ap, exists := a.packages[key]

	if !exists {
		hasAug, path := a.pathConv(f.PackagePath())
		if !hasAug {
			// store as nil to prevent trying again
			a.packages[key] = nil
			return true, nil
		}

		ap := newPackage(f.Package)
		a.packages[key] = ap

		ar := newReader(ap, errGroup, a.build)
		ar.addPackage(path)
	}
	if ap == nil {
		// path was previously skipped
		return true, nil
	}

	if con, err := ap.Modify(f, errGroup); err != nil || !con {
		return false, err
	}
	return true, nil
}

func (a *Augmenter) LoadDone(errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	keys := slices.Collect(maps.Keys(a.packages))
	sort.Strings(keys)
	for _, key := range keys {
		a.packages[key].LoadDone(errGroup)
	}
	return nil
}
