package interep

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/interep/blocks"
)

type Config struct {

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup
}

type Interep struct {
	Funcs []*blocks.Func
}

func Remodel(pkg *project.Package, cfg *Config) (ir *Interep, err error) {
	defer faults.Recover(&err)
	cfg.Logger.LogGroup(`Remodelling %q`, pkg.PkgPath())

	return nil, nil
}
