package transpiler

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs"
)

// TestGoReader is a shell for testing the GoReader.
type TestGoReader struct {
	t      *testing.T
	gr     *GoReader
	logBuf *bytes.Buffer
}

// NewTestGoReader creates a new GoReader shell for testing.
func NewTestGoReader(t *testing.T) *TestGoReader {
	logBuf := &bytes.Buffer{}
	gr := NewGoReader()
	gr.Logger().SetOutput(logBuf)
	return &TestGoReader{
		t:      t,
		gr:     gr,
		logBuf: logBuf,
	}
}

// fail writes a failure to the test session.
func (test *TestGoReader) fail(msg ...interface{}) {
	test.t.Log(fmt.Sprint(msg...))
	test.t.Fail()
}

// CheckNoErrors checks that no errors have occurred.
func (test *TestGoReader) CheckNoErrors(msg string) {
	if test.gr.Logger().ErrorCount() > 0 {
		test.fail("Errors while ", msg, ":\n", test.logBuf.String())
	}
}

// CheckErrors checks that the given expected errors have been logged.
func (test *TestGoReader) CheckErrors(expCount int, expMsg ...interface{}) {
	exp := fmt.Sprint(expMsg...)
	msg := strings.TrimSpace(test.logBuf.String())
	count := test.gr.Logger().ErrorCount()
	if (count != expCount) || (exp != msg) {
		test.fail("Unexpected or missing errors logged:",
			"\n   Expected: (", expCount, ") \"", exp, "\"",
			"\n   Result:   (", count, ") \"", msg, "\"")
	}
}

// AddCode adds another code file to this transpiler.
func (test *TestGoReader) AddCode(path string, code ...string) {
	test.gr.AddCode(path, code...)
	test.CheckNoErrors("adding code")
}

// CheckPackages checks that the expected packages have been added to the transpiler.
func (test *TestGoReader) CheckPackages(expPackages ...string) {
	packages := make([]string, len(test.gr.Program.Packages))
	index := 0
	for name := range test.gr.Program.Packages {
		packages[index] = name
		index++
	}
	if missing, extra, diff := common.DiffStringSets(packages, expPackages); diff {
		test.fail("Unexpected or missing packages:",
			"\n   Packages: ", strings.Join(packages, ", "),
			"\n   Expected: ", strings.Join(expPackages, ", "),
			"\n   Missing:  ", strings.Join(missing, ", "),
			"\n   Extra:    ", strings.Join(extra, ", "))
	}
}

// getPack gets the package for the given name.
func (test *TestGoReader) getPack(packName string) *constructs.PackageType {
	pack, exists := test.gr.Program.Packages[packName]
	if !exists {
		test.fail("Failed to find package ", packName, " for CheckImports.")
	}
	return pack
}

// CheckImports checks that the given imports are in the given package.
func (test *TestGoReader) CheckImports(packName string, expImports ...string) {
	pack := test.getPack(packName)
	imports := make([]string, len(pack.Imports))
	index := 0
	for name := range pack.Imports {
		imports[index] = name
		index++
	}
	if missing, extra, diff := common.DiffStringSets(imports, expImports); diff {
		test.fail("Unexpected or missing package imports:",
			"\n   Package:  ", packName,
			"\n   Imports:  ", strings.Join(imports, ", "),
			"\n   Expected: ", strings.Join(expImports, ", "),
			"\n   Missing:  ", strings.Join(missing, ", "),
			"\n   Extra:    ", strings.Join(extra, ", "))
	}
}

// CheckFunctions checks that the given functions are in the given package.
func (test *TestGoReader) CheckFunctions(packName string, expFunctions ...string) {
	pack := test.getPack(packName)
	functions := make([]string, len(pack.Functions))
	index := 0
	for name := range pack.Functions {
		functions[index] = name
		index++
	}
	if missing, extra, diff := common.DiffStringSets(functions, expFunctions); diff {
		test.fail("Unexpected or missing package functions:",
			"\n   Package:   ", packName,
			"\n   Functions: ", strings.Join(functions, ", "),
			"\n   Expected:  ", strings.Join(expFunctions, ", "),
			"\n   Missing:   ", strings.Join(missing, ", "),
			"\n   Extra:     ", strings.Join(extra, ", "))
	}
}

// CheckFunctionBody checks the function body in the given package's function.
func (test *TestGoReader) CheckFunctionBody(packName string, funcName string, expBody ...string) {
	pack := test.getPack(packName)
	tfunc, exists := pack.Functions[funcName]
	if !exists {
		test.fail("Failed to find function ", funcName, " in ", packName, ".")
	}
	expResult := strings.Join(expBody, "\n")
	if result := tfunc.String(); result != expResult {
		test.fail("Unexpected function construct:",
			"\n   Package:  ", packName,
			"\n   Function: ", funcName,
			"\n   Expected: ", strings.Replace(expResult, "\n", "\n             ", -1),
			"\n   Result:   ", strings.Replace(result, "\n", "\n             ", -1))
	}
}

// Transpile runs the traspilation and cheks for any errors.
func (test *TestGoReader) Transpile() {
	test.gr.Transpile()
	test.CheckNoErrors("transpiling")
}
