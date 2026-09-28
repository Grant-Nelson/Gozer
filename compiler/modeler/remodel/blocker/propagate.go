package blocker

import (
	"go/token"
	"go/types"
	"maps"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

// objectSet is a set of types.Objects.
//
// The value is an monotonically incrementing value (i.e. the length
// of the map plus one). The values must be unique and continuous.
// Those values are used to get the order of the keys in the order
// the keys were added to the set.
type objectSet map[types.Object]int

func newObjectSet() objectSet { return objectSet{} }

func (s objectSet) add(o types.Object) {
	if o != nil {
		s[o] = len(s) + 1
	}
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
		if obj := n.Object(); !def.has(obj) {
			ref.add(obj)
		}
	}

	addWrite := func(w *ir.WalkStep, n ir.Expr) {
		switch v := n.(type) {
		case *ir.VarDef:
			def.add(v.Object())
			w.Skip(v)
		case *ir.VarRef:
			def.add(v.Object())
			w.Skip(v)
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
			if n.Op.IsAssignment() {
				addWrite(w, n.X)
			}
		}
	}
	return
}

// successors returns the unique successor blocks reachable from the given block
// by any GotoBlockStmt, FuncCallStmt.Follow, etc anywhere in its body.
func successors(b *ir.Block) []*ir.Block {
	seen := map[*ir.Block]bool{}
	var out []*ir.Block
	forEachJumpTarget(b, func(ref *ir.BlockRef, _ token.Pos) {
		target := ref.Block
		if target != nil && !seen[target] {
			seen[target] = true
			out = append(out, target)
		}
	})
	return out
}

// forEachJumpTarget invokes fn for every BlockRef appearing in b's body,
// recursing into nested statements.
func forEachJumpTarget(b *ir.Block, fn func(ref *ir.BlockRef, srcPos token.Pos)) {
	walkBlockRefs(b.Body, fn)
}

// TODO: SEE IF THIS CAN USE `WalkNodes`
func walkBlockRefs(stmts []ir.Stmt, fn func(ref *ir.BlockRef, srcPos token.Pos)) {
	for _, s := range stmts {
		switch s := s.(type) {
		case *ir.GotoBlockStmt:
			if s.Block != nil {
				fn(s.Block, s.SrcPos)
			}
		case *ir.FuncCallStmt:
			if s.Follow != nil {
				fn(s.Follow, s.Pos())
			}
		case *ir.IfStmt:
			walkBlockRefs(s.Body, fn)
			walkBlockRefs(s.Else, fn)
		case *ir.StmtListStmt:
			walkBlockRefs(s.List, fn)
		case *ir.ForStmt:
			walkBlockRefs(s.Body, fn)
		case *ir.LabeledStmt:
			if s.Stmt != nil {
				walkBlockRefs([]ir.Stmt{s.Stmt}, fn)
			}
		}
	}
}

// paramObjectSet returns the set of types.Objects referred to by the
// given block params.
func paramObjectSet(params []*ir.Param) objectSet {
	out := newObjectSet()
	for _, p := range params {
		if p.Name == nil {
			continue
		}
		if obj := info.ObjectOf(p.Name); obj != nil {
			out.add(obj)
		}
	}
	return out
}

// propagateParams runs the live-variable analysis on the function's
// block graph and assigns Block.Params and BlockRef.Args so that every
// block receives exactly the variables it needs (transitively through
// its successors).
func propagateParams(fn *ir.Func, errGroup *faults.ErrGroup) {
	if fn == nil || len(fn.Blocks) == 0 {
		return
	}

	useMap := make(map[*ir.Block]objectSet, len(fn.Blocks))
	defMap := make(map[*ir.Block]objectSet, len(fn.Blocks))
	succMap := make(map[*ir.Block][]*ir.Block, len(fn.Blocks))
	liveIn := make(map[*ir.Block]objectSet, len(fn.Blocks))

	for _, b := range fn.Blocks {
		use, def := computeUseDef(b.Body)
		useMap[b] = use
		defMap[b] = def
		succMap[b] = successors(b)
		liveIn[b] = use.clone()
	}

	// Fixed-point: live_in(B) = use(B) ∪ (∪ live_in(S) for S ∈ succ(B)) − def(B)
	for changed := true; changed; {
		changed = false
		for i := len(fn.Blocks) - 1; i >= 0; i-- {
			b := fn.Blocks[i]
			def := defMap[b]
			newIn := useMap[b].clone()
			for _, s := range succMap[b] {
				for o := range liveIn[s] {
					if !def.has(o) {
						newIn.add(o)
					}
				}
			}
			if !newIn.equal(liveIn[b]) {
				liveIn[b] = newIn
				changed = true
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
