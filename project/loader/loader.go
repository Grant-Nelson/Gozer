package project

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

	// Parser is the file parser to use.
	// If nil, then the default parser is used.
	Parser parser.Parser
}

// Load reads, parses, modifies, and collects type information for a project
// based on the given configuration.
func Load(cfg Config) (*project.Project, error) {
	p := cfg.Parser
	if p == nil {
		p = parser.Default
	}

	ld := &loader{
		group:    mods.Group(cfg.Modifiers),
		fSet:     token.NewFileSet(),
		errGroup: faults.NewGroup(-1),
		parser:   p,
		overlay:  cfg.Overlay,
	}
	if err := ld.loadFileNames(cfg); err != nil {
		return nil, err
	}

	// TODO: Need to check packages for errors.
	if err := ld.parseProject(); err != nil {
		return nil, err
	}
	return ld.proj, nil
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
		packages.NeedEmbedFiles

	c := &packages.Config{
		Mode:       allNeeds,
		Dir:        cfg.Dir,
		BuildFlags: cfg.BuildFlags,
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

func (ld *loader) parseProject() error {
	// TODO: Could load these asynchronously any package that
	// has had all of its dependencies finished.
	for _, pkg := range ld.proj.AllPackages {
		if err := ld.parsePackage(pkg); err != nil {
			return err
		}
	}
	return nil
}

func (ld *loader) parsePackage(pkg *project.Package) error {
	if con, err := ld.group.PackageStart(pkg, ld.errGroup); err != nil || !con {
		return err
	}

	for _, filename := range pkg.GoFiles {
		f, err := ld.parseFile(filename)
		if err != nil {
			return err
		}
		pkg.Syntax = append(pkg.Syntax, f)
	}

	_, err := ld.group.PackageDone(pkg, ld.errGroup)
	return err
}

func (ld *loader) parseFile(filename string) (*ast.File, error) {
	var src []byte
	if over, ok := ld.overlay[filename]; ok {
		src = over
	}

	f, err := ld.parser(ld.fSet, filename, src)
	if err != nil {
		return nil, ld.errGroup.Fatal(err)
	}

	if _, err := ld.group.ModifyFile(f, ld.errGroup); err != nil {
		return nil, err
	}
	return f, nil
}
