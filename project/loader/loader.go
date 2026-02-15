package loader

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

// Config is the configuration for the loader.
type Config struct {

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
func Load(cfg Config) (*project.Project, error) {
	p := cfg.Parser
	if p == nil {
		p = parser.Default
	}

	ld := &loader{
		group:    cfg.Modifiers,
		fSet:     token.NewFileSet(),
		errGroup: faults.NewGroup(-1),
		parser:   p,
		overlay:  cfg.Overlay,
	}
	if err := ld.loadFileNames(cfg); err != nil {
		return nil, err
	}

	// TODO: Need to check packages for errors.

	if cfg.Parallel {
		err := ld.parallelParseProject()
		return ld.proj, err
	}

	err := ld.parseProject()
	return ld.proj, err
}

type loader struct {
	group    mods.Group
	fSet     *token.FileSet
	errGroup *faults.Group
	proj     *project.Project
	parser   parser.Parser
	overlay  map[string][]byte
}

type blockChan chan struct{}

func (ld *loader) loadFileNames(cfg Config) error {
	const allNeeds = packages.NeedName |
		packages.NeedFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedEmbedFiles

	c := &packages.Config{
		Mode:       allNeeds,
		Dir:        cfg.Dir,
		BuildFlags: cfg.Build,
		Tests:      cfg.Tests,
		Fset:       ld.fSet,
		Overlay:    cfg.Overlay,
	}
	roots, err := packages.Load(c, cfg.Patterns...)
	if err != nil {
		return err
	}
	ld.proj = project.New(ld.fSet, roots)
	return nil
}

func (ld *loader) parallelParseProject() error {
	cancel := make(blockChan)
	pending := make(map[string]blockChan, len(ld.proj.AllPackages))
	allDeps := make([]blockChan, len(ld.proj.AllPackages))
	for i, pkg := range ld.proj.AllPackages {
		pkgChan := make(blockChan)
		pending[pkg.PkgPath()] = pkgChan
		allDeps[i] = pkgChan

		deps := make([]blockChan, len(pkg.Ast.Imports))
		for pkgPath := range pkg.Ast.Imports {
			deps = append(deps, pending[pkgPath])
		}

		go ld.asyncParsePackage(pkg, pkgChan, cancel, deps)
	}

	waitOnDeps(cancel, deps)

	return nil
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
	con, mg, err := ld.group.StartPackage(pkg, ld.errGroup)
	if err != nil || !con {
		return err
	}

	for _, filename := range pkg.Ast.GoFiles {
		f, err := ld.parseFile(mg, filename)
		if err != nil {
			return err
		}
		pkg.Ast.Syntax = append(pkg.Ast.Syntax, f)
	}

	_, err = mg.PackageDone()
	return err
}

func (ld *loader) parseFile(mg mods.Modifier, filename string) (*ast.File, error) {
	var src []byte
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
