package targets

import (
	"os"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/interRep"
	"github.com/Grant-Nelson/Gozer/project/loader"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/pkgDropper"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/typeChecker"
)

type typeScriptTarget struct{}

func (ts typeScriptTarget) Language() string { return `typescript` }

func (ts typeScriptTarget) Aliases() []string { return []string{`typescript`, `ts`} }

func (ts typeScriptTarget) Env() []string {
	return append(os.Environ(), `GOOS=js`, `GOARCH=ecmascript`)
}

func (ts typeScriptTarget) List(cfg *ListConfig) (*project.Project, error) {
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

func (ts typeScriptTarget) Build(cfg *BuildConfig) error {
	defer cfg.Logger.LogGroup("Building %s", ts.Language())()

	// Load project into AST form and type check.
	proj, err := ts.load(cfg)
	if err != nil {
		return nil
	}

	// Remodel any packages that need to be compiled into the intermediate form.
	if cfg.Parallel {
		err = ts.asyncFinishPackages(proj, cfg)
	} else {
		err = ts.syncFinishPackages(proj, cfg)
	}
	if err != nil {
		return err
	}

	// Compile of packages that need to be compiled.
	// TODO: Finish

	return cfg.ErrGroup.AnyOrNil()
}

func (ts typeScriptTarget) load(cfg *BuildConfig) (*project.Project, error) {
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
	proj, err := loader.Load(loaderCfg)
	if err != nil {
		cfg.ErrGroup.Add(faults.New(`loading failed: %w`, err))
	}
	return proj, cfg.ErrGroup.AnyOrNil()
}

func (ts typeScriptTarget) asyncFinishPackages(proj *project.Project, cfg *BuildConfig) error {
	// TODO: Use work group to run several of these in parallel.
	for pkg := range proj.UnfinishedPackages() {
		if err := ts.finishPackage(pkg, cfg); err != nil {
			return err
		}
	}
	return nil
}

func (ts typeScriptTarget) syncFinishPackages(proj *project.Project, cfg *BuildConfig) error {
	for pkg := range proj.UnfinishedPackages() {
		if err := ts.finishPackage(pkg, cfg); err != nil {
			return err
		}
	}
	return nil
}

func (ts typeScriptTarget) finishPackage(pkg *project.Package, cfg *BuildConfig) error {
	ircCfg := &interRep.Config{
		Logger:   cfg.Logger,
		ErrGroup: cfg.ErrGroup,
		Package:  pkg,
	}
	if err := interRep.Remodel(ircCfg); err != nil {
		return err
	}

	// TODO: Finish

	return nil
}
