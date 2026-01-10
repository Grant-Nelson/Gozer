package artifacts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"iter"
	"reflect"
	"slices"
	"strconv"
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

type posVisitor func(pt PosTuple) bool

type posWalker struct {
	yield posVisitor

	// skipFileComments will skip over the comments from the file.
	// The only comments that will then be skipped will be the floating ones
	// not attached to a node.
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

type PosTuple struct {
	Node  ast.Node
	Pos   *token.Pos
	Width int
	Text  string
	Id    string
}

func (pt *PosTuple) String() string {
	id := pt.Id
	if !strings.Contains(id, `.`) {
		typStr := fmt.Sprintf("%T", pt.Node)
		typStr = strings.TrimPrefix(typStr, `*ast.`)
		id = typStr + `.` + id
	}
	text := ``
	if len(pt.Text) > 0 {
		text = strconv.Quote(pt.Text)
	}
	return fmt.Sprintf("%d:%v:%d%s", *pt.Pos, id, pt.Width, text)
}

func newPosTuple(n ast.Node, pos *token.Pos, text string, width int, id string) *PosTuple {
	if pos.IsValid() {
		return &PosTuple{
			Node:  n,
			Pos:   pos,
			Width: width,
			Text:  text,
			Id:    id,
		}
	}
	return nil
}

func commentTuple(c *ast.Comment, id string) *PosTuple {
	return newPosTuple(c, &c.Slash, c.Text, len(c.Text), id)
}

func tokTuple(n ast.Node, pos *token.Pos, tok token.Token, id string) *PosTuple {
	return newPosTuple(n, pos, tok.String(), len(tok.String()), id)
}

func appendTuple(frame []*PosTuple, pt *PosTuple) []*PosTuple {
	if pt != nil {
		return append(frame, pt)
	}
	return frame
}

func (p *posWalker) push(frame []*PosTuple) {
	p.zipStack = append(p.zipStack, frame)
}

func (p *posWalker) pushTuple(pt *PosTuple) {
	p.push(appendTuple(nil, pt))
}

func (p *posWalker) pushComments(cg *ast.CommentGroup, id string) {
	var frame []*PosTuple
	if cg != nil {
		for _, c := range cg.List {
			frame = appendTuple(frame, commentTuple(c, id))
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

func (p *posWalker) visit(n ast.Node, pos *token.Pos, text string, width int, id string) {
	p.visitTuple(newPosTuple(n, pos, text, width, id))
}

func (p *posWalker) visitTok(n ast.Node, pos *token.Pos, tok token.Token, id string) {
	p.visitTuple(tokTuple(n, pos, tok, id))
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
	if pt == nil {
		return
	}

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

func (p *posWalker) walkComment(cg *ast.CommentGroup, id string) {
	if cg != nil {
		for _, c := range cg.List {
			p.visit(c, &c.Slash, c.Text, len(c.Text), id)
		}
	}
}

func (p *posWalker) walkFieldList(n *ast.FieldList, opening, closing token.Token) {
	if n != nil {
		p.visitTok(n, &n.Opening, opening, `Opening`)
		walkPosList(p, n.List)
		p.visitTok(n, &n.Closing, closing, `Closing`)
	}
}

func (p *posWalker) walk(node ast.Node) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	switch n := any(node).(type) {

	// ======[ Comments ]======
	case *ast.Comment:
		p.visitTuple(commentTuple(n, `X.Comment`))

	case *ast.CommentGroup:
		p.walkComment(n, `X.Comment`)

	// ======[ Fields ]======
	case *ast.Field:
		p.pushComments(n.Comment, `Field.Comment`)
		p.walkComment(n.Doc, `Field.Doc`)
		walkPosList(p, n.Names)
		p.walk(n.Type)
		p.walk(n.Tag)
		p.pop()

	case *ast.FieldList:
		p.walkFieldList(n, token.LPAREN, token.RPAREN)

	// ======[ Expressions ]======
	case *ast.BadExpr:
		p.visit(n, &n.From, ``, int(n.To)-int(n.From), `From`)
		p.visit(n, &n.To, ``, 0, `To`)

	case *ast.Ident:
		p.visit(n, &n.NamePos, n.Name, len(n.Name), `Name`)

	case *ast.Ellipsis:
		p.visitTok(n, &n.Ellipsis, token.ELLIPSIS, `X.Ellipsis`)
		p.walk(n.Elt)

	case *ast.BasicLit:
		p.visit(n, &n.ValuePos, n.Value, len(n.Value), `Value`)

	case *ast.FuncLit:
		p.walk(n.Type)
		p.walk(n.Body)

	case *ast.CompositeLit:
		p.walk(n.Type)
		p.visitTok(n, &n.Lbrace, token.LBRACE, `Lbrace`)
		walkPosList(p, n.Elts)
		p.visitTok(n, &n.Rbrace, token.RBRACE, `Rbrace`)

	case *ast.ParenExpr:
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		p.walk(n.X)
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.SelectorExpr:
		p.walk(n.X)
		p.walk(n.Sel)

	case *ast.IndexExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		p.walk(n.Index)
		p.visitTok(n, &n.Rbrack, token.RBRACK, `Rbrack`)

	case *ast.IndexListExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		walkPosList(p, n.Indices)
		p.visitTok(n, &n.Rbrack, token.RBRACK, `Rbrack`)

	case *ast.SliceExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		p.walk(n.Low)
		p.walk(n.High)
		p.walk(n.Max)
		p.visitTok(n, &n.Rbrack, token.RBRACK, `Rbrack`)

	case *ast.TypeAssertExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		p.walk(n.Type)
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.CallExpr:
		p.walk(n.Fun)
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		p.pushTuple(tokTuple(n, &n.Ellipsis, token.ELLIPSIS, `CallExpr.Ellipsis`))
		walkPosList(p, n.Args)
		p.pop()
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.StarExpr:
		p.visitTok(n, &n.Star, token.MUL, `Star`)
		p.walk(n.X)

	case *ast.UnaryExpr:
		p.visitTok(n, &n.OpPos, n.Op, `Op`)
		p.walk(n.X)

	case *ast.BinaryExpr:
		p.walk(n.X)
		p.visitTok(n, &n.OpPos, n.Op, `Op`)
		p.walk(n.Y)

	case *ast.KeyValueExpr:
		p.walk(n.Key)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		p.walk(n.Value)

	// ======[ Types ]======
	case *ast.ArrayType:
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		p.walk(n.Len)
		// There is no Rbrack
		p.walk(n.Elt)

	case *ast.StructType:
		p.visitTok(n, &n.Struct, token.STRUCT, `Struct`)
		p.walkFieldList(n.Fields, token.LBRACE, token.RBRACE)

	case *ast.FuncType:
		p.visitTok(n, &n.Func, token.FUNC, `Func`)
		p.walkFieldList(n.TypeParams, token.LBRACK, token.RBRACK)
		p.walkFieldList(n.Params, token.LPAREN, token.RPAREN)
		p.walkFieldList(n.Results, token.LPAREN, token.RPAREN)

	case *ast.InterfaceType:
		p.visitTok(n, &n.Interface, token.INTERFACE, `Interface`)
		p.walkFieldList(n.Methods, token.LBRACE, token.RBRACE)

	case *ast.MapType:
		p.visitTok(n, &n.Map, token.MAP, `Map`)
		p.walk(n.Key)
		p.walk(n.Value)

	case *ast.ChanType:
		if n.Begin == n.Arrow {
			p.visit(n, &n.Arrow, ``, 0, `ArrowChan.Arrow`)
			text := token.ARROW.String() + token.CHAN.String()
			p.visit(n, &n.Begin, text, len(text), `ArrowChan.Chan`)
		} else if n.Arrow.IsValid() {
			p.visitTok(n, &n.Begin, token.CHAN, `ChanArrow.Chan`)
			p.visitTok(n, &n.Arrow, token.ARROW, `ChanArrow.Arrow`)
		} else {
			p.visitTok(n, &n.Begin, token.CHAN, `Chan.Chan`)
		}
		p.walk(n.Value)

	// ======[ Statements ]======
	case *ast.BadStmt:
		p.visit(n, &n.From, ``, int(n.To)-int(n.From), `From`)
		p.visit(n, &n.To, ``, 0, `To`)

	case *ast.DeclStmt:
		p.walk(n.Decl)

	case *ast.EmptyStmt:
		p.visitTok(n, &n.Semicolon, token.SEMICOLON, `Semicolon`)

	case *ast.LabeledStmt:
		p.walk(n.Label)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		p.walk(n.Stmt)

	case *ast.ExprStmt:
		p.walk(n.X)

	case *ast.SendStmt:
		p.walk(n.Chan)
		p.visitTok(n, &n.Arrow, token.ARROW, `Arrow`)
		p.walk(n.Value)

	case *ast.IncDecStmt:
		p.walk(n.X)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)

	case *ast.AssignStmt:
		walkPosList(p, n.Lhs)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		walkPosList(p, n.Rhs)

	case *ast.GoStmt:
		p.visitTok(n, &n.Go, token.GO, `Go`)
		p.walk(n.Call)

	case *ast.DeferStmt:
		p.visitTok(n, &n.Defer, token.DEFER, `Defer`)
		p.walk(n.Call)

	case *ast.ReturnStmt:
		p.visitTok(n, &n.Return, token.RETURN, `Return`)
		walkPosList(p, n.Results)

	case *ast.BranchStmt:
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.walk(n.Label)

	case *ast.BlockStmt:
		p.visitTok(n, &n.Lbrace, token.LBRACE, `Lbrace`)
		walkPosList(p, n.List)
		p.visitTok(n, &n.Rbrace, token.RBRACE, `Rbrace`)

	case *ast.IfStmt:
		p.visitTok(n, &n.If, token.IF, `If`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Body)
		p.walk(n.Else)

	case *ast.CaseClause:
		p.visitTok(n, &n.Case, token.CASE, `Case`)
		walkPosList(p, n.List)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		walkPosList(p, n.Body)

	case *ast.SwitchStmt:
		p.visitTok(n, &n.Switch, token.SWITCH, `Switch`)
		p.walk(n.Init)
		p.walk(n.Tag)
		p.walk(n.Body)

	case *ast.TypeSwitchStmt:
		p.visitTok(n, &n.Switch, token.SWITCH, `Switch`)
		p.walk(n.Init)
		p.walk(n.Assign)
		p.walk(n.Body)

	case *ast.CommClause:
		p.visitTok(n, &n.Case, token.CASE, `Case`)
		p.walk(n.Comm)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		walkPosList(p, n.Body)

	case *ast.SelectStmt:
		p.visitTok(n, &n.Select, token.SELECT, `Select`)
		p.walk(n.Body)

	case *ast.ForStmt:
		p.visitTok(n, &n.For, token.FOR, `For`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Post)
		p.walk(n.Body)

	case *ast.RangeStmt:
		p.visitTok(n, &n.For, token.FOR, `For`)
		p.walk(n.Key)
		p.walk(n.Value)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.visitTok(n, &n.Range, token.RANGE, `Range`)
		p.walk(n.X)
		p.walk(n.Body)

	// ======[ Declarations ]======
	case *ast.ImportSpec:
		p.pushComments(n.Comment, `ImportSpec.Comment`)
		p.walkComment(n.Doc, `ImportSpec.Doc`)
		p.walk(n.Name)
		p.walk(n.Path)
		p.pop()
		p.visit(n, &n.EndPos, ``, 0, `EndPos`)

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
		p.walkFieldList(n.TypeParams, token.LBRACK, token.RBRACK)
		p.visitTok(n, &n.Assign, token.EQL, `Assign`)
		p.walk(n.Type)
		p.pop()

	case *ast.BadDecl:
		p.visit(n, &n.From, ``, int(n.To)-int(n.From), `From`)
		p.visit(n, &n.To, ``, 0, `To`)

	case *ast.GenDecl:
		p.walkComment(n.Doc, `GenDecl.Doc`)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		walkPosList(p, n.Specs)
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.FuncDecl:
		p.walkComment(n.Doc, `FuncDecl.Doc`)
		// handle FuncType uniquely here to get the name in the correct order.
		p.visitTok(n.Type, &n.Type.Func, token.FUNC, `Func`)
		p.walkFieldList(n.Recv, token.LPAREN, token.RPAREN)
		p.walk(n.Name)
		p.walkFieldList(n.Type.TypeParams, token.LBRACK, token.RBRACK)
		p.walkFieldList(n.Type.Params, token.LPAREN, token.RPAREN)
		p.walkFieldList(n.Type.Results, token.LPAREN, token.RPAREN)
		p.walk(n.Body)

	// ======[ Files ]======
	case *ast.File:
		p.visit(n, &n.FileStart, ``, 0, `Start`)
		if !p.skipFileComments {
			for _, c := range slices.Backward(n.Comments) {
				p.pushComments(c, `File.Comment`)
			}
		}
		p.walkComment(n.Doc, `File.Doc`)
		p.visitTok(n, &n.Package, token.PACKAGE, `Package`)
		p.walk(n.Name)
		walkPosList(p, n.Decls)
		if !p.skipFileComments {
			for range n.Comments {
				p.pop()
			}
		}
		p.visit(n, &n.FileEnd, ``, 0, `End`)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]v`, n))
	}
}
