package blocker

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/Grant-Nelson/Gozer/avail/crumb"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/ir"
	"github.com/Grant-Nelson/Gozer/project/modeler/remodel"
)

type Config struct {

	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

func New(cfg *Config) remodel.RemodelFactory {
	return &blocker{errGroup: cfg.ErrGroup}
}

type blocker struct {
	errGroup *faults.ErrGroup
}

func (b *blocker) StartPackage(pkg *project.Package) (bool, remodel.Remodeler, error) {
	bb := &blockBuilder{
		errGroup: b.errGroup,
		pkg:      pkg,
	}
	return true, bb, nil
}

type blockBuilder struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
}

func (bb *blockBuilder) PackageDone() (bool, error) { return true, nil }

type funcBlockBuilder struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
	fn       *ir.Func
	curBlock *ir.Block

	forRangeItType types.Type

	stmtIndex   int
	curStmtList []ir.Stmt

	labelBlock    map[token.Pos]*ir.Block
	labelToStmt   map[token.Pos]token.Pos
	stmtLabelName map[token.Pos]string
	continueBlock map[token.Pos]*ir.Block
	breakBlock    map[token.Pos]*ir.Block
	blockPos      map[*ir.Block]token.Pos
}

func (bb *blockBuilder) RemodelFunc(fn *ir.Func) (con bool, err error) {
	bb.errGroup.Recover(&err)
	if fn.Atomic() {
		return true, nil
	}

	fbb := &funcBlockBuilder{
		errGroup: bb.errGroup,
		pkg:      bb.pkg,
		fn:       fn,

		labelBlock:    map[token.Pos]*ir.Block{},
		labelToStmt:   map[token.Pos]token.Pos{},
		stmtLabelName: map[token.Pos]string{},
		continueBlock: map[token.Pos]*ir.Block{},
		breakBlock:    map[token.Pos]*ir.Block{},
		blockPos:      map[*ir.Block]token.Pos{},
	}

	if len(fn.Blocks) <= 0 {
		return true, bb.errGroup.Add(faults.New(`function has no blocks`).
			With(`function`, fn.Name))
	}
	if len(fn.Blocks) > 1 {
		return true, bb.errGroup.Add(faults.New(`function already has multiple blocks`).
			With(`count`, len(fn.Blocks)).
			With(`function`, fn.Name))
	}

	// TODO: FIX by moving to end after blocks have been broken out.
	//if !ir.IsFlowControlStatement(fn.Blocks[0].LastStmt()) {
	//	fn.Blocks[0].Body = append(fn.Blocks[0].Body, &ir.ReturnStmt{})
	//}

	for blockIndex := 0; blockIndex < len(fn.Blocks); blockIndex++ {
		fbb.remodelBlock(fn.Blocks[blockIndex])
	}
	return true, bb.errGroup.FullOrNil()
}

func (fbb *funcBlockBuilder) pos(p token.Pos) token.Position {
	return fbb.pkg.Position(p)
}

func (fbb *funcBlockBuilder) info() *types.Info {
	return fbb.pkg.Ast.TypesInfo
}

func (fbb *funcBlockBuilder) remodelBlock(b *ir.Block) {
	fbb.curBlock = b
	fbb.curStmtList = fbb.curBlock.Body
	defer func() {
		fbb.curBlock.Body = fbb.curStmtList
		fbb.curStmtList = nil
		fbb.curBlock = nil
	}()
	for fbb.stmtIndex = 0; fbb.stmtIndex < len(fbb.curStmtList); fbb.stmtIndex++ {
		fbb.remodelStmt(fbb.curStmtList[fbb.stmtIndex])
	}
}

func (fbb *funcBlockBuilder) remodelStmtSlice(ss []ir.Stmt) {
	priorStmtIndex := fbb.stmtIndex
	priorStmtList := fbb.curStmtList
	defer func() {
		fbb.stmtIndex = priorStmtIndex
		fbb.curStmtList = priorStmtList
	}()
	fbb.curStmtList = ss
	for fbb.stmtIndex = 0; fbb.stmtIndex < len(fbb.curStmtList); fbb.stmtIndex++ {
		fbb.remodelStmt(fbb.curStmtList[fbb.stmtIndex])
	}
}

// remodelStmt remodels a statement to extract statements that cause
// flow-control changes so that blocks can be built for a scheduler.
// See: https://go.dev/ref/spec#Terminating_statements
func (fbb *funcBlockBuilder) remodelStmt(s ir.Stmt) {
	switch s := s.(type) {
	case *ir.DeclStmt, *ir.GotoBlockStmt:
		// Do Nothing
	case *ir.LabeledStmt:
		fbb.remodelLabeledStmt(s)
	case *ir.ForStmt:
		fbb.remodelForStmt(s)
	case *ir.RangeStmt:
		fbb.remodelRangeStmt(s)
	case *ir.AssignStmt:
		fbb.remodelAssignStmt(s)
	case *ir.ReturnStmt:
		fbb.remodelReturnStmt(s)
	case *ir.BranchStmt:
		fbb.remodelBranchStmt(s)
	case *ir.IfStmt:
		fbb.remodelIfStmt(s)
	case *ir.ExprStmt:
		fbb.remodelExprStmt(s)
	default:
		fbb.errGroup.Add(faults.New(`unhandled statement node in blocker`).
			With(`pos`, fbb.pos(s.Pos())).
			WithF(`type`, `%T`, s))
		return
	}
}

func (fbb *funcBlockBuilder) remodelAssignStmt(s *ir.AssignStmt) {
	fbb.remodelExprSlice(s, nil, s.Lhs)
	fbb.remodelExprSlice(s, nil, s.Rhs)
}

// splitCurBlock will break the current block into two parts.
// The current block is cut so that the current statement and everything after
// is removed. All the statements after the current statement are added to
// the next block. The current statement is not in either block but returned.
// A goto is added from the current block that jumps unconditionally to the
// next block. The goto is also returned so that it can be redirected.
// Lastly the current statement index is backed up since the current
// statement has been removed.
func (fbb *funcBlockBuilder) splitCurBlock(nextBlk *ir.Block) (ir.Stmt, *ir.GotoBlockStmt) {
	stmt := fbb.curStmtList[fbb.stmtIndex]
	nextBlk.Body = append(nextBlk.Body, fbb.curStmtList[fbb.stmtIndex+1:]...)
	fbb.curStmtList = slices.Clone(fbb.curStmtList[:fbb.stmtIndex])
	gotoLabel := ir.NewGotoBlockStmt(stmt.Pos(), nextBlk)
	fbb.curStmtList = append(fbb.curStmtList, gotoLabel)
	fbb.stmtIndex--
	return stmt, gotoLabel
}

// remodelLabeledStmt processes a label statement.
// A label can be jumped to so the code reachable from the label
// needs to be put into it's own block.
//
//	+--[Cur]-----------+     +--[Cur]-----------+
//	|    ...           |	 |     ...          |
//	|   stmt k-1       |	 | > stmt k-1       |
//	| > stmt k (label) | ==> |   goto Next      |
//	|   stmt k+1       |     +------------------+
//	|    ...           |     +--[Next]----------+
//	+------------------+     |   stmt k+1       |
//	                         |   ...            |
//	                         +------------------+
//
// See: https://go.dev/ref/spec#Labeled_statements
func (fbb *funcBlockBuilder) remodelLabeledStmt(s *ir.LabeledStmt) {
	// Check if a block was preemptively created by prior code
	// that is jumping forward to this block.
	nextBlk, ok := fbb.labelBlock[s.Label.Pos()]
	if ok {
		if len(nextBlk.Body) > 0 {
			fbb.errGroup.Add(faults.New(`preemptive label block is already populated`).
				With(`pos`, fbb.pos(s.Pos())).
				With(`statements`, len(nextBlk.Body)).
				With(`label block`, nextBlk).
				With(`current block`, fbb.curBlock).
				With(`label`, s.Label.String()))
			return
		}
	} else {
		// Create a new block for the code reachable from the label.
		nextBlk = fbb.fn.NewBlock(`Label ` + s.Label.String())
		fbb.labelBlock[s.Label.Pos()] = nextBlk
	}

	// Put the statement that the label is attached to, if there is one, as the
	// first statement in the block, then move all following statements from
	// the current block into the next block.
	if s.Stmt != nil {
		nextBlk.Body = append(nextBlk.Body, s.Stmt)
		fbb.labelToStmt[s.Label.Pos()] = s.Stmt.Pos()
		fbb.stmtLabelName[s.Stmt.Pos()] = s.Label.String()
	}
	fbb.splitCurBlock(nextBlk)
}

// remodelForStmt remodels a for-loop (without a range) into blocks.
//
//	+--[Cur]-----------------+     +--[Cur]-----------------+
//	|    ...                 |     |     ...                |
//	|   stmt k-1             |     | > stmt k-1             |
//	| > stmt k (for)         | ==> |   for-init...          |
//	|   stmt k+1             |     |   goto Body            |
//	|    ...                 |     +------------------------+
//	+------------------------+     +--[Body]----------------+
//	                               |   if !cond: goto After |
//	                               |   for-body...          |
//	                               |   goto Post            |
//	                               +------------------------+
//	                               +--[Post]----------------+
//	                               |   for-post...          |
//	                               |   goto Body            |
//	                               +------------------------+
//	                               +--[After]---------------+
//	                               |   stmt k+1             |
//	                               |   ...                  |
//	                               +------------------------+
func (fbb *funcBlockBuilder) remodelForStmt(s *ir.ForStmt) {
	var labelName string
	if name, ok := fbb.stmtLabelName[s.Pos()]; ok {
		labelName = name + `: `
	}

	// Create a block for the body of the for-loop.
	bodyBlk := fbb.fn.NewBlock(labelName + `For-loop Body`)
	fbb.blockPos[bodyBlk] = s.Pos()

	// Create a block to run post or jump to on continue.
	// If there is no post then just use the body block as the post block.
	postBlk := bodyBlk
	if s.Post != nil {
		postBlk = fbb.fn.NewBlock(labelName+`For-loop Post`, s.Post)
		postBlk.Body = append(postBlk.Body, ir.NewGotoBlockStmt(s.Pos(), bodyBlk))
	}
	fbb.continueBlock[s.Pos()] = postBlk

	// Split current block to make room for for-loop.
	afterBlk := fbb.fn.NewBlock(labelName + `After For-loop`)
	_, curJump := fbb.splitCurBlock(afterBlk)
	curJump.Block.Block = bodyBlk
	fbb.breakBlock[s.Pos()] = afterBlk

	// Insert the for-loop initialization into the current block before the jump.
	if s.Init != nil {
		index := len(fbb.curStmtList) - 1
		fbb.curStmtList = slices.Insert(fbb.curStmtList, index, s.Init)
		s.Init = nil
	}

	// Fill out the body for the for-loop including the conditional exit.
	if s.Cond != nil {
		ifCond := &ir.IfStmt{Cond: &ast.UnaryExpr{OpPos: s.Cond.Pos(), Op: token.NOT, X: s.Cond}}
		ifCond.Body = append(ifCond.Body, ir.NewGotoBlockStmt(s.Cond.Pos(), afterBlk))
		bodyBlk.Body = append(bodyBlk.Body, ifCond)
	}
	bodyBlk.Body = append(bodyBlk.Body, s.Body...)
	if !ir.IsFlowControlStatement(bodyBlk.LastStmt()) {
		bodyBlk.Body = append(bodyBlk.Body, ir.NewGotoBlockStmt(s.Pos(), postBlk))
	}
	if !ir.IsFlowControlStatement(postBlk.LastStmt()) {
		postBlk.Body = append(postBlk.Body, ir.NewGotoBlockStmt(s.Pos(), bodyBlk))
	}
}

func (fbb *funcBlockBuilder) remodelRangeStmt(s *ir.RangeStmt) {
	/*
		var labelName string
		if name, ok := fbb.stmtLabelName[s.Pos()]; ok {
			labelName = name + `: `
		}

		// Create a block for the body of the for-range.
		bodyBlk := fbb.fn.NewBlock(labelName + `For-range Body`)
		fbb.blockPos[bodyBlk] = s.Pos()
		fbb.continueBlock[s.Pos()] = bodyBlk

		// Split current block to make room for for-range.
		afterBlk := fbb.fn.NewBlock(labelName + `After For-range`)
		_, curJump := fbb.splitCurBlock(afterBlk)
		curJump.Block.Block = bodyBlk
		fbb.breakBlock[s.Pos()] = afterBlk
	*/

	// TODO: Finish
	panic(faults.New(`unimplemented`))
}

func (fbb *funcBlockBuilder) remodelReturnStmt(s *ir.ReturnStmt) {
	fbb.remodelExprSlice(s, nil, s.Results)
}

func (fbb *funcBlockBuilder) remodelBranchStmt(s *ir.BranchStmt) {
	switch s.Tok {
	case token.GOTO:
		fbb.remodelGotoBranchStmt(s)
	case token.BREAK:
		fbb.remodelBreakBranchStmt(s)
	case token.CONTINUE:
		fbb.remodelContinueBranchStmt(s)
	case token.FALLTHROUGH:
		fbb.remodelFallThroughBranchStmt(s)
	default:
		fbb.errGroup.Add(faults.New(`unhandled statement node in blocker`).
			With(`pos`, fbb.pos(s.Pos())).
			With(`branch`, s.Tok.String()).
			WithF(`type`, `%T`, s))
		return
	}
}

// remodelGotoBranchStmt handles a "goto" flow-control statement.
// See: https://go.dev/ref/spec#Goto_statements
func (fbb *funcBlockBuilder) remodelGotoBranchStmt(s *ir.BranchStmt) {
	if s.Label == nil {
		fbb.errGroup.Add(faults.New(`goto branch statement does not have a label`).
			With(`pos`, fbb.pos(s.Pos())))
		return
	}

	obj, ok := fbb.info().Uses[s.Label]
	if !ok {
		fbb.errGroup.Add(faults.New(`goto branch statement does not have object for usage`).
			With(`label`, s.Label.Name).
			With(`pos`, fbb.pos(s.Pos())))
		return
	}

	blk, ok := fbb.labelBlock[obj.Pos()]
	if !ok {
		// Create a preliminary block for the label this goes to.
		// Store this block with the label location so that any jumps to
		// this label can look up the block for this label and the actual
		// label can fill it out.
		blk = fbb.fn.NewBlock(`Label ` + s.Label.String())
		fbb.labelBlock[obj.Pos()] = blk
	}

	// Replace the branch statement with a goto block flow control
	// to jump to the label's block.
	fbb.curStmtList[fbb.stmtIndex] = ir.NewGotoBlockStmt(s.Pos(), blk)
	fbb.stmtIndex--
}

// findBlockPos attempts to determine the block that the given branch affects.
// The returned value is the position of the block that can be used to look
// up branching information for the block.
func (fbb *funcBlockBuilder) findBlockPos(s *ir.BranchStmt) token.Pos {
	if s.Label != nil {
		if obj, ok := fbb.info().Uses[s.Label]; ok {
			if stmtPos, ok := fbb.labelToStmt[obj.Pos()]; ok {
				return stmtPos
			}
			fbb.errGroup.Add(faults.New(`failed to find statement position for label position`).
				With(`label pos`, obj.Pos()).
				With(`label`, s.Label.String()).
				With(`pos`, fbb.pos(s.Pos())).
				With(`branch`, s.Tok.String()))
			return token.NoPos
		}
		fbb.errGroup.Add(faults.New(`failed to find position for labelled block`).
			With(`label`, s.Label.String()).
			With(`pos`, fbb.pos(s.Pos())).
			With(`branch`, s.Tok.String()))
		return token.NoPos
	}

	if pos, ok := fbb.blockPos[fbb.curBlock]; ok {
		return pos
	}
	fbb.errGroup.Add(faults.New(`failed to find position for block`).
		With(`pos`, fbb.pos(s.Pos())).
		With(`branch`, s.Tok.String()).
		WithF(`type`, `%T`, s))
	return token.NoPos
}

// remodelBreakBranchStmt handles a "break" flow-control statement.
// See: https://go.dev/ref/spec#Break_statements
func (fbb *funcBlockBuilder) remodelBreakBranchStmt(s *ir.BranchStmt) {
	pos := fbb.findBlockPos(s)
	blk, ok := fbb.breakBlock[pos]
	if !ok || blk == nil {
		fbb.errGroup.Add(faults.New(`failed to find break block for a pos`).
			With(`block pos`, fbb.pos(pos)).
			With(`branch`, s.Tok.String()).
			WithF(`type`, `%T`, s).
			With(`pos`, fbb.pos(s.Pos())))
		return
	}

	// Replace the branch statement with a goto block flow control
	// to jump to after the for-loop.
	fbb.curStmtList[fbb.stmtIndex] = ir.NewGotoBlockStmt(s.Pos(), blk)
	fbb.stmtIndex--
}

// remodelContinueBranchStmt handles a "continue" flow-control statement.
// See: https://go.dev/ref/spec#Continue_statements
func (fbb *funcBlockBuilder) remodelContinueBranchStmt(s *ir.BranchStmt) {
	pos := fbb.findBlockPos(s)
	blk, ok := fbb.continueBlock[pos]
	if !ok || blk == nil {
		fbb.errGroup.Add(faults.New(`failed to find continue block for a pos`).
			With(`block pos`, fbb.pos(pos)).
			With(`branch`, s.Tok.String()).
			WithF(`type`, `%T`, s).
			With(`pos`, fbb.pos(s.Pos())))
		return
	}

	// Replace the branch statement with a goto block flow control
	// to jump to after the for-loop.
	fbb.curStmtList[fbb.stmtIndex] = ir.NewGotoBlockStmt(s.Pos(), blk)
	fbb.stmtIndex--
}

// remodelFallThroughBranchStmt handles a "fallthrough" flow-control statement.
// See: https://go.dev/ref/spec#Fallthrough_statements
func (fbb *funcBlockBuilder) remodelFallThroughBranchStmt(s *ir.BranchStmt) {
	if s.Label != nil {
		fbb.errGroup.Add(faults.New(`unexpected label on a fall through branch statement`).
			With(`pos`, fbb.pos(s.Pos())).
			With(`branch`, s.Tok.String()).
			WithF(`type`, `%T`, s))
	}

	crumb.DropMsg(`Unimplemented`) //TODO: Implement
}

func (fbb *funcBlockBuilder) remodelIfStmt(s *ir.IfStmt) {
	if s.Init != nil {
		fbb.curStmtList = slices.Insert(fbb.curStmtList, fbb.stmtIndex, s.Init)
		s.Init = nil
		fbb.stmtIndex--
	}
	fbb.remodelExpr(s, nil, s.Cond)
	fbb.remodelStmtSlice(s.Body)
	fbb.remodelStmtSlice(s.Else)
}

func (fbb *funcBlockBuilder) remodelExprStmt(s *ir.ExprStmt) {
	fbb.remodelExpr(s, nil, s.X)
}

func (fbb *funcBlockBuilder) remodelExprSlice(s ir.Stmt, stack []ast.Expr, es []ast.Expr) {
	for _, e := range es {
		fbb.remodelExpr(s, stack, e)
	}
}

func (fbb *funcBlockBuilder) remodelExpr(s ir.Stmt, stack []ast.Expr, e ast.Expr) {
	switch e := e.(type) {
	case nil, *ast.BadExpr, *ast.Ident, *ast.BasicLit:
		// Do Nothing
		return
	case *ast.StarExpr:
		fbb.remodelExpr(s, append(stack, e), e.X)
	case *ast.UnaryExpr:
		fbb.remodelUnaryExpr(s, stack, e)
	case *ast.BinaryExpr:
		fbb.remodelBinaryExpr(s, stack, e)
	case *ast.CallExpr:
		fbb.remodelCallExpr(s, stack, e)
	default:
		fbb.errGroup.Add(faults.New(`unhandled expression node in blocker`).
			With(`pos`, fbb.pos(e.Pos())).
			WithF(`type`, `%T`, e))
		return
	}
}

func (fbb *funcBlockBuilder) remodelUnaryExpr(s ir.Stmt, stack []ast.Expr, e *ast.UnaryExpr) {
	switch e.Op {
	case token.INC, token.DEC, token.ADD, token.SUB, token.NOT, token.XOR, token.MUL, token.AND:
		fbb.remodelExpr(s, append(stack, e), e.X)
	case token.ARROW:
		fbb.remodelReceiveExpr(s, stack, e)
	default:
		fbb.errGroup.Add(faults.New(`unhandled unary expression node in blocker`).
			With(`pos`, fbb.pos(e.Pos())).
			With(`op`, e.Op.String()).
			WithF(`type`, `%T`, e))
		return
	}
}

// remodelReceiveExpr handles a unary expression for receiving a value from a channel.
// See: https://go.dev/ref/spec#Receive_operator
func (fbb *funcBlockBuilder) remodelReceiveExpr(s ir.Stmt, stack []ast.Expr, e *ast.UnaryExpr) {
	crumb.DropMsg(`Unimplemented`) //TODO: Implement
}

func (fbb *funcBlockBuilder) remodelBinaryExpr(s ir.Stmt, stack []ast.Expr, e *ast.BinaryExpr) {
	switch e.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.AND, token.OR,
		token.XOR, token.SHL, token.SHR, token.AND_NOT, token.ADD_ASSIGN,
		token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN,
		token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.SHL_ASSIGN,
		token.SHR_ASSIGN, token.AND_NOT_ASSIGN, token.EQL, token.LSS, token.GTR,
		token.NEQ, token.LEQ, token.GEQ:
		fbb.remodelExpr(s, append(stack, e), e.X)
		fbb.remodelExpr(s, append(stack, e), e.Y)
	case token.LAND:
		fbb.remodelLogicalAndExpr(s, stack, e)
	case token.LOR:
		fbb.remodelLogicalOrExpr(s, stack, e)
	default:
		fbb.errGroup.Add(faults.New(`unhandled binary expression node in blocker`).
			With(`pos`, fbb.pos(e.Pos())).
			With(`op`, e.Op.String()).
			WithF(`type`, `%T`, e))
		return
	}
}

func (fbb *funcBlockBuilder) remodelLogicalAndExpr(s ir.Stmt, stack []ast.Expr, e *ast.BinaryExpr) {
	crumb.DropMsg(`Unimplemented`) //TODO: Implement
}

func (fbb *funcBlockBuilder) remodelLogicalOrExpr(s ir.Stmt, stack []ast.Expr, e *ast.BinaryExpr) {
	crumb.DropMsg(`Unimplemented`) //TODO: Implement
}

func (fbb *funcBlockBuilder) remodelCallExpr(s ir.Stmt, stack []ast.Expr, e *ast.CallExpr) {
	if len(stack) <= 0 {
		follow := fbb.fn.NewBlock(`follow call`)
		_, gotoFollow := fbb.splitCurBlock(follow)
		// Replace the goto with a call
		fbb.curStmtList[fbb.stmtIndex+1] = &ir.FuncCallStmt{
			Ast:    e,
			Fun:    e.Fun,
			Args:   e.Args,
			Follow: gotoFollow.Block,
		}

		// TODO: Add assignment or return to follow
		return
	}

	fmt.Printf("(Stmt) %T\n", s) // TODO: Remove
	for i, sx := range stack {   // TODO: Remove
		fmt.Printf("(%d) %T\n", i+1, sx) // TODO: Remove
	} // TODO: Remove
	ast.Print(token.NewFileSet(), e) // TODO: Remove

	crumb.DropMsg(`Unimplemented`) //TODO: Implement
}
