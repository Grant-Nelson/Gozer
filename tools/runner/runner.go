package runner

import (
	"fmt"
	"os"
	"time"

	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/targets"
)

type Config struct {
	Usage    string         `arg:"help"`
	Lang     string         `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool           `arg:"flag, v|verbose, Indicates status information should be printed while building and running."`
	Timeout  *time.Duration `arg:"flag, t|timeout, The amount of time to let the transpiled application run until cancelling it."`
	Patterns []string       `arg:"pos, patterns, One or more patterns for the root files for a project."`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Builds a project then runs it. If needed, the run will use ` +
			`a tool specific to the transpiled language, e.g. node.js. ` +
			`This application will not end until the running program ends ` +
			`or until a set timeout it hit.`,
		Lang: targets.DefaultLang,
	}
}

func Run(cfg *Config) bool {
	// TODO: Validate configs

	// TODO: Add optional build flags to config
	// TODO: Add optional parallel flags to config

	buildCfg := &targets.BuildConfig{
		Lang:     cfg.Lang,
		Patterns: cfg.Patterns,
		Logger:   logger.New(cfg.Verbose, nil),
		Tests:    false,
	}
	if err := targets.Build(buildCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Run failed when building: %v\n", err)
	}

	// TODO: Finish implementing
	fmt.Println(`Run is not implemented yet.`)
	return false
}
