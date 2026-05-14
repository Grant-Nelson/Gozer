//go:build testApp

// This tests is a simple test of recursive function calls.
//
// An experiment was created for this test and manually written in
// experiments/exp002 to test out one possible output from compiling this code.
package main

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	println(`fib(-1) = `, fib(-1))
	println(`fib(2) = `, fib(2))
	println(`fib(5) = `, fib(5))
	println(`fib(10) = `, fib(10))
}

// Output:
// fib(-1) = -1
// fib(2) = 1
// fib(5) = 5
// fib(10) = 55
