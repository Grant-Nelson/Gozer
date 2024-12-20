package project

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"
)

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

	// Augmenters to process each file with.
	Augmenters []Augmenter
}

type Augmenter func(fSet *token.FileSet, filename string, file *ast.File) error

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
	const mode = parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fSet, filename, src, mode)
	if err != nil {
		return nil, err
	}

	for _, aug := range ld.cfg.Augmenters {
		if err := aug(fSet, filename, file); err != nil {
			return nil, err
		}
	}

	return file, nil
}
