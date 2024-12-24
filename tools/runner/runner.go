package runner

import (
	"fmt"
	"time"

	"github.com/Grant-Nelson/Gozer/tools/builder"
)

type Config struct {
	Usage    string         `arg:"help"`
	Lang     string         `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool           `arg:"flag, v|verbose, Indicates status information should be printed while building and running."`
	Timeout  *time.Duration `arg:"flag, t|timeout, The amount of time to let the transpiled application run until cancelling it."`
	Patterns []string       `arg:"pos, patterns, One or more patterns for the root files for a project."`
}

func DefaultConfig(defaultLang string) *Config {
	return &Config{
		Usage: `Builds a project then runs it. If needed, the run will use ` +
			`a tool specific to the transpiled language, e.g. node.js. ` +
			`This application will not end until the running program ends ` +
			`or until a set timeout it hit.`,
		Lang: defaultLang,
	}
}

func Run(cfg *Config) bool {
	output := `./temp` // TODO: Should use a temp directory
	if !build(cfg, output) {
		return false
	}

	// TODO: Implement
	fmt.Println(`Run is not implemented yet.`)
	fmt.Printf("\tConfig was %#v\n", cfg)
	return false
}

func build(cfg *Config, output string) bool {
	bCfg := builder.DefaultConfig(cfg.Lang)
	bCfg.Lang = cfg.Lang
	bCfg.Output = output
	bCfg.Patterns = cfg.Patterns
	bCfg.Verbose = cfg.Verbose
	bCfg.Test = false
	return builder.Build(bCfg)
}
