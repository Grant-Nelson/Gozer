package converter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
)

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

// TODO: Add `//go:linkname` information for declarations

type Converter struct {
	Info    *types.Info
	FileSet *token.FileSet
}

func (c *Converter) pos(p token.Pos) string {
	return c.FileSet.Position(p).String()
}

func (c *Converter) exprTypes(e ast.Expr) types.TypeAndValue {
	tv, ok := c.Info.Types[e]
	if !ok {
		panic(fmt.Errorf(`failed to get type for %T at %s`, e, c.pos(e.Pos())))
	}
	return tv
}

func (c *Converter) FromNodeSlice(ns []ast.Node) []ir.Node {
	result := make([]ir.Node, 0, len(ns))
	for _, n := range ns {
		result = append(result, c.FromNode(n))
	}
	return result
}

func (c *Converter) FromNode(n ast.Node) ir.Node {
	switch n := n.(type) {
	case nil:
		return nil
	case ast.Decl:
		return c.FromDecl(n)
	case ast.Stmt:
		return c.FromStmt(n)
	case ast.Expr:
		return c.FromExpr(n)
	default:
		panic(faults.New(`unexpected AST node type`).
			WithF(`type`, `%T`, n).
			With(`pos`, c.pos(n.Pos())))
	}
}

func (c *Converter) FromDecl(d ast.Decl) ir.Stmt {
	switch d := d.(type) {
	case nil, *ast.BadDecl:
		return nil
	case *ast.GenDecl:
		return c.FromGenDecl(d)
	case *ast.FuncDecl:
		return c.FromFuncDecl(d)
	default:
		panic(faults.New(`unexpected AST decl node type`).
			WithF(`type`, `%T`, d).
			With(`pos`, c.pos(d.Pos())))
	}
}

func (c *Converter) FromGenDecl(d *ast.GenDecl) ir.Stmt {
	constant := d.Tok == token.CONST
	ss := &ir.StmtListStmt{}
	for _, s := range d.Specs {
		if v, ok := s.(*ast.ValueSpec); ok {
			ss.List = append(ss.List, c.FromValueSpec(v, constant))
		}
	}
	return c.SimplifyStmt(ss)
}

func (c *Converter) FromValueSpec(s *ast.ValueSpec, constant bool) *ir.StmtListStmt {
	ss := &ir.StmtListStmt{}
	for i, n := range s.Names {
		v := &ir.ValueDecl{
			Constant: constant,
			Name:     c.FromIdent(n),
		}
		if len(s.Values) > i {
			v.Value = c.FromExpr(s.Values[i])
		}
		ss.List = append(ss.List, v)
	}
	return ss
}

func (c *Converter) FromStmtSlice(ss []ast.Stmt) []ir.Stmt {
	result := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, c.ExpandStmt(c.FromStmt(s))...)
	}
	return result
}

func (c *Converter) SimplifyStmt(s ir.Stmt) ir.Stmt {
	if b, ok := s.(*ir.StmtListStmt); ok {
		b.List = iterator.NotZero(iterator.Iterate(c.ExpandStmtSlice(b.List)...)).ToSlice()
		switch len(b.List) {
		case 0:
			return nil
		case 1:
			return b.List[0]
		default:
			return b
		}
	}
	return s
}

func (c *Converter) ExpandStmtSlice(ss []ir.Stmt) []ir.Stmt {
	st := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		for x := range iterator.NotZero(iterator.Iterate(c.ExpandStmt(s)...)) {
			st = append(st, x)
		}
	}
	return st
}

func (c *Converter) ExpandStmt(s ir.Stmt) []ir.Stmt {
	if s == nil {
		return []ir.Stmt{}
	}
	if b, ok := s.(*ir.StmtListStmt); ok {
		return c.ExpandStmtSlice(b.List)
	}
	return []ir.Stmt{s}
}

func (c *Converter) FromStmt(s ast.Stmt) ir.Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return c.FromDeclStmt(s)
	case *ast.LabeledStmt:
		return c.FromLabeledStmt(s)
	case *ast.ExprStmt:
		return c.FromExprStmt(s)
	case *ast.SendStmt:
		return c.FromSendStmt(s)
	case *ast.IncDecStmt:
		return c.FromIncDecStmt(s)
	case *ast.AssignStmt:
		return c.FromAssignStmt(s)
	case *ast.GoStmt:
		return c.FromGoStmt(s)
	case *ast.DeferStmt:
		return c.FromDeferStmt(s)
	case *ast.ReturnStmt:
		return c.FromReturnStmt(s)
	case *ast.BranchStmt:
		return c.FromBranchStmt(s)
	case *ast.BlockStmt:
		return c.FromBlockStmt(s)
	case *ast.IfStmt:
		return c.FromIfStmt(s)
	case *ast.CaseClause:
		return c.FromCaseClause(s)
	case *ast.SwitchStmt:
		return c.FromSwitchStmt(s)
	case *ast.TypeSwitchStmt:
		return c.FromTypeSwitchStmt(s)
	case *ast.CommClause:
		return c.FromCommClause(s)
	case *ast.SelectStmt:
		return c.FromSelectStmt(s)
	case *ast.ForStmt:
		return c.FromForStmt(s)
	case *ast.RangeStmt:
		return c.FromRangeStmt(s)
	default:
		panic(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s).
			With(`pos`, c.pos(s.Pos())))
	}
}

func (c *Converter) FromExprSlice(es []ast.Expr) []ir.Expr {
	result := make([]ir.Expr, 0, len(es))
	for _, e := range es {
		result = append(result, c.FromExpr(e))
	}
	return result
}

func (c *Converter) FromExpr(e ast.Expr) ir.Expr {
	switch e := e.(type) {
	case nil, *ast.BadExpr:
		return nil
	case *ast.Ident:
		return c.FromIdent(e)
	case *ast.Ellipsis:
		return c.FromEllipsis(e)
	case *ast.BasicLit:
		return c.FromBasicLit(e)
	// TODO: FIX by making func lit have an expression reference to the func
	//case *ast.FuncLit:
	//	return c.FromFuncLit(e)
	case *ast.CompositeLit:
		return c.FromCompositeLit(e)
	case *ast.ParenExpr:
		return c.FromParenExpr(e)
	case *ast.SelectorExpr:
		return c.FromSelectorExpr(e)
	case *ast.IndexExpr:
		return c.FromIndexExpr(e)
	case *ast.IndexListExpr:
		return c.FromIndexListExpr(e)
	case *ast.SliceExpr:
		return c.FromSliceExpr(e)
	case *ast.TypeAssertExpr:
		return c.FromTypeAssertExpr(e)
	case *ast.CallExpr:
		return c.FromCallExpr(e)
	case *ast.StarExpr:
		return c.FromStarExpr(e)
	case *ast.UnaryExpr:
		return c.FromUnaryExpr(e)
	case *ast.BinaryExpr:
		return c.FromBinaryExpr(e)
	case *ast.KeyValueExpr:
		return c.FromKeyValueExpr(e)
	case *ast.ArrayType:
		return c.FromArrayType(e)
	case *ast.StructType:
		return c.FromStructType(e)
	case *ast.FuncType:
		return c.FromFuncType(e)
	case *ast.InterfaceType:
		return c.FromInterfaceType(e)
	case *ast.MapType:
		return c.FromMapType(e)
	case *ast.ChanType:
		return c.FromChanType(e)
	default:
		panic(faults.New(`unexpected AST expression type`).
			WithF(`type`, `%T`, e).
			With(`pos`, c.pos(e.Pos())))
	}
}

func (c *Converter) FromAssignStmt(s *ast.AssignStmt) *ir.AssignStmt {
	if s == nil {
		return nil
	}
	return &ir.AssignStmt{
		TokPos: s.TokPos,
		Lhs:    c.FromExprSlice(s.Lhs),
		Define: s.Tok == token.DEFINE,
		Rhs:    c.FromExprSlice(s.Rhs),
	}
}

func (c *Converter) FromBranchToken(t token.Token) branchKind.BranchKind {
	switch t {
	case token.BREAK:
		return branchKind.Break
	case token.CONTINUE:
		return branchKind.Continue
	case token.GOTO:
		return branchKind.Goto
	case token.FALLTHROUGH:
		return branchKind.Fallthrough
	default:
		panic(faults.New(`unexpected token for a branch kind`).
			With(`token`, t.String()))
	}
}

func (c *Converter) FromBranchStmt(s *ast.BranchStmt) *ir.BranchStmt {
	if s == nil {
		return nil
	}
	return &ir.BranchStmt{
		TokPos: s.TokPos,
		Kind:   c.FromBranchToken(s.Tok),
		Label:  s.Label,
	}
}

func (c *Converter) FromCaseClause(s *ast.CaseClause) *ir.CaseClause {
	if s == nil {
		return nil
	}
	return &ir.CaseClause{
		Case: s.Case,
		List: c.FromExprSlice(s.List),
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCaseClauseSlice(s *ast.BlockStmt) []*ir.CaseClause {
	if s == nil {
		return []*ir.CaseClause{}
	}
	ccs := make([]*ir.CaseClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CaseClause)
		if !ok {
			panic(faults.New(`expected case clause`).
				WithF(`type`, `%T`, ct).
				With(`pos`, c.pos(ct.Pos())).
				With(`index`, i))
		}
		ccs = append(ccs, c.FromCaseClause(cc))
	}
	return ccs
}

func (c *Converter) FromCommClause(s *ast.CommClause) *ir.CommClause {
	if s == nil {
		return nil
	}
	return &ir.CommClause{
		Case: s.Case,
		Comm: c.FromStmt(s.Comm),
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCommClauseSlice(s *ast.BlockStmt) []*ir.CommClause {
	if s == nil {
		return []*ir.CommClause{}
	}
	ccs := make([]*ir.CommClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CommClause)
		if !ok {
			panic(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, ct).
				With(`pos`, c.pos(ct.Pos())).
				With(`index`, i))
		}
		ccs = append(ccs, c.FromCommClause(cc))
	}
	return ccs
}

func (c *Converter) FromDeclStmt(s *ast.DeclStmt) ir.Stmt {
	if s == nil {
		return nil
	}
	return c.FromDecl(s.Decl)
}

func (c *Converter) FromDeferStmt(s *ast.DeferStmt) *ir.DeferStmt {
	if s == nil {
		return nil
	}
	return &ir.DeferStmt{
		Defer: s.Defer,
		Call:  c.FromCallExpr(s.Call),
	}
}

func (c *Converter) FromExprStmt(s *ast.ExprStmt) *ir.ExprStmt {
	if s == nil {
		return nil
	}
	return &ir.ExprStmt{
		X: c.FromExpr(s.X),
	}
}

func (c *Converter) FromIncDecStmt(s *ast.IncDecStmt) *ir.ExprStmt {
	if s == nil {
		return nil
	}
	return &ir.ExprStmt{
		X: &ir.UnaryExpr{
			OpPos: s.Pos(),
			Op:    s.Tok,
			X:     c.FromExpr(s.X),
		},
	}
}

func (c *Converter) FromForStmt(s *ast.ForStmt) *ir.ForStmt {
	if s == nil {
		return nil
	}
	return &ir.ForStmt{
		ForPos: s.For,
		Init:   c.FromStmt(s.Init),
		Cond:   c.FromExpr(s.Cond),
		Post:   c.FromStmt(s.Post),
		Body:   c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromGoStmt(s *ast.GoStmt) *ir.GoStmt {
	if s == nil {
		return nil
	}
	return &ir.GoStmt{
		GoPos: s.Go,
		Call:  c.FromCallExpr(s.Call),
	}
}

func (c *Converter) FromIfStmt(s *ast.IfStmt) *ir.IfStmt {
	if s == nil {
		return nil
	}
	return &ir.IfStmt{
		IfPos: s.If,
		Init:  c.FromStmt(s.Init),
		Cond:  c.FromExpr(s.Cond),
		Body:  c.ExpandStmt(c.FromBlockStmt(s.Body)),
		Else:  c.ExpandStmt(c.FromStmt(s.Else)),
	}
}

func (c *Converter) FromLabeledStmt(s *ast.LabeledStmt) *ir.LabeledStmt {
	if s == nil {
		return nil
	}
	return &ir.LabeledStmt{
		Label: c.FromIdent(s.Label),
		Stmt:  c.FromStmt(s.Stmt),
	}
}

func (c *Converter) ExpandParams(ft *ast.FuncType) []*ir.Param {
	if ft == nil {
		return nil
	}
	return c.ExpandFieldList(ft.Params)
}

func (c *Converter) ExpandResults(ft *ast.FuncType) []*ir.Param {
	if ft == nil {
		return nil
	}
	return c.ExpandFieldList(ft.Results)
}

func (c *Converter) ExpandFieldList(fl *ast.FieldList) []*ir.Param {
	if fl == nil || len(fl.List) <= 0 {
		return nil
	}
	params := make([]*ir.Param, 0, len(fl.List))
	for _, field := range fl.List {
		// If field.Names is nil then this is an unnamed field
		// that should be skipped over when defining block params.
		for _, name := range field.Names {
			if name == nil {
				continue
			}
			param := &ir.Param{
				Name: c.FromIdent(name),
				Expr: c.FromExpr(field.Type),
			}
			params = append(params, param)
		}
	}
	return params
}

func (c *Converter) FromRangeStmt(s *ast.RangeStmt) *ir.RangeStmt {
	if s == nil {
		return nil
	}
	return &ir.RangeStmt{
		ForPos: s.For,
		Key:    c.FromExpr(s.Key),
		Value:  c.FromExpr(s.Value),
		X:      c.FromExpr(s.X),
		Body:   c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromReturnStmt(s *ast.ReturnStmt) *ir.ReturnStmt {
	if s == nil {
		return nil
	}
	return &ir.ReturnStmt{
		ReturnPos: s.Return,
		Results:   c.FromExprSlice(s.Results),
	}
}

func (c *Converter) FromSelectStmt(s *ast.SelectStmt) *ir.SelectStmt {
	if s == nil {
		return nil
	}
	return &ir.SelectStmt{
		SelectPos: s.Select,
		Body:      c.FromCommClauseSlice(s.Body),
	}
}

func (c *Converter) FromSendStmt(s *ast.SendStmt) *ir.SendStmt {
	if s == nil {
		return nil
	}
	return &ir.SendStmt{
		ArrowPos: s.Arrow,
		Chan:     c.FromExpr(s.Chan),
		Value:    c.FromExpr(s.Value),
	}
}

func (c *Converter) FromBlockStmt(s *ast.BlockStmt) *ir.StmtListStmt {
	if s == nil {
		return nil
	}
	return &ir.StmtListStmt{
		List: c.FromStmtSlice(s.List),
	}
}

func (c *Converter) FromSwitchStmt(s *ast.SwitchStmt) *ir.SwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.SwitchStmt{
		SwitchPos: s.Switch,
		Init:      c.FromStmt(s.Init),
		Tag:       c.FromExpr(s.Tag),
		Body:      c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromTypeSwitchStmt(s *ast.TypeSwitchStmt) *ir.TypeSwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.TypeSwitchStmt{
		SwitchPos: s.Switch,
		Init:      c.FromStmt(s.Init),
		Assign:    c.FromStmt(s.Assign),
		Body:      c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromIdent(e *ast.Ident) *ir.Ident {
	if e == nil {
		return nil
	}
	id := &ir.Ident{
		NamePos: e.NamePos,
		Name:    e.Name,
	}
	if tv, ok := c.Info.Types[e]; ok {
		id.TypeAndValue = &tv
	}
	if in, ok := c.Info.Instances[e]; ok {
		id.Instance = &in
	}
	if ds, ok := c.Info.Defs[e]; ok {
		id.Def = ds
	}
	if us, ok := c.Info.Uses[e]; ok {
		id.Use = us
	}
	return id
}

func (c *Converter) FromEllipsis(e *ast.Ellipsis) ir.Expr {
	if e == nil {
		return nil
	}
	return c.FromExpr(e.Elt)
}

func (c *Converter) FromBasicLit(e *ast.BasicLit) *ir.BasicLit {
	if e == nil {
		return nil
	}
	return &ir.BasicLit{
		ValuePos:     e.ValuePos,
		TypeAndValue: c.exprTypes(e),
	}
}

func (c *Converter) FromFuncDecl(astFunc *ast.FuncDecl) *ir.Func {
	if astFunc == nil {
		return nil
	}
	fn := &ir.Func{
		FuncPos: astFunc.Type.Func,
		Name:    c.FromIdent(astFunc.Name),
	}
	c.setInitialBlock(fn, astFunc.Body, astFunc.Type)

	if astFunc.Doc != nil {
		dv := astTools.Directives(astFunc.Doc.List, directiveGroup)
		if s, ok := dv[directiveAtomicFunc]; ok {
			// The atomic directive should have no fields.
			// Any fields will be ignored (if asserts are off).
			assert.EmptySlice(s)
			fn.Atomic = true
		}
	}

	return fn
}

func (c *Converter) FromFuncLit(astFunc *ast.FuncLit) *ir.Func {
	if astFunc == nil {
		return nil
	}
	pos := astFunc.Type.Func
	fn := &ir.Func{
		FuncPos: pos,
	}
	c.setInitialBlock(fn, astFunc.Body, astFunc.Type)
	return fn
}

// setInitialBlock creates an initial block, populate it with current statements, add params.
// Named return values are appended after input params so they are in scope
// for the function body just like declared locals.
func (c *Converter) setInitialBlock(fn *ir.Func, fnBody *ast.BlockStmt, fnType *ast.FuncType) *ir.Block {
	body := c.ExpandStmt(c.FromBlockStmt(fnBody))
	params := append(c.ExpandParams(fnType), c.ExpandResults(fnType)...)
	return fn.NewBlock(`initial`, body, params)
}

func (c *Converter) FromCompositeLit(e *ast.CompositeLit) *ir.TypeExpr {
	// TODO: Need to store the initial values
	return c.FromTypeExpr(e)
}

func (c *Converter) FromParenExpr(e *ast.ParenExpr) ir.Expr {
	if e == nil {
		return nil
	}
	// Drop the paren since the order of operators is already kept in the tree shape.
	return c.FromExpr(e.X)
}

func (c *Converter) FromSelectorExpr(e *ast.SelectorExpr) *ir.SelectorExpr {
	if e == nil {
		return nil
	}
	return &ir.SelectorExpr{
		X:       c.FromExpr(e.X),
		Sel:     c.FromIdent(e.Sel),
		SelType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromIndexExpr(e *ast.IndexExpr) *ir.IndexExpr {
	if e == nil {
		return nil
	}
	// TODO: Should determine if this is for a string, slice, or map
	//       but if it is for generics then it should be an index list.
	return &ir.IndexExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Index:      c.FromExpr(e.Index),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromIndexListExpr(e *ast.IndexListExpr) *ir.IndexListExpr {
	if e == nil {
		return nil
	}
	return &ir.IndexListExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Indices:    c.FromExprSlice(e.Indices),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromSliceExpr(e *ast.SliceExpr) *ir.SliceExpr {
	if e == nil {
		return nil
	}
	return &ir.SliceExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Low:        c.FromExpr(e.Low),
		High:       c.FromExpr(e.High),
		Max:        c.FromExpr(e.Max),
		Slice3:     e.Slice3,
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromTypeAssertExpr(e *ast.TypeAssertExpr) *ir.TypeAssertExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeAssertExpr{
		X:          c.FromExpr(e.X),
		LparenPos:  e.Lparen,
		AssertType: c.FromExpr(e.Type),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromCallExpr(e *ast.CallExpr) *ir.CallExpr {
	if e == nil {
		return nil
	}
	return &ir.CallExpr{
		Fun:          c.FromExpr(e.Fun),
		LeftParenPos: e.Lparen,
		Args:         c.FromExprSlice(e.Args),
		Expanded:     e.Ellipsis.IsValid(),
		ResultType:   c.exprTypes(e).Type,
	}
}

func (c *Converter) FromStarExpr(e *ast.StarExpr) ir.Expr {
	if e == nil {
		return nil
	}
	tv := c.exprTypes(e)
	if tv.IsType() {
		return &ir.TypeExpr{
			TypePos:      e.Star,
			TypeAndValue: tv,
		}
	}
	return &ir.UnaryExpr{
		OpPos:      e.Star,
		Op:         token.MUL,
		X:          c.FromExpr(e.X),
		ResultType: tv.Type,
	}
}

func (c *Converter) FromUnaryExpr(e *ast.UnaryExpr) *ir.UnaryExpr {
	if e == nil {
		return nil
	}
	return &ir.UnaryExpr{
		OpPos:      e.OpPos,
		Op:         e.Op,
		X:          c.FromExpr(e.X),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromBinaryExpr(e *ast.BinaryExpr) *ir.BinaryExpr {
	if e == nil {
		return nil
	}
	return &ir.BinaryExpr{
		X:          c.FromExpr(e.X),
		OpPos:      e.OpPos,
		Op:         e.Op,
		Y:          c.FromExpr(e.Y),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromKeyValueExpr(e *ast.KeyValueExpr) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromArrayType(e *ast.ArrayType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromStructType(e *ast.StructType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromFuncType(e *ast.FuncType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromInterfaceType(e *ast.InterfaceType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromMapType(e *ast.MapType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromChanType(e *ast.ChanType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromTypeExpr(e ast.Expr) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos:      e.Pos(),
		TypeAndValue: c.exprTypes(e),
	}
}
