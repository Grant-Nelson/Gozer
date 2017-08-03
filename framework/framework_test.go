package framework

import (
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestBuiltin(tt *testing.T) {
	t := common.NewTester(tt)
	prog := types.Program()
	BuiltinPrebuild(prog)
	// check double loading doesn't cause a problem
	BuiltinPrebuild(prog)

	// The current values for builtin, will change as parts are added
	CheckPackage(t, prog, "builtin",
		`import builtin{`,
		`  variant append(variant a, variant... b)`,
		`  int cap(variant a)`,
		`  void delete(variant m, variant key)`,
		`  int len(variant a)`,
		`  void print(interface{}... a)`,
		`  void println(interface{}... a)`,
		`  error{`,
		`    string Error()`,
		`  }`,
		`}`)
}

func TestFmt(tt *testing.T) {
	t := common.NewTester(tt)
	prog := types.Program()
	FmtPrebuild(prog)
	// check double loading doesn't cause a problem
	FmtPrebuild(prog)

	// The current values for fmt, will change as parts are added
	CheckPackage(t, prog, "fmt",
		`import fmt{`,
		`  error Errorf(string format, interface{}... a)`,
		`  printResult Print(interface{}... a)`,
		`  printResult Printf(string format, interface{}... a)`,
		`  printResult Println(interface{}... a)`,
		`  string Sprint(interface{}... a)`,
		`  string Sprintf(string format, interface{}... a)`,
		`  string Sprintln(interface{}... a)`,
		`  printResult{`,
		`    int n`,
		`    error err`,
		`  }`,
		`}`)
}

func TestIO(tt *testing.T) {
	t := common.NewTester(tt)
	prog := types.Program()
	IOPrebuild(prog)
	// Check double loading doesn't cause a problem
	IOPrebuild(prog)

	// The current values for io, will change as parts are added
	CheckPackage(t, prog, "io",
		`import io{}`)
}

//============================================================================

// CheckPackage checks that the package's string matches the given string.
func CheckPackage(t *common.Tester, prog *types.ProgramType, packageName string, exp ...string) {
	pack, found := prog.Packages.Find(packageName)
	t.CheckBool(found, true, "expected ", packageName, " package")
	t.CheckStr(pack.FullString(), exp...)
}
