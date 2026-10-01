package blocker

import (
	"go/token"
	"go/types"
	"maps"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/binaryOp"
)

// objectSet is a set of types.Objects.
//
// The value is an monotonically incrementing value (i.e. the length
// of the map plus one). The values must be unique and continuous.
// Those values are used to get the order of the keys in the order
// the keys were added to the set.
type objectSet map[types.Object]int

func newObjectSet() objectSet { return objectSet{} }

func (s objectSet) add(o types.Object) bool {
	if o != nil {
		if _, has := s[o]; has {
			return false
		}
		s[o] = len(s) + 1
		return true
	}
	return false
}

func (s objectSet) has(o types.Object) bool {
	return s[o] > 0
}

func (s objectSet) clone() objectSet {
	return maps.Clone(s)
}

func (s objectSet) equal(o objectSet) bool {
	return maps.Equal(s, o)
}

// orderedObjects returns a deterministically ordered slice of the objects.
// Ordered by the order that the objects were added to the set.
func orderedObjects(s objectSet) []types.Object {
	out := make([]types.Object, len(s))
	for k, v := range s {
		out[v] = k
	}
	return out
}

// computeRefDef walks the given statements in evaluation order and
// computes the set of objects referenced (read before being defined locally)
// and defined (assigned to or declared) within them.
func computeRefDef(s []ir.Stmt) (ref, def objectSet) {
	ref = newObjectSet() // referenced (reads)
	def = newObjectSet() // defined (writes)

	addRead := func(n *ir.VarRef) {
		// For a reference, we don't want to add a reference if the variable
		// was defined in this block or assigned, however if the variable was
		// referenced before the definition then we still want to keep it as
		// a reference.
		if obj := n.Object(); !def.has(obj) {
			ref.add(obj)
		}
	}

	addWrite := func(w *ir.WalkStep, n ir.Expr) {
		switch v := n.(type) {
		case *ir.VarDef:
			def.add(v.Object())
		case *ir.VarRef:
			def.add(v.Object())
		}
	}

	for w := range ir.Walk(s...) {
		switch n := w.Node.(type) {
		case *ir.VarRef:
			addRead(n)
		case *ir.MultiAssignStmt:
			for _, ln := range n.Lhs {
				addWrite(w, ln)
			}
		case *ir.BinaryExpr:
			if n.Op == binaryOp.Assign {
				// Treat `x = y` as a write to `x` but other assignments
				// of `x`, e.g. `x += y`, as a reference to `x`.
				addWrite(w, n.X)
			}
		}
	}
	return
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

	blocks := make(map[*ir.Block]int, len(fn.Blocks))
	refMap := make([]objectSet, len(fn.Blocks))
	defMap := make([]objectSet, len(fn.Blocks))
	sucMap := make([][]*ir.Block, len(fn.Blocks))
	liveIn := make([]objectSet, len(fn.Blocks))

	for i, b := range fn.Blocks {
		blocks[b] = i
		ref, def := computeRefDef(b.Body)
		refMap[i] = ref
		defMap[i] = def
		sucMap[i] = successors(b)
		liveIn[i] = ref.clone()
	}

	// Fixed-point: liveIn(B) = ref(B) ∪ (∪ liveIn(S) for S ∈ successors(B)) − def(B)
	for changed := true; changed; {
		changed = false
		for i := len(fn.Blocks) - 1; i >= 0; i-- {
			def := defMap[i]
			ref := refMap[i]
			for _, s := range sucMap[i] {
				suc := blocks[s]
				for o := range liveIn[suc] {
					if !def.has(o) {
						changed = ref.add(o) || changed
					}
				}
			}
		}
	}

	// Replace Params for every non-initial block from liveIn.
	for i, b := range fn.Blocks {
		if i == 0 {
			// Initial block params are the function's external interface.
			// Flag any extra live-in object as a free variable since
			// closures aren't yet supported.
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
		ordered := orderedObjects(liveIn[b])
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
}
