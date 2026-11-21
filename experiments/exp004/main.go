package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const hello = `
package main

import "fmt"

func main() {
	message := "Hello, World!"
	for _, r := range message {
		if r >= 'a' && r <= 'z' {
			r = r - 'a' + 'A'
		}
		fmt.Println(string(r))
	}
}
`

func main() {
	// Parse the source files.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", hello, parser.ParseComments)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}
	files := []*ast.File{f}

	// Create the type-checker's package.
	pkg := types.NewPackage("hello", "")

	// Type-check the package, load dependencies.
	// Create and build the SSA program.
	cfg := &types.Config{Importer: importer.Default()}
	hello, _, err := ssautil.BuildPackage(cfg, fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Print out the package.
	hello.WriteTo(os.Stdout)

	// Print out the package-level functions.
	//hello.Func("init").WriteTo(os.Stdout)
	hello.Func("main").WriteTo(os.Stdout)

	/*
		keys := slices.Sorted(maps.Keys(hello.Members))
		for _, key := range keys {
			fmt.Println("-------------------------------")
			member := hello.Members[key]
			switch m := member.(type) {
			case *ssa.Function:
				fmt.Print(`Function: `)
				m.WriteTo(os.Stdout)
			case *ssa.Global:
				fmt.Println(`Global:`, m)
			default:
				fmt.Printf("%s: %T\n", key, member)
			}
		}
	*/
}
