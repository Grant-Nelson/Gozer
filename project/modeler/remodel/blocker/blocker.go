package blocker

import (
	"go/ast"
	"go/token"
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
	return &blocker{
		errGroup: cfg.ErrGroup,
	}
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
		errGroup:   bb.errGroup,
		pkg:        bb.pkg,
		fn:         fn,
		labelBlock: map[token.Pos]*ir.Block{},
	}

	for blockIndex := 0; blockIndex < len(fn.Blocks); blockIndex++ {
		b := fn.Blocks[blockIndex]
		fbb.curBlock = b
		fbb.remodelBlock(b)
	}
	return true, bb.errGroup.FullOrNil()
}

type funcBlockBuilder struct {
	errGroup  *faults.ErrGroup
	pkg       *project.Package
	fn        *ir.Func
	curBlock  *ir.Block
	stmtIndex int

	labelBlock map[token.Pos]*ir.Block
}

func (fbb *funcBlockBuilder) remodelBlock(b *ir.Block) {
	for fbb.stmtIndex = 0; fbb.stmtIndex < len(b.Body); fbb.stmtIndex++ {
		fbb.remodelIrcStmt(b.Body[fbb.stmtIndex])
	}
}

func (fbb *funcBlockBuilder) remodelIrcStmt(s ir.Stmt) {
	switch s := s.(type) {
	case *ir.BaseStmt:
		fbb.remodelAstStmt(s.Stmt)
	}
}

func (fbb *funcBlockBuilder) remodelAstStmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.DeclStmt, *ast.EmptyStmt:
		// Do Nothing
	case *ast.LabeledStmt:
		fbb.remodelLabeledStmt(s)
	//case *ast.ExprStmt:   // TODO: Implement
	//case *ast.SendStmt:   // TODO: Implement
	//case *ast.IncDecStmt: // TODO: Implement
	case *ast.AssignStmt:
		fbb.remodelAssignStmt(s)
	//case *ast.GoStmt:    // TODO: Implement
	//case *ast.DeferStmt: // TODO: Implement
	case *ast.ReturnStmt:
		fbb.remodelReturnStmt(s)
	case *ast.BranchStmt:
		fbb.remodelBranchStmt(s)
	case *ast.BlockStmt:
		fbb.remodelBlockStmt(s)
	case *ast.IfStmt:
		fbb.remodelIfStmt(s)
	//case *ast.SwitchStmt:     // TODO: Implement
	//case *ast.TypeSwitchStmt: // TODO: Implement
	//case *ast.SelectStmt:     // TODO: Implement
	//case *ast.ForStmt:        // TODO: Implement
	//case *ast.RangeStmt:      // TODO: Implement
	default:
		panic(faults.New(`unhandled statement node in blocker`).
			WithF(`type`, `%T`, s))
	}
}

func (fbb *funcBlockBuilder) remodelAssignStmt(s *ast.AssignStmt) {
	fbb.remodelExpSlice(s, s.Lhs)
	fbb.remodelExpSlice(s, s.Rhs)
}

// remodelLabeledStmt processes a label statement.
// A label can be jumped to so the code reachable from the label
// needs to be put into it's own block.
func (fbb *funcBlockBuilder) remodelLabeledStmt(s *ast.LabeledStmt) {
	// Check if a block was preemptively created by prior code
	// that is jumping forward to this block.
	blk, ok := fbb.labelBlock[s.Label.Pos()]
	if ok {
		if len(blk.Body) > 0 {
			fbb.errGroup.Add(faults.New(`preemptive label block is already populated`).
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
		fbb.labelBlock[s.Label.Pos()] = blk
	}

	// TODO: Need to handle for-loop or special targeted statement for the label
	// so that things like the initialization of the for-loop is in cur block
	// and the comparator and body get put into another block, that is the
	// labelled block.

	// Put the statement that the label is attached to as the first statement
	// in the block, then move all following statements from the current block
	// into the new block.
	blk.Body = []ir.Stmt{&ir.BaseStmt{Stmt: s.Stmt}}
	blk.Body = append(blk.Body, fbb.curBlock.Body[fbb.stmtIndex+1:]...)
	fbb.curBlock.Body = fbb.curBlock.Body[:fbb.stmtIndex]

	// Add a goto in the current block to jump to the label block since the
	// code flow goes from the current block into the label unconditionally.
	jump := &ir.GotoStmt{Block: &ir.BlockRef{Block: blk}}
	fbb.curBlock.Body = append(fbb.curBlock.Body, jump)

	// Step back th statement index so that the new goto statement that
	// took the index space of the current statement is processed.
	fbb.stmtIndex--
}

func (fbb *funcBlockBuilder) remodelReturnStmt(s *ast.ReturnStmt) {
	fbb.remodelExpSlice(s, s.Results)
}

func (fbb *funcBlockBuilder) remodelBranchStmt(s *ast.BranchStmt) {
	//TODO: Implement
	crumb.DropMsg(`Unimplemented`)
}

func (fbb *funcBlockBuilder) remodelBlockStmt(s *ast.BlockStmt) {

	//TODO: Implement
	crumb.Drop()
}

func (fbb *funcBlockBuilder) remodelIfStmt(s *ast.IfStmt) {
	if s.Init != nil {
		init := &ir.BaseStmt{Stmt: s.Init}
		fbb.curBlock.Body = slices.Insert(fbb.curBlock.Body, fbb.stmtIndex, ir.Stmt(init))
		s.Init = nil
		fbb.stmtIndex--
	}
	fbb.remodelExp(s, s.Cond)
	fbb.remodelBlockStmt(s.Body)
	if s.Else != nil {
		fbb.remodelAstStmt(s.Else)
	}
}

func (fbb *funcBlockBuilder) remodelExpSlice(s ast.Stmt, es []ast.Expr) {
	for _, e := range es {
		fbb.remodelExp(s, e)
	}
}

func (fbb *funcBlockBuilder) remodelExp(s ast.Stmt, e ast.Expr) {
	//TODO: Implement
	crumb.DropMsg(`Unimplemented`)
}
