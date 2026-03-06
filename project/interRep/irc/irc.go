package irc

import (
	"fmt"
	"go/token"
	"strings"
)

func csvString[E any, S []E](s S) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = fmt.Sprintf(`%v`, elem)
	}
	return strings.Join(elems, `, `)
}

func linesString[E any, S []E](s S, indent string) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := fmt.Sprintf(`%s%v`, indent, elem)
		elems[i] = strings.ReplaceAll(eStr, "\n", "\n"+indent)
	}
	return strings.Join(elems, "\n")
}

func endOfSlice[E interface{ End() token.Pos }, S []E](s S) token.Pos {
	if n := len(s); n > 0 {
		return s[n-1].End()
	}
	return token.NoPos
}
