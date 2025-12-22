package artifacts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"iter"
)

type ComparableNode interface {
	ast.Node
	comparable
}

var errEndWalkPos = errors.New(`end WalkPos`)

func WalkPos[N ComparableNode](node N) iter.Seq2[ast.Node, *token.Pos] {
	return func(yield func(n ast.Node, off *token.Pos) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndWalkPos {
				panic(r)
			}
		}()

		walkPos(node, yield)
	}
}

type posVisitor func(n ast.Node, off *token.Pos) bool

func walkPosVisit(n ast.Node, off *token.Pos, yield posVisitor) {
	if off.IsValid() {
		if !yield(n, off) {
			panic(errEndWalkPos)
		}
	}
}

func walkPosList[N ComparableNode](list []N, yield posVisitor) {
	for _, node := range list {
		walkPos(node, yield)
	}
}

func walkPos[N ComparableNode](node N, yield posVisitor) {
	var zero N
	if node == zero {
		return
	}

	switch n := any(node).(type) {
	case nil:
		// do nothing

	// Comments and fields
	case *ast.Comment:
		walkPosVisit(n, &n.Slash, yield)

	case *ast.CommentGroup:
		walkPosList(n.List, yield)

	case *ast.Field:
		walkPos(n.Doc, yield)
		walkPosList(n.Names, yield)
		walkPos(n.Type, yield)
		walkPos(n.Tag, yield)
		walkPos(n.Comment, yield)

	case *ast.FieldList:
		walkPosVisit(n, &n.Opening, yield)
		walkPosList(n.List, yield)
		walkPosVisit(n, &n.Closing, yield)

	// Expressions
	case *ast.BadExpr:
		walkPosVisit(n, &n.From, yield)
		walkPosVisit(n, &n.To, yield)

	case *ast.Ident:
		walkPosVisit(n, &n.NamePos, yield)

	case *ast.Ellipsis:
		walkPosVisit(n, &n.Ellipsis, yield)
		walkPos(n.Elt, yield)

	case *ast.BasicLit:
		walkPosVisit(n, &n.ValuePos, yield)

	case *ast.FuncLit:
		walkPos(n.Type, yield)
		walkPos(n.Body, yield)

	case *ast.CompositeLit:
		walkPos(n.Type, yield)
		walkPosVisit(n, &n.Lbrace, yield)
		walkPosList(n.Elts, yield)
		walkPosVisit(n, &n.Rbrace, yield)

	case *ast.ParenExpr:
		walkPosVisit(n, &n.Lparen, yield)
		walkPos(n.X, yield)
		walkPosVisit(n, &n.Rparen, yield)

	case *ast.SelectorExpr:
		walkPos(n.X, yield)
		walkPos(n.Sel, yield)

	case *ast.IndexExpr:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.Lbrack, yield)
		walkPos(n.Index, yield)
		walkPosVisit(n, &n.Rbrack, yield)

	case *ast.IndexListExpr:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.Lbrack, yield)
		walkPosList(n.Indices, yield)
		walkPosVisit(n, &n.Rbrack, yield)

	case *ast.SliceExpr:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.Lbrack, yield)
		walkPos(n.Low, yield)
		walkPos(n.High, yield)
		walkPos(n.Max, yield)
		walkPosVisit(n, &n.Rbrack, yield)

	case *ast.TypeAssertExpr:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.Lparen, yield)
		walkPos(n.Type, yield)
		walkPosVisit(n, &n.Rparen, yield)

	case *ast.CallExpr:
		walkPos(n.Fun, yield)
		walkPosVisit(n, &n.Lparen, yield)
		walkPosList(n.Args, yield)
		// TODO: Determine where the ellipsis needs to go
		walkPosVisit(n, &n.Ellipsis, yield)
		walkPosVisit(n, &n.Rparen, yield)

	case *ast.StarExpr:
		walkPosVisit(n, &n.Star, yield)
		walkPos(n.X, yield)

	case *ast.UnaryExpr:
		walkPosVisit(n, &n.OpPos, yield)
		walkPos(n.X, yield)

	case *ast.BinaryExpr:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.OpPos, yield)
		walkPos(n.Y, yield)

	case *ast.KeyValueExpr:
		walkPos(n.Key, yield)
		walkPosVisit(n, &n.Colon, yield)
		walkPos(n.Value, yield)

	// Types
	case *ast.ArrayType:
		walkPosVisit(n, &n.Lbrack, yield)
		// TODO: Determine why no Rbrack?
		walkPos(n.Len, yield)
		walkPos(n.Elt, yield)

	case *ast.StructType:
		walkPosVisit(n, &n.Struct, yield)
		walkPos(n.Fields, yield)

	case *ast.FuncType:
		walkPosVisit(n, &n.Func, yield)
		walkPos(n.TypeParams, yield)
		walkPos(n.Params, yield)
		walkPos(n.Results, yield)

	case *ast.InterfaceType:
		walkPosVisit(n, &n.Interface, yield)
		walkPos(n.Methods, yield)

	case *ast.MapType:
		walkPosVisit(n, &n.Map, yield)
		walkPos(n.Key, yield)
		walkPos(n.Value, yield)

	case *ast.ChanType:
		// TODO: Check the order is correct for the Begin and Arrow
		walkPosVisit(n, &n.Begin, yield)
		walkPosVisit(n, &n.Arrow, yield)
		walkPos(n.Value, yield)

	// Statements
	case *ast.BadStmt:
		walkPosVisit(n, &n.From, yield)
		walkPosVisit(n, &n.To, yield)

	case *ast.DeclStmt:
		walkPos(n.Decl, yield)

	case *ast.EmptyStmt:
		walkPosVisit(n, &n.Semicolon, yield)

	case *ast.LabeledStmt:
		walkPos(n.Label, yield)
		walkPosVisit(n, &n.Colon, yield)
		walkPos(n.Stmt, yield)

	case *ast.ExprStmt:
		walkPos(n.X, yield)

	case *ast.SendStmt:
		walkPos(n.Chan, yield)
		walkPosVisit(n, &n.Arrow, yield)
		walkPos(n.Value, yield)

	case *ast.IncDecStmt:
		walkPos(n.X, yield)
		walkPosVisit(n, &n.TokPos, yield)

	case *ast.AssignStmt:
		walkPosList(n.Lhs, yield)
		walkPosVisit(n, &n.TokPos, yield)
		walkPosList(n.Rhs, yield)

	case *ast.GoStmt:
		walkPosVisit(n, &n.Go, yield)
		walkPos(n.Call, yield)

	case *ast.DeferStmt:
		walkPosVisit(n, &n.Defer, yield)
		walkPos(n.Call, yield)

	case *ast.ReturnStmt:
		walkPosVisit(n, &n.Return, yield)
		walkPosList(n.Results, yield)

	case *ast.BranchStmt:
		walkPosVisit(n, &n.TokPos, yield)
		walkPos(n.Label, yield)

	case *ast.BlockStmt:
		walkPosVisit(n, &n.Lbrace, yield)
		walkPosList(n.List, yield)
		walkPosVisit(n, &n.Rbrace, yield)

	case *ast.IfStmt:
		walkPosVisit(n, &n.If, yield)
		walkPos(n.Init, yield)
		walkPos(n.Cond, yield)
		walkPos(n.Body, yield)
		walkPos(n.Else, yield)

	case *ast.CaseClause:
		walkPosVisit(n, &n.Case, yield)
		walkPosList(n.List, yield)
		walkPosVisit(n, &n.Colon, yield)
		walkPosList(n.Body, yield)

	case *ast.SwitchStmt:
		walkPosVisit(n, &n.Switch, yield)
		walkPos(n.Init, yield)
		walkPos(n.Tag, yield)
		walkPos(n.Body, yield)

	case *ast.TypeSwitchStmt:
		walkPosVisit(n, &n.Switch, yield)
		walkPos(n.Init, yield)
		walkPos(n.Assign, yield)
		walkPos(n.Body, yield)

	case *ast.CommClause:
		walkPosVisit(n, &n.Case, yield)
		walkPos(n.Comm, yield)
		walkPosVisit(n, &n.Colon, yield)
		walkPosList(n.Body, yield)

	case *ast.SelectStmt:
		walkPosVisit(n, &n.Select, yield)
		walkPos(n.Body, yield)

	case *ast.ForStmt:
		walkPosVisit(n, &n.For, yield)
		walkPos(n.Init, yield)
		walkPos(n.Cond, yield)
		walkPos(n.Post, yield)
		walkPos(n.Body, yield)

	case *ast.RangeStmt:
		walkPosVisit(n, &n.For, yield)
		walkPos(n.Key, yield)
		walkPos(n.Value, yield)
		walkPosVisit(n, &n.TokPos, yield)
		walkPosVisit(n, &n.Range, yield)
		walkPos(n.X, yield)
		walkPos(n.Body, yield)

	// Declarations
	case *ast.ImportSpec:
		walkPos(n.Doc, yield)
		walkPos(n.Name, yield)
		walkPos(n.Path, yield)
		walkPos(n.Comment, yield)
		walkPosVisit(n, &n.EndPos, yield)

	case *ast.ValueSpec:
		walkPos(n.Doc, yield)
		walkPosList(n.Names, yield)
		walkPos(n.Type, yield)
		walkPosList(n.Values, yield)
		walkPos(n.Comment, yield)

	case *ast.TypeSpec:
		walkPos(n.Doc, yield)
		walkPos(n.Name, yield)
		walkPos(n.TypeParams, yield)
		walkPosVisit(n, &n.Assign, yield)
		walkPos(n.Type, yield)
		walkPos(n.Comment, yield)

	case *ast.BadDecl:
		walkPosVisit(n, &n.From, yield)
		walkPosVisit(n, &n.To, yield)

	case *ast.GenDecl:
		walkPos(n.Doc, yield)
		walkPosVisit(n, &n.TokPos, yield)
		walkPosVisit(n, &n.Lparen, yield)
		walkPosList(n.Specs, yield)
		walkPosVisit(n, &n.Rparen, yield)

	case *ast.FuncDecl:
		walkPos(n.Doc, yield)
		walkPos(n.Recv, yield)
		// handle FuncType uniquely here to get the name in the correct order.
		walkPosVisit(n, &n.Type.Func, yield)
		walkPos(n.Name, yield)
		walkPos(n.Type.TypeParams, yield)
		walkPos(n.Type.Params, yield)
		walkPos(n.Type.Results, yield)
		walkPos(n.Body, yield)

	// Files
	case *ast.File:
		walkPosVisit(n, &n.FileStart, yield)
		walkPos(n.Doc, yield)
		walkPosVisit(n, &n.Package, yield)
		walkPos(n.Name, yield)
		walkPosList(n.Decls, yield)
		walkPosVisit(n, &n.FileEnd, yield)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]yield`, n))
	}
}
