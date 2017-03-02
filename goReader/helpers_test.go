package transpiler

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
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

// CheckImports checks that the imports in the given package.
func (test *TestGoReader) CheckImports(packName string, expImports ...string) {
	pack, exists := test.gr.Program.Packages[packName]
	if !exists {
		test.fail("Failed to find package ", packName, ".")
		return
	}
	imports := make([]string, len(pack.Imports))
	index := 0
	for name := range pack.Imports {
		imports[index] = name
		index++
	}
	if missing, extra, diff := common.DiffStringSets(imports, expImports); diff {
		test.fail("Unexpected or missing packages:",
			"\n   Package:  ", packName,
			"\n   Imports:  ", strings.Join(imports, ", "),
			"\n   Expected: ", strings.Join(expImports, ", "),
			"\n   Missing:  ", strings.Join(missing, ", "),
			"\n   Extra:    ", strings.Join(extra, ", "))
	}
}

// Transpile runs the traspilation and cheks for any errors.
func (test *TestGoReader) Transpile() {
	test.gr.Transpile()
	test.CheckNoErrors("transpiling")
}
