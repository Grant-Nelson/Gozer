package ssaConverter

import (
	"fmt"
	"go/token"
	"go/types"
	"maps"
	"slices"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func ConvertPackage(astPkg *packages.Package, ssaPkg *ssa.Package, errGroup *faults.ErrGroup) (p *ir.Package, err error) {
	defer errGroup.Recover(&err)

	assert.NotNil(astPkg)
	assert.NotNil(ssaPkg)
	assert.NotNil(ssaPkg.Pkg)
	assert.NotNil(ssaPkg.Prog)
	assert.NotNil(ssaPkg.Prog.Fset)

	c := &converter{
		AstPkg:  astPkg,
		SsaPkg:  ssaPkg,
		FileSet: ssaPkg.Prog.Fset,
		Info:    astPkg.TypesInfo,
		Sizes:   astPkg.TypesSizes,
		Errors:  errGroup,
	}
	return c.FromPackage(ssaPkg), nil
}

type converter struct {
	AstPkg  *packages.Package
	SsaPkg  *ssa.Package
	FileSet *token.FileSet
	Info    *types.Info
	Sizes   types.Sizes
	Errors  *faults.ErrGroup
}

func (c *converter) pos(p token.Pos) string {
	return c.FileSet.Position(p).String()
}

func (c *converter) addFault(f *faults.Fault) {
	if c == nil || c.Errors == nil {
		panic(f)
	}
	if err := c.Errors.Add(f); err != nil {
		panic(err)
	}
}

func (c *converter) FromPackage(pkg *ssa.Package) *ir.Package {
	p := &ir.Package{
		Name:    pkg.Pkg.Name(),
		Path:    pkg.Pkg.Path(),
		FileSet: c.FileSet,
		Sizes:   c.Sizes,
	}

	names := slices.Sorted(maps.Keys(pkg.Members))
	for _, name := range names {
		mem := c.FromMember(pkg.Members[name])
		switch mem := mem.(type) {
		case nil:
		case *ir.ConstDef:
			if mem != nil {
				p.Consts = append(p.Consts, mem)
			}
		case *ir.VarDef:
			if mem != nil {
				p.Vars = append(p.Vars, mem)
			}
		case *ir.FuncDef:
			if mem != nil {
				p.Funcs = append(p.Funcs, mem)
			}
		case *ir.TypeDef:
			if mem != nil {
				p.Types = append(p.Types, mem)
			}
		default:
			c.addFault(faults.New(`unexpected AST node type`).
				WithTypeOf(`type`, mem).
				With(`statement`, mem).
				With(`pos`, c.pos(mem.Pos())))
		}
	}

	return p
}

func (c *converter) FromMember(m ssa.Member) ir.Def {
	switch m := m.(type) {
	case *ssa.Function:
		return c.FromFunction(m)
	case *ssa.Global:
		return c.FromGlobal(m)
	case *ssa.NamedConst:
		return c.FromNamedConst(m)
	case *ssa.Type:
		return c.FromType(m)
	default:
		panic(faults.New(`Unexpected Member type`).
			With(`package`, c.SsaPkg.Pkg.Path()).
			WithF(`member`, `%v`, m).
			With(`pos`, c.pos(m.Pos())).
			WithTypeOf(`type`, m))
	}
}

func (c *converter) FromGlobal(g *ssa.Global) *ir.VarDef {
	//printMemberTree(g) // TODO: REMOVE

	// TODO: Implement
	fmt.Printf("VarDef: %v\n", g)

	return &ir.VarDef{}
}

func (c *converter) FromNamedConst(nc *ssa.NamedConst) *ir.ConstDef {
	//printMemberTree(nc) // TODO: REMOVE

	// TODO: Implement
	fmt.Printf("ConstDef: %v\n", nc)

	return &ir.ConstDef{}
}

func (c *converter) FromType(t *ssa.Type) *ir.TypeDef {
	//printMemberTree(t) // TODO: REMOVE

	// TODO: Implement
	fmt.Printf("TypeDef: %v\n", t)

	return &ir.TypeDef{}
}

func (c *converter) FromFunction(f *ssa.Function) *ir.FuncDef {
	synthetic := len(f.Synthetic) > 0
	if synthetic {
		// TODO: Determine if any of the synthetic methods needs to be kept,
		//       for example, check if test main is synthesized and can be used.
		//       Note that synthetic methods have a nil Object.
		fmt.Printf("Skipping Synthetic Function: %s\n", f.Synthetic)
		return nil
	}

	printFunctionDebug(f) // TODO: REMOVE
	//printMemberTree(f)    // TODO: REMOVE

	if f.Object() == nil {
		panic(faults.New(`Function had null object`).
			With(`pos`, c.pos(f.Pos())).
			With(`name`, f.Name()))
	}

	obj, isObjFunc := f.Object().(*types.Func)
	if !isObjFunc {
		panic(faults.New(`Function had null object`).
			With(`pos`, c.pos(f.Pos())).
			With(`name`, f.Name()).
			WithTypeOf(`object type`, f.Object()).
			With(`object`, f.Object().String()))
	}

	blocks := make([]*ir.Block, len(f.Blocks))
	for i, blk := range f.Blocks {
		blocks[i] = c.convertBlock(blk)
	}

	// TODO: Update when blocks are injected for obstructions
	getIrBlock := func(bb *ssa.BasicBlock) *ir.Block { return blocks[bb.Index] }
	for i, blk := range f.Blocks {
		b := blocks[i]
		b.Prior = iterator.Iterate(blk.Preds).Select(getIrBlock).ToSlice()
		b.Follow = iterator.Iterate(blk.Succs).Select(getIrBlock).ToSlice()
	}

	return &ir.FuncDef{
		FuncObj: obj,
		Func: &ir.Func{
			FuncPos:   f.Pos(),
			Signature: f.Signature,
			Blocks:    blocks,
			// TODO: Set ReturnBlocks
		},
	}
}

func (c *converter) convertBlock(blk *ssa.BasicBlock) *ir.Block {
	body := []ir.Stmt{}
	for _, in := range blk.Instrs {
		if stmt := c.fromInstruction(in); stmt != nil {
			body = append(body, stmt)
		}
	}
	return &ir.Block{
		Index: blk.Index,
		Hint:  blk.Comment,
		Body:  body,
		//Params []*Param
	}
}

func (c *converter) fromInstruction(in ssa.Instruction) ir.Stmt {
	switch x := in.(type) {
	case ssa.Value:
		return c.fromValueInstruction(x)
	// TODO: Add more as needed
	default:
		panic(faults.New(`Unexpected instruction type`).
			With(`package`, c.SsaPkg.Pkg.Path()).
			With(`instruction`, in.String()).
			With(`pos`, c.pos(in.Pos())).
			WithTypeOf(`type`, in))
	}
}

func (c *converter) fromValueInstruction(v ssa.Value) ir.Stmt {
	if embeddable(v) {
		return nil
	}
	expr := c.fromValue(v)
	if expr == nil {
		return expr
	}
	return &ir.VarDef{
		VarObj: types.NewVar(v.Pos(), c.SsaPkg.Pkg, v.Name(), v.Type()),
		Value:  expr,
	}
}

func (c *converter) fromValue(v ssa.Value) ir.Expr {
	switch x := v.(type) {
	case *ssa.BinOp:
		return c.fromBinOp(x)
	// TODO: Add more as needed
	default:
		panic(faults.New(`Unexpected value type`).
			With(`package`, c.SsaPkg.Pkg.Path()).
			With(`instruction`, v.String()).
			With(`pos`, c.pos(v.Pos())).
			WithTypeOf(`type`, v))
	}
}

func (c *converter) fromBinOp(in *ssa.BinOp) ir.Expr {

}

func (c *converter) fromOperand(v ssa.Value) ir.Expr {
	if embeddable(v) {
		return c.fromValue(v)
	}

	// TODO: Finish
}

func embeddable(v ssa.Value) bool {
	if len(*v.Referrers()) != 1 {
		return false
	}
	switch x := v.(type) {
	case *ssa.BinOp:
		return true
	case *ssa.UnOp:
		return x.Op != token.ARROW && !x.CommaOk
		// TODO: Add more as needed
	}
	return false
}

func obstructionInstr(in ssa.Instruction) bool {
	switch x := in.(type) {
	case *ssa.Call, *ssa.Send:
		return true
	case *ssa.UnOp:
		return x.Op == token.ARROW && !x.CommaOk
		// TODO: Add more as needed
	}
	return false
}
