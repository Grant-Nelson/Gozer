// This package contains bread crumbs. Hansel and Gretel dropped bread crumbs
// to help them find there way back. The bread crumbs indicate where they
// had been. This package is the same idea but for tracing code.
//
// While debugging, sometimes stepping through code isn't very useful
// if you don't know where you should put a break or stepping through the code
// isn't possible. In those cases leaving crumbs will help you narrow down
// where the code has flown. You drop crumbs around the code you're debugging
// an it will output the location whenever a crumb is reached.
// Some devs will do this by dropping in a something like `println("Flag 1")`,
// `println("Flag 2.A.x")`, `println(">> Boom!")`, etc. The crumbs is
// a similar idea except that it will output the file path and function name.
//
// The crumbs are easily disabled so that you can run without them on but
// without removing them. All bread crumbs should be removed when done
// debugging. Since the crumbs aren't just random print statements but are in
// this package, they are much easier to find and remove.
//
// If you want to preserve standard out or standard out get cleared, the
// crumbs can be redirected to a different output.
package crumb

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	prefix  = `Crumb @ `
	unknown = `unknown`
)

// Disable will disable dropping crumbs.
var Disable = false

// Out is the output to write the crumbs to. Defaults to standard output.
// If Out is nil, the crumb will be skipped the same as if crumbs were disabled.
var Out = os.Stdout

// Drop puts a "bread crumb" at the location this is called at
// such that it will print the location of the method calling
// this crumb.
//
// This is useful for debugging as a quick method for tracing the code.
func Drop() {
	if Disable || Out == nil {
		return
	}
	fmt.Fprintf(Out, "%s%s\n", prefix, Location(1, true))
	Out.Sync()
}

// DropMsg puts a "bread crumb" at the location this is called at
// such that it will print the location of the method calling
// this crumb followed by the given custom message.
//
// This is useful for debugging as a quick method for tracing the code
// and outputting specific values or additional context with the crumb.
func DropMsg(format string, args ...any) {
	if Disable || Out == nil {
		return
	}
	if msg := fmt.Sprintf(format, args...); len(msg) > 0 {
		fmt.Fprintf(Out, "%s%s: %s\n", prefix, Location(1, true), msg)
	} else {
		fmt.Fprintf(Out, "%s%s\n", prefix, Location(1, true))
	}
	Out.Sync()
}

// Location gets the location that a bread crumb would output.
// Skip is the number of stack frames to skip where 0 will return the
// function that called this function, 1 will return the function that
// called the function that called this function, and so on.
func Location(skip int, includeFuncName bool) string {
	if pc, file, line, ok := runtime.Caller(skip + 1); ok {
		if includeFuncName {
			if fn := runtime.FuncForPC(pc); fn != nil {
				fnName := fn.Name()
				if index := strings.LastIndex(fnName, "/"); index >= 0 {
					fnName = fnName[index+1:]
				}
				return fmt.Sprintf("%s:%d in %s", file, line, fnName)
			}
		}
		return fmt.Sprintf("%s:%d", file, line)
	}
	return unknown
}
