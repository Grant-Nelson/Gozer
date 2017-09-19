package goReader

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs/types"
	"github.com/grant-nelson/Gozer/msg"
	"github.com/grant-nelson/Gozer/tests"
)

// TestGoReader is a shell for testing the GoReader.
type TestGoReader struct {
	t     *common.Tester
	gr    *GoReader
	logIO *msg.LogIO
}

// NewTestGoReader creates a new GoReader shell for testing.
func NewTestGoReader(t *testing.T) *TestGoReader {
	logIO := msg.NewLogIO(&bytes.Buffer{})
	logIO.Debug = true
	logIO.Info = true
	logIO.Warnings = true
	logIO.Errors = true
	gr := NewGoReader()
	gr.Logger().Push(logIO)
	return &TestGoReader{
		t:     common.NewTester(t),
		gr:    gr,
		logIO: logIO,
	}
}

// MainMethodBodyTest checks transpilation using code in only the main method.
func (test *TestGoReader) MainMethodBodyTest(input string, exp string) {
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`func main() {`,
		`  `+common.Indent(input, "  "),
		`}`)
	test.Transpile()
	test.CheckFunction("test", "main",
		`void main() {`,
		`  `+common.Indent(exp, "  "),
		`}`)
}

// MainMethodBodyError checks that the main method transpilation returns an error.
func (test *TestGoReader) MainMethodBodyError(input string, expErrCount int, expErr string) {
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`func main() {`,
		`  `+common.Indent(input, "  "),
		`}`)
	test.gr.Transpile()
	test.CheckErrors(expErrCount, expErr)
}

// CheckNoErrors checks that no errors have occurred.
func (test *TestGoReader) CheckNoErrors(arg string) {
	if test.gr.Logger().ErrorCount() > 0 {
		test.t.Fatal(msg.NewError("Errors while ", arg).
			Add("Log", test.logIO.Output.(*bytes.Buffer).String()))
	}
}

// CheckErrors checks that the given expected errors have been logged.
func (test *TestGoReader) CheckErrors(expCount int, expMsg ...interface{}) {
	exp := fmt.Sprint(expMsg...)
	result := strings.TrimSpace(test.logIO.Output.(*bytes.Buffer).String())
	count := test.gr.Logger().ErrorCount()
	if (count != expCount) || (exp != result) {
		test.t.Failed("Unexpected or missing errors logged", common.NewMap().
			Add("Expected", "(", expCount, ") \"", exp, "\"").
			Add("Result", "(", count, ") \"", result, "\""))
	}
}

// AddCode adds another code file to this transpiler.
func (test *TestGoReader) AddCode(path string, code ...string) {
	test.gr.AddCode(path, code...)
	test.CheckNoErrors("adding code")
}

// CheckPackages checks that the expected packages have been added to the transpiler.
func (test *TestGoReader) CheckPackages(expPackages ...string) {
	packages := make([]string, test.gr.Results().Packages.Len())
	for i, pack := range test.gr.Results().Packages.Packages() {
		packages[i] = pack.GetName()
	}
	if missing, extra, diff := common.DiffStringSets(packages, expPackages); diff {
		test.t.Fatal(msg.NewError("Unexpected or missing packages").
			Add("Packages", strings.Join(packages, ", ")).
			Add("Expected", strings.Join(expPackages, ", ")).
			Add("Missing", strings.Join(missing, ", ")).
			Add("Extra", strings.Join(extra, ", ")))
	}
}

// getPack gets the package for the given name.
func (test *TestGoReader) getPack(packName string) *types.PackageType {
	pack, exists := test.gr.Results().Packages.Find(packName)
	if !exists {
		test.t.Fatal(msg.NewError("Failed to find package ", packName, " for CheckImports."))
	}
	return pack
}

// CheckImports checks that the given imports are in the given package.
func (test *TestGoReader) CheckImports(packName string, expImports ...string) {
	pack := test.getPack(packName)
	imports := make([]string, pack.Imports.Len())
	for i, importType := range pack.Imports.Packages() {
		imports[i] = importType.GetName()
	}
	if missing, extra, diff := common.DiffStringSets(imports, expImports); diff {
		test.t.Fatal(msg.NewError("Unexpected or missing package imports").
			Add("Package", packName).
			Add("Imports", strings.Join(imports, ", ")).
			Add("Expected", strings.Join(expImports, ", ")).
			Add("Missing", strings.Join(missing, ", ")).
			Add("Extra", strings.Join(extra, ", ")))
	}
}

// CheckPackage checks that the given package is in the given package.
func (test *TestGoReader) CheckPackage(packName string, expBody ...string) {
	pack := test.getPack(packName)
	expResult := tests.Lines(expBody...)
	if result := pack.FullString(); result != expResult {
		test.t.Fatal(msg.NewError("Unexpected package body").
			Add("Package", packName).
			Add("Result", result).
			Add("Expected", expResult))
	}
}

// CheckFunctions checks that the given functions are in the given package.
func (test *TestGoReader) CheckFunctions(packName string, expFunctions ...string) {
	pack := test.getPack(packName)
	functions := make([]string, pack.Functions.Len())
	for i, funcType := range pack.Functions.Functions {
		functions[i] = funcType.GetName()
	}
	if missing, extra, diff := common.DiffStringSets(functions, expFunctions); diff {
		test.t.Fatal(msg.NewError("Unexpected or missing package functions").
			Add("Package", packName).
			Add("Functions", strings.Join(functions, ", ")).
			Add("Expected", strings.Join(expFunctions, ", ")).
			Add("Missing", strings.Join(missing, ", ")).
			Add("Extra", strings.Join(extra, ", ")))
	}
}

// CheckFunction checks the function body in the given package's function.
func (test *TestGoReader) CheckFunction(packName string, funcName string, expBody ...string) {
	pack := test.getPack(packName)
	tfunc, exists := pack.Functions.Find(funcName)
	if !exists {
		test.t.Fatal(msg.NewError("Failed to find function ", funcName, " in ", packName, "."))
	}
	indent := func(s string) string {
		s = strings.Replace(s, " ", "\u00B7", -1)
		//s = "`" + strings.Replace(s, "\n", "`,\n`", -1) + "`"
		return s
	}
	expResult := tests.Lines(expBody...)
	if result := tfunc.FullBodyString(); result != expResult {
		test.t.Fatal(msg.NewError("Unexpected function construct").
			Add("Package", packName).
			Add("Function", funcName).
			Add("Expected", indent(expResult)).
			Add("Result", indent(result)))
	}
}

// Transpile runs the traspilation and cheks for any errors.
func (test *TestGoReader) Transpile() {
	test.gr.Transpile()
	test.CheckNoErrors("transpiling")
}
