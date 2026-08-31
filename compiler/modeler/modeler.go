package modeler

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel"
	"github.com/Grant-Nelson/Gozer/compiler/project"
	"github.com/Grant-Nelson/Gozer/compiler/project/enums/buildState"
)

type Config struct {
	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup

	// Package is the package to create a IR for.
	Package *project.Package

	// Remodelers are the tools to apply to the IR to transform from
	// very similar to AST over to a form needed by the target.
	Remodelers remodel.Group
}

// Model converts the AST in the packages in the given configs
// into the IR form and sets that to the package.
func Model(cfg *Config) (err error) {
	defer faults.Recover(&err)
	cfg.Logger.LogGroup(`Remodelling %q`, cfg.Package.PkgPath())

	pkg := cfg.Package
	pkg.State = buildState.Remodelling

	rm := &modeler{
		logger:   cfg.Logger,
		errGroup: cfg.ErrGroup,
		pkg:      pkg,
	}

	if err := rm.packageStart(cfg.Remodelers); err != nil {
		return err
	}

	if err := rm.remodelPackage(); err != nil {
		return err
	}

	if err := rm.packageDone(); err != nil {
		return err
	}

	pkg.State = buildState.Remodelled
	return rm.errGroup.FullOrNil()
}

type modeler struct {
	logger   *logger.Logger
	errGroup *faults.ErrGroup
	pkg      *project.Package
	group    remodel.Remodeler
}

func (rm *modeler) packageStart(group remodel.Group) error {
	rm.logger.Printf(`Package Starting`)
	defer rm.logger.Indent()()

	con, rg, err := group.StartPackage(rm.pkg)
	if err != nil || !con {
		return rm.errGroup.Add(err)
	}
	rm.group = rg
	return rm.errGroup.FullOrNil()
}

func (rm *modeler) remodelPackage() error {
	if rm.group == nil {
		return nil
	}

	rm.logger.Printf(`Remodelling Package`)
	defer rm.logger.Indent()()

	if _, err := rm.group.Remodel(); err != nil {
		return rm.errGroup.Add(err)
	}
	return rm.errGroup.FullOrNil()
}

func (rm *modeler) packageDone() error {
	if rm.group == nil {
		return nil
	}

	rm.logger.Printf(`Package Finishing`)
	defer rm.logger.Indent()()

	if pd, ok := rm.group.(remodel.ProjectDoneExt); ok {
		if _, err := pd.PackageDone(); err != nil {
			return rm.errGroup.Add(err)
		}
	}
	return rm.errGroup.FullOrNil()
}
