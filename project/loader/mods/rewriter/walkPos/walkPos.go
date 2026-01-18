package walkPos

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"iter"
	"reflect"
	"slices"
)

var errEndWalkPos = errors.New(`end WalkPos`)

type WalkPosOption string

const (
	SkipFileComments WalkPosOption = `SkipFileComments`
	SkipPseudoPos    WalkPosOption = `SkipPseudoPos`
)

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
func WalkPos(fs *token.File, node ast.Node, options ...WalkPosOption) iter.Seq[PosTuple] {
	return func(yield func(PosTuple) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndWalkPos {
				panic(r)
			}
		}()
		op := pwOptions{
			skipFileComments: slices.Contains(options, SkipFileComments),
			skipPseudoPos:    slices.Contains(options, SkipPseudoPos),
		}
		p := &posWalker{
			fs:        fs,
			yield:     yield,
			handled:   map[*token.Pos]struct{}{},
			pwOptions: op,
		}
		p.walk(node)
	}
}

type posVisitor func(pt PosTuple) bool

type pwOptions struct {
	// skipFileComments will skip over the comments from the file.
	// The only comments that will then be skipped will be the floating ones
	// not attached to a node.
	skipFileComments bool

	// skipPseudoPos will skip over pseudo positions for glyphs that are not
	// given a specific actual location in the AST but will appear in the source.
	skipPseudoPos bool
}

type posWalker struct {
	pwOptions
	fs    *token.File
	yield posVisitor

	// lastEndPos stores the next unused position after the currently yielded
	// position information.
	lastEndPos token.Pos

	// cond, if not nil, is a conditional pseudo position that may or
	// may not be added to the output before the next position.
	cond *conditionalPseudo

	// zipStack is a stack of positions that need to be interwoven into other
	// positions. These are typically comments. Once a frame of positions is
	// done (popped), all unused positions will be outputted.
	zipStack [][]*PosTuple

	// handled is used for prevent outputting the same value multiple times.
	// This is only used for positions that might be outputted multiple times
	// such as those put into the zipStack.
	handled map[*token.Pos]struct{}
}

type conditionalPseudo struct {
	// newLine indicates that:
	// - if true, the conditionalPseudo needs to be added
	//   iff a newLine is between the prior position and the next position
	// - if false, the conditional Pseudo needs to be added
	//   iff a newline is NOT between the prior position and the next position
	newLine bool

	prior token.Pos
	n     ast.Node
	tok   token.Token
	id    string
}

func isNodeNil(node ast.Node) bool {
	return node == nil || reflect.ValueOf(node).IsNil()
}

func tokWidth(tok token.Token) int {
	return len(tok.String())
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

// takeTuple will get all the tokens before but not equal to the given
// [next] position. The next position could be a pending tuple but with more
// information so will go in place of the pending tuple.
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

func (p *posWalker) setConditionalPseudo(newLine bool, prior token.Pos, n ast.Node, tok token.Token, id string) {
	if p.cond != nil {
		panic(fmt.Errorf(`can not set a conditional when one is already set`))
	}
	p.cond = &conditionalPseudo{
		newLine: newLine,
		prior:   prior,
		n:       n,
		tok:     tok,
		id:      id,
	}
}

func (p *posWalker) visitConditionalPseudo(next token.Pos) {
	if p.cond == nil {
		return
	}

	hasNL := p.fs.Line(p.cond.prior) != p.fs.Line(next)
	if hasNL != p.cond.newLine {
		p.cond = nil
		return
	}
	p.visitPseudo(p.cond.n, p.cond.tok, p.cond.id)
	p.cond = nil
}

func (p *posWalker) visit(n ast.Node, pos *token.Pos, text string, width int, id string) {
	p.visitTuple(newPosTuple(n, pos, text, width, id))
}

func (p *posWalker) visitTok(n ast.Node, pos *token.Pos, tok token.Token, id string) {
	p.visitTuple(tokTuple(n, pos, tok, id))
}

func (p *posWalker) visitPseudo(n ast.Node, tok token.Token, id string) {
	if p.skipPseudoPos {
		return
	}

	// If any pending is at the lastEndPos output them and
	// bump the pseudo until after those pending.
	width := token.Pos(tokWidth(tok))
	for p.visitAllPending(p.lastEndPos + width) {
		// Do Nothing
	}

	// make a copy of the pos so that it can have a pointer for it
	// without any concern of the lastEndPos of being modified.
	pos := p.lastEndPos
	pt := tokTuple(n, &pos, tok, id)
	pt.Pseudo = true
	p.visitTuple(pt)
}

func (p *posWalker) visitAllPending(next token.Pos) bool {
	hadAny := false
	for p.visitPending(next) {
		hadAny = true
	}
	return hadAny
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
	p.lastEndPos = pt.End()
	p.setAsHandled(pt.Pos)
	return true
}

func (p *posWalker) visitTuple(pt *PosTuple) {
	if pt == nil {
		return
	}

	p.visitConditionalPseudo(*pt.Pos)

	// Zip in all pending posTuples that come before [pt].
	p.visitAllPending(*pt.Pos)

	if p.beenHandled(pt.Pos) {
		return
	}
	if !p.yield(*pt) {
		panic(errEndWalkPos)
	}
	p.lastEndPos = pt.End()
	if _, ok := pt.Node.(*ast.Comment); ok {
		p.setAsHandled(pt.Pos)
	}
}

func walkSemicolonSepList[N ast.Node](p *posWalker, parent ast.Node, list []N, semicolonId string) {
	for i, node := range list {
		if i > 0 {
			p.setConditionalPseudo(false, p.lastEndPos, parent, token.SEMICOLON, semicolonId)
		}
		p.walk(node)
	}
}

func walkCommaSepList[N ast.Node](p *posWalker, parent ast.Node, list []N, commaId string, enclosed bool) {
	for i, node := range list {
		if i > 0 {
			p.visitPseudo(parent, token.COMMA, commaId)
		}
		p.walk(node)
	}
	if enclosed {
		p.setConditionalPseudo(true, p.lastEndPos, parent, token.COMMA, commaId)
	}
}

func (p *posWalker) walkComments(cg *ast.CommentGroup, id string) {
	if cg != nil {
		for _, c := range cg.List {
			p.visit(c, &c.Slash, c.Text, len(c.Text), id)
		}
	}
}

func (p *posWalker) walkFieldList(n *ast.FieldList) {
	if n != nil {
		p.visitTok(n, &n.Opening, token.LBRACE, `LBrace`)
		walkSemicolonSepList(p, n, n.List, `FieldSep`)
		p.visitTok(n, &n.Closing, token.RBRACE, `RBrace`)
	}
}

func (p *posWalker) walkTypeParams(n *ast.FieldList) {
	if n != nil {
		p.visitTok(n, &n.Opening, token.LBRACK, `LBrack`)
		walkCommaSepList(p, n, n.List, `Comma`, true)
		p.visitTok(n, &n.Closing, token.RBRACK, `RBrack`)
	}
}

func (p *posWalker) walkParams(n *ast.FieldList) {
	if n != nil {
		p.visitTok(n, &n.Opening, token.LPAREN, `LParen`)
		walkCommaSepList(p, n, n.List, `Comma`, true)
		p.visitTok(n, &n.Closing, token.RPAREN, `RParen`)
	}
}

func (p *posWalker) walk(node ast.Node) {
	if isNodeNil(node) {
		return
	}

	switch n := any(node).(type) {

	// ======[ Comments ]======
	case *ast.Comment:
		p.visitTuple(commentTuple(n, `X.Comment`))

	case *ast.CommentGroup:
		p.walkComments(n, `X.Comment`)

	// ======[ Fields ]======
	case *ast.Field:
		p.pushComments(n.Comment, `Field.Comment`)
		p.walkComments(n.Doc, `Field.Doc`)
		walkCommaSepList(p, n, n.Names, `NameComma`, false)
		p.walk(n.Type)
		p.walk(n.Tag)
		p.pop()

	case *ast.FieldList:
		p.walkParams(n) // Without context, just treat like parameters

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
		walkSemicolonSepList(p, n, n.Elts, `EltSep`)
		p.visitTok(n, &n.Rbrace, token.RBRACE, `Rbrace`)

	case *ast.ParenExpr:
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		p.walk(n.X)
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.SelectorExpr:
		p.walk(n.X)
		p.visitPseudo(n, token.PERIOD, `Dot`)
		p.walk(n.Sel)

	case *ast.IndexExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		p.walk(n.Index)
		p.visitTok(n, &n.Rbrack, token.RBRACK, `Rbrack`)

	case *ast.IndexListExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		walkCommaSepList(p, n, n.Indices, `Comma`, true)
		p.visitTok(n, &n.Rbrack, token.RBRACK, `Rbrack`)

	case *ast.SliceExpr:
		p.walk(n.X)
		p.visitTok(n, &n.Lbrack, token.LBRACK, `Lbrack`)
		p.walk(n.Low)
		p.visitPseudo(n, token.COLON, `FirstColon`)
		p.walk(n.High)
		if !isNodeNil(n.Max) {
			p.visitPseudo(n, token.COLON, `SecondColon`)
			p.walk(n.Max)
		}
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
		walkCommaSepList(p, n, n.Args, `ArgComma`, true)
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
		p.visitPseudo(n, token.RBRACK, `Rbrack`)
		p.walk(n.Elt)

	case *ast.StructType:
		p.visitTok(n, &n.Struct, token.STRUCT, `Struct`)
		p.walkFieldList(n.Fields)

	case *ast.FuncType:
		p.visitTok(n, &n.Func, token.FUNC, `Func`)
		p.walkTypeParams(n.TypeParams)
		p.walkParams(n.Params)
		p.walkParams(n.Results)

	case *ast.InterfaceType:
		p.visitTok(n, &n.Interface, token.INTERFACE, `Interface`)
		p.walkFieldList(n.Methods)

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
		walkCommaSepList(p, n, n.Lhs, `LhsComma`, false)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		walkCommaSepList(p, n, n.Rhs, `RhsComma`, false)

	case *ast.GoStmt:
		p.visitTok(n, &n.Go, token.GO, `Go`)
		p.walk(n.Call)

	case *ast.DeferStmt:
		p.visitTok(n, &n.Defer, token.DEFER, `Defer`)
		p.walk(n.Call)

	case *ast.ReturnStmt:
		p.visitTok(n, &n.Return, token.RETURN, `Return`)
		walkCommaSepList(p, n, n.Results, `Comma`, false)

	case *ast.BranchStmt:
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.walk(n.Label)

	case *ast.BlockStmt:
		p.visitTok(n, &n.Lbrace, token.LBRACE, `Lbrace`)
		walkSemicolonSepList(p, n, n.List, `Sep`)
		p.visitTok(n, &n.Rbrace, token.RBRACE, `Rbrace`)

	case *ast.IfStmt:
		p.visitTok(n, &n.If, token.IF, `If`)
		p.walk(n.Init)
		p.walk(n.Cond)
		p.walk(n.Body)
		p.walk(n.Else)

	case *ast.CaseClause:
		p.visitTok(n, &n.Case, token.CASE, `Case`)
		walkCommaSepList(p, n, n.List, `CaseComma`, false)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		walkSemicolonSepList(p, n, n.Body, `BodySep`)

	case *ast.SwitchStmt:
		p.visitTok(n, &n.Switch, token.SWITCH, `Switch`)
		if !isNodeNil(n.Init) {
			p.walk(n.Init)
			p.visitPseudo(n, token.SEMICOLON, `InitComma`)
		}
		p.walk(n.Tag)
		p.walk(n.Body)

	case *ast.TypeSwitchStmt:
		p.visitTok(n, &n.Switch, token.SWITCH, `Switch`)
		if !isNodeNil(n.Init) {
			p.walk(n.Init)
			p.visitPseudo(n, token.SEMICOLON, `InitComma`)
		}
		p.walk(n.Assign)
		p.walk(n.Body)

	case *ast.CommClause:
		p.visitTok(n, &n.Case, token.CASE, `Case`)
		p.walk(n.Comm)
		p.visitTok(n, &n.Colon, token.COLON, `Colon`)
		walkSemicolonSepList(p, n, n.Body, `BodySep`)

	case *ast.SelectStmt:
		p.visitTok(n, &n.Select, token.SELECT, `Select`)
		p.walk(n.Body)

	case *ast.ForStmt:
		p.visitTok(n, &n.For, token.FOR, `For`)
		if !isNodeNil(n.Init) {
			p.walk(n.Init)
			p.visitPseudo(n, token.SEMICOLON, `InitSemicolon`)
		}
		p.walk(n.Cond)
		if !isNodeNil(n.Post) {
			p.visitPseudo(n, token.SEMICOLON, `PostSemicolon`)
			p.walk(n.Post)
		}
		p.walk(n.Body)

	case *ast.RangeStmt:
		p.visitTok(n, &n.For, token.FOR, `For`)
		p.walk(n.Key)
		if !isNodeNil(n.Value) {
			p.visitPseudo(n, token.COMMA, `Comma`)
			p.walk(n.Value)
		}
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.visitTok(n, &n.Range, token.RANGE, `Range`)
		p.walk(n.X)
		p.walk(n.Body)

	// ======[ Declarations ]======
	case *ast.ImportSpec:
		p.pushComments(n.Comment, `ImportSpec.Comment`)
		p.walkComments(n.Doc, `ImportSpec.Doc`)
		p.walk(n.Name)
		p.walk(n.Path)
		p.pop()
		p.visit(n, &n.EndPos, ``, 0, `EndPos`)

	case *ast.ValueSpec:
		p.pushComments(n.Comment, `ValueSpec.Comment`)
		p.walkComments(n.Doc, `ValueSpec.Doc`)
		walkCommaSepList(p, n, n.Names, `NameComma`, false)
		p.walk(n.Type)
		walkCommaSepList(p, n, n.Values, `ValueComma`, false)
		p.pop()

	case *ast.TypeSpec:
		p.pushComments(n.Comment, `TypeSpec.Comment`)
		p.walkComments(n.Doc, `TypeSpec.Doc`)
		p.walk(n.Name)
		p.walkTypeParams(n.TypeParams)
		p.visitTok(n, &n.Assign, token.EQL, `Assign`)
		p.walk(n.Type)
		p.pop()

	case *ast.BadDecl:
		p.visit(n, &n.From, ``, int(n.To)-int(n.From), `From`)
		p.visit(n, &n.To, ``, 0, `To`)

	case *ast.GenDecl:
		p.walkComments(n.Doc, `GenDecl.Doc`)
		p.visitTok(n, &n.TokPos, n.Tok, `Tok`)
		p.visitTok(n, &n.Lparen, token.LPAREN, `Lparen`)
		walkSemicolonSepList(p, n, n.Specs, `SpecSep`)
		p.visitTok(n, &n.Rparen, token.RPAREN, `Rparen`)

	case *ast.FuncDecl:
		p.walkComments(n.Doc, `FuncDecl.Doc`)
		// handle FuncType uniquely here to get the name in the correct order.
		p.visitTok(n.Type, &n.Type.Func, token.FUNC, `Func`)
		p.walkParams(n.Recv)
		p.walk(n.Name)
		p.walkTypeParams(n.Type.TypeParams)
		p.walkParams(n.Type.Params)
		p.walkParams(n.Type.Results)
		p.walk(n.Body)

	// ======[ Files ]======
	case *ast.File:
		p.visit(n, &n.FileStart, ``, 0, `Start`)
		if !p.skipFileComments {
			for _, c := range slices.Backward(n.Comments) {
				p.pushComments(c, `File.Comment`)
			}
		}
		p.walkComments(n.Doc, `File.Doc`)
		p.visitTok(n, &n.Package, token.PACKAGE, `Package`)
		p.walk(n.Name)
		walkSemicolonSepList(p, n, n.Decls, `DeclSep`)
		if !p.skipFileComments {
			for range n.Comments {
				p.pop()
			}
		}
		p.visit(n, &n.FileEnd, ``, 0, `End`)

	default:
		panic(fmt.Errorf(`unexpected node: (%[1]T) %[1]v`, n))
	}
}
