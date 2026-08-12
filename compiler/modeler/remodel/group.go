package remodel

import (
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type (
	Group        []RemodelFactory
	RemodelGroup []Remodeler
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

func (rg RemodelGroup) PackageDone() (bool, error) {
	for _, m := range rg {
		if con, err := m.PackageDone(); err != nil || !con {
			return false, err
		}
	}
	return true, nil
}

func (rg RemodelGroup) RemodelFunc(f *ir.Func) (bool, error) {
	for _, factory := range rg {
		if m, ok := factory.(RemodelFuncExt); ok {
			if con, err := m.RemodelFunc(f); err != nil || !con {
				return false, err
			}
		}
	}
	return true, nil
}
