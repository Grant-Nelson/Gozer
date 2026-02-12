package astTools

import (
	"fmt"
	"go/ast"
	"strings"
)

// Directives finds all the directives with the given prefix.
func Directives(comments []*ast.Comment, prefix string) map[string][]string {
	prefix = `//` + prefix + `:`
	result := map[string][]string{}
	for _, c := range comments {
		if tail, ok := strings.CutPrefix(c.Text, prefix); ok {
			var key, value string
			if i := strings.Index(tail, ` `); i > 0 {
				key = strings.TrimSpace(tail[:i])
				value = strings.TrimSpace(tail[i:])
			} else {
				key = tail
				value = ``
			}
			result[key] = append(result[key], value)
		}
	}
	return result
}

func RemoveDirectives(cg *ast.CommentGroup, prefix string) {
	if cg == nil || len(cg.List) <= 0 {
		return
	}
	prefix = `//` + prefix + `:`
	result := make([]*ast.Comment, 0, len(cg.List))
	for _, c := range cg.List {
		if !strings.HasPrefix(c.Text, prefix) {
			result = append(result, c)
		}
	}
	cg.List = result
}

func CommentsAttachedToNode(n ast.Node) []*ast.CommentGroup {
	comments := []*ast.CommentGroup{}
	ast.Inspect(n, func(n ast.Node) bool {
		switch n := n.(type) {
		case nil:
			return true
		case *ast.CommentGroup:
			comments = append(comments, n)
			return false
		}
		return true
	})
	return comments
}

func FileCommentsForNode(f *ast.File, n ast.Node) []*ast.CommentGroup {
	comments := []*ast.CommentGroup{}
	for _, cg := range f.Comments {
		if cg != nil {
			fmt.Printf("::: %d > %d (%t) && %d < %d (%t) ::: %q\n",
				cg.End(), n.Pos(), cg.End() > n.Pos(),
				cg.Pos(), n.End(), cg.Pos() < n.End(),
				cg.Text()) // TODO: REMOVE

			if cg.End() > n.Pos() && cg.Pos() < n.End() {
				comments = append(comments, cg)
			}
		}
	}
	return comments
}
