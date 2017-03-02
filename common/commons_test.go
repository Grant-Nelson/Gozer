package common

import (
	"fmt"
	"strings"
	"testing"
)

// TestDiffStringSets test of the DiffStringSets method.
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
		t.Fatal(fmt.Sprint("Unexpected result from DiffStringSets:",
			"\n   Set 1:            [", strings.Join(set1, ", "), "]",
			"\n   Set 2:            [", strings.Join(set2, ", "), "]",
			"\n   Diff:             ", diff,
			"\n   Exp Diff:         ", expDiff,
			"\n   Not in Set 1:     [", strings.Join(notInSet1, ", "), "]",
			"\n   Exp Not in Set 1: [", strings.Join(expNotInSet1, ", "), "]",
			"\n   Not in Set 2:     [", strings.Join(notInSet2, ", "), "]",
			"\n   Exp Not in Set 2: [", strings.Join(expNotInSet2, ", "), "]"))
	}
}
