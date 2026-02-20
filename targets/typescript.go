package targets

import (
	"fmt"
	"os"

	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/interep"
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

	loaderCfg := loader.Config{
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
		cfg.ErrGroup.Add(fmt.Errorf(`Listing failed: %w`, err))
	}
	return proj, cfg.ErrGroup.ErrorOrNil()
}

func (ts typeScriptTarget) load(cfg *BuildConfig) (*project.Project, error) {
	defer cfg.Logger.LogGroup("Loading %s", ts.Language())()

	mods := mods.Group{
		pkgDropper.New(&pkgDropper.Config{
			ErrGroup: cfg.ErrGroup,
		}),
		//cache.New(&cache.Config{
		//	Build: build,
		//	ErrGroup: cfg.ErrGroup,
		//	//Converter: , // TODO:
		//}),
		//augmenter.New(&augmenter.Config{
		//	Build: build,
		//	ErrGroup: cfg.ErrGroup,
		//	//Converter: , // TODO:
		//}),
		typeChecker.New(&typeChecker.Config{
			ErrGroup: cfg.ErrGroup,
		}),
	}

	// Load all packages for this project.
	loaderCfg := loader.Config{
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
		cfg.ErrGroup.Add(fmt.Errorf(`Loading failed: %w`, err))
	}
	return proj, cfg.ErrGroup.ErrorOrNil()
}

func (ts typeScriptTarget) Build(cfg *BuildConfig) error {
	defer cfg.Logger.LogGroup("Building %s", ts.Language())()

	proj, err := ts.load(cfg)
	if err != nil {
		return nil
	}

	// Remodel any packages that need to be compiled into the intermediate form.
	remodelCfg := &interep.Config{
		Logger:   cfg.Logger,
		ErrGroup: cfg.ErrGroup,
	}
	for pkg := range proj.UnfinishedPackages() {
		// TODO: Use work group to run several of these in parallel when [cfg.Parallel] is `true`.
		if err := interep.Remodel(pkg, remodelCfg); err != nil {
			return err
		}
	}

	// Compile of packages that need to be compiled.
	// TODO: Finish

	return cfg.ErrGroup.ErrorOrNil()
}
