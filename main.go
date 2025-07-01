package main

import (
	"fmt"
	"os"

	"github.com/Grant-Nelson/Gozer/internal/args"
	"github.com/Grant-Nelson/Gozer/tools/builder"
	"github.com/Grant-Nelson/Gozer/tools/listLangs"
	"github.com/Grant-Nelson/Gozer/tools/runner"
	"github.com/Grant-Nelson/Gozer/tools/tester"
	"github.com/Grant-Nelson/Gozer/tools/version"
)

type mainConfig struct {
	Usage     string            `arg:"help"`
	Builder   *builder.Config   `arg:"tool, b|build, Builds a project."`
	Runner    *runner.Config    `arg:"tool, r|run, Builds a project then runs it."`
	Tester    *tester.Config    `arg:"tool, t|test, Builds and tests all or part of a project"`
	Version   *version.Config   `arg:"tool, v|version, Shows the version."`
	ListLangs *listLangs.Config `arg:"tool, list, Shows the list of languages available to transpile into."`
}

func defaultConfig(defaultLang string) *mainConfig {
	return &mainConfig{
		Usage: `Gozer cross compiles from Go into other languages, transpilation. ` +
			`Complicated code should be tested, vetted, and written once ` +
			`so that updates and fixes are properly integrated. By letting ` +
			`Gozer transpile allows complicated code be written once in Go ` +
			`and reused in several other languages.` +
			`To use Gozer, select a tool to build, run, test, etc a project.`,
		Builder:   builder.DefaultConfig(defaultLang),
		Runner:    runner.DefaultConfig(defaultLang),
		Tester:    tester.DefaultConfig(defaultLang),
		Version:   version.DefaultConfig(),
		ListLangs: listLangs.DefaultConfig(),
	}
}

func main() {
	cfg := defaultConfig(`dart`)
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

	case cfg.ListLangs != nil:
		listLangs.ListLangs(cfg.ListLangs)

	default:
		fmt.Fprintln(os.Stderr, `Must select a tool to use.`)
		fmt.Fprintf(os.Stderr, "Use %q to print help.\n", os.Args[0]+` -h`)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
