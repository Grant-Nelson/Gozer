package common

import (
	"fmt"
	"strings"
)

// TestInterface is an interface used for testing.
type TestInterface interface {

	// Fatal will fail a test and print the given msg.
	Fatal(msg ...interface{})
}

// Tester is a tool to help unit-testing.
type Tester struct {

	// t the interface to the testing instance.
	t TestInterface
}

// NewTester creates a new tester.
func NewTester(t TestInterface) *Tester {
	return &Tester{
		t: t,
	}
}

// Fatal will fail a test and print the given msg.
func (t *Tester) Fatal(msg ...interface{}) {
	t.t.Fatal(msg...)
}

// Failed indicates a test has failed.
func (t *Tester) Failed(text string, m Map) {
	if m == nil {
		m = NewMap()
	}
	if !m.Contains("Stack") {
		m.Add("Stack", StackTrace(0, 5))
	}
	result := ""
	if !m.Empty() {
		result = ":\n  " + m.FormatMap("  ")
	}
	t.Fatal(text + result)
}

// CheckStr checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func (t *Tester) CheckStr(result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		m := NewMap().
			Add("Expected", expStr).
			Add("Gotten", result).
			Add("Stack", StackTrace(0, 5))
		t.Failed("Unexpected string", m)
	}
}

// CheckInt checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func (t *Tester) CheckInt(result int, exp int, msg ...interface{}) {
	if result != exp {
		m := NewMap().
			Add("Expected", exp).
			Add("Gotten", result).
			Add("Stack", StackTrace(0, 5))
		if len(msg) > 0 {
			m.Add("Message", fmt.Sprint(msg...))
		}
		t.Failed("Unexpected integer", m)
	}
}

// CheckBool checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func (t *Tester) CheckBool(result bool, exp bool, msg ...interface{}) {
	if result != exp {
		m := NewMap().
			Add("Expected", exp).
			Add("Gotten", result).
			Add("Stack", StackTrace(0, 5))
		if len(msg) > 0 {
			m.Add("Message", fmt.Sprint(msg...))
		}
		t.Failed("Unexpected boolean", m)
	}
}
