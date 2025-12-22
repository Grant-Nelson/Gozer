package project

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
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

func Load(cfg Config) (*project.Project, error) {
	finalFileSet := token.NewFileSet()
	ld := &loader{
		errGroup:     faults.NewGroup(-1),
		group:        mods.Group(cfg.Modifiers),
		curPkg:       nil,
		tempFileSet:  artifacts.NewFileSet(),
		finalFileSet: artifacts.NewFileSet(),
	}
	c := &packages.Config{
		Mode:       allNeeds,
		Dir:        cfg.Dir,
		BuildFlags: cfg.BuildFlags,
		ParseFile:  ld.parseFile,
		Fset:       finalFileSet,
		Tests:      cfg.Tests,
		Overlay:    cfg.Overlay,
	}
	packages, err := packages.Load(c, cfg.Patterns...)
	if err != nil {
		return nil, err
	}

	ld.packageDone()
	if err := ld.group.LoadDone(ld.errGroup); err != nil {
		return nil, err
	}

	proj := project.New(finalFileSet, packages)
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
	errGroup     *faults.Group
	group        mods.Group
	curPkg       *artifacts.Package
	tempFileSet  *artifacts.FileSet
	finalFileSet *artifacts.FileSet
}

func (ld *loader) parseFile(fs *token.FileSet, filename string, src []byte) (*ast.File, error) {
	f, err := artifacts.Load(ld.tempFileSet, filename, src)
	if err != nil {
		return nil, ld.errGroup.Fatal(err)
	}

	pkgName, pkgPath := f.PackageName(), f.PackagePath()
	if ld.packageChanged(pkgName, pkgPath) {
		ld.packageDone()
		ld.packageStart(pkgName, pkgName)
	}

	if err := ld.group.Modify(f, ld.errGroup); err != nil {
		return nil, err
	}

	//f.Remap(ld.finalFileSet) // TODO: Fix

	final, err := f.Reload(ld.finalFileSet)
	if err != nil {
		return nil, ld.errGroup.Fatal(err)
	}
	return final.File, nil
}

func (ld *loader) packageChanged(pkgName, pkgPath string) bool {
	return ld.curPkg == nil || (ld.curPkg.Name != pkgName && ld.curPkg.Path != pkgPath)
}

func (ld *loader) packageDone() {
	if ld.curPkg != nil {
		ld.group.PackageDone(ld.curPkg, ld.errGroup)
		ld.curPkg = nil
	}
}

func (ld *loader) packageStart(pkgName, pkgPath string) {
	ld.curPkg = &artifacts.Package{
		Name:        pkgName,
		Path:        pkgPath,
		TempFileSet: ld.tempFileSet.FileSet(),
	}
	ld.group.PackageStart(ld.curPkg, ld.errGroup)
}
