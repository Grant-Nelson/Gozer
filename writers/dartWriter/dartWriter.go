package dartWriter

import (
	"bytes"

	"github.com/grant-nelson/Gozer/constructs/types"
)

// DartWriter
type DartWriter struct {

	// program is the constructs being written out by this writer.
	program *types.ProgramType

	// testing indicates that the unit-test should be written out.
	testing bool

	// files is the list of file names to file information to write.
	files map[string]interface{}
}

func NewDartWriter(program *types.ProgramType) *DartWriter {
	return &DartWriter{
		program: program,
		testing: false,
		files:   map[string]interface{}{},
	}
}

// TestingEnabled indicates if testing is enabled.
func (dw *DartWriter) TestingEnabled() bool {
	return dw.testing
}

// EnableTesting enables or disables testing for outputting
// and adding unit-testing while reading packages.
func (dw *DartWriter) EnableTesting(test bool) {
	dw.testing = test
}

// Transpile will convert the program into dart data.
func (dw *DartWriter) Transpile() {
	dw.files = map[string]interface{}{}
	for _, pack := range dw.program.Packages.Packages() {
		dw.files[pack.GetName()] = 0
		// TODO: Fill out types
	}
}

// FileNames gets the file names for the transpiled program.
func (dw *DartWriter) FileNames() []string {
	files := make([]string, len(dw.files))
	i := 0
	for file := range dw.files {
		files[i] = file
		i++
	}
	return files
}

func (dw *DartWriter) WriteFile(fileName string, buf *bytes.Buffer) {

}
