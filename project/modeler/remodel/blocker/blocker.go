package blocker

import (
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

func (bb *blockBuilder) RemodelFunc(fn *ir.Func) (con bool, err error) {
	bb.errGroup.Recover(&err)
	if fn.Atomic() {
		return true, nil
	}

	fbb := &funcBlockBuilder{
		errGroup: bb.errGroup,
		pkg:      bb.pkg,
		fn:       fn,

		beforeBlock: map[token.Pos]*ir.Block{},
		innerBlock:  map[token.Pos]*ir.Block{},
		afterBlock:  map[token.Pos]*ir.Block{},
		forBlock:    map[*ir.Block]token.Pos{},
	}

	for blockIndex := 0; blockIndex < len(fn.Blocks); blockIndex++ {
		fbb.remodelBlock(fn.Blocks[blockIndex])
	}
	return true, bb.errGroup.FullOrNil()
}

type funcBlockBuilder struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
	fn       *ir.Func
	curBlock *ir.Block

	stmtIndex   int
	curStmtList []ir.Stmt

	beforeBlock map[token.Pos]*ir.Block
	innerBlock  map[token.Pos]*ir.Block
	afterBlock  map[token.Pos]*ir.Block
	forBlock    map[*ir.Block]token.Pos
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

func (fbb *funcBlockBuilder) remodelStmt(s ir.Stmt) {
	switch s := s.(type) {
	case *ir.DeclStmt, *ir.GotoBlockStmt:
		// Do Nothing
	case *ir.LabeledStmt:
		fbb.remodelLabeledStmt(s)
	case *ir.AssignStmt:
		fbb.remodelAssignStmt(s)
	case *ir.ReturnStmt:
		fbb.remodelReturnStmt(s)
	case *ir.BranchStmt:
		fbb.remodelBranchStmt(s)
	case *ir.IfStmt:
		fbb.remodelIfStmt(s)
	case *ir.IncDecStmt:
		fbb.remodelIncDecStmt(s)
	default:
		fbb.errGroup.Add(faults.New(`unhandled statement node in blocker`).
			With(`pos`, fbb.pos(s.Pos())).
			WithF(`type`, `%T`, s))
		return
	}
}

func (fbb *funcBlockBuilder) remodelAssignStmt(s *ir.AssignStmt) {
	fbb.remodelExpSlice(s, s.Lhs)
	fbb.remodelExpSlice(s, s.Rhs)
}

// remodelLabeledStmt processes a label statement.
// A label can be jumped to so the code reachable from the label
// needs to be put into it's own block.
func (fbb *funcBlockBuilder) remodelLabeledStmt(s *ir.LabeledStmt) {
	// Check if a block was preemptively created by prior code
	// that is jumping forward to this block.
	nextBlk, ok := fbb.beforeBlock[s.Label.Pos()]
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
		nextBlk = fbb.fn.NewBlock()
		nextBlk.Hint = `Label ` + s.Label.String()
		fbb.beforeBlock[s.Label.Pos()] = nextBlk
	}

	// Remove all the following statements from the current block for a later
	// block. Add a goto in the current block to jump to the label block since
	// the code flow goes from the current block into the label unconditionally.
	follow := fbb.curStmtList[fbb.stmtIndex+1:]
	fbb.curStmtList = slices.Clone(fbb.curStmtList[:fbb.stmtIndex])
	fbb.curStmtList = append(fbb.curStmtList, ir.NewGotoBlockStmt(s.Stmt.Pos(), nextBlk))
	fbb.stmtIndex--

	// Handle for-loop or special targeted statement for the label so
	// that things like the initialization of the for-loop is in cur block
	// and the comparator and body get put into another block, that is the
	// labelled block.
	if sf, ok := s.Stmt.(*ir.ForStmt); ok {
		if sf.Init == nil {
			// Without a for-loop initialization, the loop on the body of the
			// for-loop is the same as the jump for the label and body of the loop.
			fbb.innerBlock[s.Label.Pos()] = nextBlk
		} else {
			// Change the next block into the for-loop initialization (init block),
			// create a new block for the body of the for-loop, then add the
			// for-loop init statement and a goto to the init block.
			initBlk := nextBlk
			initBlk.Hint = `For-loop Init for ` + s.Label.String()
			nextBlk = fbb.fn.NewBlock()
			initBlk.Body = append(initBlk.Body, sf.Init, ir.NewGotoBlockStmt(sf.Init.Pos(), nextBlk))
			sf.Init = nil
		}

		// Create a new block for after the for-loop, fill out the body of the
		// for-loop, and add the basic for-loop goto's to loop correctly.
		bodyBlk := nextBlk
		bodyBlk.Hint = `For-loop Body for ` + s.Label.String()
		fbb.innerBlock[s.Label.Pos()] = bodyBlk
		fbb.forBlock[bodyBlk] = s.Label.Pos()
		nextBlk = fbb.fn.NewBlock()
		nextBlk.Hint = `After For-loop for ` + s.Label.String()
		fbb.afterBlock[s.Label.Pos()] = nextBlk

		// Fill out the body for the for-loop including the conditional exit.
		if sf.Cond != nil {
			ifCond := &ir.IfStmt{Cond: &ast.UnaryExpr{OpPos: sf.Cond.Pos(), Op: token.NOT, X: sf.Cond}}
			ifCond.Body = append(ifCond.Body, ir.NewGotoBlockStmt(sf.Cond.Pos(), nextBlk))
			bodyBlk.Body = append(bodyBlk.Body, ifCond)
		}
		bodyBlk.Body = append(bodyBlk.Body, sf.Body...)
		if !ir.IsFlowControlStatement(bodyBlk.LastStmt()) {
			// TODO: Move Post to it's own block for end of loop and continues.
			if sf.Post != nil {
				bodyBlk.Body = append(bodyBlk.Body, sf.Post)
			}
			bodyBlk.Body = append(bodyBlk.Body, ir.NewGotoBlockStmt(sf.Pos(), bodyBlk))
		}

		// Put all following statements into the block after the for-loop.
		nextBlk.Body = append(nextBlk.Body, follow...)
		return
	}

	// Put the statement that the label is attached to as the first statement
	// in the block, then move all following statements from the current block
	// into the new block.
	nextBlk.Body = []ir.Stmt{s.Stmt}
	nextBlk.Body = append(nextBlk.Body, follow...)
}

func (fbb *funcBlockBuilder) remodelReturnStmt(s *ir.ReturnStmt) {
	fbb.remodelExpSlice(s, s.Results)
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

	blk, ok := fbb.beforeBlock[obj.Pos()]
	if !ok {
		// Create a preliminary block for the label this goes to.
		blk = fbb.fn.NewBlock()
		blk.Hint = `Label ` + s.Label.String()

		// Store this block with the label location so that any jumps to
		// this label can look up the block for this label and the actual
		// label can fill it out.
		fbb.beforeBlock[obj.Pos()] = blk
	}

	// Replace the branch statement with a goto block flow control
	// to jump to the label's block.
	fbb.curStmtList[fbb.stmtIndex] = ir.NewGotoBlockStmt(s.Pos(), blk)
	fbb.stmtIndex--
}

func (fbb *funcBlockBuilder) findBlockPos(s *ir.BranchStmt) token.Pos {
	if s.Label != nil {
		obj, ok := fbb.info().Uses[s.Label]
		if !ok {
			fbb.errGroup.Add(faults.New(`failed to find for-loop position for labelled block`).
				With(`label`, s.Label.String()).
				With(`pos`, fbb.pos(s.Pos())).
				With(`branch`, s.Tok.String()))
			return token.NoPos
		}
		return obj.Pos()
	}
	if pos, ok := fbb.forBlock[fbb.curBlock]; ok {
		return pos
	}
	fbb.errGroup.Add(faults.New(`failed to find for-loop position for block`).
		With(`pos`, fbb.pos(s.Pos())).
		With(`branch`, s.Tok.String()).
		WithF(`type`, `%T`, s))
	return token.NoPos
}

func (fbb *funcBlockBuilder) remodelBreakBranchStmt(s *ir.BranchStmt) {
	pos := fbb.findBlockPos(s)
	blk, ok := fbb.afterBlock[pos]
	if !ok || blk == nil {
		fbb.errGroup.Add(faults.New(`failed to find after block for a pos`).
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

func (fbb *funcBlockBuilder) remodelContinueBranchStmt(s *ir.BranchStmt) {
	pos := fbb.findBlockPos(s)
	blk, ok := fbb.innerBlock[pos]
	if !ok || blk == nil {
		fbb.errGroup.Add(faults.New(`failed to find after block for a pos`).
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

func (fbb *funcBlockBuilder) remodelFallThroughBranchStmt(s *ir.BranchStmt) {
	if s.Label != nil {
		fbb.errGroup.Add(faults.New(`unexpected label on a fall through branch statement`).
			With(`pos`, fbb.pos(s.Pos())).
			With(`branch`, s.Tok.String()).
			WithF(`type`, `%T`, s))
	}

	//TODO: Implement
	crumb.DropMsg(`Unimplemented`)
}

func (fbb *funcBlockBuilder) remodelIfStmt(s *ir.IfStmt) {
	if s.Init != nil {
		fbb.curStmtList = slices.Insert(fbb.curStmtList, fbb.stmtIndex, s.Init)
		s.Init = nil
		fbb.stmtIndex--
	}
	fbb.remodelExp(s, s.Cond)
	fbb.remodelStmtSlice(s.Body)
	fbb.remodelStmtSlice(s.Else)
}

func (fbb *funcBlockBuilder) remodelIncDecStmt(s *ir.IncDecStmt) {
	fbb.remodelExp(s, s.X)
}

func (fbb *funcBlockBuilder) remodelExpSlice(s ir.Stmt, es []ast.Expr) {
	for _, e := range es {
		fbb.remodelExp(s, e)
	}
}

func (fbb *funcBlockBuilder) remodelExp(s ir.Stmt, e ast.Expr) {
	//TODO: Implement
	crumb.DropMsg(`Unimplemented`)
}
