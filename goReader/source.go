package transpiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"reflect"
	"strconv"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs"
)

// Source is the information about an input source file.
type Source struct {

	// log is the logger for reporting transpile issues with.
	log *common.Logger

	// fileSet is the full set of loaded  files.
	fileSet *token.FileSet

	// Path is the location of this source file for debugging purposes.
	Path string

	// Data is the data for this source file.
	Data *ast.File

	// Imports whieh were local to this method.
	// The key is the short name or path for the import.
	Imports map[string]*constructs.PackageType

	// Package is the package which this method belongs to.
	Package *constructs.PackageType
}

// NewSource creates a new source file descriptions.
func NewSource(log *common.Logger, fileSet *token.FileSet) *Source {
	return &Source{
		log:     log,
		fileSet: fileSet,
		Path:    "",
		Data:    nil,
		Imports: map[string]*constructs.PackageType{},
		Package: nil,
	}
}

// ProcessTypes determines all the class signatures, interfaces, functions, and handles.
func (src *Source) ProcessTypes() {
	scope := src.fillOutScope()
	for _, decl := range src.Data.Decls {
		switch data := decl.(type) {
		case *ast.GenDecl:
			// TODO: handle
			fmt.Println("GenDecl: ", data)
		case *ast.FuncDecl:
			src.readFunctionType(scope, data)
		default:
			common.ThrowError("Unhandled type declaration: ", decl, " (", reflect.TypeOf(decl), ")")
		}
	}
}

// ProcessBodies transpiles the bodies of the
func (src *Source) ProcessBodies() {
	// scope := src.fillOutScope()
	// TODO: fill out Constansts
	// TODO: fill out Variables
	// for _, named := range src.Functions.AllNamed() {
	// 	src.transpileFunction(scope, named.Data.(*FunctionType))
	// }

	// TODO: fill out Library Functions
	// TODO: fill out Class Functions
}

// addBasedPackage adds a package at the base level, not under a name.
func (src *Source) addBasedPackage(scope *Scope, pack *constructs.PackageType) {
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
func (src *Source) readType(scope *Scope, desc ast.Expr) (constructs.Type, bool) {
	if desc == nil {
		common.ThrowError("Nil type expression")
		return nil, false
	}
	switch id := desc.(type) {
	case *ast.Ident:
		return src.lookupType(scope, id.Name), false
	case *ast.Ellipsis:
		desc, _ := src.readType(scope, id.Elt)
		return constructs.List(desc), true
	default:
		common.ThrowError("Unhandled type expression: ", desc, " (", reflect.TypeOf(desc), ")")
		return nil, false
	}
}

// lookupType gets the type for the given Go type name.
func (src *Source) lookupType(scope *Scope, typeName string) constructs.Type {
	switch typeName {
	case "string":
		return constructs.String()
	default:
		common.ThrowError("Unhandled type name: ", typeName)
		return nil
	}
}

// readFieldList reads a list of parameters or returns from the given field.
func (src *Source) readFieldList(scope *Scope, fields *ast.FieldList) ([]string, []constructs.Type, bool) {
	names := []string{}
	types := []constructs.Type{}
	var typeDesc constructs.Type
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

// readFunction reads the given method into the given library
// and adds it to a class if the class is defined.
func (src *Source) readFunctionType(scope *Scope, data *ast.FuncDecl) {
	//ast.Print(src.fileSet, data)
	fn := constructs.Function()
	fn.Name = data.Name.Name

	// Receiver
	receiverNames, receiverTypes, _ := src.readFieldList(scope, data.Recv)
	if len(receiverNames) > 0 {
		fn.ReceiverName = receiverNames[0]
		class := receiverTypes[0].(*constructs.ClassType)
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
		fn.Body = src.parseBlock(scope, data.Body)
	}
}

// parseBlock reads a block statement.
// https://golang.org/pkg/go/ast/#BlockStmt
func (src *Source) parseBlock(scope *Scope, block *ast.BlockStmt) *constructs.BlockStatement {
	blockScope := NewScope(scope)
	stats := make([]constructs.Statement, len(block.List))
	for i, stmt := range block.List {
		stats[i] = src.parseStatement(blockScope, stmt)
	}
	return constructs.Block(stats...)
}

// parseStatement reads a statement.
// https://golang.org/pkg/go/ast/#Stmt
func (src *Source) parseStatement(scope *Scope, statement ast.Stmt) constructs.Statement {
	statScope := NewScope(scope)
	switch stat := statement.(type) {
	case *ast.ExprStmt:
		return src.parseExpression(statScope, stat.X)
	// case *ast.ReturnStmt: dw.writeReturn(st)
	default:
		src.log.Error("Unhandled statement type ", reflect.TypeOf(stat))
		return nil
	}
}

// parseExpression reads in some kind of expression.
// https://golang.org/pkg/go/ast/#Expr
func (src *Source) parseExpression(scope *Scope, expr ast.Expr) constructs.Expression {
	switch ex := expr.(type) {
	case *ast.BasicLit:
		return src.parseLiteral(scope, ex)
	// case *ast.BinaryExpr: src.parseBinary(ex)
	case *ast.CallExpr:
		return src.parseCall(scope, ex)
	case *ast.Ident:
		return src.parseIdentifier(scope, ex)
	case *ast.SelectorExpr:
		return src.parseSelector(scope, ex)
	default:
		src.log.Error("Unhandled expression type ", reflect.TypeOf(ex))
		return nil
	}
}

// parseLiteral reads a code literal.
// https://golang.org/pkg/go/ast/#BasicLit
// e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
func (src *Source) parseLiteral(scope *Scope, lit *ast.BasicLit) *constructs.LiteralExp {
	switch lit.Kind {
	case token.INT:
		val, err := strconv.ParseInt(lit.Value, 0, 64)
		if err != nil {
			src.log.Error("Error reading integer literal: ", err)
			return nil
		}
		if val > math.MaxInt32 {
			return constructs.Literal(lit.Value, constructs.Int64())
		}
		return constructs.Literal(lit.Value, constructs.Int())
	case token.FLOAT:
		return constructs.Literal(lit.Value, constructs.Float64())
	case token.IMAG:
		return constructs.Literal(lit.Value, constructs.Imaginary())
	case token.CHAR:
		return constructs.Literal(lit.Value, constructs.Rune())
	case token.STRING:
		// TODO: Check back tick strings and normalize string.
		return constructs.Literal(lit.Value, constructs.String())
	default:
		src.log.Error("Unhandled literal kind ", lit.Kind)
		return nil
	}
}

// parseIdentifier reads an identifier.
// https://golang.org/pkg/go/ast/#Ident
func (src *Source) parseIdentifier(scope *Scope, sel *ast.Ident) *constructs.IdentifierExp {
	return constructs.Identifier(sel.Name)
}

// parseSelector reads an identifier selector.
// https://golang.org/pkg/go/ast/#SelectorExpr
func (src *Source) parseSelector(scope *Scope, sel *ast.SelectorExpr) *constructs.SelectorExp {
	exp := src.parseExpression(scope, sel.X)
	return constructs.Selector(exp, sel.Sel.Name)
}

// parseCall reads a code literal.
// https://golang.org/pkg/go/ast/#CallExpr
func (src *Source) parseCall(scope *Scope, call *ast.CallExpr) *constructs.CallExp {
	fnHandl := src.parseExpression(scope, call.Fun)
	paramLen := len(call.Args)
	params := make([]constructs.Expression, paramLen)
	for i, param := range call.Args {
		params[i] = src.parseExpression(scope, param)
	}
	// Set the function to nil for now. It will be set when determining types.
	return constructs.Call(nil, fnHandl, params)
}
