package builder

import (
	"github.com/Grant-Nelson/Gozer/project/interep"
	"github.com/Grant-Nelson/Gozer/project/loader"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/typeChecker"
)

func buildTS(cfg *Config) error {
	defer cfg.Logger.LogGroup("Building typescript")()

	build := append(cfg.Build, `ts`)

	// Prepared modifiers to parse project for typescript.
	mods := mods.Group{
		//cache.New(&cache.Config{
		//	Build: build,
		//	//Converter: , // TODO:
		//}),
		//augmenter.New(&augmenter.Config{
		//	Build: build,
		//	//Converter: , // TODO:
		//}),
		typeChecker.New(nil),
	}

	// Load all packages for this project.
	loaderCfg := loader.Config{
		Logger:    cfg.Logger,
		Dir:       cfg.Dir,
		Build:     build,
		Patterns:  cfg.Patterns,
		Tests:     cfg.Tests,
		Overlay:   cfg.Overlay,
		Modifiers: mods,
		Parallel:  cfg.Parallel,
	}
	proj, err := loader.Load(loaderCfg)
	if err != nil {
		return err
	}

	// Remodel any packages that need to be compiled into the intermediate form.
	remodelCfg := &interep.Config{}
	for pkg := range proj.UnfinishedPackages() {
		// TODO: Use work group to run several of these in parallel when [cfg.Parallel] is `true`.
		if err := interep.Remodel(pkg, remodelCfg); err != nil {
			return err
		}
	}

	// Compile of packages that need to be compiled.
	// TODO: Finish

	return nil
}
