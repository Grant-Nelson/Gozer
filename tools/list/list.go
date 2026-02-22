// This package defines the list tool run via "gozer list".
// This tool is similar to "go list" but not required to be identical.
//
// See: https://pkg.go.dev/cmd/go#hdr-List_packages_or_modules
// See: https://pkg.go.dev/cmd/go/internal/list
package list

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/targets"
	"golang.org/x/tools/go/packages"
)

type Config struct {
	Usage    string   `arg:"help"`
	Lang     string   `arg:"flag, l|lang|language, The language that would be transpiled to."`
	Patterns []string `arg:"pos, patterns, One or more patterns for the root files for a project."`
	Tests    bool     `arg:"flag, t|test, Indicate the transpile would be for a test."`
	Deps     bool     `arg:"flag, d|deps, Indicate that dependencies should be listed as well."`
	Json     bool     `arg:"flag, j|json, Outputs as json. Json is ignored if a format is given."`
	Format   string   `arg:"flag, f|fmt|format, Outputs with the given template format."`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Prints the list of files that would be used for the given build. ` +
			`This will not include augmentation files. ` +
			`This tool is similar to but is not the same as "go list".`,
		Lang:     targets.DefaultLang,
		Patterns: []string{},
		Tests:    false,
	}
}

func List(cfg *Config) {
	listCfg := &targets.ListConfig{
		Lang:     cfg.Lang,
		Patterns: cfg.Patterns,
		Tests:    cfg.Tests,
	}
	proj, err := targets.List(listCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "List failed: %v\n", err)
		return
	}

	pkgs := proj.Roots
	if cfg.Deps {
		pkgs = proj.AllPackages
	}
	out := os.Stdout

	switch {
	case len(cfg.Format) > 0:
		outputFormatted(out, cfg.Format, pkgs)
	case cfg.Json:
		outputJson(out, pkgs)
	default:
		outputDefault(out, pkgs)
	}
}

func outputFormatted(w io.Writer, format string, pkgs []*project.Package) {
	data := getAllData(pkgs)
	fm := template.FuncMap{
		"join": strings.Join,
	}
	tmpl, err := template.New("main").Funcs(fm).Parse(format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in template format: %v\n", err)
		return
	}
	for _, d := range data {
		buf := &strings.Builder{}
		if err := tmpl.Execute(buf, d); err != nil {
			fmt.Fprintf(os.Stderr, "Error templating data: %v\n", err)
			return
		}
		res := buf.String()
		if !strings.HasSuffix(res, "\n") {
			res += "\n"
		}
		fmt.Fprint(w, res)
	}
}

func outputJson(w io.Writer, pkgs []*project.Package) {
	data := getAllData(pkgs)
	en := json.NewEncoder(w)
	en.SetIndent(``, "\t")
	if err := en.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error list marshalling packages: %v\n", err)
		return
	}
}

func outputDefault(w io.Writer, pkgs []*project.Package) {
	paths := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		paths[i] = pkg.PkgPath()
	}
	slices.Sort(paths)
	for _, p := range paths {
		fmt.Fprintln(w, p)
	}
}

type packageData struct {
	Dir        string `json:",omitempty"` // directory containing package sources
	ImportPath string `json:",omitempty"` // import path of package in dir
	Name       string `json:",omitempty"` // package name
	Target     string `json:",omitempty"` // installed target for this package (may be executable)
	ForTest    string `json:",omitempty"` // package is only for use in named test
	Export     string `json:",omitempty"` // file containing export data (set by go list -export)
	DepOnly    bool   `json:",omitempty"` // package is only as a dependency, not explicitly listed

	//Goroot     bool `json:",omitempty"` // is this package found in the Go root?
	//Standard   bool `json:",omitempty"` // is this package part of the standard Go library?

	GoFiles       []string `json:",omitempty"` // .go source files
	IgnoredFiles  []string `json:",omitempty"` // source files ignored due to build constraints
	OtherFiles    []string `json:",omitempty"` // other source files
	EmbedPatterns []string `json:",omitempty"` // //go:embed patterns
	EmbedFiles    []string `json:",omitempty"` // files matched by EmbedPatterns
	Imports       []string `json:",omitempty"` // import paths used by this package
	Deps          []string `json:",omitempty"` // all (recursively) imported dependencies
}

func getData(pkg *project.Package) *packageData {
	data := &packageData{
		Dir:        pkg.Ast.Dir,
		ImportPath: pkg.Ast.PkgPath,
		Name:       pkg.Ast.Name,
		Target:     pkg.Ast.Target,
		ForTest:    pkg.Ast.ForTest,
		Export:     pkg.Ast.ExportFile,
		DepOnly:    !pkg.Root,

		GoFiles:       pkg.Ast.GoFiles,
		IgnoredFiles:  pkg.Ast.IgnoredFiles,
		OtherFiles:    pkg.Ast.OtherFiles,
		EmbedPatterns: pkg.Ast.EmbedPatterns,
		EmbedFiles:    pkg.Ast.EmbedFiles,
	}

	imports := slices.Collect(maps.Keys(pkg.Ast.Imports))
	slices.Sort(imports)
	data.Imports = imports

	depMap := map[string]bool{}
	for dep := range packages.Postorder([]*packages.Package{pkg.Ast}) {
		depMap[dep.PkgPath] = true
	}
	deps := slices.Collect(maps.Keys(depMap))
	slices.Sort(deps)
	data.Deps = deps

	return data
}

func getAllData(pkgs []*project.Package) []*packageData {
	data := make([]*packageData, len(pkgs))
	for i, pkg := range pkgs {
		data[i] = getData(pkg)
	}
	slices.SortFunc(data, func(a, b *packageData) int {
		return cmp.Compare(a.ImportPath, b.ImportPath)
	})
	return data
}
