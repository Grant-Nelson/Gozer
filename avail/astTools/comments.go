package astTools

import (
	"go/ast"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

func Directives(comments []*ast.Comment) iterator.Iterator[*ast.Directive] {
	return func(yield func(*ast.Directive) bool) {
		for _, cm := range comments {
			if cm != nil {
				d := DirectivesFromComment(cm)
				if d != nil && !yield(d) {
					return
				}
			}
		}
	}
}

func DirectivesFromGroup(cg *ast.CommentGroup) iterator.Iterator[*ast.Directive] {
	if cg == nil {
		return iterator.Empty[*ast.Directive]()
	}
	return Directives(cg.List)
}

func DirectivesFromComment(cm *ast.Comment) *ast.Directive {
	if d, ok := ast.ParseDirective(cm.Pos(), cm.Text); ok {
		return &d
	}
	return nil
}

// TODO: Continue updating below

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
			if cg.End() > n.Pos() && cg.Pos() < n.End() {
				comments = append(comments, cg)
			}
		}
	}
	return comments
}
