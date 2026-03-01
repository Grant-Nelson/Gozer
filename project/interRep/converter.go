package interRep

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/enums/basicType"
	"github.com/Grant-Nelson/Gozer/project/enums/binaryOp"
	"github.com/Grant-Nelson/Gozer/project/interRep/irc"
)

type converter struct {
	logger   *logger.Logger
	errGroup *faults.ErrGroup
	pkg      *project.Package
	fn       *irc.Func
}

func (cv *converter) pos(n interface{ Pos() token.Pos }) token.Position {
	if n == nil {
		return token.Position{}
	}
	return cv.pkg.Ast.Fset.Position(n.Pos())
}

func (cv *converter) Info() *types.Info {
	return cv.pkg.Ast.TypesInfo
}

// TODO: Add code to set Arg values for Goto either during or after.
// TODO: Add code to remove unreachable blocks such as if [afterIf] is never used.
// TODO: Add code to fill out the block links to other blocks for the graph.

func (cv *converter) statement(block *irc.Block, stmt ast.Stmt) (*irc.Block, error) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return cv.blockStmt(block, stmt)
	case *ast.ExprStmt:
		return cv.exprStmt(block, stmt)
	case *ast.IfStmt:
		return cv.ifStmt(block, stmt)
	case *ast.ReturnStmt:
		return cv.returnStmt(block, stmt)
	default:
		err := faults.New(`unexpected statement node type`).
			With(`pos`, cv.pos(stmt)).
			WithF(`type`, `%T`, stmt)
		return block, cv.errGroup.Add(err)
	}
}

func (cv *converter) blockStmt(block *irc.Block, stmt *ast.BlockStmt) (*irc.Block, error) {
	var err error
	for _, s := range stmt.List {
		if block, err = cv.statement(block, s); err != nil {
			return block, err
		}
	}
	return block, nil
}

func (cv *converter) exprStmt(block *irc.Block, stmt *ast.ExprStmt) (*irc.Block, error) {
	exp, err := cv.expression(stmt.X)
	if err != nil {
		return block, err
	}
	block.Body = append(block.Body, &irc.ExprStmt{Expr: exp})
	return block, nil
}

func (cv *converter) ifStmt(block *irc.Block, stmt *ast.IfStmt) (*irc.Block, error) {
	var err error

	// Add if-statement initialization to current block since the scoping was already handled by Go.
	if stmt.Init != nil {
		if block, err = cv.statement(block, stmt.Init); err != nil {
			return block, err
		}
	}

	cond, err := cv.expression(stmt.Cond)
	if err != nil {
		return block, err
	}

	ircIf := &irc.IfStmt{
		IfPos: stmt.If,
		Cond:  cond,
	}
	block.Body = append(block.Body, ircIf)

	afterIf := cv.fn.NewBlock()

	thenStart := cv.fn.NewBlock()
	thenEnd, err := cv.blockStmt(thenStart, stmt.Body)
	if err != nil {
		return block, err
	}

	cv.addFollowIfNeeded(thenEnd, afterIf)
	ircIf.Then = []irc.Stmt{irc.NewGotoBlock(thenStart)}

	if stmt.Else != nil {
		elseStart := cv.fn.NewBlock()
		elseEnd, err := cv.statement(elseStart, stmt.Else)
		if err != nil {
			return block, err
		}

		cv.addFollowIfNeeded(elseEnd, afterIf)
		ircIf.Else = []irc.Stmt{irc.NewGotoBlock(elseStart)}

	} else {
		// No else so skip to after block.
		block.Body = append(block.Body, irc.NewGotoBlock(afterIf))
	}

	return afterIf, nil
}

func (cv *converter) returnStmt(block *irc.Block, stmt *ast.ReturnStmt) (*irc.Block, error) {
	results := make([]irc.Expr, 0, len(stmt.Results))
	for _, e := range stmt.Results {
		exp, err := cv.expression(e)
		if err != nil {
			return block, err
		}
		results = append(results, exp)
	}

	var result irc.Expr
	if len(results) == 1 {
		result = results[0]
	} else if len(results) > 1 {
		result = &irc.TupleExpr{
			OpenPos: stmt.Return,
			Values:  results,
		}
	}

	ret := &irc.RetStmt{
		KeyPos: stmt.Return,
		Result: result,
	}
	block.Body = append(block.Body, ret)
	return block, nil
}

func (cv *converter) addFollowIfNeeded(block, after *irc.Block) {
	if stmt := block.LastStmt(); stmt != nil && !irc.IsFlowControl(stmt) {
		block.Body = append(block.Body, irc.NewGotoBlock(after))
	}
}

func (cv *converter) expression(expr ast.Expr) (irc.Expr, error) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return cv.basicLit(expr)
	case *ast.BinaryExpr:
		return cv.binaryExpr(expr)
	case *ast.CallExpr:
		return cv.callExpr(expr)
	case *ast.Ident:
		return cv.ident(expr)

	case *ast.UnaryExpr:

	default:
		err := faults.New(`unexpected expression node type`).
			With(`pos`, cv.pos(expr)).
			WithF(`type`, `%T`, expr)
		return nil, cv.errGroup.Add(err)
	}
}

func (cv *converter) basicLit(expr *ast.BasicLit) (irc.Expr, error) {
	typ, err := cv.typeForExpr(expr)
	if err != nil {
		return nil, err
	}

	return &irc.BasicLit{
		ValuePos: expr.ValuePos,
		Value:    expr.Value,
		Type:     typ.(*irc.BasicType),
	}, nil
}

func (cv *converter) binaryExpr(expr *ast.BinaryExpr) (irc.Expr, error) {
	left, err := cv.expression(expr.X)
	if err != nil {
		return nil, err
	}

	right, err := cv.expression(expr.Y)
	if err != nil {
		return nil, err
	}

	typ, err := cv.typeForExpr(expr)
	if err != nil {
		return nil, err
	}

	return &irc.BinaryExpr{
		OpPos:  expr.OpPos,
		Op:     binaryOp.FromToken(expr.Op),
		Left:   left,
		Right:  right,
		Result: typ,
	}, nil
}

func (cv *converter) callExpr(expr *ast.CallExpr) (irc.Expr, error) {
	fn, err := cv.expression(expr.Fun)
	if err != nil {
		return nil, err
	}

	typ, err := cv.typeForExpr(expr)
	if err != nil {
		return nil, err
	}

	args := make([]irc.Expr, len(expr.Args))
	for i, a := range expr.Args {
		arg, err := cv.expression(a)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	return &irc.CallExpr{
		Func:      fn,
		LeftParen: expr.Lparen,
		Args:      args,
		Result:    typ,
	}, nil

}

func (cv *converter) ident(id *ast.Ident) (irc.Expr, error) {
	typ, err := cv.typeForExpr(id)
	if err != nil {
		return nil, err
	}

	return &irc.Ident{
		NamePos: id.NamePos,
		Name:    id.Name,
		Type:    typ,
	}, nil
}

func (cv *converter) unaryExpr(expr *ast.CallExpr) (irc.Expr, error) {

	// TODO: Implement

}

func (cv *converter) typeForExpr(expr ast.Expr) (irc.Type, error) {
	tv, ok := cv.Info().Types[expr]
	if !ok {
		err := faults.New(`unable to find type for expression`).
			With(`pos`, cv.pos(expr)).
			WithF(`type`, `%T`, expr)
		return nil, cv.errGroup.Add(err)
	}
	return cv.convType(tv.Type)
}

func (cv *converter) convType(t types.Type) (irc.Type, error) {
	switch t := t.(type) {
	case *types.Basic:
		return cv.basic(t)
	default:
		err := faults.New(`unexpected type node type`).
			WithF(`type`, `%T`, t)
		return nil, cv.errGroup.Add(err)
	}
}

func (cv *converter) basic(b *types.Basic) (irc.Type, error) {
	return &irc.BasicType{
		Kind: basicType.FromKind(b.Kind()),
	}, nil
}
