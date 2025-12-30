package artifacts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"iter"
	"reflect"
)

var errEndWalkPos = errors.New(`end WalkPos`)

// WalkPos will walk the branch of AST nodes with the given node as the root.
// It will return all the positions in the branch.
//
// The position is a pointer to the actual field so that the positions
// can be updated with care. Invalid positions are not returned.
func WalkPos(node ast.Node) iter.Seq[PosTuple] {
	return func(yield func(PosTuple) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndWalkPos {
				panic(r)
			}
		}()

		p := &posWalker{
			yield:   yield,
			handled: map[*token.Pos]struct{}{},
		}
		p.walk(node)
	}
}

type posVisitor func(pt PosTuple) bool

type posWalker struct {
	yield posVisitor

	// zipStack is a stack of positions that need to be interwoven into other
	// positions. These are typically comments. Once a frame of positions is
	// done (popped), all unused positions will be outputted.
	zipStack [][]*PosTuple

	// handled is used for prevent outputting the same value multiple times.
	// This is only used for positions that might be outputted multiple times
	// such as those put into the zipStack.
	handled map[*token.Pos]struct{}
}

type PosTuple struct {
	Node ast.Node
	Pos  *token.Pos
	Name string
}

func (p *posWalker) push(pos []*PosTuple) {
	p.zipStack = append(p.zipStack, pos)
}

func (p *posWalker) pushComments(cg *ast.CommentGroup, name string) {
	frame := []*PosTuple{}
	if cg != nil {
		for _, comment := range cg.List {
			frame = append(frame, &PosTuple{
				Node: comment,
				Pos:  &comment.Slash,
				Name: name,
			})
		}
	}
	p.push(frame)
}

func (p *posWalker) pop() {
	max := len(p.zipStack) - 1
	frame := p.zipStack[max]
	p.zipStack = p.zipStack[:max]
	for _, pt := range frame {
		p.visitTuple(pt)
	}
}

func (p *posWalker) takeTuple(next token.Pos) *PosTuple {
	index := -1
	var min *PosTuple
	for i := len(p.zipStack) - 1; i >= 0; i-- {
		if len(p.zipStack[i]) > 0 {
			top := p.zipStack[i][0]
			if *top.Pos > next {
				continue
			}
			if min == nil || *min.Pos < *top.Pos {
				min = top
				index = i
			}
		}
	}
	if index >= 0 {
		// remove min from stack
		p.zipStack[index] = p.zipStack[index][1:]
	}
	return min
}

func (p *posWalker) beenHandled(pos *token.Pos) bool {
	_, handled := p.handled[pos]
	return handled
}

func (p *posWalker) setAsHandled(pos *token.Pos) {
	p.handled[pos] = struct{}{}
}

func (p *posWalker) visit(n ast.Node, pos *token.Pos, name string) {
	if !pos.IsValid() {
		return
	}
	p.visitTuple(&PosTuple{
		Node: n,
		Pos:  pos,
		Name: name,
	})
}

func (p *posWalker) visitPending(next token.Pos) bool {
	tuple := p.takeTuple(next)
	if tuple == nil {
		return false
	}
	if p.beenHandled(tuple.Pos) {
		return true
	}
	if !p.yield(*tuple) {
		panic(errEndWalkPos)
	}
	p.setAsHandled(tuple.Pos)
	return true
}

func (p *posWalker) visitTuple(pt *PosTuple) {
	// Zip in all pending posTuples that come before `off`.
	for p.visitPending(*pt.Pos) {
		// Do Nothing
	}

	if p.beenHandled(pt.Pos) {
		return
	}
	if !p.yield(*pt) {
		panic(errEndWalkPos)
	}
	if _, ok := pt.Node.(*ast.Comment); ok {
		p.setAsHandled(pt.Pos)
	}
}

func walkPosList[N ast.Node](p *posWalker, list []N) {
	for _, node := range list {
		p.walk(node)
	}
}

func (p *posWalker) walk(node ast.Node) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	switch n := any(node).(type) {

	// ======[ Comments ]======
	case *ast.Comment:
		p.visit(n, &n.Slash, `Comment.Slash`)

	case *ast.CommentGroup:
		walkPosList(p, n.List)

	// ======[ Fields ]======
	case *ast.Field:
		p.pushComments(n.Comment, `Field.Comment`)
		p.walk(n.Doc)
		walkPosList(p, n.Names)
		p.walk(n.Type)
		p.walk(n.Tag)
		p.pop()

	case *ast.FieldList:
		p.visit(n, &n.Opening, `FieldList.Opening`)
		walkPosList(p, n.List)
		p.visit(n, &n.Closing, `FieldList.Closing`)

	// ======[ Expressions ]======
	case *ast.BadExpr:
		p.visit(n, &n.From, `BadExpr.From`)
		p.visit(n, &n.To, `BadExpr.To`)

	case *ast.Ident:
		p.visit(n, &n.NamePos, `Ident.NamePos`)

	case *ast.Ellipsis:
		p.visit(n, &n.Ellipsis, `Ellipsis`)
		p.walk(n.Elt)

	case *ast.BasicLit:
		p.visit(n, &n.ValuePos, `BasicLit.ValuePos`)

	case *ast.FuncLit:
		p.walk(n.Type)
		p.walk(n.Body)

	case *ast.CompositeLit:
		p.walk(n.Type)
		p.visit(n, &n.Lbrace, `CompositeLit.LeftBrace`)
		walkPosList(p, n.Elts)
		p.visit(n, &n.Rbrace, `CompositeLit.RightBrace`)

	case *ast.ParenExpr:
		p.visit(n, &n.Lparen, `ParenExpr.LeftParen`)
		p.walk(n.X)
		p.visit(n, &n.Rparen, `ParenExpr.RightParen`)

	case *ast.SelectorExpr:
		p.walk(n.X)
		p.walk(n.Sel)

	case *ast.IndexExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `IndexExpr.LeftBrace`)
		p.walk(n.Index)
		p.visit(n, &n.Rbrack, `IndexExpr.RightBrace`)

	case *ast.IndexListExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `IndexListExpr.LeftBrace`)
		walkPosList(p, n.Indices)
		p.visit(n, &n.Rbrack, `IndexListExpr.RightBrace`)

	case *ast.SliceExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `SliceExpr.LeftBracket`)
		p.walk(n.Low)
		p.walk(n.High)
		p.walk(n.Max)
		p.visit(n, &n.Rbrack, `SliceExpr.RightBracket`)

	case *ast.TypeAssertExpr:
		p.walk(n.X)
		p.visit(n, &n.Lparen, `TypeAssertExpr.LeftParen`)
		p.walk(n.Type)
		p.visit(n, &n.Rparen, `TypeAssertExpr.RightParen`)

	case *ast.CallExpr:
		p.walk(n.Fun)
		p.visit(n, &n.Lparen, `CallExpr.LeftParen`)
		p.push([]*PosTuple{{
			Node: n,
			Pos:  &n.Ellipsis,
		}})
		walkPosList(p, n.Args)
		p.visit(n, &n.Rparen, `CallExpr.RightParen`)
		p.pop()

	case *ast.StarExpr:
		p.visit(n, &n.Star, `StarExpr.Star`)
		p.walk(n.X)

	case *ast.UnaryExpr:
		p.visit(n, &n.OpPos, `UnaryExpr.`+n.Op.String())
		p.walk(n.X)

	case *ast.BinaryExpr:
		p.walk(n.X)
		p.visit(n, &n.OpPos, `UnaryExpr.`+n.Op.String())
		p.walk(n.Y)

	case *ast.KeyValueExpr:
		p.walk(n.Key)
		p.visit(n, &n.Colon, `KeyValueExpr.Colon`)
		p.walk(n.Value)

	// ======[ Types ]======
	case *ast.ArrayType:
		p.visit(n, &n.Lbrack, `ArrayType.LeftBracket`)
		// TODO: Determine why no Rbrack?
		p.walk(n.Len)
		p.walk(n.Elt)

	case *ast.StructType:
		p.visit(n, &n.Struct, `Struct`)
		p.walk(n.Fields)

	case *ast.FuncType:
		p.visit(n, &n.Func, `Func`)
		p.walk(n.TypeParams)
		p.walk(n.Params)
		p.walk(n.Results)

	case *ast.InterfaceType:
		p.visit(n, &n.Interface, `Interface`)
		p.walk(n.Methods)

	case *ast.MapType:
		p.visit(n, &n.Map, `Map`)
		p.walk(n.Key)
		p.walk(n.Value)

	case *ast.ChanType:
		// TODO: Check the order is correct for the Begin and Arrow
		p.visit(n, &n.Begin, `ChanType.Begin`)
		p.visit(n, &n.Arrow, `ChanType.Arrow`)
		p.walk(n.Value)

	// ======[ Statements ]======
	case *ast.BadStmt:
		p.visit(n, &n.From, `BadStmt.From`)
		p.visit(n, &n.To, `BadStmt.To`)

	case *ast.DeclStmt:
		p.walk(n.Decl)

	case *ast.EmptyStmt:
		p.visit(n, &n.Semicolon, `EmptyStmt.Semicolon`)

	case *ast.LabeledStmt:
		p.walk(n.Label)
		p.visit(n, &n.Colon, `LabeledStmt.Colon`)
		p.walk(n.Stmt)

	case *ast.ExprStmt:
		p.walk(n.X)

	case *ast.SendStmt:
		p.walk(n.Chan)
		p.visit(n, &n.Arrow, `SendStmt.Arrow`)
		p.walk(n.Value)

	case *ast.IncDecStmt:
		p.walk(n.X)
		p.visit(n, &n.TokPos, `IncDecStmt.`+n.Tok.String())

	case *ast.AssignStmt:
		walkPosList(p, n.Lhs)
		p.visit(n, &n.TokPos, `AssignStmt.`+n.Tok.String())
		walkPosList(p, n.Rhs)

	case *ast.GoStmt:
		p.visit(n, &n.Go, `Go`)
		p.walk(n.Call)

	case *ast.DeferStmt:
		p.visit(n, &n.Defer, `Defer`)
		p.walk(n.Call)

	case *ast.ReturnStmt:
		p.visit(n, &n.Return, `Return`)
		walkPosList(p, n.Results)

	case *ast.BranchStmt:
		p.visit(n, &n.TokPos, `Branch.`+n.Tok.String())
		p.walk(n.Label)

	case *ast.BlockStmt:
		p.visit(n, &n.Lbrace, `BlockStmt.LeftBrace`)
		walkPosList(p, n.List)
		p.visit(n, &n.Rbrace, `BlockStmt.RightBrace`)

	case *ast.IfStmt:
		p.visit(n, &n.If, `If`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Body)
		p.walk(n.Else)

	case *ast.CaseClause:
		p.visit(n, &n.Case, `CaseClause.Case`)
		walkPosList(p, n.List)
		p.visit(n, &n.Colon, `CaseClause.Colon`)
		walkPosList(p, n.Body)

	case *ast.SwitchStmt:
		p.visit(n, &n.Switch, `SwitchStmt.Switch`)
		p.walk(n.Init)
		p.walk(n.Tag)
		p.walk(n.Body)

	case *ast.TypeSwitchStmt:
		p.visit(n, &n.Switch, `TypeSwitchStmt.Switch`)
		p.walk(n.Init)
		p.walk(n.Assign)
		p.walk(n.Body)

	case *ast.CommClause:
		p.visit(n, &n.Case, `CommClause.Case`)
		p.walk(n.Comm)
		p.visit(n, &n.Colon, `CommClause.Colon`)
		walkPosList(p, n.Body)

	case *ast.SelectStmt:
		p.visit(n, &n.Select, `Select`)
		p.walk(n.Body)

	case *ast.ForStmt:
		p.visit(n, &n.For, `ForStmt.For`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Post)
		p.walk(n.Body)

	case *ast.RangeStmt:
		p.visit(n, &n.For, `RangeStmt.For`)
		p.walk(n.Key)
		p.walk(n.Value)
		p.visit(n, &n.TokPos, `RangeStmt.`+n.Tok.String())
		p.visit(n, &n.Range, `RangeStmt.Range`)
		p.walk(n.X)
		p.walk(n.Body)

	// ======[ Declarations ]======
	case *ast.ImportSpec:
		p.pushComments(n.Comment, `ImportSpec.Comment`)
		p.walk(n.Doc)
		p.walk(n.Name)
		p.walk(n.Path)
		p.pop()
		p.visit(n, &n.EndPos, `ImportSpec.EndPos`)

	case *ast.ValueSpec:
		p.pushComments(n.Comment, `ValueSpec.Comment`)
		p.walk(n.Doc)
		walkPosList(p, n.Names)
		p.walk(n.Type)
		walkPosList(p, n.Values)
		p.pop()

	case *ast.TypeSpec:
		p.pushComments(n.Comment, `TypeSpec.Comment`)
		p.walk(n.Doc)
		p.walk(n.Name)
		p.walk(n.TypeParams)
		p.visit(n, &n.Assign, `TypeSpec.Assign`)
		p.walk(n.Type)
		p.pop()

	case *ast.BadDecl:
		p.visit(n, &n.From, `BadDecl.From`)
		p.visit(n, &n.To, `BadDecl.To`)

	case *ast.GenDecl:
		p.walk(n.Doc)
		p.visit(n, &n.TokPos, `GenDecl.`+n.Tok.String()+`.Pos`)
		p.visit(n, &n.Lparen, `GenDecl.`+n.Tok.String()+`.LeftParen`)
		walkPosList(p, n.Specs)
		p.visit(n, &n.Rparen, `GenDecl.`+n.Tok.String()+`.RightParen`)

	case *ast.FuncDecl:
		p.walk(n.Doc)
		p.walk(n.Recv)
		// handle FuncType uniquely here to get the name in the correct order.
		p.visit(n, &n.Type.Func, `FuncDecl.Func`)
		p.walk(n.Name)
		p.walk(n.Type.TypeParams)
		p.walk(n.Type.Params)
		p.walk(n.Type.Results)
		p.walk(n.Body)

	// ======[ Files ]======
	case *ast.File:
		p.visit(n, &n.FileStart, `File.Start`)
		for i := len(n.Comments) - 1; i >= 0; i-- {
			p.pushComments(n.Comments[i], `File.Comment`)
		}
		p.walk(n.Doc)
		p.visit(n, &n.Package, `File.Package`)
		p.walk(n.Name)
		walkPosList(p, n.Decls)
		for range n.Comments {
			p.pop()
		}
		p.visit(n, &n.FileEnd, `File.End`)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]yield`, n))
	}
}
