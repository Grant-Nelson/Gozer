package common

import (
	"runtime/debug"
	"sort"
	"strings"
)

// DiffStringSets determines the difference between the two given sets.
func DiffStringSets(set1 []string, set2 []string) (notInSet1 []string, notInSet2 []string, diff bool) {
	sort.Strings(set1)
	sort.Strings(set2)
	len1, len2 := len(set1), len(set2)
	if len1 <= 0 {
		return set2, []string{}, true
	}
	if len2 <= 0 {
		return []string{}, set1, true
	}
	notInSet1, notInSet2, diff = []string{}, []string{}, false
	i, j := 0, 0
	for (i < len1) && (j < len2) {
		cmp := strings.Compare(set1[i], set2[j])
		if cmp > 0 {
			diff = true
			notInSet1 = append(notInSet1, set2[j])
			j++
		} else if cmp < 0 {
			diff = true
			notInSet2 = append(notInSet2, set1[i])
			i++
		} else {
			i++
			j++
		}
	}
	if i < len1 {
		notInSet2 = append(notInSet2, set1[i:]...)
		diff = true
	} else if j < len2 {
		notInSet1 = append(notInSet1, set2[j:]...)
		diff = true
	}
	return
}

// Indent returns the given text indented.
func Indent(text string, indent string) string {
	return strings.Replace(text, "\n", "\n"+indent, -1)
}

// getStack gets the stack trace.
// Made available for testing.
var getStack = debug.Stack

// StackTrace gets the current stack trace.
func StackTrace(offset int, count int) string {
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
