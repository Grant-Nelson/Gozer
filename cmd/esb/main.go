package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/Grant-Nelson/Gozer/avail/args"
)

var Config = &struct {
	Help   string `arg:"help"`
	Output string `arg:"optional,flag,o|out,Output file to write to instead of the input with *.js extension"`
	Input  string `arg:"required,pos,Input,The input *.ts file to transform into *.js."`
	Minify bool   `arg:"optional,flag,m|mini,Minifies any *.js being built"`
}{
	Help:   `This compiles a *.ts file into a *.js file.`,
	Output: ``,
	Input:  ``,
	Minify: false,
}

func main() {
	if !args.Parse(Config) {
		return
	}

	if len(Config.Output) <= 0 {
		Config.Output = strings.TrimSuffix(Config.Input, filepath.Ext(Config.Input)) + `*.js`
	}

	data, err := os.ReadFile(Config.Input)
	if err != nil {
		log.Fatalf(`Failed reading source: %v`, err)
	}

	result := transformToJS(string(data))

	err = os.WriteFile(Config.Output, result, 0644)
	if err != nil {
		log.Fatalf(`Failed writing output: %v`, err)
	}
}

func transformToJS(source string) []byte {
	options := api.TransformOptions{
		Target:      api.ES2015,
		Loader:      api.LoaderTS,
		TreeShaking: api.TreeShakingFalse,
		Sourcemap:   api.SourceMapInline,
		Sourcefile:  Config.Input,
	}

	if Config.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.KeepNames = true
	}

	result := api.Transform(source, options)
	for _, w := range result.Warnings {
		log.Printf("Warning: %d:%d: %s\n%s\n", w.Location.Line, w.Location.Column, w.Text, w.Location.LineText)
	}
	if errCount := len(result.Errors); errCount > 0 {
		for _, e := range result.Errors {
			log.Printf("Error: %d:%d: %s\n%s\n", e.Location.Line, e.Location.Column, e.Text, e.Location.LineText)
		}
		log.Fatalf(`JS transform failed with %d errors`, errCount)
	}

	return result.Code
}
