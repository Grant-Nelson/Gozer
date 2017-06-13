package common

import "strings"

// TestInterface is an interface used for testing.
type TestInterface interface {

	// Fatal will fail a test and print the given msg.
	Fatal(msg ...interface{})
}

// Failed indicates a test has failed.
func Failed(t TestInterface, text string, m Map) {
	result := ""
	if !m.Empty() {
		result = ":\n   " + m.FormatMap("   ")
	}
	t.Fatal(text + result)
}

// CheckString checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func CheckString(t TestInterface, result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		Failed(t, "Unexpected construct string", NewMap().
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}
