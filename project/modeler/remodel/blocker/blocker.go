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

		jumpBlock:  map[token.Pos]*ir.Block{},
		afterBlock: map[token.Pos]*ir.Block{},
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

	jumpBlock  map[token.Pos]*ir.Block
	afterBlock map[token.Pos]*ir.Block
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
	blk, ok := fbb.jumpBlock[s.Label.Pos()]
	if ok {
		if len(blk.Body) > 0 {
			fbb.errGroup.Add(faults.New(`preemptive label block is already populated`).
				With(`pos`, fbb.pos(s.Pos())).
				With(`statements`, len(blk.Body)).
				With(`label block`, blk).
				With(`current block`, fbb.curBlock).
				With(`label`, s.Label.String()))
			return
		}
	} else {
		// Create a new block for the code reachable from the label.
		blk = fbb.fn.NewBlock()
		blk.Hint = `Label ` + s.Label.String()

		// Store this block with the label location so that any jumps to
		// this label can look up the block for this label.
		fbb.jumpBlock[s.Label.Pos()] = blk
	}

	// TODO: Need to handle for-loop or special targeted statement for the label
	// so that things like the initialization of the for-loop is in cur block
	// and the comparator and body get put into another block, that is the
	// labelled block.

	// Put the statement that the label is attached to as the first statement
	// in the block, then move all following statements from the current block
	// into the new block.
	blk.Body = []ir.Stmt{s.Stmt}
	blk.Body = append(blk.Body, fbb.curStmtList[fbb.stmtIndex+1:]...)
	fbb.curStmtList = fbb.curStmtList[:fbb.stmtIndex]

	// Add a goto in the current block to jump to the label block since the
	// code flow goes from the current block into the label unconditionally.
	jump := &ir.GotoBlockStmt{SrcPos: s.Stmt.Pos(), Block: &ir.BlockRef{Block: blk}}
	fbb.curStmtList = append(fbb.curStmtList, jump)

	// Step back the statement index so that the new goto statement,
	// that took the index space of the current statement, is processed.
	fbb.stmtIndex--
}

func (fbb *funcBlockBuilder) remodelReturnStmt(s *ir.ReturnStmt) {
	fbb.remodelExpSlice(s, s.Results)
}

func (fbb *funcBlockBuilder) remodelBranchStmt(s *ir.BranchStmt) {
	switch s.Tok {
	case token.GOTO:
		fbb.remodelGotoBranchStmt(s)
	default:
		fbb.errGroup.Add(faults.New(`unhandled statement node in blocker`).
			With(`pos`, fbb.pos(s.Pos())).
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

	blk, ok := fbb.jumpBlock[obj.Pos()]
	if !ok {
		// Create a preliminary block for the label this goes to.
		blk = fbb.fn.NewBlock()
		blk.Hint = `Label ` + s.Label.String()

		// Store this block with the label location so that any jumps to
		// this label can look up the block for this label and the actual
		// label can fill it out.
		fbb.jumpBlock[obj.Pos()] = blk
	}

	// Replace the branch statement with a goto block flow control
	// to jump to the label's block.
	jump := &ir.GotoBlockStmt{SrcPos: s.Pos(), Block: &ir.BlockRef{Block: blk}}
	fbb.curStmtList[fbb.stmtIndex] = jump

	// Step back the statement index so that the new goto statement,
	// that took the index space of the current statement, is processed.
	fbb.stmtIndex--
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

func (fbb *funcBlockBuilder) remodelExpSlice(s ir.Stmt, es []ast.Expr) {
	for _, e := range es {
		fbb.remodelExp(s, e)
	}
}

func (fbb *funcBlockBuilder) remodelExp(s ir.Stmt, e ast.Expr) {
	//TODO: Implement
	crumb.DropMsg(`Unimplemented`)
}
