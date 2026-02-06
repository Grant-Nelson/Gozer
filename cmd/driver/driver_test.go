package main

import (
	"fmt"
	"go/token"
	"log"
	"os"
	"testing"

	"golang.org/x/tools/go/packages"
)

func Test_Driver_Basic(t *testing.T) {
	const allNeeds = packages.NeedName |
		packages.NeedFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedExportFile
	driverEnv := `GOPACKAGESDRIVER=cmd/driver/driver.exe`
	logEnv := `GOPACKAGESDEBUG=true`

	fileSet := token.NewFileSet()
	c := &packages.Config{
		Mode:       allNeeds,
		Dir:        `../../`,
		BuildFlags: []string{},
		Fset:       fileSet,
		Logf:       log.Printf,
		Env:        append(os.Environ(), driverEnv, logEnv),
	}

	pkgs, err := packages.Load(c, `math`)
	if err != nil {
		t.Fatalf(`Failed to load: %v`, err)
	}

	for pkg := range packages.Postorder(pkgs) {
		fmt.Printf("Package: %q\n", pkg.PkgPath)
		printOneFile(`ID`, pkg.ID)
		printOneFile(`Name`, pkg.Name)
		printOneFile(`Dir`, pkg.Dir)
		printFiles(`Go Files`, pkg.GoFiles)
		printFiles(`Compiled Go Files`, pkg.CompiledGoFiles)
		printFiles(`Other Files`, pkg.OtherFiles)
		printFiles(`Embed Files`, pkg.EmbedFiles)
		printFiles(`Embed Patterns`, pkg.EmbedPatterns)
		printFiles(`IgnoredFiles`, pkg.IgnoredFiles)
		printOneFile(`Export Files`, pkg.ExportFile)
		printOneFile(`Target`, pkg.Target)
		printOneFile(`For Test`, pkg.ForTest)
	}

	t.Errorf(`Just an error`)
}

func printOneFile(name, file string) {
	if len(file) > 0 {
		fmt.Printf("\t%s: %q\n", name, file)
	}
}

func printFiles(name string, files []string) {
	if len(files) > 0 {
		fmt.Printf("\t%s (%d)\n", name, len(files))
		for i, f := range files {
			fmt.Printf("\t\t(%d) %q\n", i, f)
		}
	}
}
