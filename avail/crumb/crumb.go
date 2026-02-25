package crumb

import (
	"fmt"
	"runtime"
)

const (
	disable = false
	prefix  = `Crumb @ `
)

// Drop puts a "bread crumb" at this locations such that it will
// print the file and line number that the Drop was called from.
//
// This is useful for debugging as a quick method for tracing the code.
func Drop() {
	if disable {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		fmt.Printf("%s%s:%d\n", prefix, file, line)
	} else {
		fmt.Printf("%sUnknown\n", prefix)
	}
}
