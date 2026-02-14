package tester

import (
	"fmt"
	"time"

	"github.com/Grant-Nelson/Gozer/presets"
	"github.com/Grant-Nelson/Gozer/tools/builder"
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
	output := `./temp` // TODO: Should use a temp directory
	if !build(cfg, output) {
		return false
	}

	// TODO: Implement
	fmt.Println(`Test is not implemented yet.`)
	fmt.Printf("\tConfig was %#v\n", cfg)
	return false
}

func build(cfg *Config, output string) bool {
	bCfg := builder.DefaultConfig()
	bCfg.Lang = cfg.Lang
	bCfg.Output = output
	bCfg.Patterns = cfg.Patterns
	bCfg.Verbose = cfg.Verbose
	bCfg.Test = true
	return builder.Build(bCfg)
}
