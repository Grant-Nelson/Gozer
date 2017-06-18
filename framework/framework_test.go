package framework

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestBuiltin(t *testing.T) {
	prog := types.Program()
	BuiltinPrebuild(prog)
	// check double loading doesn't cause a problem
	BuiltinPrebuild(prog)

	// The current values for builtin, will change as parts are added
	checkType(t, prog.Packages["builtin"],
		`{`,
		`  variant append(variant a, variant... b)`,
		`  int cap(variant a)`,
		`  int len(variant a)`,
		`  void print(interface{}... a)`,
		`  void println(interface{}... a)`,
		`  error{`,
		`    string Error()`,
		`  }`,
		`}`)
}

func TestFmt(t *testing.T) {
	prog := types.Program()
	FmtPrebuild(prog)
	// check double loading doesn't cause a problem
	FmtPrebuild(prog)

	// The current values for fmt, will change as parts are added
	checkType(t, prog.Packages["fmt"],
		`{`,
		`  string Sprintln(interface{}... a)`,
		`  error Errorf(string format, interface{}... a)`,
		`  printResult Print(interface{}... a)`,
		`  printResult Printf(string format, interface{}... a)`,
		`  printResult Println(interface{}... a)`,
		`  string Sprint(interface{}... a)`,
		`  string Sprintf(string format, interface{}... a)`,
		`  printResult{`,
		`    int n`,
		`    error err`,
		`  }`,
		`}`)
}

func TestIO(t *testing.T) {
	prog := types.Program()
	IOPrebuild(prog)
	// Check double loading doesn't cause a problem
	IOPrebuild(prog)

	// The current values for io, will change as parts are added
	checkType(t, prog.Packages["io"],
		`{}`)
}

//============================================================================

// checkType checks that the type's string matches the given string.
func checkType(t *testing.T, ty types.Type, exp ...string) {
	checkString(t, types.ToString(ty), exp...)
}
