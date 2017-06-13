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

func TestTestTools(t *testing.T) {
	ftc := &FatalTestCall{}
	CheckString(ftc, "Now", "Never")
	CheckString(t, ftc.Result,
		`Unexpected construct string:`,
		`   Expected: Never`,
		`   Gotten:   Now`)
}

func TestDiffStringSets(t *testing.T) {
	checkDiffStringSets(t, "a|b|c", "a|b|c", "", "", false)
	checkDiffStringSets(t, "a|b|c", "b|c|a", "", "", false)
	checkDiffStringSets(t, "b|c|a", "a|b|c", "", "", false)
	checkDiffStringSets(t, "a|b|c", "b|c", "", "a", true)
	checkDiffStringSets(t, "b|c", "a|b|c", "a", "", true)
	checkDiffStringSets(t, "a|b|c", "a|c", "", "b", true)
	checkDiffStringSets(t, "a|c", "a|b|c", "b", "", true)
	checkDiffStringSets(t, "a|b|c", "a|b", "", "c", true)
	checkDiffStringSets(t, "a|b", "a|b|c", "c", "", true)
	checkDiffStringSets(t, "a|b|c", "", "", "a|b|c", true)
	checkDiffStringSets(t, "", "a|b|c", "a|b|c", "", true)
	checkDiffStringSets(t, "a|b", "c|d", "c|d", "a|b", true)
	checkDiffStringSets(t, "a|b|c", "c|d|e", "d|e", "a|b", true)
	checkDiffStringSets(t, "a|c|d|e", "a|b|c|e", "b", "d", true)
}

func TestIndent(t *testing.T) {
	result := Indent("No Indent\nIndent\nAlso Indented", "   ")
	exp := "No Indent\n   Indent\n   Also Indented"
	if result != exp {
		Failed(t, "Unexpected result from Indent", NewMap().
			Add("Result", result).
			Add("Expected", exp))
	}
}

func TestMapFormatting(t *testing.T) {
	checkMap(t, NewMap(),
		"")
	checkMap(t, NewMap().
		Add("A", 1).
		Add("B", 2),
		"A: 1",
		"B: 2")
	checkMap(t, NewMap().
		Add("Horse", "Neh").
		Add("Bird", "Tweet").
		Add("Dog", "Woof").
		Add("Cat", "Meow"),
		"Bird:  Tweet",
		"Cat:   Meow",
		"Dog:   Woof",
		"Horse: Neh")
	checkMap(t, NewMap().
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
	checkMap(t, NewMap().
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
}

//============================================================================

// splitLists converts a bar seperated list onto a slice of strings.
func splitLists(listStr string) []string {
	if len(listStr) <= 0 {
		return []string{}
	}
	return strings.Split(listStr, "|")
}

// StringSlicesEqual checks if the two slices of strings are equal.
func stringSlicesEqual(set1 []string, set2 []string) bool {
	for i, str := range set1 {
		if set2[i] != str {
			return false
		}
	}
	return true
}

// checkDiffStringSets checks the DiffStringSets inputs and expected results.
func checkDiffStringSets(t *testing.T, set1Str string, set2Str string, expNotInSet1Str string, expNotInSet2Str string, expDiff bool) {
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
		Failed(t, "Unexpected result from DiffStringSets", NewMap().
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

// checkMap checks that a map matches the expected result.
func checkMap(t *testing.T, m Map, expLines ...string) {
	result := m.String()
	exp := strings.Join(expLines, "\n")
	if result != exp {
		Failed(t, "Unexpected string result from Map", NewMap().
			Add("Expected", exp).
			Add("Result", result))
	}
	if (len(exp) == 0) != m.Empty() {
		Failed(t, "Unexpected empty result from Map", NewMap().
			Add("Expected", len(exp) == 0).
			Add("Result", m.Empty()).
			Add("String", result))
	}
}
