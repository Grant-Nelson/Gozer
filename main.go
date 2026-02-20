package main

import (
	"fmt"
	"os"

	"github.com/Grant-Nelson/Gozer/avail/args"
	"github.com/Grant-Nelson/Gozer/tools/builder"
	"github.com/Grant-Nelson/Gozer/tools/langs"
	"github.com/Grant-Nelson/Gozer/tools/list"
	"github.com/Grant-Nelson/Gozer/tools/runner"
	"github.com/Grant-Nelson/Gozer/tools/tester"
	"github.com/Grant-Nelson/Gozer/tools/version"
)

type mainConfig struct {
	Usage   string          `arg:"help"`
	Builder *builder.Config `arg:"tool, b|build, Builds a project."`
	Runner  *runner.Config  `arg:"tool, r|run, Builds a project then runs it."`
	Tester  *tester.Config  `arg:"tool, t|test, Builds and tests all or part of a project"`
	Version *version.Config `arg:"tool, v|version, Shows the version."`
	List    *list.Config    `arg:"tool, l|list, Shows the list of files used when building."`
	Langs   *langs.Config   `arg:"tool, langs, Shows the list of languages available to transpile into."`
}

// TODO: Add 'env' for printing environment variables that are used.
// TODO: Add 'serve' to serve a website

func defaultConfig() *mainConfig {
	return &mainConfig{
		Usage: `Gozer cross compiles from Go into other languages, transpilation. ` +
			`Complicated code should be tested, vetted, and written once ` +
			`so that updates and fixes are properly integrated. By letting ` +
			`Gozer transpile allows complicated code be written once in Go ` +
			`and reused in several other languages.` +
			`To use Gozer, select a tool to build, run, test, etc a project.`,
		Builder: builder.DefaultConfig(),
		Runner:  runner.DefaultConfig(),
		Tester:  tester.DefaultConfig(),
		Version: version.DefaultConfig(),
		List:    list.DefaultConfig(),
		Langs:   langs.DefaultConfig(),
	}
}

func main() {
	cfg := defaultConfig()
	if !args.Parse(cfg) {
		os.Exit(1)
	}

	switch {
	case cfg.Builder != nil:
		builder.Build(cfg.Builder)

	case cfg.Runner != nil:
		runner.Run(cfg.Runner)

	case cfg.Tester != nil:
		tester.Test(cfg.Tester)

	case cfg.Version != nil:
		version.Version(cfg.Version)

	case cfg.List != nil:
		list.List(cfg.List)

	case cfg.Langs != nil:
		langs.Langs(cfg.Langs)

	default:
		fmt.Fprintln(os.Stderr, `Must select a tool to use.`)
		fmt.Fprintf(os.Stderr, "Use %q to print help.\n", os.Args[0]+` -h`)
		os.Exit(1)
	}
}
