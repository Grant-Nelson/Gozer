package goReader

import (
	"fmt"
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
	Imports *types.PackageSet

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
		Imports:      types.NewPackageSet(),
		Package:      nil,
		pendingFuncs: map[*types.FunctionType]*ast.BlockStmt{},
	}
}

// ProcessTypes determines all the class signatures, interfaces, functions, and handles.
func (src *Source) ProcessTypes() {
	defer src.log.RecoverError()
	defer src.log.PushPretext("Error occurred while processing types").Pop()

	scope := src.fillOutScope()
	for _, decl := range src.Data.Decls {
		switch data := decl.(type) {
		case *ast.GenDecl:
			src.readGenericDeclaration(scope, data)
		case *ast.FuncDecl:
			src.readFunctionType(scope, data)
		default:
			src.log.ThrowError("Unhandled type declaration: ", decl, " (", reflect.TypeOf(decl), ")")
		}
	}
}

// IsUnderscore determines if the given value is an underscore no-op identifier.
func (src *Source) IsUnderscore(x interface{}) bool {
	if stmt, ok := x.(*ast.Ident); ok && stmt.Name == "_" {
		return true
	}
	return false
}

// ProcessBodies transpiles the bodies of the functions and expressions of constants.
func (src *Source) ProcessBodies() {
	defer src.log.RecoverError()
	defer src.log.PushPretext("Error occurred while processing bodies").Pop()

	scope := src.fillOutScope()
	// TODO: fill out Constansts
	// TODO: fill out Variables
	for fn, body := range src.pendingFuncs {
		src.processPendingFunc(scope, fn, body)
	}
	// TODO: fill out Library Functions
	// TODO: fill out Class Functions
}

// processPendingFunc processes the pending function body.
func (src *Source) processPendingFunc(scope *Scope, fn *types.FunctionType, body *ast.BlockStmt) {
	defer src.log.RethrowError()
	defer src.log.PushData("Stage", "Processing pending function body").Pop()
	defer src.log.PushData("Path", src.getPath(body.Pos())).Pop()
	defer src.log.PushData("Mathod", fn.GetName()).Pop()

	fn.Body = src.parseBlock(scope, body)
}

// getPath gets the string for the path given the ast token position.
func (src *Source) getPath(pos token.Pos) string {
	loc := src.fileSet.Position(pos)
	return fmt.Sprint(loc.Filename, ":", loc.Line, ":", loc.Column)
}

// addBasedPackage adds a package at the base level, not under a name.
func (src *Source) addBasedPackage(scope *Scope, pack *types.PackageType) {
	shorts := pack.Imports.Shorts()
	for i, t := range pack.Imports.Packages() {
		name := t.GetName()
		if short := shorts[i]; len(short) > 0 {
			name = short
		}
		scope.Add(name, t)
	}
	for _, t := range pack.Declarations.Declarations {
		scope.Add(t.GetName(), t)
	}
	for _, t := range pack.Functions.Functions {
		scope.Add(t.GetName(), t)
	}
	for _, t := range pack.Interfaces.Interfaces {
		scope.Add(t.GetName(), t)
	}
	for _, t := range pack.Classes.Classes {
		scope.Add(t.GetName(), t)
	}
}

// fillOutScope fills out the scope for the containing package.
func (src *Source) fillOutScope() *Scope {
	defer src.log.RethrowError()
	defer src.log.PushPretext("Error occurred while filling out the scrope").Pop()

	scope := NewScope(nil)
	shorts := src.Imports.Shorts()
	for i, pack := range src.Imports.Packages() {
		name := pack.GetName()
		if short := shorts[i]; len(short) > 0 {
			name = short
		}
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
	defer src.log.PushData("Path", src.getPath(desc.Pos())).Pop()

	if desc == nil {
		src.log.ThrowError("Nil type expression")
		return nil, false
	}
	switch id := desc.(type) {
	case *ast.ArrayType:
		if id.Len != nil {
			src.log.ThrowError("Unhandled array length expression: ", id, " (", reflect.TypeOf(desc), ")")
		}
		element, _ := src.readType(scope, id.Elt)
		return types.List(element), false
	case *ast.MapType:
		key, _ := src.readType(scope, id.Key)
		value, _ := src.readType(scope, id.Value)
		return types.Map(key, value), false
	case *ast.Ident:
		return src.lookupType(scope, id.Name), false
	case *ast.Ellipsis:
		desc, _ := src.readType(scope, id.Elt)
		return types.List(desc), true
	default:
		src.log.ThrowError("Unhandled type expression: ", desc, " (", reflect.TypeOf(desc), ")")
		return nil, false
	}
}

// lookupType gets the type for the given Go type name.
func (src *Source) lookupType(scope *Scope, typeName string) types.Type {
	result := types.LookupType(typeName)
	if result == nil {
		src.log.ThrowError("Unhandled type name: ", typeName)
		return nil
	}
	return result
}

// readFieldList reads a list of parameters or returns from the given field.
func (src *Source) readFieldList(scope *Scope, fields *ast.FieldList) ([]string, []types.Type, bool) {
	names := []string{}
	typeList := []types.Type{}
	var typeDesc types.Type
	ellipsis := false
	if fields != nil {
		for _, field := range fields.List {
			typeDesc, ellipsis = src.readType(scope, field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					names = append(names, name.Name)
					typeList = append(typeList, typeDesc)
				}
			} else {
				names = append(names, "")
				typeList = append(typeList, typeDesc)
			}
		}
	}
	return names, typeList, ellipsis
}

// readGenericDeclaration reads the given generic declaration into the library.
func (src *Source) readGenericDeclaration(scope *Scope, data *ast.GenDecl) {
	defer src.log.RethrowError()
	defer src.log.PushPretext("Error occurred while reading a generic declaration").Pop()
	defer src.log.PushData("Stage", "Reading generic declaration").Pop()
	defer src.log.PushData("Path", src.getPath(data.Pos())).Pop()

	switch data.Tok {
	case token.IMPORT:
		// Ignore imports while reading generic declarations.
		return
	default:
		src.log.ThrowError("Unhandled generic declaration: ", data, " (", reflect.TypeOf(data), ")")
	}
}

// readFunctionType reads the given method into the library
// and adds it to a class if the class is defined.
func (src *Source) readFunctionType(scope *Scope, data *ast.FuncDecl) {
	defer src.log.RethrowError()
	defer src.log.PushPretext("Error occurred while reading a function type").Pop()
	defer src.log.PushData("Stage", "Reading function declaration").Pop()
	defer src.log.PushData("Path", src.getPath(data.Pos())).Pop()

	//ast.Print(src.fileSet, data)
	fn := types.Function()
	fn.Name = data.Name.Name

	// Receiver
	receiverNames, receiverTypes, _ := src.readFieldList(scope, data.Recv)
	if len(receiverNames) > 0 {
		fn.ReceiverName = receiverNames[0]
		class := receiverTypes[0].(*types.ClassType)
		class.Functions.Add(fn)
		fn.ReceiverClass = class
	} else {
		src.Package.Functions.Add(fn)
	}

	// Input parameters
	paramNames, paramTypes, ellipsis := src.readFieldList(scope, data.Type.Params)
	fn.ParamNames = paramNames
	fn.ParamTypes = paramTypes
	fn.Ellipsis = ellipsis

	// Return paramters
	returnNames, returnTypes, _ := src.readFieldList(scope, data.Type.Results)
	if len(returnTypes) <= 0 {
		fn.SetReturn(types.Void())
	} else if len(returnTypes) == 1 {
		if len(returnNames[0]) > 0 {
			ret := types.NewReturnSet()
			ret.AddMember(returnNames[0], returnTypes[0])
			fn.SetReturn(ret)
		} else {
			fn.SetReturn(returnTypes[0])
		}
	} else {
		ret := types.NewReturnSet()
		for i, retType := range returnTypes {
			name := returnNames[i]
			if len(name) <= 0 {
				name = fmt.Sprintf("val%d", i)
			}
			ret.AddMember(name, retType)
		}
		fn.SetReturn(ret)
	}

	// Read the body for the function
	if data.Body != nil {
		src.pendingFuncs[fn] = data.Body
	}
}

// parseBlock reads a block statement.
// https://golang.org/pkg/go/ast/#BlockStmt
func (src *Source) parseBlock(scope *Scope, block *ast.BlockStmt) *statements.BlockStat {
	defer src.log.RethrowError()
	defer src.log.PushPretext("Error occurred while parsing a block").Pop()

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
	case *ast.BlockStmt:
		return []statements.Statement{src.parseBlock(scope, stat)}
	case *ast.BranchStmt:
		return []statements.Statement{src.parseBranchStatement(scope, stat)}
	case *ast.ExprStmt:
		return []statements.Statement{src.parseExpression(scope, stat.X)}
	case *ast.ForStmt:
		return []statements.Statement{src.parseForStatement(scope, stat)}
	case *ast.IfStmt:
		return []statements.Statement{src.parseIfStatement(scope, stat)}
	case *ast.IncDecStmt:
		return []statements.Statement{src.parseIncDecStatement(scope, stat)}
	case *ast.RangeStmt:
		return []statements.Statement{src.parseRangeStatement(scope, stat)}

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

// parseForStatement reads a for-statement
// https://golang.org/pkg/go/ast/#ForStmt
func (src *Source) parseForStatement(scope *Scope, forStat *ast.ForStmt) statements.Statement {
	container := statements.Block()

	var init statements.Statement
	if forStat.Init != nil {
		inits := src.parseStatement(scope, forStat.Init)
		if len(inits) == 1 {
			init = inits[0]
		} else if len(inits) > 1 {
			container.Statements = append(container.Statements, inits...)
		}
	}

	var cond expressions.Expression
	if forStat.Cond != nil {
		cond = src.parseExpression(scope, forStat.Cond)
	}

	var post statements.Statement
	if forStat.Post != nil {
		posts := src.parseStatement(scope, forStat.Post)
		if len(posts) == 1 {
			post = posts[0]
		} else if len(posts) > 1 {
			fn := types.Function()
			fn.Body = statements.Block(posts...)
			id := scope.AddTemp(fn)
			def := expressions.Definition(id, expressions.Lambda(fn))
			container.Statements = append(container.Statements, def)
			post = expressions.Call(fn, id, nil)
		}
	}

	body := src.parseBlock(scope, forStat.Body)
	result := statements.For(init, cond, post, body)
	container.Statements = append(container.Statements, result)

	if len(container.Statements) == 1 {
		return container.Statements[0]
	}
	return container
}

// parseRangeStatement reads a range-statement (for-each-statement) for a list, string, or map.
// https://golang.org/pkg/go/ast/#RangeStmt
func (src *Source) parseRangeStatement(scope *Scope, stmt *ast.RangeStmt) statements.Statement {
	rangeExp := src.parseExpression(scope, stmt.X)
	switch t := rangeExp.ReturnType().(type) {
	case *types.StringType, *types.ListType:
		return src.parseListRangeStatement(scope, stmt, rangeExp)
	case *types.MapType:
		return src.parseMapRangeStatement(scope, stmt, rangeExp)
	default:
		src.log.Error("Unhandled foreach type ", reflect.TypeOf(t))
		return nil
	}
}

// parseListRangeStatement reads a range-statement (for-each-statement) for a list or string.
func (src *Source) parseListRangeStatement(scope *Scope, stmt *ast.RangeStmt, rangeExp expressions.Expression) statements.Statement {
	definition := stmt.Tok == token.DEFINE
	innerScope := NewScope(scope)

	tempRangeID := false
	var rangeID *expressions.IdentifierExp
	if idExp, ok := rangeExp.(*expressions.IdentifierExp); ok {
		rangeID = idExp
	} else {
		tempRangeID = true
		rangeID = innerScope.AddTemp(rangeExp.ReturnType())
	}

	var indexID *expressions.IdentifierExp
	var init expressions.Expression
	if (stmt.Key != nil) && !src.IsUnderscore(stmt.Key) {
		if definition {
			indexID = innerScope.Add(stmt.Key.(*ast.Ident).Name, types.Int())
			init = expressions.Definition(indexID, expressions.Literal("0", types.Int()))
		} else {
			indexID = src.parseExpression(innerScope, stmt.Key).(*expressions.IdentifierExp)
			init = expressions.Assignment(indexID, expressions.Literal("0", types.Int()))
		}
	} else {
		indexID = innerScope.AddTemp(types.Int())
		init = expressions.Definition(indexID, expressions.Literal("0", types.Int()))
	}

	var value expressions.Expression
	if (stmt.Value != nil) && !src.IsUnderscore(stmt.Value) {
		var valueID *expressions.IdentifierExp
		if definition {
			var valueType types.Type
			switch exp := rangeExp.ReturnType().(type) {
			case types.IndexableType:
				valueType = exp.ElementType()
			default:
				src.log.Error("Unhandled list range type ", reflect.TypeOf(exp))
			}
			valueID = innerScope.Add(stmt.Value.(*ast.Ident).Name, valueType)
			value = expressions.Definition(valueID, expressions.Indexer(rangeID, indexID))
		} else {
			valueID = src.parseExpression(innerScope, stmt.Value).(*expressions.IdentifierExp)
			value = expressions.Assignment(valueID, expressions.Indexer(rangeID, indexID))
		}
	}

	body := src.parseBlock(innerScope, stmt.Body)
	if value != nil {
		body.Statements = append([]statements.Statement{value}, body.Statements...)
	}

	lenFunc := scope.Get("len")
	condLen := expressions.Call(lenFunc.Type.(*types.FunctionType), lenFunc, []expressions.Expression{rangeID})
	cond := expressions.BinaryOp(indexID, condLen, expressions.LessThanOp, types.Bool())
	post := statements.IncDecOp(indexID, true)
	result := statements.For(init, cond, post, body)

	if tempRangeID {
		return statements.Block(
			expressions.Definition(rangeID, rangeExp),
			result)
	}
	return result
}

// parseMapRangeStatement reads a range-statement (for-each-statement) for maps.
func (src *Source) parseMapRangeStatement(scope *Scope, stmt *ast.RangeStmt, rangeExp expressions.Expression) statements.Statement {
	definition := stmt.Tok == token.DEFINE
	innerScope := NewScope(scope)

	var key expressions.Expression
	if (stmt.Key != nil) && !src.IsUnderscore(stmt.Key) {
		if definition {
			var keyType types.Type
			switch exp := rangeExp.ReturnType().(type) {
			case *types.StringType:
				keyType = types.Int()
			case *types.ListType:
				keyType = types.Int()
			case *types.MapType:
				keyType = exp.Key
			}
			keyID := innerScope.Add(stmt.Key.(*ast.Ident).Name, keyType)
			key = expressions.Definition(keyID, nil)
		} else {
			key = src.parseExpression(innerScope, stmt.Key)
		}
	}

	var value expressions.Expression
	if (stmt.Value != nil) && !src.IsUnderscore(stmt.Value) {
		if definition {
			var valueType types.Type
			switch exp := rangeExp.ReturnType().(type) {
			case *types.StringType:
				valueType = exp.ElementType()
			case *types.ListType:
				valueType = exp.ElementType()
			case *types.MapType:
				valueType = exp.Value
			}
			valueID := innerScope.Add(stmt.Value.(*ast.Ident).Name, valueType)
			value = expressions.Definition(valueID, nil)
		} else {
			value = src.parseExpression(innerScope, stmt.Value)
		}
	}

	body := src.parseBlock(innerScope, stmt.Body)

	// rangeIsID := false
	// switch rangeExp.(type) {
	// case *expressions.IdentifierExp:
	// 	rangeIsID = true
	// }

	// TODO: Finish
	return statements.Foreach(key, value, rangeExp, body)
}

// parseIncDecStatement reads an increment or decrement statement.
// https://golang.org/pkg/go/ast/#IncDecStmt
func (src *Source) parseIncDecStatement(scope *Scope, stmt *ast.IncDecStmt) statements.Statement {
	exp := src.parseExpression(scope, stmt.X)
	return statements.IncDecOp(exp, stmt.Tok == token.INC)
}

// parseBranchStatement reads a branch statement.
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

	// TODO: Handle underscores
	// TODO: Confirm type compatability

	// For single assignment assign it directly
	if (len(assign.Lhs) == 1) && (len(assign.Rhs) == 1) {
		left := src.parseExpression(scope, assign.Lhs[0])
		right := src.parseExpression(scope, assign.Rhs[0])
		return []expressions.Expression{expressions.Assignment(left, right)}
	}

	// Get all left expressions
	lefts := make([]expressions.Expression, len(assign.Lhs))
	for i, exp := range assign.Lhs {
		lefts[i] = src.parseExpression(scope, exp)
	}

	// Store all results to temporary locations so that if a swap is occuring
	// the correct value is still used on the right after the initial assignment.
	results := make([]expressions.Expression, 0, len(assign.Rhs))
	tempIDs := make([]expressions.Expression, 0, len(assign.Lhs))
	for _, exp := range assign.Rhs {
		right := src.parseExpression(scope, exp)

		// Assign result to single temporary.
		ret := right.ReturnType()
		tempID := scope.AddTemp(ret)
		results = append(results, expressions.Definition(tempID, right))

		// If multi-return then select on temp structure members so results are assigned correctly.
		if retSet, ok := ret.(*types.ReturnSet); ok {
			retCount := retSet.Members.Len()
			for i := 0; i < retCount; i++ {
				decl := retSet.Members.Declarations[i]
				tempIDs = append(tempIDs, expressions.Selector(tempID, decl.Name, decl.Data))
			}
		} else {
			tempIDs = append(tempIDs, tempID)
		}
	}

	// Write all temporary values to the final values.
	for i, tempID := range tempIDs {
		results = append(results, expressions.Assignment(lefts[i], tempID))
	}
	return results
}

// parseDefinition reads in an assigment definition statement.
func (src *Source) parseDefinition(scope *Scope, assign *ast.AssignStmt) []expressions.Expression {
	leftOffset := 0
	results := make([]expressions.Expression, 0, len(assign.Rhs))
	for _, exp := range assign.Rhs {
		right := src.parseExpression(scope, exp)

		// If multi-return then store to temp then select on temp structure
		// members so results are assigned correctly.
		ret := right.ReturnType()
		if retSet, ok := ret.(*types.ReturnSet); ok {
			tempID := scope.AddTemp(ret)
			results = append(results, expressions.Definition(tempID, right))
			retCount := retSet.Members.Len()
			for j := 0; j < retCount; j++ {
				decl := retSet.Members.Declarations[j]
				rightDef := expressions.Selector(tempID, decl.Name, decl.Data)

				leftExp := assign.Lhs[leftOffset]
				leftOffset++
				id := scope.Add(leftExp.(*ast.Ident).Name, decl.Data)
				results = append(results, expressions.Definition(id, rightDef))
			}
		} else {
			leftExp := assign.Lhs[leftOffset]
			leftOffset++
			id := scope.Add(leftExp.(*ast.Ident).Name, ret)
			results = append(results, expressions.Definition(id, right))
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
	case *ast.KeyValueExpr:
		return src.parseKeyValue(scope, ex)
	case *ast.ParenExpr:
		return src.parseExpression(scope, ex.X)
	case *ast.SelectorExpr:
		return src.parseSelector(scope, ex)
	case *ast.SliceExpr:
		return src.parseSlice(scope, ex)
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

// parseCompositeLit reads a composite literial value for a slice.
// https://golang.org/pkg/go/ast/#CompositeLit
func (src *Source) parseCompositeLit(scope *Scope, exp *ast.CompositeLit) *expressions.CompositeLiteralExp {
	litType, _ := src.readType(scope, exp.Type)
	elements := make([]expressions.Expression, len(exp.Elts))
	for i, elts := range exp.Elts {
		elements[i] = src.parseExpression(scope, elts)
	}
	return expressions.CompositeLiteral(elements, litType)
}

// parseMapLit reads a composite literial value for a map.
// https://golang.org/pkg/go/ast/#KeyValueExpr
func (src *Source) parseKeyValue(scope *Scope, exp *ast.KeyValueExpr) *expressions.KeyValueExp {
	key := src.parseExpression(scope, exp.Key)
	value := src.parseExpression(scope, exp.Value)
	return expressions.KeyValue(key, value)
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
		resultType = left.ReturnType()
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
	} else if name == "make" {
		return expressions.Make()
	}
	src.log.Error("Unable to find ", name, " in scope")
	return expressions.Identifier(name, types.Variant())
}

// parseIndexer reads an index expression.
// https://golang.org/pkg/go/ast/#IndexExpr
func (src *Source) parseIndexer(scope *Scope, ind *ast.IndexExpr) expressions.Expression {
	exp := src.parseExpression(scope, ind.X)
	index := src.parseExpression(scope, ind.Index)

	collectionType := exp.ReturnType()
	indexType, ok := collectionType.(types.IndexableType)
	if !ok {
		src.log.Error("Unhandled indexed type: ", indexType)
		return nil
	}

	return expressions.Indexer(exp, index)
}

// parseSelector reads an identifier selector.
// https://golang.org/pkg/go/ast/#SelectorExpr
func (src *Source) parseSelector(scope *Scope, sel *ast.SelectorExpr) *expressions.SelectorExp {
	exp := src.parseExpression(scope, sel.X)
	name := sel.Sel.Name
	if t, exists := types.FindSubtype(exp.ReturnType(), name); exists {
		return expressions.Selector(exp, name, t)
	}
	src.log.Error("Failed to find subtype for selector").
		Add("Expression", exp).
		Add("Selector", name)
	return expressions.Selector(exp, name, types.Variant())
}

// parseSlice reads a subslice creation call.
// https://golang.org/pkg/go/ast/#SliceExpr
// https://blog.golang.org/go-slices-usage-and-internals
func (src *Source) parseSlice(scope *Scope, sel *ast.SliceExpr) *expressions.SubsliceExp {
	exp := src.parseExpression(scope, sel.X)
	var low, high, max expressions.Expression
	if sel.Low != nil {
		low = src.parseExpression(scope, sel.Low)
	}
	if sel.High != nil {
		high = src.parseExpression(scope, sel.High)
	}
	if sel.Max != nil {
		max = src.parseExpression(scope, sel.Max)
	}
	return expressions.Subslice(exp, low, high, max)
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
		resultType = exp.ReturnType()
	}

	// https://golang.org/pkg/go/token/#Token
	operand := ""
	switch una.Op {
	case token.ADD: // +
		operand = expressions.PosOp
	case token.SUB: // -
		operand = expressions.NegateOp
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

// parseCall reads a code method call expression.
// https://golang.org/pkg/go/ast/#CallExpr
func (src *Source) parseCall(scope *Scope, call *ast.CallExpr) expressions.Expression {
	fnExp := src.parseExpression(scope, call.Fun)
	if m, ok := fnExp.(*expressions.MakeExp); ok {
		return src.parseMakeCall(scope, call, m)
	}

	paramLen := len(call.Args)
	params := make([]expressions.Expression, paramLen)
	for i, param := range call.Args {
		params[i] = src.parseExpression(scope, param)
	}
	returnType := fnExp.ReturnType()
	if fnHndl, exists := returnType.(*types.FunctionType); exists {
		return expressions.Call(fnHndl, fnExp, params)
	}
	src.log.Error("Called type is not a function handle:",
		"\n   Expression: ", fnExp)
	return expressions.Call(nil, fnExp, params)
}

// parseMakeCall reads a make method call.
func (src *Source) parseMakeCall(scope *Scope, call *ast.CallExpr, makeExp *expressions.MakeExp) expressions.Expression {
	paramLen := len(call.Args)
	if (paramLen >= 1) && (paramLen <= 3) {
		typeDef, _ := src.readType(scope, call.Args[0])
		makeExp.Type = typeDef
		if paramLen >= 2 {
			makeExp.Length = src.parseExpression(scope, call.Args[1])
		}
		if paramLen >= 3 {
			makeExp.Capacity = src.parseExpression(scope, call.Args[2])
		}
	} else {
		src.log.Error("Make call must have 1 to 3 arguments but got ", paramLen)
	}
	return makeExp
}

// expSliceToStatSlice converts a slice of expressions into a slice of statements.
func (src *Source) expSliceToStatSlice(exps []expressions.Expression) []statements.Statement {
	parts := make([]statements.Statement, len(exps))
	for i, exp := range exps {
		parts[i] = exp
	}
	return parts
}
