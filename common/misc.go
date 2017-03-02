package common

import (
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
