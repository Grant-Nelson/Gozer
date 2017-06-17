package common

import (
	"fmt"
	"runtime/debug"
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

// getStack gets the stack trace
var getStack func() []byte = debug.Stack

// Stack gets the current stack trace.
func (t *Tester) Stack(offset int, count int) string {
	if count <= 0 {
		return ""
	}
	if offset < 0 {
		offset = 0
	}
	stack := strings.Split(string(getStack()), "\n")
	length := len(stack)
	start := offset*2 + 5
	if start >= length {
		start = length - 1
	}
	stop := start + count*2
	if stop >= length {
		stop = length - 1
	}
	return strings.Join(stack[start:stop], "\n")
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
		m.Add("Stack", t.Stack(0, 5))
	}
	result := ""
	if !m.Empty() {
		result = ":\n   " + m.FormatMap("   ")
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
			Add("Stack", t.Stack(0, 5))
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
			Add("Stack", t.Stack(0, 5))
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
			Add("Stack", t.Stack(0, 5))
		if len(msg) > 0 {
			m.Add("Message", fmt.Sprint(msg...))
		}
		t.Failed("Unexpected boolean", m)
	}
}
