//go:build experiment

package main

import (
	"fmt"
	"iter"
	"runtime/debug"
)

func Counter(max int) iter.Seq[int] {
	return func(yield func(int) bool) { // funcLit named "main.Counter.func1"
		defer func() {
			fmt.Println(`counter`, max, `defer`)
		}()
		fmt.Println(`counter`, max, `started`)
		for i := range max {
			if !yield(i) {
				fmt.Println(`counter`, max, `exited on`, i)
				return
			}
		}
		fmt.Println(`counter`, max, `done`)
	}
}

func Multiplier(scalar int, in iter.Seq[int]) iter.Seq[int] {
	return func(yield func(int) bool) { // funcLit named "main.Multiplier.func1"
		defer func() {
			fmt.Println(`multiplier`, scalar, `defer`)
		}()
		/*
			fmt.Println(`multiplier`, scalar, `started`)
			for i := range in {
				if v := i * scalar; !yield(v) {
					fmt.Println(`multiplier`, scalar, `exited on`, v)
					return
				}
			}
			fmt.Println(`multiplier`, scalar, `done`)
		*/

		// This is equivalent to the above code except the return and break are
		// the same. This will not work if there was a return value on the
		// outer since then return and break don't work the same and the
		// return value needs to be stashed away until the iteration finishes.
		body := func(i int) bool { // funcLit named "main.Multiplier.func1.2"
			if v := i * scalar; !yield(v) {
				fmt.Println(`multiplier`, scalar, `exited on`, v)
				return false
			}
			return true
		}

		fmt.Println(`multiplier`, scalar, `started`)
		in(body)
		fmt.Println(`multiplier`, scalar, `done`)
	}
}

func main() {
	runFull()
	runFullWonky()
	runBreak()
	runBreakWonky()
	runContinue()
	runReturn()
	runPanic()
	runValueReturn()
	runValueReturnWonky()
	runCallStack()
}

func runFull() {
	fmt.Println()
	fmt.Println(`====[ Full ]==============================`)
	it := Multiplier(3, Multiplier(2, Counter(3)))
	for i := range it {
		fmt.Println(`body got`, i)
	}
	fmt.Println(`body done`)
}

func runFullWonky() {
	fmt.Println()
	fmt.Println(`====[ Full Wonky ]========================`)
	body := func(i int) bool {
		fmt.Println(`body got`, i)
		return true
	}
	Multiplier(3, Multiplier(2, Counter(3)))(body)
	fmt.Println(`body done`)
}

func runBreak() {
	fmt.Println()
	fmt.Println(`====[ Break ]=============================`)
	it := Multiplier(3, Multiplier(2, Counter(3)))
	for i := range it {
		fmt.Println(`body got`, i)
		if i == 6 {
			fmt.Println(`body break`)
			break
		}
	}
	fmt.Println(`body done`)
}

func runBreakWonky() {
	fmt.Println()
	fmt.Println(`====[ Break Wonky ]=======================`)
	body := func(i int) bool {
		fmt.Println(`body got`, i)
		if i == 6 {
			return false
		}
		return true
	}
	Multiplier(3, Multiplier(2, Counter(3)))(body)
	fmt.Println(`body done`)
}

func runContinue() {
	fmt.Println()
	fmt.Println(`====[ Continue ]==========================`)
	it := Multiplier(3, Multiplier(2, Counter(3)))
	for i := range it {
		fmt.Println(`body got`, i)
		if i == 6 {
			fmt.Println(`body continue`)
			continue
		}
	}
	fmt.Println(`body done`)
}

func runReturn() {
	fmt.Println()
	fmt.Println(`====[ Return ]============================`)
	it := Multiplier(3, Multiplier(2, Counter(3)))
	for i := range it {
		fmt.Println(`body got`, i)
		if i == 6 {
			fmt.Println(`body return`)
			return
		}
	}
	fmt.Println(`body done`)
}

func runPanic() {
	fmt.Println()
	fmt.Println(`====[ Panic ]=============================`)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(`got panic:`, r)
		}
	}()
	it := Multiplier(3, Multiplier(2, Counter(3)))
	for i := range it {
		fmt.Println(`body got`, i)
		if i == 6 {
			fmt.Println(`body panic`)
			panic(`oops`)
		}
	}
	fmt.Println(`body done`)
}

func runValueReturn() {
	fmt.Println()
	fmt.Println(`====[ Value Return ]======================`)
	bounds := func(doRet bool) string {
		fmt.Println(`----[ Value Return ]----------------------`)
		fmt.Println(`do return:`, doRet)
		it := Multiplier(3, Multiplier(2, Counter(3)))
		for i := range it {
			fmt.Println(`body got`, i)
			if i == 6 {
				if doRet {
					fmt.Println(`body return`)
					return `cat`
				}
				break
			}
		}
		fmt.Println(`body done`)
		return `dog`
	}
	fmt.Println(`result:`, bounds(false))
	fmt.Println(`result:`, bounds(true))
}

func runValueReturnWonky() {
	fmt.Println()
	fmt.Println(`====[ Value Return Wonky ]================`)
	bounds := func(doRet bool) string {
		fmt.Println(`----[ Value Return Wonky ]----------------`)
		fmt.Println(`do return:`, doRet)

		// The body result is only for this experiment. In the actual
		// blocker, the result will have to handle multiple returns and
		// any other information more elegantly.
		var bodyResult *string
		body := func(i int) bool {
			fmt.Println(`body got`, i)
			if i == 6 {
				if doRet {
					fmt.Println(`body return`)
					bodyResult = new(`cat`)
					return false // return `cat`
				}
				return false // break
			}
			return true
		}
		// Note: bodyResult needs to be reset each time body is used
		// and the values of bodyResult needs to be checked. This shouldn't
		// be a problem since in the non-wonky form, the for-range body is
		// only used once inside of that specific for-range.
		Multiplier(3, Multiplier(2, Counter(3)))(body)
		if bodyResult != nil {
			return *bodyResult
		}

		fmt.Println(`body done`)
		return `dog`
	}
	fmt.Println(`result:`, bounds(false))
	fmt.Println(`result:`, bounds(true))
}

func runCallStack() {
	fmt.Println()
	fmt.Println(`====[ Call Stack ]========================`)
	fmt.Println(`----[ Call Stack ]----------------`)
	for range Multiplier(3, Multiplier(2, Counter(3))) { // body funcLit named "main.runCallStack-range1"
		debug.PrintStack()
		break
	}
	fmt.Println(`----[ Call Stack ]----------------`)
	for range Multiplier(3, Multiplier(2, Counter(3))) { // body funcLit named "main.runCallStack-range2"
		debug.PrintStack()
		break
	}
}
