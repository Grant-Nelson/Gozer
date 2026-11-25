package ssaBuilder

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type Config struct {
	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

func New(cfg *Config) *SsaBuilder {
	return &SsaBuilder{errGroup: cfg.ErrGroup}
}

var _ remodel.RemodelFactory = (*SsaBuilder)(nil)
var _ remodel.Remodeler = (*ssaRemodel)(nil)

type SsaBuilder struct {
	errGroup *faults.ErrGroup
}

func (sb *SsaBuilder) StartPackage(pkg *project.Package) (bool, remodel.Remodeler, error) {
	sr := &ssaRemodel{
		errGroup: sb.errGroup,
		pkg:      pkg,
	}
	return true, sr, nil
}

type ssaRemodel struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
}

func (sr *ssaRemodel) Remodel() (bool, error) {
	sr.pkg.Ssa.Build()
	return true, nil
}
