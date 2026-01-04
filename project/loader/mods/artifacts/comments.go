package artifacts

import (
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

func RemoveDirectives(comments []*ast.Comment, prefix string) []*ast.Comment {
	prefix = `//` + prefix + `:`
	result := make([]*ast.Comment, 0, len(comments))
	for _, c := range comments {
		if !strings.HasPrefix(c.Text, prefix) {
			result = append(result, c)
		}
	}
	return result
}

func CommentsForNode(f *ast.File, n ast.Node) []*ast.Comment {
	comments := []*ast.Comment{}
	for _, cg := range f.Comments {
		if cg != nil {
			for _, c := range cg.List {
				if c.End() > n.Pos() && c.Pos() < n.End() {
					comments = append(comments, c)
				}
			}
		}
	}
	return comments
}
