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
		case *ir.ConstDef:
			p.Consts = append(p.Consts, mem)
		case *ir.VarDef:
			p.Vars = append(p.Vars, mem)
		case *ir.FuncDef:
			p.Funcs = append(p.Funcs, mem)
		case *ir.TypeDef:
			p.Types = append(p.Types, mem)
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

func (c *converter) FromFunction(f *ssa.Function) *ir.FuncDef {

	// TODO: Implement
	fmt.Printf("FuncDef: %v\n", f)

	blocks := make([]*ir.Block, len(f.Blocks))
	for i, blk := range f.Blocks {
		b := c.FromBlock(blk)
		b.Index = i
		blocks[i] = b
	}

	return &ir.FuncDef{
		//FuncObj: fnObj,
		Func: &ir.Func{
			FuncPos:   f.Pos(),
			Signature: f.Signature,
			Blocks:    blocks,
		},
	}
}

func (c *converter) FromGlobal(g *ssa.Global) *ir.VarDef {

	// TODO: Implement
	fmt.Printf("VarDef: %v\n", g)

	return &ir.VarDef{}
}

func (c *converter) FromNamedConst(nc *ssa.NamedConst) *ir.ConstDef {

	// TODO: Implement
	fmt.Printf("ConstDef: %v\n", nc)

	return &ir.ConstDef{}
}

func (c *converter) FromType(t *ssa.Type) *ir.TypeDef {

	// TODO: Implement
	fmt.Printf("TypeDef: %v\n", t)

	return &ir.TypeDef{}
}

func (c *converter) FromBlock(blk *ssa.BasicBlock) *ir.Block {

	// TODO: Implement
	fmt.Printf("Block: %v\n", blk)

	return &ir.Block{}
}
