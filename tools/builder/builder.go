package builder

import (
	"fmt"
	"os"
	"strings"

	"github.com/Grant-Nelson/Gozer/presets"
	"github.com/Grant-Nelson/Gozer/project/interep"
	"github.com/Grant-Nelson/Gozer/project/loader"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/typeChecker"
)

type Config struct {
	Usage    string   `arg:"help"`
	Lang     string   `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool     `arg:"flag, v|verbose, Indicates status information should be printed while building."`
	Output   string   `arg:"flag, o|out, The directory to write the resulting application out to."`
	Patterns []string `arg:"pos, patterns, One or more patterns for the root files for a project."`
	Test     bool     `arg:"skip"` // Set via tester tool
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Builds a project into a specific language at a specific` +
			`output location. The outputted code may be a package, library, ` +
			`or full application that may be used in other projects.`,
		Lang:     presets.DefaultLang,
		Verbose:  false,
		Output:   `./out`,
		Patterns: []string{},
		Test:     false,
	}
}

func Build(cfg *Config) bool {
	switch strings.ToLower(cfg.Lang) {
	case `ts`, `typescript`:
		return buildTs(cfg)
	default:
		fmt.Fprintln(os.Stderr, `Unknown language selected: %q`, cfg.Lang)
		return false
	}
}

func buildTs(cfg *Config) bool {
	build := []string{`ts`}

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

	loaderCfg := loader.Config{
		Build:     build,
		Patterns:  cfg.Patterns,
		Tests:     cfg.Test,
		Modifiers: mods,
	}

	proj, err := loader.Load(loaderCfg)
	if err != nil {
		// TODO: log error
		return false
	}

	interep.Remodel(proj)

}
