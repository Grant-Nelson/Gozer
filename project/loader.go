package project

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// TODO: Add Modifiers:
//  - that performs overlays like gopherjs
//  - that runs per function or func lit
//  - to simplify constants
//  - to remove defers into a `deferBlock` call
//  - to remove Goto and labels (aka flatten)
//  - to join initialization for a package
//  - to generate return structures for multiple returns
//  - to replace multiple assignments with a `multiAssign` call
//  - to flatten select statements and switches as needed
//  - to adjust imports
// TODO: Need post processing for determining things like inheritance

type Config struct {

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string

	// BuildFlags is a list of command-line flags to be passed through to
	// the build system's query tool.
	BuildFlags []string

	// Patterns are the file patterns to use for the project root.
	Patterns []string

	// Tests indicates the test files will also be read for the patterns.
	Tests bool

	// Overlay is a mapping from absolute file paths to file contents.
	Overlay map[string][]byte

	// Modifiers to process each file with.
	Modifiers []Modifier
}

// Modifier performs a set changes to the given file.
type Modifier interface {
	Modify(file *File) error
}

func Load(cfg Config) (*Project, error) {
	fSet := &token.FileSet{}
	ld := &loader{cfg: cfg}
	c := &packages.Config{
		Mode:       allNeeds,
		Dir:        cfg.Dir,
		BuildFlags: cfg.BuildFlags,
		ParseFile:  ld.parseFile,
		Fset:       fSet,
		Tests:      cfg.Tests,
		Overlay:    cfg.Overlay,
	}
	packages, err := packages.Load(c, cfg.Patterns...)
	if err != nil {
		return nil, err
	}
	proj := &Project{
		fSet:     fSet,
		packages: packages,
	}
	return proj, nil
}

const allNeeds = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedExportFile |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

type loader struct {
	cfg Config
}

func (ld *loader) parseFile(fSet *token.FileSet, filename string, src []byte) (*ast.File, error) {
	file, err := initFile(filename, src)
	if err != nil {
		return nil, err
	}

	for _, mod := range ld.cfg.Modifiers {
		if err := mod.Modify(file); err != nil {
			return nil, err
		}
	}

	return file.finalize(fSet)
}
