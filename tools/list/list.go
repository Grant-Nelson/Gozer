package list

import (
	"fmt"
	"os"

	"github.com/Grant-Nelson/Gozer/targets"
)

type Config struct {
	Usage    string   `arg:"help"`
	Lang     string   `arg:"flag, l|lang|language, The language to transpile into."`
	Patterns []string `arg:"pos, patterns, One or more patterns for the root files for a project."`
	Tests    bool     `arg:"flag, t|test, Indicate the transpile is for a test."`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Prints the list of files that would be used for the given build. ` +
			`This will not include augmentation files.`,
		Lang:     targets.DefaultLang,
		Patterns: []string{},
		Tests:    false,
	}
}

func List(cfg *Config) {
	listCfg := &targets.ListConfig{
		Lang:     cfg.Lang,
		Patterns: cfg.Patterns,
		Tests:    cfg.Tests,
	}
	proj, err := targets.List(listCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "List failed: %v\n", err)
	}

	nodes := make([]*node, len(proj.AllPackages))
	for i, pkg := range proj.AllPackages {
		nodes[i] = &node{
			Path: pkg.Ast.PkgPath,
			Name: pkg.Ast.Name,

			// TODO: Finish
		}

	}

	// TODO: Implement
}

type node struct {
	Root    bool
	Path    string
	Name    string
	Imports []string
	Files   []string
}
