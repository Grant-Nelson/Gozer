package interep

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
)

type Config struct {
}

func Remodel(pkg *project.Package, cfg *Config) (err error) {
	errGroup := faults.NewGroup(-1)
	defer faults.Recover(&err)

	return errGroup.Wrap()
}
