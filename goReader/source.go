package transpiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

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

	// Declarations gets the set of constant and variable declarations from this source.
	Declarations map[string]constructs.Type

	// Functions is the set of functions from this source.
	Functions map[string]*constructs.FunctionType

	// Interfaces is the set of interfaces from this source.
	Interfaces map[string]*constructs.InterfaceType

	// Classes is the set of classes from this source.
	Classes map[string]*constructs.ClassType
}

// NewSource creates a new source file descriptions.
func NewSource(log *common.Logger, fileSet *token.FileSet) *Source {
	return &Source{
		log:          log,
		fileSet:      fileSet,
		Path:         "",
		Data:         nil,
		Imports:      map[string]*constructs.PackageType{},
		Package:      nil,
		Declarations: map[string]constructs.Type{},
		Functions:    map[string]*constructs.FunctionType{},
		Interfaces:   map[string]*constructs.InterfaceType{},
		Classes:      map[string]*constructs.ClassType{},
	}
}

// ProcessTypes determines all the class signatures, interfaces, functions, and handles.
func (src *Source) ProcessTypes() {
	scope := src.fillOutScope()
	for _, decl := range src.Data.Decls {
		switch data := decl.(type) {
		case *ast.GenDecl:
			// TODO: Impelement
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
	name := data.Name.Name

	// Receiver
	receiverNames, receiverTypes, _ := src.readFieldList(scope, data.Recv)
	if len(receiverNames) > 0 {
		fn.ReceiverName = receiverNames[0]
		class := receiverTypes[0].(*constructs.ClassType)
		class.Interface.Functions[name] = fn
		fn.ReceiverClass = class
	} else {
		src.Package.Functions[name] = fn
	}

	// Input parameters
	paramNames, paramTypes, ellipsis := src.readFieldList(scope, data.Type.Params)
	fn.ParamNames = paramNames
	fn.ParamTypes = paramTypes
	fn.Ellipsis = ellipsis

	// Return paramters
	resultNames, resultTypes, _ := src.readFieldList(scope, data.Type.Results)
	fn.ResultNames = resultNames
	fn.ResultTypes = resultTypes
}
