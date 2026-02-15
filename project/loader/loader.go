package loader

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

// Config is the configuration for the loader.
type Config struct {

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger logger.Logger

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string

	// Build is a list of command-line flags to be passed through to
	// the build system's query tool.
	Build []string

	// Patterns are the file patterns to use for the project root.
	Patterns []string

	// Tests indicates the test files will also be read for the patterns.
	Tests bool

	// Overlay is a mapping from absolute file paths to file contents.
	Overlay map[string][]byte

	// Modifiers to process each file with.
	Modifiers mods.Group

	// Parser is the file parser to use.
	// If nil, then the default parser is used.
	Parser parser.Parser

	// Parallel indicates that the packages should be loaded
	// in parallel when possible based on dependencies.
	// Otherwise, the packages are loaded one at a time in the order
	// that the packages are defined in the project.
	Parallel bool
}

// Load reads, parses, modifies, and collects type information for a project
// based on the given configuration.
func Load(cfg Config) (proj *project.Project, err error) {
	errGroup := faults.NewGroup(-1)
	//defer faults.Recover(&err) // TODO: Make this use the errGroup
	defer cfg.Logger.LogGroup("Loading")()

	// TODO: Add verbose logs (which packages loaded from cache, etc)

	p := cfg.Parser
	if p == nil {
		p = parser.Default
	}

	ld := &loader{
		group:    cfg.Modifiers,
		fSet:     token.NewFileSet(),
		errGroup: errGroup,
		parser:   p,
		overlay:  cfg.Overlay,
	}
	if err := ld.loadFileNames(cfg); err != nil {
		return nil, errGroup.Fatal(err)
	}

	// TODO: Need to check packages for errors.

	if cfg.Parallel {
		err = ld.parallelParseProject()
		return ld.proj, errGroup.Fatal(err)
	}
	err = ld.parseProject()
	return ld.proj, errGroup.Fatal(err)
}

type loader struct {
	group    mods.Group
	fSet     *token.FileSet
	errGroup *faults.Group
	proj     *project.Project
	parser   parser.Parser
	overlay  map[string][]byte
}

func (ld *loader) loadFileNames(cfg Config) error {
	const allNeeds = packages.NeedName |
		packages.NeedFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedEmbedFiles |
		packages.NeedForTest

	loadCfg := &packages.Config{
		Mode: allNeeds,
		Dir:  cfg.Dir,
		//BuildFlags: cfg.Build, // TODO: Fix the build constraints not working
		Tests:   cfg.Tests,
		Fset:    ld.fSet,
		Overlay: cfg.Overlay,
	}

	roots, err := packages.Load(loadCfg, cfg.Patterns...)
	if err != nil {
		return fmt.Errorf(`Listing files failed: %w`, err)
	}
	ld.proj = project.New(ld.fSet, roots)
	return nil
}

// parallelParseProject loads all the packages as parallel as possible.
// If there are any errors between groups of parallel packages, then
// the errors will be returned since after that group is a group that
// depends on the group with errors.
func (ld *loader) parallelParseProject() error {
	prev, depth := 0, 0
	for i, pkg := range ld.proj.AllPackages {
		if pkg.Depth != depth {
			ld.parallelParseGroup(ld.proj.AllPackages[prev:i])
			if err := ld.errGroup.Wrap(); err != nil {
				return err
			}
			prev, depth = i, pkg.Depth
		}
	}
	ld.parallelParseGroup(ld.proj.AllPackages[prev:])
	return ld.errGroup.Wrap()
}

// parallelParseGroup is a group of packages that all have their imports already
// loaded and do not depend on each other, so they can be loaded in parallel.
func (ld *loader) parallelParseGroup(pkgs []*project.Package) {
	wg := &sync.WaitGroup{}
	wg.Add(len(pkgs))
	for _, pkg := range pkgs {
		go func(pkg *project.Package) {
			// errors should already be added to the error group.
			// Adding the error again here may cause it to be added a second
			// time since another package loading could have added an error
			// between this error and the prior error so it wouldn't be skipped.
			_ = ld.parsePackage(pkg)
			wg.Done()
		}(pkg)
	}
	wg.Wait()
}

func (ld *loader) parseProject() error {
	for _, pkg := range ld.proj.AllPackages {
		if err := ld.parsePackage(pkg); err != nil {
			return err
		}
	}
	return nil
}

func (ld *loader) parsePackage(pkg *project.Package) error {
	pkg.State = buildState.Loading

	con, mg, err := ld.group.StartPackage(pkg, ld.errGroup)
	if err != nil || !con {
		return ld.errGroup.Add(err)
	}

	for _, filename := range pkg.Ast.GoFiles {
		f, err := ld.parseFile(mg, filename)
		if err != nil {
			fmt.Printf("ERROR: %v", err) // TODO: REMOVE
			if err2 := ld.errGroup.Add(err); err2 != nil {
				return err2
			}
			continue
		}
		pkg.Ast.Syntax = append(pkg.Ast.Syntax, f)
	}

	if _, err = mg.PackageDone(); err != nil {
		return ld.errGroup.Add(err)
	}

	pkg.State = buildState.Loaded
	return nil
}

func (ld *loader) parseFile(mg mods.Modifier, filename string) (*ast.File, error) {
	var src any = nil
	if over, ok := ld.overlay[filename]; ok {
		src = over
	}

	f, err := ld.parser(ld.fSet, filename, src)
	if err != nil {
		return nil, ld.errGroup.Fatal(err)
	}

	if m, ok := mg.(mods.ModifyAstFileExt); ok {
		if _, err := m.ModifyAstFile(f); err != nil {
			return nil, err
		}
	}
	return f, nil
}
