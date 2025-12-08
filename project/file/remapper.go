package file

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func (f *File) Remap(fileSet *token.FileSet) {
	frm := &fileRemapper{
		f:      f,
		offset: 1,
	}
	fmt.Printf("===============\n")
	walkPos(f.File, frm.mapPos)
	frm.finish(fileSet)
}

type fileRemapper struct {
	f      *File
	offset int
	lines  []token.Position

	priorOff token.Pos
	priorPos token.Position
}

func (frm *fileRemapper) finish(fileSet *token.FileSet) {
	p := frm.f.FileSet.Position(frm.f.File.FileStart)
	f := fileSet.AddFile(p.Filename, 1, frm.offset)
	var prior token.Position
	for _, ln := range frm.lines {
		if prior.Filename != ln.Filename {
			f.AddLineColumnInfo(ln.Offset, ln.Filename, ln.Line, ln.Column)
		} else if prior.Line != ln.Line {
			f.AddLine(ln.Line)
		}
		prior = ln
	}
	frm.f.FileSet = fileSet
}

func (frm *fileRemapper) mapPos(n ast.Node, off *token.Pos) {
	if !off.IsValid() {
		return
	}

	pos := frm.f.FileSet.Position(*off)
	if pos.Filename != frm.priorPos.Filename || pos.Line < frm.priorPos.Line {
		fmt.Printf(">%q:\n", pos.Filename)

		// TODO: Implement
	}

	// TODO: Implement

	fmt.Printf("\t[offset: %2d]  [offset: %2d, line: %2d, column: %2d]  %T\n",
		frm.offset, int(*off), pos.Line, pos.Column, n)

	frm.priorOff = *off
	frm.priorPos = pos
}

type ComparableNode interface {
	ast.Node
	comparable
}

type posVisitor func(n ast.Node, off *token.Pos)

func walkPosList[N ComparableNode](list []N, v posVisitor) {
	for _, node := range list {
		walkPos(node, v)
	}
}

func walkPos[N ComparableNode](node N, v posVisitor) {
	var zero N
	if node == zero {
		return
	}

	switch n := any(node).(type) {
	case nil:
		// do nothing

	// Comments and fields
	case *ast.Comment:
		v(n, &n.Slash)

	case *ast.CommentGroup:
		walkPosList(n.List, v)

	case *ast.Field:
		walkPos(n.Doc, v)
		walkPosList(n.Names, v)
		walkPos(n.Type, v)
		walkPos(n.Tag, v)
		walkPos(n.Comment, v)

	case *ast.FieldList:
		v(n, &n.Opening)
		walkPosList(n.List, v)
		v(n, &n.Closing)

	// Expressions
	case *ast.BadExpr:
		v(n, &n.From)
		v(n, &n.To)

	case *ast.Ident:
		v(n, &n.NamePos)

	case *ast.Ellipsis:
		v(n, &n.Ellipsis)
		walkPos(n.Elt, v)

	case *ast.BasicLit:
		v(n, &n.ValuePos)

	case *ast.FuncLit:
		walkPos(n.Type, v)
		walkPos(n.Body, v)

	case *ast.CompositeLit:
		walkPos(n.Type, v)
		v(n, &n.Lbrace)
		walkPosList(n.Elts, v)
		v(n, &n.Rbrace)

	case *ast.ParenExpr:
		v(n, &n.Lparen)
		walkPos(n.X, v)
		v(n, &n.Rparen)

	case *ast.SelectorExpr:
		walkPos(n.X, v)
		walkPos(n.Sel, v)

	case *ast.IndexExpr:
		walkPos(n.X, v)
		v(n, &n.Lbrack)
		walkPos(n.Index, v)
		v(n, &n.Rbrack)

	case *ast.IndexListExpr:
		walkPos(n.X, v)
		v(n, &n.Lbrack)
		walkPosList(n.Indices, v)
		v(n, &n.Rbrack)

	case *ast.SliceExpr:
		walkPos(n.X, v)
		v(n, &n.Lbrack)
		walkPos(n.Low, v)
		walkPos(n.High, v)
		walkPos(n.Max, v)
		v(n, &n.Rbrack)

	case *ast.TypeAssertExpr:
		walkPos(n.X, v)
		v(n, &n.Lparen)
		walkPos(n.Type, v)
		v(n, &n.Rparen)

	case *ast.CallExpr:
		walkPos(n.Fun, v)
		v(n, &n.Lparen)
		walkPosList(n.Args, v)
		// TODO: Determine where the ellipsis needs to go
		v(n, &n.Ellipsis)
		v(n, &n.Rparen)

	case *ast.StarExpr:
		v(n, &n.Star)
		walkPos(n.X, v)

	case *ast.UnaryExpr:
		v(n, &n.OpPos)
		walkPos(n.X, v)

	case *ast.BinaryExpr:
		walkPos(n.X, v)
		v(n, &n.OpPos)
		walkPos(n.Y, v)

	case *ast.KeyValueExpr:
		walkPos(n.Key, v)
		v(n, &n.Colon)
		walkPos(n.Value, v)

	// Types
	case *ast.ArrayType:
		v(n, &n.Lbrack)
		// TODO: Determine why no Rbrack?
		walkPos(n.Len, v)
		walkPos(n.Elt, v)

	case *ast.StructType:
		v(n, &n.Struct)
		walkPos(n.Fields, v)

	case *ast.FuncType:
		v(n, &n.Func)
		walkPos(n.TypeParams, v)
		walkPos(n.Params, v)
		walkPos(n.Results, v)

	case *ast.InterfaceType:
		v(n, &n.Interface)
		walkPos(n.Methods, v)

	case *ast.MapType:
		v(n, &n.Map)
		walkPos(n.Key, v)
		walkPos(n.Value, v)

	case *ast.ChanType:
		// TODO: Check the order is correct for the Begin and Arrow
		v(n, &n.Begin)
		v(n, &n.Arrow)
		walkPos(n.Value, v)

	// Statements
	case *ast.BadStmt:
		v(n, &n.From)
		v(n, &n.To)

	case *ast.DeclStmt:
		walkPos(n.Decl, v)

	case *ast.EmptyStmt:
		v(n, &n.Semicolon)

	case *ast.LabeledStmt:
		walkPos(n.Label, v)
		v(n, &n.Colon)
		walkPos(n.Stmt, v)

	case *ast.ExprStmt:
		walkPos(n.X, v)

	case *ast.SendStmt:
		walkPos(n.Chan, v)
		v(n, &n.Arrow)
		walkPos(n.Value, v)

	case *ast.IncDecStmt:
		walkPos(n.X, v)
		v(n, &n.TokPos)

	case *ast.AssignStmt:
		walkPosList(n.Lhs, v)
		v(n, &n.TokPos)
		walkPosList(n.Rhs, v)

	case *ast.GoStmt:
		v(n, &n.Go)
		walkPos(n.Call, v)

	case *ast.DeferStmt:
		v(n, &n.Defer)
		walkPos(n.Call, v)

	case *ast.ReturnStmt:
		v(n, &n.Return)
		walkPosList(n.Results, v)

	case *ast.BranchStmt:
		v(n, &n.TokPos)
		walkPos(n.Label, v)

	case *ast.BlockStmt:
		v(n, &n.Lbrace)
		walkPosList(n.List, v)
		v(n, &n.Rbrace)

	case *ast.IfStmt:
		v(n, &n.If)
		walkPos(n.Init, v)
		walkPos(n.Cond, v)
		walkPos(n.Body, v)
		walkPos(n.Else, v)

	case *ast.CaseClause:
		v(n, &n.Case)
		walkPosList(n.List, v)
		v(n, &n.Colon)
		walkPosList(n.Body, v)

	case *ast.SwitchStmt:
		v(n, &n.Switch)
		walkPos(n.Init, v)
		walkPos(n.Tag, v)
		walkPos(n.Body, v)

	case *ast.TypeSwitchStmt:
		v(n, &n.Switch)
		walkPos(n.Init, v)
		walkPos(n.Assign, v)
		walkPos(n.Body, v)

	case *ast.CommClause:
		v(n, &n.Case)
		walkPos(n.Comm, v)
		v(n, &n.Colon)
		walkPosList(n.Body, v)

	case *ast.SelectStmt:
		v(n, &n.Select)
		walkPos(n.Body, v)

	case *ast.ForStmt:
		v(n, &n.For)
		walkPos(n.Init, v)
		walkPos(n.Cond, v)
		walkPos(n.Post, v)
		walkPos(n.Body, v)

	case *ast.RangeStmt:
		v(n, &n.For)
		walkPos(n.Key, v)
		walkPos(n.Value, v)
		v(n, &n.TokPos)
		v(n, &n.Range)
		walkPos(n.X, v)
		walkPos(n.Body, v)

	// Declarations
	case *ast.ImportSpec:
		walkPos(n.Doc, v)
		walkPos(n.Name, v)
		walkPos(n.Path, v)
		walkPos(n.Comment, v)
		v(n, &n.EndPos)

	case *ast.ValueSpec:
		walkPos(n.Doc, v)
		walkPosList(n.Names, v)
		walkPos(n.Type, v)
		walkPosList(n.Values, v)
		walkPos(n.Comment, v)

	case *ast.TypeSpec:
		walkPos(n.Doc, v)
		walkPos(n.Name, v)
		walkPos(n.TypeParams, v)
		v(n, &n.Assign)
		walkPos(n.Type, v)
		walkPos(n.Comment, v)

	case *ast.BadDecl:
		v(n, &n.From)
		v(n, &n.To)

	case *ast.GenDecl:
		walkPos(n.Doc, v)
		v(n, &n.TokPos)
		v(n, &n.Lparen)
		walkPosList(n.Specs, v)
		v(n, &n.Rparen)

	case *ast.FuncDecl:
		walkPos(n.Doc, v)
		walkPos(n.Recv, v)
		walkPos(n.Name, v)
		walkPos(n.Type, v)
		walkPos(n.Body, v)

	// Files
	case *ast.File:
		v(n, &n.FileStart)
		walkPos(n.Doc, v)
		v(n, &n.Package)
		walkPos(n.Name, v)
		walkPosList(n.Decls, v)
		v(n, &n.FileEnd)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]v`, n))
	}
}
