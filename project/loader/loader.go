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
//  - to preload cached packages then shortcut modified files
//  - to store modified files in a cached package that saves when load is done
//  - to simplify constants (except concatenated strings that have separate variables)
//  - to remove defers into a `deferBlock` call
//  - to remove Goto and labels (aka flatten)
//  - to inject Jumps and labels to replace other flow-controls
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
		packages:     map[string]*artifacts.Package{},
		errGroup:     faults.NewGroup(-1),
		group:        mods.Group(cfg.Modifiers),
		tempFileSet:  token.NewFileSet(),
		finalFileSet: token.NewFileSet(),
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
	packages     map[string]*artifacts.Package
	errGroup     *faults.Group
	group        mods.Group
	tempFileSet  *token.FileSet
	finalFileSet *token.FileSet
}

func (ld *loader) parseFile(fs *token.FileSet, filename string, src []byte) (*ast.File, error) {
	fm, err := artifacts.DefaultFileParser.Parse(ld.tempFileSet, filename, src)
	if err != nil {
		return nil, ld.errGroup.Fatal(err)
	}
	f := artifacts.NewFile(ld.tempFileSet, fm)

	pkgKey := f.PackageKey()
	pkg, exists := ld.packages[pkgKey]
	if exists {
		// replace the temporary package with a shared one
		f.Package = pkg
	} else {
		// use the temporary package as a shared one
		ld.packages[pkgKey] = f.Package
	}

	if _, err := ld.group.Modify(f, ld.errGroup); err != nil {
		return nil, err
	}
	return f.File, nil
}
