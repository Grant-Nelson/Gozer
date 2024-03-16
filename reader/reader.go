package reader

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

// Read reads a project and all its packages and files.
func Read(config *Config) (proj *Project, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	proj = &Project{}
	cfg := getParseConfigs(config)
	proj.Packages, err = packages.Load(cfg, config.Path)
	if err != nil {
		return proj, err
	}
	return proj, proj.Errors()
}

const allNeeds = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedExportFile |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes |
	packages.NeedModule |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns

func getParseConfigs(config *Config) *packages.Config {
	cfg := &packages.Config{
		Dir:        config.Path,
		BuildFlags: config.BuildFlags,
		Context:    config.Context,
		Mode:       allNeeds,
		Tests:      config.Tests,
	}

	if config.AugmentFile != nil {
		fa := fileAugmenter{AugmentFile: config.AugmentFile}
		cfg.ParseFile = fa.parseFile
	}

	if config.Verbose {
		cfg.Logf = func(format string, args ...any) {
			_, err := fmt.Printf(format, args...)
			panic(err)
		}
	}

	return cfg
}

type fileAugmenter struct {
	AugmentFile func(args *AugmentFileArgs) error
}

func (fa *fileAugmenter) parseFile(fileSet *token.FileSet, filename string, src []byte) (*ast.File, error) {
	const mode = parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fileSet, filename, src, mode)
	if err != nil {
		return nil, err
	}

	err = fa.AugmentFile(&AugmentFileArgs{
		Filename: filename,
		FileSet:  fileSet,
		File:     file,
	})
	return file, err
}
