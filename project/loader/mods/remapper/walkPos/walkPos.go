package walkPos

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"iter"
	"reflect"
	"strings"
)

var errEndWalkPos = errors.New(`end WalkPos`)

// WalkPos will walk the branch of AST nodes with the given node as the root.
// It will return all the positions in the branch.
//
// skipFileComments will skip trying to position floating file level comments.
// If the file is unmodified or has been remapped into the same fileSet space,
// then the floating file level comments can be properly positioned,
// otherwise they may not be positioned in the correct locations.
//
// The position is a pointer to the actual field so that the positions
// can be updated with care. Invalid positions are not returned.
func WalkPos(node ast.Node, skipFileComments bool) iter.Seq[PosTuple] {
	return func(yield func(PosTuple) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndWalkPos {
				panic(r)
			}
		}()

		p := &posWalker{
			yield:            yield,
			skipFileComments: skipFileComments,
			handled:          map[*token.Pos]struct{}{},
		}
		p.walk(node)
	}
}

type PosTuple struct {
	Node ast.Node
	Pos  *token.Pos
	Name string
}

type posVisitor func(pt PosTuple) bool

type posWalker struct {
	yield posVisitor

	skipFileComments bool

	// zipStack is a stack of positions that need to be interwoven into other
	// positions. These are typically comments. Once a frame of positions is
	// done (popped), all unused positions will be outputted.
	zipStack [][]*PosTuple

	// handled is used for prevent outputting the same value multiple times.
	// This is only used for positions that might be outputted multiple times
	// such as those put into the zipStack.
	handled map[*token.Pos]struct{}
}

func (pt PosTuple) String() string {
	typStr := fmt.Sprintf("%T", pt.Node)
	typStr = strings.TrimPrefix(typStr, `*ast.`)
	return fmt.Sprintf("%d:%s:%v", *pt.Pos, typStr, pt.Name)
}

func (p *posWalker) push(frame []*PosTuple) {
	p.zipStack = append(p.zipStack, frame)
}

func (p *posWalker) pushSingle(n ast.Node, pos *token.Pos) {
	var frame []*PosTuple
	if pos.IsValid() {
		pt := &PosTuple{
			Node: n,
			Pos:  pos,
		}
		frame = append(frame, pt)
	}
	p.push(frame)
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
			if (*top.Pos < next) && (min == nil || *min.Pos > *top.Pos) {
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
	pt := p.takeTuple(next)
	if pt == nil {
		return false
	}
	if p.beenHandled(pt.Pos) {
		return true
	}
	if !p.yield(*pt) {
		panic(errEndWalkPos)
	}
	p.setAsHandled(pt.Pos)
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

func (p *posWalker) walkComment(cg *ast.CommentGroup, name string) {
	if cg != nil {
		for _, c := range cg.List {
			p.visit(c, &c.Slash, name)
		}
	}
}

func (p *posWalker) walk(node ast.Node) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	switch n := any(node).(type) {

	// ======[ Comments ]======
	case *ast.Comment:
		p.visit(n, &n.Slash, `Slash`)

	case *ast.CommentGroup:
		walkPosList(p, n.List)

	// ======[ Fields ]======
	case *ast.Field:
		p.pushComments(n.Comment, `Field.Comment`)
		p.walkComment(n.Doc, `Field.Doc`)
		walkPosList(p, n.Names)
		p.walk(n.Type)
		p.walk(n.Tag)
		p.pop()

	case *ast.FieldList:
		p.visit(n, &n.Opening, `Opening`)
		walkPosList(p, n.List)
		p.visit(n, &n.Closing, `Closing`)

	// ======[ Expressions ]======
	case *ast.BadExpr:
		p.visit(n, &n.From, `From`)
		p.visit(n, &n.To, `To`)

	case *ast.Ident:
		p.visit(n, &n.NamePos, n.Name)

	case *ast.Ellipsis:
		p.visit(n, &n.Ellipsis, `Ellipsis`)
		p.walk(n.Elt)

	case *ast.BasicLit:
		p.visit(n, &n.ValuePos, `ValuePos`)

	case *ast.FuncLit:
		p.walk(n.Type)
		p.walk(n.Body)

	case *ast.CompositeLit:
		p.walk(n.Type)
		p.visit(n, &n.Lbrace, `LeftBrace`)
		walkPosList(p, n.Elts)
		p.visit(n, &n.Rbrace, `RightBrace`)

	case *ast.ParenExpr:
		p.visit(n, &n.Lparen, `LeftParen`)
		p.walk(n.X)
		p.visit(n, &n.Rparen, `RightParen`)

	case *ast.SelectorExpr:
		p.walk(n.X)
		p.walk(n.Sel)

	case *ast.IndexExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `LeftBrace`)
		p.walk(n.Index)
		p.visit(n, &n.Rbrack, `RightBrace`)

	case *ast.IndexListExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `LeftBrace`)
		walkPosList(p, n.Indices)
		p.visit(n, &n.Rbrack, `RightBrace`)

	case *ast.SliceExpr:
		p.walk(n.X)
		p.visit(n, &n.Lbrack, `LeftBracket`)
		p.walk(n.Low)
		p.walk(n.High)
		p.walk(n.Max)
		p.visit(n, &n.Rbrack, `RightBracket`)

	case *ast.TypeAssertExpr:
		p.walk(n.X)
		p.visit(n, &n.Lparen, `LeftParen`)
		p.walk(n.Type)
		p.visit(n, &n.Rparen, `RightParen`)

	case *ast.CallExpr:
		p.walk(n.Fun)
		p.visit(n, &n.Lparen, `LeftParen`)
		p.pushSingle(n, &n.Ellipsis)
		walkPosList(p, n.Args)
		p.visit(n, &n.Rparen, `RightParen`)
		p.pop()

	case *ast.StarExpr:
		p.visit(n, &n.Star, `Star`)
		p.walk(n.X)

	case *ast.UnaryExpr:
		p.visit(n, &n.OpPos, n.Op.String())
		p.walk(n.X)

	case *ast.BinaryExpr:
		p.walk(n.X)
		p.visit(n, &n.OpPos, n.Op.String())
		p.walk(n.Y)

	case *ast.KeyValueExpr:
		p.walk(n.Key)
		p.visit(n, &n.Colon, `Colon`)
		p.walk(n.Value)

	// ======[ Types ]======
	case *ast.ArrayType:
		p.visit(n, &n.Lbrack, `LeftBracket`)
		p.walk(n.Len)
		// There is no Rbrack
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
		p.visit(n, &n.Begin, `Begin`)
		p.visit(n, &n.Arrow, `Arrow`)
		p.walk(n.Value)

	// ======[ Statements ]======
	case *ast.BadStmt:
		p.visit(n, &n.From, `From`)
		p.visit(n, &n.To, `To`)

	case *ast.DeclStmt:
		p.walk(n.Decl)

	case *ast.EmptyStmt:
		p.visit(n, &n.Semicolon, `Semicolon`)

	case *ast.LabeledStmt:
		p.walk(n.Label)
		p.visit(n, &n.Colon, `Colon`)
		p.walk(n.Stmt)

	case *ast.ExprStmt:
		p.walk(n.X)

	case *ast.SendStmt:
		p.walk(n.Chan)
		p.visit(n, &n.Arrow, `Arrow`)
		p.walk(n.Value)

	case *ast.IncDecStmt:
		p.walk(n.X)
		p.visit(n, &n.TokPos, n.Tok.String())

	case *ast.AssignStmt:
		walkPosList(p, n.Lhs)
		p.visit(n, &n.TokPos, n.Tok.String())
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
		p.visit(n, &n.TokPos, n.Tok.String())
		p.walk(n.Label)

	case *ast.BlockStmt:
		p.visit(n, &n.Lbrace, `LeftBrace`)
		walkPosList(p, n.List)
		p.visit(n, &n.Rbrace, `RightBrace`)

	case *ast.IfStmt:
		p.visit(n, &n.If, `If`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Body)
		p.walk(n.Else)

	case *ast.CaseClause:
		p.visit(n, &n.Case, `Case`)
		walkPosList(p, n.List)
		p.visit(n, &n.Colon, `Colon`)
		walkPosList(p, n.Body)

	case *ast.SwitchStmt:
		p.visit(n, &n.Switch, `Switch`)
		p.walk(n.Init)
		p.walk(n.Tag)
		p.walk(n.Body)

	case *ast.TypeSwitchStmt:
		p.visit(n, &n.Switch, `Switch`)
		p.walk(n.Init)
		p.walk(n.Assign)
		p.walk(n.Body)

	case *ast.CommClause:
		p.visit(n, &n.Case, `Case`)
		p.walk(n.Comm)
		p.visit(n, &n.Colon, `Colon`)
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
		p.visit(n, &n.For, `For`)
		p.walk(n.Key)
		p.walk(n.Value)
		p.visit(n, &n.TokPos, n.Tok.String())
		p.visit(n, &n.Range, `Range`)
		p.walk(n.X)
		p.walk(n.Body)

	// ======[ Declarations ]======
	case *ast.ImportSpec:
		p.pushComments(n.Comment, `ImportSpec.Comment`)
		p.walkComment(n.Doc, `ImportSpec.Doc`)
		p.walk(n.Name)
		p.walk(n.Path)
		p.pop()
		p.visit(n, &n.EndPos, `EndPos`)

	case *ast.ValueSpec:
		p.pushComments(n.Comment, `ValueSpec.Comment`)
		p.walkComment(n.Doc, `ValueSpec.Doc`)
		walkPosList(p, n.Names)
		p.walk(n.Type)
		walkPosList(p, n.Values)
		p.pop()

	case *ast.TypeSpec:
		p.pushComments(n.Comment, `TypeSpec.Comment`)
		p.walkComment(n.Doc, `TypeSpec.Doc`)
		p.walk(n.Name)
		p.walk(n.TypeParams)
		p.visit(n, &n.Assign, `Assign`)
		p.walk(n.Type)
		p.pop()

	case *ast.BadDecl:
		p.visit(n, &n.From, `From`)
		p.visit(n, &n.To, `To`)

	case *ast.GenDecl:
		p.walkComment(n.Doc, `GenDecl.Doc`)
		p.visit(n, &n.TokPos, n.Tok.String()+`.Pos`)
		p.visit(n, &n.Lparen, n.Tok.String()+`.LeftParen`)
		walkPosList(p, n.Specs)
		p.visit(n, &n.Rparen, n.Tok.String()+`.RightParen`)

	case *ast.FuncDecl:
		p.walkComment(n.Doc, `FuncDecl.Doc`)
		// handle FuncType uniquely here to get the name in the correct order.
		p.visit(n, &n.Type.Func, `Func`)
		p.walk(n.Recv)
		p.walk(n.Name)
		p.walk(n.Type.TypeParams)
		p.walk(n.Type.Params)
		p.walk(n.Type.Results)
		p.walk(n.Body)

	// ======[ Files ]======
	case *ast.File:
		p.visit(n, &n.FileStart, `Start`)
		if !p.skipFileComments {
			for i := len(n.Comments) - 1; i >= 0; i-- {
				p.pushComments(n.Comments[i], `File.Comment`)
			}
		}
		p.walkComment(n.Doc, `File.Doc`)
		p.visit(n, &n.Package, `Package`)
		p.walk(n.Name)
		walkPosList(p, n.Decls)
		if !p.skipFileComments {
			for range n.Comments {
				p.pop()
			}
		}
		p.visit(n, &n.FileEnd, `End`)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]v`, n))
	}
}
