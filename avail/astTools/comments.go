package astTools

import (
	"go/ast"
	"go/token"
	"slices"

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

func RemoveDirectives(cg *ast.CommentGroup, directives ...*ast.Directive) {
	if cg == nil || len(cg.List) <= 0 || len(directives) <= 0 {
		return
	}
	pos := map[token.Pos]bool{}
	for _, d := range directives {
		pos[d.Pos()] = true
	}
	cg.List = slices.DeleteFunc(cg.List,
		func(c *ast.Comment) bool { return pos[c.Pos()] })
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
