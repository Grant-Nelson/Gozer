package artifacts

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func (f *File) Remap(fileSet *FileSet) {
	frm := &fileRemapper{
		f:      f,
		offset: 1,
	}
	walkPos(f.File, frm.mapPos)
	frm.finish(fileSet)
}

type fileRemapper struct {
	f       *File
	offset  int
	edits   []remapperEdit
	expNext int
}

type remapperEdit func(f *token.File)

func (frm *fileRemapper) finish(fileSet *FileSet) {
	p := frm.f.FileSet.Position(frm.f.File.FileStart)
	f := fileSet.fileSet.AddFile(p.Filename, 1, frm.offset)
	for _, e := range frm.edits {
		e(f)
	}
	frm.f.FileSet = fileSet
}

func (frm *fileRemapper) mapPos(n ast.Node, off *token.Pos) {
	if !off.IsValid() {
		return
	}

	if int(*off) != frm.expNext {
		pos := frm.f.FileSet.Position(*off)
		frm.edits = append(frm.edits, func(f *token.File) {
			f.AddLineColumnInfo(frm.offset, pos.Filename, pos.Line, pos.Column)
		})
	}

	*off = token.Pos(frm.offset)
	total, lines := frm.f.FileSet.Widths(*off)
	frm.expNext = int(*off) + total

	for i, ln := range lines {
		if i > 1 {
			offset := frm.offset
			frm.edits = append(frm.edits, func(f *token.File) {
				f.AddLine(offset)
			})
		}
		frm.offset += ln
	}
}

type ComparableNode interface {
	ast.Node
	comparable
}

type posVisitor func(n ast.Node, off *token.Pos)

func walkPosVisit(n ast.Node, off *token.Pos, v posVisitor) {
	if off.IsValid() {
		v(n, off)
	}
}

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
		walkPosVisit(n, &n.Slash, v)

	case *ast.CommentGroup:
		walkPosList(n.List, v)

	case *ast.Field:
		walkPos(n.Doc, v)
		walkPosList(n.Names, v)
		walkPos(n.Type, v)
		walkPos(n.Tag, v)
		walkPos(n.Comment, v)

	case *ast.FieldList:
		walkPosVisit(n, &n.Opening, v)
		walkPosList(n.List, v)
		walkPosVisit(n, &n.Closing, v)

	// Expressions
	case *ast.BadExpr:
		walkPosVisit(n, &n.From, v)
		walkPosVisit(n, &n.To, v)

	case *ast.Ident:
		walkPosVisit(n, &n.NamePos, v)

	case *ast.Ellipsis:
		walkPosVisit(n, &n.Ellipsis, v)
		walkPos(n.Elt, v)

	case *ast.BasicLit:
		walkPosVisit(n, &n.ValuePos, v)

	case *ast.FuncLit:
		walkPos(n.Type, v)
		walkPos(n.Body, v)

	case *ast.CompositeLit:
		walkPos(n.Type, v)
		walkPosVisit(n, &n.Lbrace, v)
		walkPosList(n.Elts, v)
		walkPosVisit(n, &n.Rbrace, v)

	case *ast.ParenExpr:
		walkPosVisit(n, &n.Lparen, v)
		walkPos(n.X, v)
		walkPosVisit(n, &n.Rparen, v)

	case *ast.SelectorExpr:
		walkPos(n.X, v)
		walkPos(n.Sel, v)

	case *ast.IndexExpr:
		walkPos(n.X, v)
		walkPosVisit(n, &n.Lbrack, v)
		walkPos(n.Index, v)
		walkPosVisit(n, &n.Rbrack, v)

	case *ast.IndexListExpr:
		walkPos(n.X, v)
		walkPosVisit(n, &n.Lbrack, v)
		walkPosList(n.Indices, v)
		walkPosVisit(n, &n.Rbrack, v)

	case *ast.SliceExpr:
		walkPos(n.X, v)
		walkPosVisit(n, &n.Lbrack, v)
		walkPos(n.Low, v)
		walkPos(n.High, v)
		walkPos(n.Max, v)
		walkPosVisit(n, &n.Rbrack, v)

	case *ast.TypeAssertExpr:
		walkPos(n.X, v)
		walkPosVisit(n, &n.Lparen, v)
		walkPos(n.Type, v)
		walkPosVisit(n, &n.Rparen, v)

	case *ast.CallExpr:
		walkPos(n.Fun, v)
		walkPosVisit(n, &n.Lparen, v)
		walkPosList(n.Args, v)
		// TODO: Determine where the ellipsis needs to go
		walkPosVisit(n, &n.Ellipsis, v)
		walkPosVisit(n, &n.Rparen, v)

	case *ast.StarExpr:
		walkPosVisit(n, &n.Star, v)
		walkPos(n.X, v)

	case *ast.UnaryExpr:
		walkPosVisit(n, &n.OpPos, v)
		walkPos(n.X, v)

	case *ast.BinaryExpr:
		walkPos(n.X, v)
		walkPosVisit(n, &n.OpPos, v)
		walkPos(n.Y, v)

	case *ast.KeyValueExpr:
		walkPos(n.Key, v)
		walkPosVisit(n, &n.Colon, v)
		walkPos(n.Value, v)

	// Types
	case *ast.ArrayType:
		walkPosVisit(n, &n.Lbrack, v)
		// TODO: Determine why no Rbrack?
		walkPos(n.Len, v)
		walkPos(n.Elt, v)

	case *ast.StructType:
		walkPosVisit(n, &n.Struct, v)
		walkPos(n.Fields, v)

	case *ast.FuncType:
		walkPosVisit(n, &n.Func, v)
		walkPos(n.TypeParams, v)
		walkPos(n.Params, v)
		walkPos(n.Results, v)

	case *ast.InterfaceType:
		walkPosVisit(n, &n.Interface, v)
		walkPos(n.Methods, v)

	case *ast.MapType:
		walkPosVisit(n, &n.Map, v)
		walkPos(n.Key, v)
		walkPos(n.Value, v)

	case *ast.ChanType:
		// TODO: Check the order is correct for the Begin and Arrow
		walkPosVisit(n, &n.Begin, v)
		walkPosVisit(n, &n.Arrow, v)
		walkPos(n.Value, v)

	// Statements
	case *ast.BadStmt:
		walkPosVisit(n, &n.From, v)
		walkPosVisit(n, &n.To, v)

	case *ast.DeclStmt:
		walkPos(n.Decl, v)

	case *ast.EmptyStmt:
		walkPosVisit(n, &n.Semicolon, v)

	case *ast.LabeledStmt:
		walkPos(n.Label, v)
		walkPosVisit(n, &n.Colon, v)
		walkPos(n.Stmt, v)

	case *ast.ExprStmt:
		walkPos(n.X, v)

	case *ast.SendStmt:
		walkPos(n.Chan, v)
		walkPosVisit(n, &n.Arrow, v)
		walkPos(n.Value, v)

	case *ast.IncDecStmt:
		walkPos(n.X, v)
		walkPosVisit(n, &n.TokPos, v)

	case *ast.AssignStmt:
		walkPosList(n.Lhs, v)
		walkPosVisit(n, &n.TokPos, v)
		walkPosList(n.Rhs, v)

	case *ast.GoStmt:
		walkPosVisit(n, &n.Go, v)
		walkPos(n.Call, v)

	case *ast.DeferStmt:
		walkPosVisit(n, &n.Defer, v)
		walkPos(n.Call, v)

	case *ast.ReturnStmt:
		walkPosVisit(n, &n.Return, v)
		walkPosList(n.Results, v)

	case *ast.BranchStmt:
		walkPosVisit(n, &n.TokPos, v)
		walkPos(n.Label, v)

	case *ast.BlockStmt:
		walkPosVisit(n, &n.Lbrace, v)
		walkPosList(n.List, v)
		walkPosVisit(n, &n.Rbrace, v)

	case *ast.IfStmt:
		walkPosVisit(n, &n.If, v)
		walkPos(n.Init, v)
		walkPos(n.Cond, v)
		walkPos(n.Body, v)
		walkPos(n.Else, v)

	case *ast.CaseClause:
		walkPosVisit(n, &n.Case, v)
		walkPosList(n.List, v)
		walkPosVisit(n, &n.Colon, v)
		walkPosList(n.Body, v)

	case *ast.SwitchStmt:
		walkPosVisit(n, &n.Switch, v)
		walkPos(n.Init, v)
		walkPos(n.Tag, v)
		walkPos(n.Body, v)

	case *ast.TypeSwitchStmt:
		walkPosVisit(n, &n.Switch, v)
		walkPos(n.Init, v)
		walkPos(n.Assign, v)
		walkPos(n.Body, v)

	case *ast.CommClause:
		walkPosVisit(n, &n.Case, v)
		walkPos(n.Comm, v)
		walkPosVisit(n, &n.Colon, v)
		walkPosList(n.Body, v)

	case *ast.SelectStmt:
		walkPosVisit(n, &n.Select, v)
		walkPos(n.Body, v)

	case *ast.ForStmt:
		walkPosVisit(n, &n.For, v)
		walkPos(n.Init, v)
		walkPos(n.Cond, v)
		walkPos(n.Post, v)
		walkPos(n.Body, v)

	case *ast.RangeStmt:
		walkPosVisit(n, &n.For, v)
		walkPos(n.Key, v)
		walkPos(n.Value, v)
		walkPosVisit(n, &n.TokPos, v)
		walkPosVisit(n, &n.Range, v)
		walkPos(n.X, v)
		walkPos(n.Body, v)

	// Declarations
	case *ast.ImportSpec:
		walkPos(n.Doc, v)
		walkPos(n.Name, v)
		walkPos(n.Path, v)
		walkPos(n.Comment, v)
		walkPosVisit(n, &n.EndPos, v)

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
		walkPosVisit(n, &n.Assign, v)
		walkPos(n.Type, v)
		walkPos(n.Comment, v)

	case *ast.BadDecl:
		walkPosVisit(n, &n.From, v)
		walkPosVisit(n, &n.To, v)

	case *ast.GenDecl:
		walkPos(n.Doc, v)
		walkPosVisit(n, &n.TokPos, v)
		walkPosVisit(n, &n.Lparen, v)
		walkPosList(n.Specs, v)
		walkPosVisit(n, &n.Rparen, v)

	case *ast.FuncDecl:
		walkPos(n.Doc, v)
		walkPos(n.Recv, v)
		walkPos(n.Name, v)
		walkPos(n.Type, v)
		walkPos(n.Body, v)

	// Files
	case *ast.File:
		walkPosVisit(n, &n.FileStart, v)
		walkPos(n.Doc, v)
		walkPosVisit(n, &n.Package, v)
		walkPos(n.Name, v)
		walkPosList(n.Decls, v)
		walkPosVisit(n, &n.FileEnd, v)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]v`, n))
	}
}
