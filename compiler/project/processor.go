package project

import (
	"slices"
	"sync"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type ProcessorHandle func(*Package) error

type Processor struct {
	Parallel bool
	Packages []*Package
	ErrGroup *faults.ErrGroup
	Handle   ProcessorHandle
}

func (pr Processor) Process() (err error) {
	defer pr.ErrGroup.Recover(&err)

	ph := &processorHelper{
		Packages: slices.Clone(pr.Packages),
		ErrGroup: pr.ErrGroup,
		Handle:   pr.Handle,
	}

	sortPackagesByDepth(ph.Packages)
	if pr.Parallel {
		ph.parallel()
	} else {
		ph.serial()
	}
	return pr.ErrGroup.AnyOrNil()
}

type processorHelper struct {
	Packages []*Package
	ErrGroup *faults.ErrGroup
	Handle   ProcessorHandle
}

func (pr *processorHelper) parallel() {
	prev, depth := 0, 0
	for i, pkg := range pr.Packages {
		if pkg.Depth != depth {
			pr.parallelGroup(pr.Packages[prev:i])
			if !pr.ErrGroup.Empty() {
				return
			}
			prev, depth = i, pkg.Depth
		}
	}
	pr.parallelGroup(pr.Packages[prev:])
}

func (pr *processorHelper) parallelGroup(packages []*Package) {
	wg := &sync.WaitGroup{}
	wg.Add(len(packages))
	for _, pkg := range packages {
		go func(pkg *Package) {
			defer wg.Done()
			defer pr.ErrGroup.Recover(nil)
			if err := pr.Handle(pkg); err != nil {
				pr.ErrGroup.Add(err)
			}
		}(pkg)
	}
	wg.Wait()
}

func (pr *processorHelper) serial() {
	for _, pkg := range pr.Packages {
		if err := pr.Handle(pkg); err != nil {
			pr.ErrGroup.Add(err)
			return
		}
	}
}
