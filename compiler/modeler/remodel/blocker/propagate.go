package blocker

import (
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/binaryOp"
)

// varSet is a set of `*types.Var`'s`.
//
// The value is an monotonically incrementing value (i.e. the length
// of the map plus one). The values must be unique and continuous.
// Those values are used to get the order of the keys in the order
// the keys were added to the set.
type varSet map[*types.Var]int

func (s varSet) add(v *types.Var) bool {
	if v != nil {
		if _, has := s[v]; !has {
			s[v] = len(s) + 1
			return true
		}
	}
	return false
}

func (s varSet) has(v *types.Var) bool {
	_, ok := s[v]
	return ok
}

// toSlice returns a deterministically ordered slice of the vars.
// Ordered by the order that the vars were added to the set.
func (s varSet) toSlice() []*types.Var {
	out := make([]*types.Var, len(s))
	for v, i := range s {
		out[i] = v
	}
	return out
}

type refDef struct {
	ref varSet // referenced (reads)
	def varSet // defined (writes)
}

func newRefDef() *refDef {
	return &refDef{
		ref: varSet{},
		def: varSet{},
	}
}

func (rd *refDef) addRead(v *ir.VarRef) {
	// For a reference, we don't want to add a reference if the variable
	// was defined in this block or assigned, however if the variable was
	// referenced before the definition then we still want to keep it as
	// a reference.
	if !rd.def.has(v.VarObj) {
		rd.ref.add(v.VarObj)
	}
}

func (rd *refDef) addWrite(n ir.Expr) {
	switch v := n.(type) {
	case *ir.VarDef:
		rd.def.add(v.VarObj)
	case *ir.VarRef:
		rd.def.add(v.VarObj)
	}
	// otherwise, continue so that this expression is walked and
	// the variables are treated as a read since it is like
	// `*s=`, `s[i]=`, or `s.value=`.
}

// compute walks the given statements in evaluation order and
// computes the set of vars referenced (read before being defined locally)
// and defined (assigned to or declared) within them.
func (rd *refDef) compute[T ir.Node](roots ...T) {
	for w := range ir.Walk(roots...) {
		switch n := w.Node.(type) {
		case *ir.VarRef:
			rd.addRead(n)

		case *ir.MultiAssignStmt:
			// Walk the rhs first to collect references before assignment.
			w.Skip = true
			rd.compute(n.Rhs...)
			for _, ln := range n.Lhs {
				rd.addWrite(ln)
			}
			rd.compute(n.Lhs...)

		case *ir.BinaryExpr:
			if n.Op.IsAssignment() {
				// Walk the rhs first to collect references before assignment.
				w.Skip = true
				rd.compute(n.Y)
				// Treat `x = y` as a write to `x` but other assignments of `x`,
				// e.g. `x += y`, as a reference and write to `x`,
				// the same as if it was `x = x + y`.
				if n.Op != binaryOp.Assign {
					if expr, ok := n.X.(*ir.VarRef); ok {
						rd.addRead(expr)
					}
				}
				rd.addWrite(n.X)
				rd.compute(n.X)
			}
		}
	}
}

// successors returns the unique successor blocks reachable from the given block
// by any GotoBlockStmt, FuncCallStmt.Follow, etc anywhere in its body.
func successors(b *ir.Block) []*ir.Block {
	return iterator.Unique(
		iterator.NotZero(
			ir.WalkNodes(b).OfType[*ir.BlockRef]().
				Select(func(ref *ir.BlockRef) *ir.Block { return ref.Block }),
		),
	).ToSlice()
}

// propagateParams runs the live-variable analysis on the function's
// block graph and assigns Block.Params and BlockRef.Args so that every
// block receives exactly the variables it needs (transitively through
// its successors).
func propagateParams(fn *ir.Func, errGroup *faults.ErrGroup) {
	if fn == nil || len(fn.Blocks) <= 0 {
		return
	}

	type blockInfo struct {
		block *ir.Block
		index int
		ref   varSet
		def   varSet
		suc   []*ir.Block
	}

	infoByBlock := make(map[*ir.Block]*blockInfo, len(fn.Blocks))
	infoByIndex := make([]*blockInfo, len(fn.Blocks))
	for i, b := range fn.Blocks {
		rd := newRefDef()
		rd.compute(b.Body...)
		info := &blockInfo{
			block: b,
			index: i,
			ref:   rd.ref,
			def:   rd.def,
			suc:   successors(b),
		}
		infoByBlock[b] = info
		infoByIndex[i] = info
	}

	for changed := true; changed; {
		changed = false
		for _, info := range infoByIndex {
			for _, s := range info.suc {
				suc := infoByBlock[s]
				for o := range suc.ref {
					if !info.def.has(o) {
						changed = info.ref.add(o) || changed
					}
				}
			}
		}
	}

	/*
		for i, b := range fn.Blocks {
			if i == 0 {
				// Initial block params are the function's external interface.
				existing := paramObjectSet(b.Params)
				for o := range liveIn[b] {
					if !existing.has(o) {
						errGroup.Add(faults.New(`function block has free variable not declared as parameter`).
							With(`function`, fn.Name).
							With(`variable`, o.Name()))
					}
				}
				continue
			}
			ordered := orderedVars(liveIn[b])
			newParams := make([]*ir.Param, 0, len(ordered))
			for _, o := range ordered {
				newParams = append(newParams, makeParam(o))
			}
			b.Params = newParams
		}

		// Rebuild Args at every jump site so each ref matches its target's Params.
		for _, b := range fn.Blocks {
			forEachJumpTarget(b, func(ref *ir.BlockRef, srcPos token.Pos) {
				target := ref.Block
				if target == nil {
					ref.Args = nil
					return
				}
				newArgs := make([]ir.Expr, 0, len(target.Params))
				for _, p := range target.Params {
					obj := info.ObjectOf(p.Name)
					if obj == nil {
						continue
					}
					newArgs = append(newArgs, makeArg(obj, srcPos))
				}
				ref.Args = newArgs
			})
		}
	*/
}
