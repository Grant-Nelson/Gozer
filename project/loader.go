package project

import (
	"errors"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
	"github.com/Grant-Nelson/Gozer/project/mods"
)

// TODO: Add Modifiers:
//  - that runs per function or func lit
//  - to simplify constants
//  - to remove defers into a `deferBlock` call
//  - to remove Goto and labels (aka flatten)
//  - to inject Jumps and labels to replace other flow-controls
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
	Modifiers []mods.Modifier
}

func Load(cfg Config) (*Project, error) {
	fSet := &token.FileSet{}
	ld := &loader{
		group: mods.Group(cfg.Modifiers),
	}
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

	if len(ld.curPkgName) >= 0 {
		ld.group.PackageDone(ld.curPkgName, ld.curPkgPath)
	}
	if err := ld.group.LoadDone(); err != nil {
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
	group      mods.Group
	curPkgPath string
	curPkgName string
}

func (ld *loader) parseFile(fSet *token.FileSet, filename string, src []byte) (*ast.File, error) {
	fm := fileMod.New(filename)
	if err := fm.AddFile(filename, src); err != nil {
		return nil, err
	}

	pkgName, pkgPath := fm.PackageName(), fm.PackagePath()
	if ld.curPkgName != pkgName && ld.curPkgPath != pkgPath {
		if len(ld.curPkgName) >= 0 {
			ld.group.PackageDone(ld.curPkgName, ld.curPkgPath)
		}
		ld.curPkgName, ld.curPkgPath = pkgName, pkgPath
		ld.group.PackageStart(ld.curPkgName, ld.curPkgPath)
	}

	if err := ld.group.Modify(fm); err != nil {
		if !errors.Is(err, mods.ErrFileModDone) {
			return nil, err
		}
	}

	return fm.Finalize(fSet)
}
