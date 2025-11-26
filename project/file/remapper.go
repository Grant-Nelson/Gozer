package file

import (
	"go/ast"
	"go/token"
)

type fileSetFrame struct {
	offset   int
	filename string
	line     int
	column   int
}

type fileRemapper struct {
	f      *File
	offset int
	frames []fileSetFrame
}

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func (f *File) Remap() {
	frm := &fileRemapper{f: f}
	for n := range f.Nodes() {
		if n.Node != nil {
			mapPos(n.Node, frm.posMapper)
		}
	}
	frm.finish()
}

func (frm *fileRemapper) posMapper(pos token.Pos) token.Pos {
	// TODO: Implement
	return pos
}

func (frm *fileRemapper) finish() {
	fileSet := token.NewFileSet()

	// TODO: write to fileSet.
	//f2 := &token.File{}
	//f2.AddLineColumnInfo(offset, filename, line, column)

	frm.f.FileSet = fileSet
}

func mapPos(n ast.Node, handle func(pos token.Pos) token.Pos) {
	switch t := n.(type) {

	// Comments
	case *ast.Comment:
		t.Slash = handle(t.Slash)
	case *ast.CommentGroup:
		// none

	// Expression Parts
	case *ast.Field:
		// none
	case *ast.FieldList:
		t.Opening = handle(t.Opening)
		t.Closing = handle(t.Closing)

	// Expressions
	case *ast.BadExpr:
		t.From = handle(t.From)
		t.To = handle(t.To)
	case *ast.Ident:
		t.NamePos = handle(t.NamePos)
	case *ast.Ellipsis:
		t.Ellipsis = handle(t.Ellipsis)
	case *ast.BasicLit:
		t.ValuePos = handle(t.ValuePos)
	case *ast.FuncLit:
		// none
	case *ast.CompositeLit:
		t.Lbrace = handle(t.Lbrace)
		t.Rbrace = handle(t.Rbrace)
	case *ast.ParenExpr:
		t.Lparen = handle(t.Lparen)
		t.Rparen = handle(t.Rparen)
	case *ast.SelectorExpr:
		// none
	case *ast.IndexExpr:
		t.Lbrack = handle(t.Lbrack)
		t.Rbrack = handle(t.Rbrack)
	case *ast.IndexListExpr:
		t.Lbrack = handle(t.Lbrack)
		t.Rbrack = handle(t.Rbrack)
	case *ast.SliceExpr:
		t.Lbrack = handle(t.Lbrack)
		t.Rbrack = handle(t.Rbrack)
	case *ast.TypeAssertExpr:
		t.Lparen = handle(t.Lparen)
		t.Rparen = handle(t.Rparen)
	case *ast.CallExpr:
		t.Lparen = handle(t.Lparen)
		t.Ellipsis = handle(t.Ellipsis)
		t.Rparen = handle(t.Rparen)
	case *ast.StarExpr:
		t.Star = handle(t.Star)
	case *ast.UnaryExpr:
		t.OpPos = handle(t.OpPos)
	case *ast.BinaryExpr:
		t.OpPos = handle(t.OpPos)
	case *ast.KeyValueExpr:
		t.Colon = handle(t.Colon)

		// Types
	case *ast.ArrayType:

	}
}
