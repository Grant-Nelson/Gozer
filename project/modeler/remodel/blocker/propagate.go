package blocker

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/modeler/ir"
)

// TODO: Agent created this file, do a full check of the code since there
// are some things that need to be updated to use newer Go patterns.

// objectSet is a set of types.Objects.
type objectSet map[types.Object]struct{}

func newObjectSet() objectSet { return objectSet{} }

func (s objectSet) add(o types.Object) {
	if o != nil {
		s[o] = struct{}{}
	}
}

func (s objectSet) has(o types.Object) bool {
	_, ok := s[o]
	return ok
}

func (s objectSet) clone() objectSet {
	c := make(objectSet, len(s))
	for k := range s {
		c[k] = struct{}{}
	}
	return c
}

func (s objectSet) equal(o objectSet) bool {
	if len(s) != len(o) {
		return false
	}
	for k := range s {
		if _, ok := o[k]; !ok {
			return false
		}
	}
	return true
}

// orderedObjects returns a deterministically ordered slice of the objects.
// Ordered by declaration position then by name.
func orderedObjects(s objectSet) []types.Object {
	out := make([]types.Object, 0, len(s))
	for o := range s {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Pos() != out[j].Pos() {
			return out[i].Pos() < out[j].Pos()
		}
		return out[i].Name() < out[j].Name()
	})
	return out
}

// computeUseDef walks the given statements in evaluation order and
// computes the set of objects used (read before being defined locally)
// and defined (assigned to or declared) within them.
func computeUseDef(stmts []ir.Stmt, info *types.Info) (use, def objectSet) {
	use = newObjectSet()
	def = newObjectSet()
	for _, s := range stmts {
		visitStmtIdents(s, info, use, def)
	}
	return use, def
}

// visitStmtIdents walks a statement, populating use and def in
// evaluation order.
func visitStmtIdents(s ir.Stmt, info *types.Info, use, def objectSet) {
	if s == nil {
		return
	}
	switch s := s.(type) {
	case *ir.AssignStmt:
		for _, e := range s.Rhs {
			visitExprIdents(e, info, use, def)
		}
		// For compound assigns the LHS is also read before being written.
		if s.Tok != token.ASSIGN && s.Tok != token.DEFINE {
			for _, e := range s.Lhs {
				visitExprIdents(e, info, use, def)
			}
		}
		for _, e := range s.Lhs {
			markLhsDef(e, info, use, def)
		}
	case *ir.ExprStmt:
		visitExprIdents(s.X, info, use, def)
	case *ir.SendStmt:
		visitExprIdents(s.Chan, info, use, def)
		visitExprIdents(s.Value, info, use, def)
	case *ir.ReturnStmt:
		for _, e := range s.Results {
			visitExprIdents(e, info, use, def)
		}
	case *ir.IfStmt:
		visitStmtIdents(s.Init, info, use, def)
		visitExprIdents(s.Cond, info, use, def)
		for _, b := range s.Body {
			visitStmtIdents(b, info, use, def)
		}
		for _, b := range s.Else {
			visitStmtIdents(b, info, use, def)
		}
	case *ir.StmtListStmt:
		for _, b := range s.List {
			visitStmtIdents(b, info, use, def)
		}
	case *ir.ForStmt:
		visitStmtIdents(s.Init, info, use, def)
		visitExprIdents(s.Cond, info, use, def)
		for _, b := range s.Body {
			visitStmtIdents(b, info, use, def)
		}
		visitStmtIdents(s.Post, info, use, def)
	case *ir.LabeledStmt:
		visitStmtIdents(s.Stmt, info, use, def)
	case *ir.FuncCallStmt:
		visitExprIdents(s.Fun, info, use, def)
		for _, e := range s.Args {
			visitExprIdents(e, info, use, def)
		}
		// Follow.Args are synthesized by the blocker and intentionally skipped;
		// they will be rederived from the successor's Params.
	case *ir.GotoBlockStmt, *ir.BranchStmt, *ir.DeclStmt:
		// GotoBlockStmt.Block.Args are synthesized by the blocker and
		// intentionally skipped. BranchStmt has no variable refs.
		// DeclStmt is deferred for now.
	}
}

// visitExprIdents walks an expression and records identifier reads
// into use (unless they have already been defined locally).
func visitExprIdents(e ast.Expr, info *types.Info, use, def objectSet) {
	if e == nil {
		return
	}
	ast.Inspect(e, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		obj := info.Uses[id]
		if obj == nil {
			return true
		}
		if _, isVar := obj.(*types.Var); !isVar {
			return true
		}
		if !def.has(obj) {
			use.add(obj)
		}
		return true
	})
}

// markLhsDef marks an LHS expression as defining or redefining variables.
// Non-identifier LHS expressions (e.g. *p, a[i], s.f) contribute reads
// of their target.
func markLhsDef(e ast.Expr, info *types.Info, use, def objectSet) {
	id, ok := e.(*ast.Ident)
	if !ok {
		visitExprIdents(e, info, use, def)
		return
	}
	if obj := info.Defs[id]; obj != nil {
		if _, isVar := obj.(*types.Var); isVar {
			def.add(obj)
		}
		return
	}
	if obj := info.Uses[id]; obj != nil {
		if _, isVar := obj.(*types.Var); isVar {
			def.add(obj)
		}
	}
}

// successors returns the unique successor blocks reachable from b by
// any GotoBlockStmt or FuncCallStmt.Follow anywhere in its body.
func successors(b *ir.Block) []*ir.Block {
	seen := map[*ir.Block]bool{}
	var out []*ir.Block
	forEachJumpTarget(b, func(ref *ir.BlockRef, _ token.Pos) {
		target := ref.Block
		if target == nil || seen[target] {
			return
		}
		seen[target] = true
		out = append(out, target)
	})
	return out
}

// forEachJumpTarget invokes fn for every BlockRef appearing in b's body,
// recursing into nested statements.
func forEachJumpTarget(b *ir.Block, fn func(ref *ir.BlockRef, srcPos token.Pos)) {
	walkBlockRefs(b.Body, fn)
}

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
func paramObjectSet(params []*ir.Param, info *types.Info) objectSet {
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

// makeParam creates a synthetic block parameter for the given object.
// The synthetic identifier is registered in info.Defs.
func makeParam(obj types.Object, info *types.Info) *ir.Param {
	id := &ast.Ident{Name: obj.Name()}
	info.Defs[id] = obj
	return &ir.Param{
		Name: id,
		Type: obj.Type(),
	}
}

// makeArg creates a synthetic argument expression referring to the given
// object. The synthetic identifier is registered in info.Uses.
func makeArg(obj types.Object, srcPos token.Pos, info *types.Info) ast.Expr {
	id := &ast.Ident{NamePos: srcPos, Name: obj.Name()}
	info.Uses[id] = obj
	return id
}

// propagateParams runs the live-variable analysis on the function's
// block graph and assigns Block.Params and BlockRef.Args so that every
// block receives exactly the variables it needs (transitively through
// its successors).
func propagateParams(fn *ir.Func, info *types.Info, errGroup *faults.ErrGroup) {
	if fn == nil || len(fn.Blocks) == 0 {
		return
	}

	useMap := make(map[*ir.Block]objectSet, len(fn.Blocks))
	defMap := make(map[*ir.Block]objectSet, len(fn.Blocks))
	succMap := make(map[*ir.Block][]*ir.Block, len(fn.Blocks))
	liveIn := make(map[*ir.Block]objectSet, len(fn.Blocks))

	for _, b := range fn.Blocks {
		use, def := computeUseDef(b.Body, info)
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
			existing := paramObjectSet(b.Params, info)
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
			newParams = append(newParams, makeParam(o, info))
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
			newArgs := make([]ast.Expr, 0, len(target.Params))
			for _, p := range target.Params {
				obj := info.ObjectOf(p.Name)
				if obj == nil {
					continue
				}
				newArgs = append(newArgs, makeArg(obj, srcPos, info))
			}
			ref.Args = newArgs
		})
	}
}
