package tester

import (
	"fmt"
	"os"
	"time"

	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/presets"
	"github.com/Grant-Nelson/Gozer/project/builder"
)

type Config struct {
	Usage    string         `arg:"help"`
	Lang     string         `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool           `arg:"flag, v|verbose, Indicates status information should be printed while building and testing."`
	Timeout  *time.Duration `arg:"flag, t|timeout, The amount of time to let the tests run until cancelling it."`
	Patterns []string       `arg:"pos, patterns, One or more patterns for the root files for a project."`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Builds and tests all or part of a project. If needed, the ` +
			`tests will be run using a tool specific to the transpiled ` +
			`language. The results of the tests will be outputted to the ` +
			`console similarly to how Go's tests are outputted.`,
		Lang: presets.DefaultLang,
	}
}

func Test(cfg *Config) bool {
	// TODO: Validate configs

	// TODO: Add optional build flags to config
	// TODO: Add optional parallel flags to config

	buildCfg := &builder.Config{
		Lang:     cfg.Lang,
		Patterns: cfg.Patterns,
		Logger:   logger.New(cfg.Verbose),
		Tests:    true,
	}
	if err := builder.Build(buildCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Tests failed when building: %v\n", err)
	}

	// TODO: Finish implementing
	fmt.Println(`Run is not implemented yet.`)
	return false
}
