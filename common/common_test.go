package common

import (
	"fmt"
	"strings"
	"testing"
)

type FatalTestCall struct {
	Result string
}

func (ftc *FatalTestCall) Fatal(msg ...interface{}) {
	ftc.Result = fmt.Sprint(msg...)
}

func TestTestTools(tt *testing.T) {
	t := NewTester(tt)
	ftc := &FatalTestCall{}
	t2 := NewTester(ftc)
	tempStack := getStack
	defer func() {
		getStack = tempStack
	}()
	getStack = func() []byte {
		return []byte(fmt.Sprint(
			"1. fake stack trace\n",
			"2. fake stack trace\n",
			"3. fake stack trace\n",
			"4. fake stack trace\n",
			"5. fake stack trace\n",
			"6. fake stack trace\n",
			"7. fake stack trace\n",
			"8. fake stack trace"))
	}

	t2.CheckStr("Now", "Never")
	t.CheckStr(ftc.Result,
		`Unexpected string:`,
		`  Expected: Never`,
		`  Gotten:   Now`,
		`  Stack:    6. fake stack trace`,
		`            7. fake stack trace`)

	t2.CheckInt(12, 34, "Blooop")
	t.CheckStr(ftc.Result,
		`Unexpected integer:`,
		`  Expected: 34`,
		`  Gotten:   12`,
		`  Message:  Blooop`,
		`  Stack:    6. fake stack trace`,
		`            7. fake stack trace`)

	t2.CheckBool(true, false, "Bleep")
	t.CheckStr(ftc.Result,
		`Unexpected boolean:`,
		`  Expected: false`,
		`  Gotten:   true`,
		`  Message:  Bleep`,
		`  Stack:    6. fake stack trace`,
		`            7. fake stack trace`)

	t.CheckStr(StackTrace(1, -1), ``)
	t.CheckStr(StackTrace(-1, 2),
		`6. fake stack trace`,
		`7. fake stack trace`)
	t.CheckStr(StackTrace(2, 4), ``)

	t2.Failed("Panda", nil)
	t.CheckStr(ftc.Result,
		`Panda:`,
		`  Stack: 6. fake stack trace`,
		`         7. fake stack trace`)
}

func TestDiffStringSets(tt *testing.T) {
	t := NewTester(tt)
	t.CheckDiffStringSets("a|b|c", "a|b|c", "", "", false)
	t.CheckDiffStringSets("a|b|c", "b|c|a", "", "", false)
	t.CheckDiffStringSets("b|c|a", "a|b|c", "", "", false)
	t.CheckDiffStringSets("a|b|c", "b|c", "", "a", true)
	t.CheckDiffStringSets("b|c", "a|b|c", "a", "", true)
	t.CheckDiffStringSets("a|b|c", "a|c", "", "b", true)
	t.CheckDiffStringSets("a|c", "a|b|c", "b", "", true)
	t.CheckDiffStringSets("a|b|c", "a|b", "", "c", true)
	t.CheckDiffStringSets("a|b", "a|b|c", "c", "", true)
	t.CheckDiffStringSets("a|b|c", "", "", "a|b|c", true)
	t.CheckDiffStringSets("", "a|b|c", "a|b|c", "", true)
	t.CheckDiffStringSets("a|b", "c|d", "c|d", "a|b", true)
	t.CheckDiffStringSets("a|b|c", "c|d|e", "d|e", "a|b", true)
	t.CheckDiffStringSets("a|c|d|e", "a|b|c|e", "b", "d", true)
}

func TestIndent(tt *testing.T) {
	t := NewTester(tt)
	result := Indent("No Indent\nIndent\nAlso Indented", "  ")
	exp := "No Indent\n  Indent\n  Also Indented"
	t.CheckStr(result, exp)
}

func TestMapFormatting(tt *testing.T) {
	t := NewTester(tt)
	t.CheckMap(NewMap(),
		"")
	t.CheckMap(NewMap().
		Add("A", 1).
		Add("B", 2),
		"A: 1",
		"B: 2")
	t.CheckMap(NewMap().
		Add("Horse", "Neh").
		Add("Bird", "Tweet").
		Add("Dog", "Woof").
		Add("Cat", "Meow"),
		"Bird:  Tweet",
		"Cat:   Meow",
		"Dog:   Woof",
		"Horse: Neh")
	t.CheckMap(NewMap().
		Add("People", NewMap().
			Add("Bob", 453).
			Add("Bill", 123).
			Add("Jill", 8787)).
		Add("Dogs", NewMap().
			Add("Gizmo", 736).
			Add("Spot", 6656)),
		"Dogs:   Gizmo: 736",
		"        Spot:  6656",
		"People: Bill: 123",
		"        Bob:  453",
		"        Jill: 8787")
	t.CheckMap(NewMap().
		Add("Animals", "Cat\nDog\nHorse\nCow").
		Add("Messages", "Hello\nWorld\nGoodBye\nMoon"),
		"Animals:  Cat",
		"          Dog",
		"          Horse",
		"          Cow",
		"Messages: Hello",
		"          World",
		"          GoodBye",
		"          Moon")
	t.CheckBool(NewMap().
		Add("Bird", 1248).
		Contains("Fly"), false)
	t.CheckBool(NewMap().
		Add("Bird", 1248).
		Contains("Bird"), true)
}

//============================================================================

// splitLists converts a bar seperated list onto a slice of strings.
func splitLists(listStr string) []string {
	if len(listStr) <= 0 {
		return []string{}
	}
	return strings.Split(listStr, "|")
}

// stringSlicesEqual checks if the two slices of strings are equal.
func stringSlicesEqual(set1 []string, set2 []string) bool {
	for i, str := range set1 {
		if set2[i] != str {
			return false
		}
	}
	return true
}

// CheckDiffStringSets checks the DiffStringSets inputs and expected results.
func (t *Tester) CheckDiffStringSets(set1Str string, set2Str string, expNotInSet1Str string, expNotInSet2Str string, expDiff bool) {
	set1, set2, expNotInSet1, expNotInSet2 := splitLists(set1Str), splitLists(set2Str), splitLists(expNotInSet1Str), splitLists(expNotInSet2Str)
	notInSet1, notInSet2, diff := DiffStringSets(set1, set2)
	failed := (diff != expDiff) || (len(notInSet1) != len(expNotInSet1)) || (len(notInSet2) != len(expNotInSet2))
	if !failed {
		failed = !stringSlicesEqual(notInSet1, expNotInSet1)
		if !failed {
			failed = !stringSlicesEqual(notInSet2, expNotInSet2)
		}
	}
	if failed {
		t.Failed("Unexpected result from DiffStringSets", NewMap().
			Add("Set 1", "[", strings.Join(set1, ", "), "]").
			Add("Set 2", "[", strings.Join(set2, ", "), "]").
			Add("Diff", diff).
			Add("Exp Diff", expDiff).
			Add("Not in Set 1", "[", strings.Join(notInSet1, ", "), "]").
			Add("Exp Not in Set 1", "[", strings.Join(expNotInSet1, ", "), "]").
			Add("Not in Set 2", "[", strings.Join(notInSet2, ", "), "]").
			Add("Exp Not in Set 2", "[", strings.Join(expNotInSet2, ", "), "]"))
	}
}

// CheckMap checks that a map matches the expected result.
func (t *Tester) CheckMap(m Map, expLines ...string) {
	result := m.String()
	exp := strings.Join(expLines, "\n")
	if result != exp {
		t.Failed("Unexpected string result from Map", NewMap().
			Add("Expected", exp).
			Add("Result", result))
	}
	if (len(exp) == 0) != m.Empty() {
		t.Failed("Unexpected empty result from Map", NewMap().
			Add("Expected", len(exp) == 0).
			Add("Result", m.Empty()).
			Add("String", result))
	}
}
