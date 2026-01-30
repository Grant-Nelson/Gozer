package builder

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fetcher"
)

func New(verbose, minify bool, next fetcher.Fetcher) fetcher.Fetcher {
	return &builderImp{
		verbose: verbose,
		minify:  minify,
		next:    next,
	}
}

type builderImp struct {
	verbose bool
	minify  bool
	next    fetcher.Fetcher
}

func (b *builderImp) Fetch(path string) ([]byte, error) {
	content, err := b.next.Fetch(path)
	if err == nil {
		return content, nil
	}

	// TODO: Handle index.html

	if noExt, has := strings.CutSuffix(path, `.js`); has {
		noExt = filepath.ToSlash(noExt)
		if ok, content, err := b.buildJS(noExt+`.ts`, api.LoaderTS); ok {
			return content, err
		}
		if ok, content, err := b.buildJS(noExt+`.tsx`, api.LoaderTSX); ok {
			return content, err
		}
		if ok, content, err := b.buildJS(noExt+`.jsx`, api.LoaderJSX); ok {
			return content, err
		}
	}
	return content, err
}

func (b *builderImp) buildJS(path string, loader api.Loader) (bool, []byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, nil, err
	}

	options := api.TransformOptions{
		Target:      api.ES2015,
		TreeShaking: api.TreeShakingFalse,
		Sourcemap:   api.SourceMapInline,
		Loader:      loader,
	}
	if b.minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.KeepNames = true
	}
	result := api.Transform(string(content), options)
	if b.verbose {
		for _, w := range result.Warnings {
			log.Printf("Warning: %d:%d: %s\n%s\n", w.Location.Line, w.Location.Column, w.Text, w.Location.LineText)
		}
	}
	if errCount := len(result.Errors); errCount > 0 {
		for _, e := range result.Errors {
			log.Printf("Error: %d:%d: %s\n%s\n", e.Location.Line, e.Location.Column, e.Text, e.Location.LineText)
		}
		log.Fatalf(`JS transform failed with %d errors.`, errCount)
	}
	return true, result.Code, nil
}
