package loader

import (
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

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup

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

	// SkipFileParsing indicates that the loader should get the file names,
	// run the package start and stop modifiers but not modify files.
	SkipFileParsing bool
}

// Load reads, parses, modifies, and collects type information for a project
// based on the given configuration.
func Load(cfg Config) (proj *project.Project, err error) {
	defer cfg.ErrGroup.Recover(&err)
	defer cfg.Logger.LogGroup(`Loading`)()

	// TODO: Add verbose logs (which packages loaded from cache, etc)

	p := cfg.Parser
	if p == nil {
		p = parser.Default
	}

	ld := &loader{
		logger:          cfg.Logger,
		errGroup:        cfg.ErrGroup,
		group:           cfg.Modifiers,
		fSet:            token.NewFileSet(),
		parser:          p,
		overlay:         cfg.Overlay,
		skipFileParsing: cfg.SkipFileParsing,
	}
	if err := ld.loadFileNames(cfg); err != nil {
		return nil, err
	}

	if cfg.Parallel {
		return ld.proj, ld.parallelParseProject()
	}
	return ld.proj, ld.parseProject()
}

type loader struct {
	logger          logger.Logger
	errGroup        *faults.ErrGroup
	group           mods.Group
	fSet            *token.FileSet
	proj            *project.Project
	parser          parser.Parser
	overlay         map[string][]byte
	skipFileParsing bool
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
		ld.errGroup.Add(err)
		// Always return the errGroup to stop building.
		return ld.errGroup
	}

	ld.proj = project.New(ld.fSet, roots)
	ld.proj.CollectErrors(ld.errGroup)
	return ld.errGroup.ErrorOrNil()
}

// parallelParseProject loads all the packages as parallel as possible.
// If there are any errors between groups of parallel packages, then
// the errors will be returned since after that group is a group that
// depends on the group with errors.
func (ld *loader) parallelParseProject() error {
	prev, depth := 0, 0
	for i, pkg := range ld.proj.AllPackages {
		if pkg.Depth != depth {
			err := ld.parallelParseGroup(ld.proj.AllPackages[prev:i])
			if err != nil {
				return err
			}
			prev, depth = i, pkg.Depth
		}
	}
	return ld.parallelParseGroup(ld.proj.AllPackages[prev:])
}

// parallelParseGroup is a group of packages that all have their imports already
// loaded and do not depend on each other, so they can be loaded in parallel.
func (ld *loader) parallelParseGroup(pkgs []*project.Package) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(pkgs))
	for _, pkg := range pkgs {
		go func(pkg *project.Package) {
			_ = ld.parsePackage(pkg)
			wg.Done()
		}(pkg)
	}
	wg.Wait()
	return ld.errGroup.ErrorOrNil()
}

func (ld *loader) parseProject() error {
	for _, pkg := range ld.proj.AllPackages {
		if err := ld.parsePackage(pkg); err != nil {
			return err
		}
	}
	return ld.errGroup.ErrorOrNil()
}

func (ld *loader) parsePackage(pkg *project.Package) error {
	pkg.State = buildState.Loading

	con, mg, err := ld.group.StartPackage(pkg)
	if err != nil || !con {
		return err
	}

	for _, filename := range pkg.Ast.GoFiles {
		f, err := ld.parseFile(mg, filename)
		if err != nil {
			return err
		}
		if f != nil {
			pkg.Ast.Syntax = append(pkg.Ast.Syntax, f)
		}
	}

	if _, err = mg.PackageDone(); err != nil {
		return err
	}

	pkg.State = buildState.Loaded
	return nil
}

func (ld *loader) parseFile(mg mods.Modifier, filename string) (*ast.File, error) {
	if ld.skipFileParsing {
		return nil, nil
	}

	var src any = nil
	if over, ok := ld.overlay[filename]; ok {
		src = over
	}

	f, err := ld.parser(ld.fSet, filename, src)
	if err != nil {
		return nil, ld.errGroup.Add(err)
	}

	if m, ok := mg.(mods.ModifyAstFileExt); ok {
		if _, err := m.ModifyAstFile(f); err != nil {
			return nil, err
		}
	}
	return f, nil
}
