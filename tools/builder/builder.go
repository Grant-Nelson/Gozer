package builder

import (
	"fmt"
	"os"

	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/targets"
)

type Config struct {
	Usage    string   `arg:"help"`
	Lang     string   `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool     `arg:"flag, v|verbose, Indicates status information should be printed while building."`
	Output   string   `arg:"flag, o|out, The optional directory to copy the resulting application out to."`
	Patterns []string `arg:"pos, patterns, One or more patterns for the root files for a project."`
	Tests    bool     `arg:"skip"` // Set via tester tool
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Builds a project into a specific language at a specific` +
			`output location. The outputted code may be a package, library, ` +
			`or full application that may be used in other projects.`,
		Lang:     targets.DefaultLang,
		Verbose:  false,
		Output:   `./out`,
		Patterns: []string{},
		Tests:    false,
	}
}

func Build(cfg *Config) {
	// TODO: Validate configs

	// TODO: Add optional build flags to config
	// TODO: Add optional parallel flags to config

	buildCfg := &targets.BuildConfig{
		Lang:     cfg.Lang,
		Logger:   logger.New(cfg.Verbose),
		Output:   cfg.Output,
		Patterns: cfg.Patterns,
		Tests:    cfg.Tests,
	}
	if err := targets.Build(buildCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
	}
}
