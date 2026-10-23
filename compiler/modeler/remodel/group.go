package remodel

import (
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type (
	Group        []RemodelFactory
	RemodelGroup []Remodeler
)

var (
	_ RemodelFactory = Group{}
	_ Remodeler      = RemodelGroup{}
	_ ProjectDoneExt = RemodelGroup{}
)

func (group Group) StartPackage(pkg *project.Package) (bool, Remodeler, error) {
	rg := make(RemodelGroup, 0, len(group))
	for _, factory := range group {
		con, m, err := factory.StartPackage(pkg)
		if err != nil || !con {
			return false, nil, err
		}
		if m != nil {
			rg = append(rg, m)
		}
	}
	if len(rg) <= 0 {
		return true, nil, nil
	}
	return true, rg, nil

}

func (rg RemodelGroup) Remodel() (bool, error) {
	for _, m := range rg {
		if con, err := m.Remodel(); err != nil || !con {
			return false, err
		}
	}
	return true, nil
}

func (rg RemodelGroup) PackageDone() (bool, error) {
	for _, factory := range rg {
		if m, ok := factory.(ProjectDoneExt); ok {
			if con, err := m.PackageDone(); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}
