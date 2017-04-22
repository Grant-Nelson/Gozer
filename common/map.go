package common

import (
	"fmt"
	"sort"
	"strings"
)

// Map is a set of key to generic values.
type Map map[string]interface{}

// NewMap creates a new generic map.
func NewMap() Map {
	return map[string]interface{}{}
}

// Add sets the entry for a map.
func (m Map) Add(key string, val ...interface{}) Map {
	m[key] = fmt.Sprint(val...)
	return m
}

// String creates a string for this map.
func (m Map) String() string {
	return m.FormatMap("")
}

// FormatMap creates a formatted string for the
// given map to make the map easily readable.
func (m Map) FormatMap(indent string) string {
	maxKeyLen := 0
	for key := range m {
		if keyLen := len(key); keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}
	parts := make([]string, len(m))
	offset := "\n" + indent + strings.Repeat(" ", maxKeyLen) + "  "
	i := 0
	for key, val := range m {
		padding := strings.Repeat(" ", maxKeyLen-len(key))
		valStr := strings.Replace(fmt.Sprint(val), "\n", offset, -1)
		parts[i] = fmt.Sprint(key, ": ", padding, valStr)
		i++
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n"+indent)
}
