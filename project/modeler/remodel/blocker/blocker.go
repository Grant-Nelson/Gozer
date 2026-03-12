package blocker

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/irc"
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

func (bb *blockBuilder) RemodelFunc(fn *irc.Func) (con bool, err error) {
	bb.errGroup.Recover(&err)
	if fn.Atomic() {
		return true, nil
	}

	fbb := &funcBlockBuilder{
		errGroup:   bb.errGroup,
		pkg:        bb.pkg,
		fn:         fn,
		labelBlock: map[token.Pos]*irc.Block{},
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
	fn        *irc.Func
	curBlock  *irc.Block
	stmtIndex int

	labelBlock map[token.Pos]*irc.Block
}

func (fbb *funcBlockBuilder) remodelBlock(b *irc.Block) {
	for fbb.stmtIndex = 0; fbb.stmtIndex < len(b.Body); fbb.stmtIndex++ {
		fbb.remodelIrcStmt(b.Body[fbb.stmtIndex])
	}
}

func (fbb *funcBlockBuilder) remodelIrcStmt(s irc.Stmt) {
	switch s := s.(type) {
	case *irc.BaseStmt:
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
	//case *ast.BranchStmt: // TODO: Implement
	//case *ast.BlockStmt:  // TODO: Implement
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

func (fbb *funcBlockBuilder) remodelLabeledStmt(s *ast.LabeledStmt) {
	blk := fbb.fn.NewBlock()
	blk.Hint = `Label ` + s.Label.String()
	fbb.labelBlock[s.Label.Pos()] = blk

	blk.Body = []irc.Stmt{&irc.BaseStmt{Stmt: s.Stmt}}
	blk.Body = append(blk.Body, fbb.curBlock.Body[fbb.stmtIndex+1:]...)
	fbb.curBlock.Body = fbb.curBlock.Body[:fbb.stmtIndex]
	fbb.stmtIndex--

	ref := &irc.BlockRef{Block: blk}
	fbb.curBlock.Body = append(fbb.curBlock.Body, &irc.GotoStmt{Block: ref})
}

func (fbb *funcBlockBuilder) remodelReturnStmt(s *ast.ReturnStmt) {
	fbb.remodelExpSlice(s, s.Results)
}

func (fbb *funcBlockBuilder) remodelIfStmt(s *ast.IfStmt) {
	s.Body
	//TODO: Implement
}

func (fbb *funcBlockBuilder) remodelExpSlice(s ast.Stmt, es []ast.Expr) {
	for _, e := range es {
		fbb.remodelExp(s, e)
	}
}

func (fbb *funcBlockBuilder) remodelExp(s ast.Stmt, e ast.Expr) {
	//TODO: Implement
}
