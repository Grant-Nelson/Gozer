# SSA Notes

These are results from some runs of `experiments\exp004\main.go`.

References:

- [SSA tools](https://pkg.go.dev/golang.org/x/tools/go/ssa)
- [SSA in linters](https://nick-we.com/251119_ssa-in-go-linters/)
- [Wiki](https://en.wikipedia.org/wiki/Static_single-assignment_form)
- [Research](https://www.dcs.gla.ac.uk/~jsinger/ssa.html)

## Example 1

Running the following program in SSA analysis and printing `main()`
results in the output below.

```Go
package main

import "fmt"

func main() {
    message := "Hello, World!"
    for _, r := range message {
        if r >= 'a' && r <= 'z' {
            r = r - 'a' + 'A'
        }
        fmt.Println(string(r))
    }
}
```

```Plain
# Name: hello.main
# Package: hello
# Location: hello.go:6:6
func main():
0:                                                                entry P:0 S:1
        t0 = range "Hello, World!":string                        $ssa.rangeIter
        jump 1
1:                                                rangeiter.loop P:2 S:2 idom:0
        t1 = next t0                          (ok bool, k invalid type, v rune)
        t2 = extract t1 #0                                                 bool
        if t2 goto 2 else 3
2:                                                rangeiter.body P:1 S:2 idom:1
        t3 = extract t1 #2                                                 rune
        t4 = t3 >= 97:rune                                                 bool
        if t4 goto 6 else 5
3:                                                rangeiter.done P:1 S:0 idom:1
        return
4:                                                       if.then P:1 S:1 idom:6
        t5 = t3 - 97:rune                                                  rune
        t6 = t5 + 65:rune                                                  rune
        jump 5
5:                                                       if.done P:3 S:1 idom:2
        t7 = phi [2: t3, 6: t3, 4: t6] #r                                  rune
        t8 = convert string <- rune (t7)                                 string
        t9 = new [1]any (varargs)                                       *[1]any
        t10 = &t9[0:int]                                                   *any
        t11 = make any <- string (t8)                                       any
        *t10 = t11
        t12 = slice t9[:]                                                 []any
        t13 = fmt.Println(t12...)                            (n int, err error)
        jump 1
6:                                                     cond.true P:1 S:2 idom:2
        t14 = t3 <= 122:rune                                               bool
        if t14 goto 4 else 5
```
