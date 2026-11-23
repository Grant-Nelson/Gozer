package typescript

import (
	"os"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/loader"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods/pkgDropper"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods/ssaBuilder"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods/typeChecker"
	"github.com/Grant-Nelson/Gozer/compiler/modeler"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel/blocker"
	"github.com/Grant-Nelson/Gozer/compiler/project"
	"github.com/Grant-Nelson/Gozer/targets/target"
	"github.com/Grant-Nelson/Gozer/targets/typescript/writer"
)

type tsTarget struct{}

func Target() target.Target {
	return tsTarget{}
}

func (ts tsTarget) Language() string { return `typescript` }

func (ts tsTarget) Aliases() []string { return []string{`typescript`, `ts`} }

func (ts tsTarget) Env() []string {
	return append(os.Environ(), `GOOS=js`, `GOARCH=ecmascript`)
}

func (ts tsTarget) List(cfg *target.ListConfig) (*project.Project, error) {
	defer cfg.Logger.LogGroup("Listing %s", ts.Language())()
	loaderCfg := &loader.Config{
		Logger:          cfg.Logger,
		ErrGroup:        cfg.ErrGroup,
		Dir:             cfg.Dir,
		Build:           cfg.Build,
		Patterns:        cfg.Patterns,
		Tests:           cfg.Tests,
		Env:             ts.Env(),
		SkipFileParsing: true,
	}
	proj, err := loader.Load(loaderCfg)
	if err != nil {
		cfg.ErrGroup.Add(faults.New(`listing failed`, err))
	}
	return proj, cfg.ErrGroup.AnyOrNil()
}

func (ts tsTarget) Build(cfg *target.BuildConfig) error {
	defer cfg.Logger.LogGroup("Building %s", ts.Language())()

	// Load project into AST form and type check.
	proj, err := ts.load(cfg)
	if err != nil {
		return err
	}

	// Model any packages that need to be compiled into the intermediate form.
	if err := ts.model(proj, cfg); err != nil {
		return err
	}

	// Compile the packages that need to be compiled.
	// TODO: Finish

	return cfg.ErrGroup.AnyOrNil()
}

func (ts tsTarget) load(cfg *target.BuildConfig) (proj *project.Project, err error) {
	defer cfg.ErrGroup.Recover(&err)
	defer cfg.Logger.LogGroup("Loading %s", ts.Language())()

	mods := mods.Group{
		pkgDropper.New(&pkgDropper.Config{
			Logger:   cfg.Logger,
			ErrGroup: cfg.ErrGroup,
			PkgPathPatterns: []string{
				`runtime`,
			},
		}),

		//cache.New(&cache.Config{
		//	Build: config.Build,
		//	ErrGroup: config.ErrGroup,
		//	//Converter: , // TODO:
		//}),

		//augmenter.New(&augmenter.Config{
		//	Build:    config.Build,
		//	ErrGroup: config.ErrGroup,
		//	//Converter: , // TODO:
		//}),

		typeChecker.New(&typeChecker.Config{
			ErrGroup: cfg.ErrGroup,
		}),

		ssaBuilder.New(&ssaBuilder.Config{
			ErrGroup: cfg.ErrGroup,
		}),
	}

	// Load all packages for this project.
	loaderCfg := &loader.Config{
		Logger:    cfg.Logger,
		ErrGroup:  cfg.ErrGroup,
		Dir:       cfg.Dir,
		Build:     cfg.Build,
		Patterns:  cfg.Patterns,
		Tests:     cfg.Tests,
		Overlay:   cfg.Overlay,
		Modifiers: mods,
		Parallel:  cfg.Parallel,
		Env:       ts.Env(),
	}
	proj, err = loader.Load(loaderCfg)
	if err != nil {
		cfg.ErrGroup.Add(faults.New(`loading failed: %w`, err))
	}
	return proj, cfg.ErrGroup.AnyOrNil()
}

func (ts tsTarget) model(proj *project.Project, cfg *target.BuildConfig) (err error) {
	defer cfg.ErrGroup.Recover(&err)
	defer cfg.Logger.LogGroup("Modelling %s", ts.Language())()

	remodelers := remodel.Group{
		blocker.New(&blocker.Config{
			ErrGroup: cfg.ErrGroup,
		}),

		// TODO: Add optimizers

		writer.New(&writer.Config{
			ErrGroup: cfg.ErrGroup,
		}),

		// TODO: Write cache with result
	}

	modelerCfg := &modeler.Config{
		Logger:     cfg.Logger,
		ErrGroup:   cfg.ErrGroup,
		Proj:       proj,
		Remodelers: remodelers,
		Parallel:   cfg.Parallel,
	}
	err = modeler.Model(modelerCfg)
	if err != nil {
		cfg.ErrGroup.Add(faults.New(`modelling failed: %w`, err))
	}
	return cfg.ErrGroup.AnyOrNil()
}
