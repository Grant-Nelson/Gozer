package interRep

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/interRep/irc"
)

type converter struct {
	logger   *logger.Logger
	errGroup *faults.ErrGroup
	pkg      *project.Package
	fn       *irc.Func
}

func (cv *converter) pos(p token.Pos) token.Position {
	return cv.pkg.Ast.Fset.Position(p)
}

// TODO: Add code to set Arg values for Goto either during or after.
// TODO: Add code to remove unreachable blocks such as if [afterIf] is never used.
// TODO: Add code to fill out the block links to other blocks for the graph.

// TODO: NOW!!! Make each convert method take and return a block, if the block
// passed in is nil and a statement needs to be added to it, then create one.
// Otherwise if no statements are added, no block is created.

func (cv *converter) statement(block *irc.Block, stmt ast.Stmt) (*irc.Block, error) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return cv.blockStmt(block, stmt)
	case *ast.IfStmt:
		return cv.ifStmt(block, stmt)
	default:
		err := faults.New(`unexpected statement node type`).
			With(`pos`, cv.pos(stmt.Pos())).
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

func (cv *converter) ifStmt(block *irc.Block, stmt *ast.IfStmt) (*irc.Block, error) {
	// Add if-statement initialization to current block since the scoping was already handled by Go.
	if err := cv.statement(block, stmt.Init); err != nil {
		return err
	}

	cond, err := cv.expression(stmt.Cond)
	if err != nil {
		return err
	}

	ircIf := &irc.IfStmt{
		IfPos: stmt.If,
		Cond:  cond,
	}

	thenBlock := cv.fn.NewBlock()
	afterIf := cv.fn.NewBlock()
	if err := cv.blockStmt(thenBlock, stmt.Body); err != nil {
		return err
	}
	cv.addFollowIfNeeded(thenBlock, afterIf)
	ircIf.Then = []irc.Stmt{&irc.GotoStmt{
		Goto: &irc.BlockRef{Block: thenBlock},
	}}

	if stmt.Else != nil {
		elseBlock := cv.fn.NewBlock()
		if err := cv.statement(elseBlock, stmt.Else); err != nil {
			return err
		}
		cv.addFollowIfNeeded(elseBlock, afterIf)
		ircIf.Else = []irc.Stmt{&irc.GotoStmt{
			Goto: &irc.BlockRef{Block: elseBlock},
		}}

	} else {
		// No else so skip to after block.
		ircIf.Else = []irc.Stmt{&irc.GotoStmt{
			Goto: &irc.BlockRef{Block: afterIf},
		}}
	}

	return nil
}

func (cv *converter) addFollowIfNeeded(block, after *irc.Block) {

}

func (cv *converter) expression(expr ast.Expr) (irc.Expr, error) {

	return nil, nil
}
