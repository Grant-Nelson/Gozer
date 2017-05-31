package transpiler

import (
	"go/ast"
	"go/token"
	"math"
	"reflect"
	"strconv"

	"github.com/grant-nelson/Gozer/constructs/expressions"
	"github.com/grant-nelson/Gozer/constructs/statements"
	"github.com/grant-nelson/Gozer/constructs/types"
	"github.com/grant-nelson/Gozer/msg"
)

// Source is the information about an input source file.
type Source struct {

	// log is the logger for reporting transpile issues with.
	log *msg.Logger

	// fileSet is the full set of loaded  files.
	fileSet *token.FileSet

	// Path is the location of this source file for debugging purposes.
	Path string

	// Data is the data for this source file.
	Data *ast.File

	// Imports whieh were local to this method.
	// The key is the short name or path for the import.
	Imports map[string]*types.PackageType

	// Package is the package which this method belongs to.
	Package *types.PackageType

	// pendingFuncs is the set of functions which need the body of the function read.
	pendingFuncs map[*types.FunctionType]*ast.BlockStmt
}

// NewSource creates a new source file descriptions.
func NewSource(log *msg.Logger, fileSet *token.FileSet) *Source {
	return &Source{
		log:          log,
		fileSet:      fileSet,
		Path:         "",
		Data:         nil,
		Imports:      map[string]*types.PackageType{},
		Package:      nil,
		pendingFuncs: map[*types.FunctionType]*ast.BlockStmt{},
	}
}

// ProcessTypes determines all the class signatures, interfaces, functions, and handles.
func (src *Source) ProcessTypes() {
	defer func() {
		if r := recover(); r != nil {
			src.log.Error("Error occurred while processing types: ", r)
		}
	}()
	scope := src.fillOutScope()
	for _, decl := range src.Data.Decls {
		switch data := decl.(type) {
		case *ast.GenDecl:
			src.readGenericDeclaration(scope, data)
		case *ast.FuncDecl:
			src.readFunctionType(scope, data)
		default:
			msg.ThrowError("Unhandled type declaration: ", decl, " (", reflect.TypeOf(decl), ")")
		}
	}
}

// ProcessBodies transpiles the bodies of the functions and expressions of constants.
func (src *Source) ProcessBodies() {
	defer func() {
		if r := recover(); r != nil {
			src.log.Error("Error occurred while processing bodies: ", r)
		}
	}()
	scope := src.fillOutScope()
	// TODO: fill out Constansts
	// TODO: fill out Variables
	for fn, body := range src.pendingFuncs {
		fn.Body = src.parseBlock(scope, body)
	}
	// TODO: fill out Library Functions
	// TODO: fill out Class Functions
}

// addBasedPackage adds a package at the base level, not under a name.
func (src *Source) addBasedPackage(scope *Scope, pack *types.PackageType) {
	for id, t := range pack.Imports {
		scope.Add(id, t)
	}
	for id, t := range pack.Declarations {
		scope.Add(id, t)
	}
	for id, t := range pack.Functions {
		scope.Add(id, t)
	}
	for id, t := range pack.Interfaces {
		scope.Add(id, t)
	}
	for id, t := range pack.Classes {
		scope.Add(id, t)
	}
}

// fillOutScope fills out the scope for the containing package.
func (src *Source) fillOutScope() *Scope {
	defer func() {
		if r := recover(); r != nil {
			msg.ThrowError("Error occurred while filling out the scrope: ", r)
		}
	}()
	scope := NewScope(nil)
	for name, pack := range src.Imports {
		if name == "builtin" {
			src.addBasedPackage(scope, pack)
		} else {
			scope.Add(name, pack)
		}
	}
	src.addBasedPackage(scope, src.Package)
	return scope
}

// readType reads a type from the given expression.
func (src *Source) readType(scope *Scope, desc ast.Expr) (types.Type, bool) {
	if desc == nil {
		msg.ThrowError("Nil type expression")
		return nil, false
	}
	switch id := desc.(type) {
	case *ast.ArrayType:
		if id.Len != nil {
			msg.ThrowError("Unhandled array length expression: ", id, " (", reflect.TypeOf(desc), ")")
		}
		desc, _ := src.readType(scope, id.Elt)
		return types.List(desc), true
	case *ast.Ident:
		return src.lookupType(scope, id.Name), false
	case *ast.Ellipsis:
		desc, _ := src.readType(scope, id.Elt)
		return types.List(desc), true
	default:
		msg.ThrowError("Unhandled type expression: ", desc, " (", reflect.TypeOf(desc), ")")
		return nil, false
	}
}

// lookupType gets the type for the given Go type name.
func (src *Source) lookupType(scope *Scope, typeName string) types.Type {
	result := types.LookupType(typeName)
	if result == nil {
		msg.ThrowError("Unhandled type name: ", typeName)
		return nil
	}
	return result
}

// readFieldList reads a list of parameters or returns from the given field.
func (src *Source) readFieldList(scope *Scope, fields *ast.FieldList) ([]string, []types.Type, bool) {
	names := []string{}
	types := []types.Type{}
	var typeDesc types.Type
	ellipsis := false
	if fields != nil {
		for _, field := range fields.List {
			typeDesc, ellipsis = src.readType(scope, field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					names = append(names, name.Name)
					types = append(types, typeDesc)
				}
			} else {
				names = append(names, "")
				types = append(types, typeDesc)
			}
		}
	}
	return names, types, ellipsis
}

// readGenericDeclaration reads the given generic declaration into the library.
func (src *Source) readGenericDeclaration(scope *Scope, data *ast.GenDecl) {
	defer func() {
		if r := recover(); r != nil {
			msg.ThrowError("Error occurred while reading a generic declaration: ", r)
		}
	}()
	switch data.Tok {
	case token.IMPORT:
		// Ignore imports while reading generic declarations.
		return
	default:
		msg.ThrowError("Unhandled generic declaration: ", data, " (", reflect.TypeOf(data), ")")
	}
}

// readFunction reads the given method into the library
// and adds it to a class if the class is defined.
func (src *Source) readFunctionType(scope *Scope, data *ast.FuncDecl) {
	defer func() {
		if r := recover(); r != nil {
			msg.ThrowError("Error occurred while reading a function type: ", r)
		}
	}()
	//ast.Print(src.fileSet, data)
	fn := types.Function()
	fn.Name = data.Name.Name

	// Receiver
	receiverNames, receiverTypes, _ := src.readFieldList(scope, data.Recv)
	if len(receiverNames) > 0 {
		fn.ReceiverName = receiverNames[0]
		class := receiverTypes[0].(*types.ClassType)
		class.Interface.Functions[fn.Name] = fn
		fn.ReceiverClass = class
	} else {
		src.Package.Functions[fn.Name] = fn
	}

	// Input parameters
	paramNames, paramTypes, ellipsis := src.readFieldList(scope, data.Type.Params)
	fn.ParamNames = paramNames
	fn.ParamTypes = paramTypes
	fn.Ellipsis = ellipsis

	// Return paramters
	returnNames, returnTypes, _ := src.readFieldList(scope, data.Type.Results)
	fn.ReturnNames = returnNames
	fn.ReturnTypes = returnTypes

	// Read the body for the function
	if data.Body != nil {
		src.pendingFuncs[fn] = data.Body
	}
}

// parseBlock reads a block statement.
// https://golang.org/pkg/go/ast/#BlockStmt
func (src *Source) parseBlock(scope *Scope, block *ast.BlockStmt) *statements.BlockStat {
	defer func() {
		if r := recover(); r != nil {
			msg.ThrowError("Error occurred while parsing a block: ", r)
		}
	}()
	blockScope := NewScope(scope)
	stats := []statements.Statement{}
	for _, stmt := range block.List {
		newStates := src.parseStatement(blockScope, stmt)
		stats = append(stats, newStates...)
	}
	return statements.Block(stats...)
}

// parseStatement reads a statement.
// https://golang.org/pkg/go/ast/#Stmt
func (src *Source) parseStatement(scope *Scope, statement ast.Stmt) []statements.Statement {
	switch stat := statement.(type) {
	case *ast.AssignStmt:
		return src.expSliceToStatSlice(src.parseAssignment(scope, stat))
	case *ast.ExprStmt:
		return []statements.Statement{src.parseExpression(scope, stat.X)}
	case *ast.IfStmt:
		return []statements.Statement{src.parseIfStatement(scope, stat)}
	case *ast.BlockStmt:
		return []statements.Statement{src.parseBlock(scope, stat)}
	case *ast.ForStmt:
		return []statements.Statement{src.parseForStatement(scope, stat)}
	case *ast.IncDecStmt:
		return []statements.Statement{src.parseIncDecStatement(scope, stat)}
	case *ast.BranchStmt:
		return []statements.Statement{src.parseBranchStatement(scope, stat)}

	// case *ast.ReturnStmt: dw.writeReturn(st)
	default:
		src.log.Error("Unhandled statement type ", reflect.TypeOf(stat))
		return nil
	}
}

// parseIfStatement reads in an if-statement.
// https://golang.org/pkg/go/ast/#IfStmt
func (src *Source) parseIfStatement(scope *Scope, ifStat *ast.IfStmt) statements.Statement {
	init := []statements.Statement{}
	if ifStat.Init != nil {
		init = src.parseStatement(scope, ifStat.Init)
	}

	cond := src.parseExpression(scope, ifStat.Cond)
	bodyStat := src.parseBlock(scope, ifStat.Body)
	var elseStat statements.Statement
	if ifStat.Else != nil {
		stats := src.parseStatement(scope, ifStat.Else)
		if len(stats) == 1 {
			elseStat = stats[0]
		} else {
			elseStat = statements.Block(stats...)
		}
	}

	if len(init) > 0 {
		init = append(init, statements.If(cond, bodyStat, elseStat))
		return statements.Block(init...)
	}
	return statements.If(cond, bodyStat, elseStat)
}

// parseForStatement reads a for-statment
// https://golang.org/pkg/go/ast/#ForStmt
func (src *Source) parseForStatement(scope *Scope, forStat *ast.ForStmt) statements.Statement {
	var inits []statements.Statement
	if forStat.Init != nil {
		inits = src.parseStatement(scope, forStat.Init)
	}
	var cond expressions.Expression
	if forStat.Cond != nil {
		cond = src.parseExpression(scope, forStat.Cond)
	}
	var post []statements.Statement
	if forStat.Post != nil {
		post = src.parseStatement(scope, forStat.Post)
	}
	body := src.parseBlock(scope, forStat.Body)
	var init statements.Statement
	if len(inits) == 1 {
		init = inits[0]
	}
	result := statements.For(init, cond, post, body)
	if len(inits) > 1 {
		return statements.Block(append(inits, result)...)
	}
	return result
}

// parseIncDecStatement reads an increment or decrement statment.
// https://golang.org/pkg/go/ast/#IncDecStmt
func (src *Source) parseIncDecStatement(scope *Scope, stmt *ast.IncDecStmt) statements.Statement {
	exp := src.parseExpression(scope, stmt.X)
	return statements.IncDecOp(exp, stmt.Tok == token.INC)
}

// parseBranchStatement reads a branch statment.
// https://golang.org/pkg/go/ast/#BranchStmt
func (src *Source) parseBranchStatement(scope *Scope, stmt *ast.BranchStmt) statements.Statement {
	return statements.Branch(stmt.Tok == token.BREAK)
}

// parseAssignment reads in an assigment statement.
// https://golang.org/pkg/go/ast/#AssignStmt
func (src *Source) parseAssignment(scope *Scope, assign *ast.AssignStmt) []expressions.Expression {
	if assign.Tok == token.DEFINE {
		return src.parseDefinition(scope, assign)
	}

	results := make([]expressions.Expression, 0, len(assign.Rhs))
	lefts := make([]expressions.Expression, len(assign.Lhs))
	for i := len(assign.Lhs) - 1; i >= 0; i-- {
		lefts[i] = src.parseExpression(scope, assign.Lhs[i])
	}

	if len(assign.Rhs) == 1 {
		right := src.parseExpression(scope, assign.Rhs[0])
		results = append(results, expressions.Assignment(lefts, right))
	} else {
		leftOffset := 0
		tempIDs := make([]*expressions.IdentifierExp, len(assign.Lhs))
		for _, exp := range assign.Rhs {
			right := src.parseExpression(scope, exp)
			rets := right.ReturnTypes()
			retCount := len(rets)
			if retCount == 1 {
				tempID := scope.AddTemp(rets[0])
				tempIDs[leftOffset] = tempID
				leftOffset++
				results = append(results, expressions.Definition(tempID, right))
			} else {
				leftExps := make([]expressions.Expression, retCount)
				for j := 0; j < retCount; j++ {
					tempID := scope.AddTemp(rets[j])
					leftExps[j] = tempID
					tempIDs[leftOffset] = tempID
					leftOffset++
					results = append(results, expressions.Definition(tempID, nil))
				}
				results = append(results, expressions.Assignment(leftExps, right))
			}
		}
		for i, tempID := range tempIDs {
			leftExps := []expressions.Expression{lefts[i]}
			results = append(results, expressions.Assignment(leftExps, tempID))
		}
	}
	return results
}

// parseDefinition reads in an assigment definition statement.
func (src *Source) parseDefinition(scope *Scope, assign *ast.AssignStmt) []expressions.Expression {
	leftOffset := 0
	results := make([]expressions.Expression, 0, len(assign.Rhs))
	for _, exp := range assign.Rhs {
		right := src.parseExpression(scope, exp)
		rets := right.ReturnTypes()
		retCount := len(rets)
		if retCount == 1 {
			leftExp := assign.Lhs[leftOffset]
			leftOffset++
			tempID := scope.Add(leftExp.(*ast.Ident).Name, rets[0])
			results = append(results, expressions.Definition(tempID, right))
		} else {
			leftExps := make([]expressions.Expression, retCount)
			for j := 0; j < retCount; j++ {
				leftExp := assign.Lhs[leftOffset]
				leftOffset++
				tempID := scope.Add(leftExp.(*ast.Ident).Name, rets[j])
				leftExps[j] = tempID
				results = append(results, expressions.Definition(tempID, nil))
			}
			results = append(results, expressions.Assignment(leftExps, right))
		}
	}
	return results
}

// parseExpression reads in some kind of expression.
// https://golang.org/pkg/go/ast/#Expr
func (src *Source) parseExpression(scope *Scope, expr ast.Expr) expressions.Expression {
	switch ex := expr.(type) {
	case *ast.BasicLit:
		return src.parseLiteral(scope, ex)
	case *ast.BinaryExpr:
		return src.parseBinary(scope, ex)
	case *ast.CallExpr:
		return src.parseCall(scope, ex)
	case *ast.CompositeLit:
		return src.parseCompositeLit(scope, ex)
	case *ast.Ident:
		return src.parseIdentifier(scope, ex)
	case *ast.IndexExpr:
		return src.parseIndexer(scope, ex)
	case *ast.ParenExpr:
		return src.parseExpression(scope, ex.X)
	case *ast.SelectorExpr:
		return src.parseSelector(scope, ex)
	case *ast.UnaryExpr:
		return src.parseUnary(scope, ex)
	default:
		src.log.Error("Unhandled expression type ", reflect.TypeOf(ex))
		return nil
	}
}

// parseLiteral reads a code literal.
// https://golang.org/pkg/go/ast/#BasicLit
// e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
func (src *Source) parseLiteral(scope *Scope, lit *ast.BasicLit) *expressions.LiteralExp {
	switch lit.Kind {
	case token.INT:
		val, err := strconv.ParseInt(lit.Value, 0, 64)
		if err != nil {
			src.log.Error("Error reading integer literal: ", err)
			return nil
		}
		if val > math.MaxInt32 {
			return expressions.Literal(lit.Value, types.Int64())
		}
		return expressions.Literal(lit.Value, types.Int())
	case token.FLOAT:
		return expressions.Literal(lit.Value, types.Float64())
	case token.IMAG:
		return expressions.Literal(lit.Value, types.Complex128())
	case token.CHAR:
		return expressions.Literal(lit.Value, types.Rune())
	case token.STRING:
		str, err := strconv.Unquote(lit.Value)
		if err != nil {
			src.log.Error("Unable to unquote literal string: ", lit.Value)
			return expressions.Literal(lit.Value, types.String())
		}
		str = strconv.QuoteToASCII(str)
		return expressions.Literal(str, types.String())
	default:
		src.log.Error("Unhandled literal kind ", lit.Kind)
		return nil
	}
}

// parseCompositeLit reads a composite literial value.
// https://golang.org/pkg/go/ast/#CompositeLit
func (src *Source) parseCompositeLit(scope *Scope, exp *ast.CompositeLit) *expressions.CompoundLiteralExp {
	litType, _ := src.readType(scope, exp.Type)
	elements := make([]expressions.Expression, len(exp.Elts))
	for i, elts := range exp.Elts {
		elements[i] = src.parseExpression(scope, elts)
	}
	return expressions.CompoundLiteral(elements, litType)
}

// parseBinary reads a binary operation.
// https://golang.org/pkg/go/ast/#BinaryExpr
func (src *Source) parseBinary(scope *Scope, bin *ast.BinaryExpr) *expressions.BinaryOpExp {
	left := src.parseExpression(scope, bin.X)
	right := src.parseExpression(scope, bin.Y)

	var resultType types.Type
	if left == nil {
		src.log.Debug("Left expression in a binary is nil ", bin.Op, ".")
		resultType = types.Variant()
	} else {
		rets := left.ReturnTypes()
		if len(rets) != 1 {
			src.log.Error("Binary expected only one return type for ", bin.Op, ".")
			resultType = types.Variant()
		} else {
			resultType = rets[0]
		}
	}

	// https://golang.org/pkg/go/token/#Token
	operand := ""
	switch bin.Op {
	case token.ADD: // +
		operand = expressions.AddOp
	case token.SUB: // -
		operand = expressions.SubtractOp
	case token.MUL: // *
		operand = expressions.MultiplyOp
	case token.QUO: // /
		operand = expressions.QuotentOp
	case token.REM: // %
		operand = expressions.RemainderOp
	case token.AND: // &
		operand = expressions.AndOp
	case token.OR: // |
		operand = expressions.OrOp
	case token.XOR: // ^
		operand = expressions.ExclusiveOrOp
	case token.SHL: // <<
		operand = expressions.LeftShiftOp
	case token.SHR: // >>
		operand = expressions.RightShiftOp
	case token.AND_NOT: // &^
		operand = expressions.AndNotOp
	case token.LAND: // &&
		operand = expressions.LogicalAndOp
	case token.LOR: // ||
		operand = expressions.LogicalOrOp
	case token.EQL: // ==
		resultType = types.Bool()
		operand = expressions.EqualOp
	case token.NEQ: // !=
		resultType = types.Bool()
		operand = expressions.NotEqualOp
	case token.LSS: // <
		resultType = types.Bool()
		operand = expressions.LessThanOp
	case token.LEQ: // <=
		resultType = types.Bool()
		operand = expressions.LessThanEqualOp
	case token.GTR: // >
		resultType = types.Bool()
		operand = expressions.GreaterThanOp
	case token.GEQ: // >=
		resultType = types.Bool()
		operand = expressions.GreaterThanEqualOp
	default:
		src.log.Error("Unhandled binary operand ", bin.Op)
		return nil
	}
	return expressions.BinaryOp(left, right, operand, resultType)
}

// parseIdentifier reads an identifier.
// https://golang.org/pkg/go/ast/#Ident
func (src *Source) parseIdentifier(scope *Scope, id *ast.Ident) expressions.Expression {
	name := id.Name
	if (name == "true") || (name == "false") {
		return expressions.Literal(name, types.Bool())
	} else if id := scope.Get(name); id != nil {
		return id
	}
	src.log.Error("Unable to find ", name, " in scope.")
	return expressions.Identifier(name, types.Variant())
}

// parseIndexer reads an index expression.
// https://golang.org/pkg/go/ast/#IndexExpr
func (src *Source) parseIndexer(scope *Scope, ind *ast.IndexExpr) expressions.Expression {
	exp := src.parseExpression(scope, ind.X)
	index := src.parseExpression(scope, ind.Index)

	collectionTypes := exp.ReturnTypes()
	if len(collectionTypes) != 1 {
		src.log.Error("Expected only one return value for an index.")
		return nil
	}
	indexType, ok := collectionTypes[0].(expressions.IndexableType)
	if !ok {
		src.log.Error("Unhandled indexed type: ", indexType)
		return nil
	}

	retType := indexType.Subtype()
	return expressions.Indexer(exp, index, retType)
}

// parseSelector reads an identifier selector.
// https://golang.org/pkg/go/ast/#SelectorExpr
func (src *Source) parseSelector(scope *Scope, sel *ast.SelectorExpr) *expressions.SelectorExp {
	exp := src.parseExpression(scope, sel.X)
	name := sel.Sel.Name
	returns := exp.ReturnTypes()
	if len(returns) != 1 {
		src.log.Error("Too many return values to select from:",
			"\n   Expression: ", exp,
			"\n   Selector:   ", name)
		return expressions.Selector(exp, name, types.Variant())
	}
	if t, exists := expressions.FindSubtype(returns[0], name); exists {
		return expressions.Selector(exp, name, t)
	}
	src.log.Error("Failed to find subtype for selector:",
		"\n   Expression: ", exp,
		"\n   Selector:   ", name)
	return expressions.Selector(exp, name, types.Variant())
}

// parseUnary reads a unary operation.
// https://golang.org/pkg/go/ast/#UnaryExpr
func (src *Source) parseUnary(scope *Scope, una *ast.UnaryExpr) *expressions.UnaryOpExp {
	exp := src.parseExpression(scope, una.X)

	var resultType types.Type
	if exp == nil {
		src.log.Debug("The expression in a unary is nil ", una.Op, ".")
		resultType = types.Variant()
	} else {
		rets := exp.ReturnTypes()
		if len(rets) != 1 {
			src.log.Error("Unary expected only one return type for ", una.Op, ".")
			resultType = types.Variant()
		} else {
			resultType = rets[0]
		}
	}

	// https://golang.org/pkg/go/token/#Token
	operand := ""
	switch una.Op {
	case token.ADD: // +
		operand = expressions.PosOp
	case token.SUB: // -
		operand = expressions.NegateOp
	case token.INC: // ++
		operand = expressions.IncrementOp
	case token.DEC: // --
		operand = expressions.DecrementOp
	case token.NOT: // !
		operand = expressions.NotOp
	case token.MUL: // *
		operand = expressions.DereferanceOp
	case token.AND: // &
		operand = expressions.ReferanceOp
	default:
		src.log.Error("Unhandled unary operand ", una.Op)
		return nil
	}
	return expressions.UnaryOp(exp, operand, resultType)
}

// parseCall reads a code literal.
// https://golang.org/pkg/go/ast/#CallExpr
func (src *Source) parseCall(scope *Scope, call *ast.CallExpr) *expressions.CallExp {
	fnExp := src.parseExpression(scope, call.Fun)
	paramLen := len(call.Args)
	params := make([]expressions.Expression, paramLen)
	returns := fnExp.ReturnTypes()
	if len(returns) != 1 {
		src.log.Error("Too many return values for a method call:",
			"\n   Expression: ", fnExp)
		return expressions.Call(nil, fnExp, params)
	}
	for i, param := range call.Args {
		params[i] = src.parseExpression(scope, param)
	}
	if fnHndl, exists := returns[0].(*types.FunctionType); exists {
		return expressions.Call(fnHndl, fnExp, params)
	}
	src.log.Error("Called type is not a function handle:",
		"\n   Expression: ", fnExp)
	return expressions.Call(nil, fnExp, params)
}

// expSliceToStatSlice converts a slice of expressions into a slice of statments.
func (src *Source) expSliceToStatSlice(exps []expressions.Expression) []statements.Statement {
	parts := make([]statements.Statement, len(exps))
	for i, exp := range exps {
		parts[i] = exp
	}
	return parts
}
