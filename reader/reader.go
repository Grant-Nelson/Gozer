package reader

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

type Config struct {
	Verbose         bool
	MainPackagePath string
	Context         context.Context
	Tests           bool
	BuildFlags      []string
	AugmentFile     func(args *AugmentFileArgs) error
}

type AugmentFileArgs struct {
	Filename string
	FileSet  *token.FileSet
	File     *ast.File
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

func Read(config *Config) (_ *Project, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	cfg := getParseConfigs(config)
	mainPackages, err := packages.Load(cfg, config.MainPackagePath)
	if err != nil {
		panic(err)
	}

	p := &Project{Packages: mainPackages}
	err = newError(p.Errors())
	return p, err
}

func getParseConfigs(config *Config) *packages.Config {
	cfg := &packages.Config{
		Dir:        config.MainPackagePath,
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
	f, err := parser.ParseFile(fileSet, filename, src, mode)
	if err != nil {
		return nil, err
	}

	err = fa.AugmentFile(&AugmentFileArgs{
		Filename: filename,
		FileSet:  fileSet,
		File:     f,
	})
	return f, err
}
