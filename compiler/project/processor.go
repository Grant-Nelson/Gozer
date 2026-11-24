package project

import (
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

func (pr *Processor) Process() (err error) {
	defer pr.ErrGroup.Recover(&err)
	sortPackagesByDepth(pr.Packages)
	if pr.Parallel {
		pr.parallelProcessProject()
	} else {
		pr.serialProcessProject()
	}
	return pr.ErrGroup.AnyOrNil()
}

func (pr *Processor) parallelProcessProject() {
	prev, depth := 0, 0
	for i, pkg := range pr.Packages {
		if pkg.Depth != depth {
			pr.parallelProcessGroup(pr.Packages[prev:i])
			if !pr.ErrGroup.Empty() {
				return
			}
			prev, depth = i, pkg.Depth
		}
	}
	pr.parallelProcessGroup(pr.Packages[prev:])
}

func (pr *Processor) parallelProcessGroup(packages []*Package) {
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

func (pr *Processor) serialProcessProject() {
	for _, pkg := range pr.Packages {
		if err := pr.Handle(pkg); err != nil {
			pr.ErrGroup.Add(err)
			return
		}
	}
}
