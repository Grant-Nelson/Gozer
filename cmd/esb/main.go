package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

const (
	input  = "../../experiments/exp001/scheduler.ts"
	output = "../../experiments/exp001/scheduler.js"
)

func main() {
	data, err := os.ReadFile(input)
	if err != nil {
		log.Fatalf(`Failed reading source: %v`, err)
	}

	result := transformToJS(string(data), false)

	err = os.WriteFile(output, result, 0644)
	if err != nil {
		log.Fatalf(`Failed writing output: %v`, err)
	}
}

func transformToJS(source string, minify bool) []byte {
	options := api.TransformOptions{
		Target:      api.ES2015,
		Loader:      api.LoaderTS,
		TreeShaking: api.TreeShakingFalse,
	}

	if minify {
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
