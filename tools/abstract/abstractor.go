package abstract

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/Gozer/internal/reader"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

// TODO:
// - Figure out implemented interfaces.
// - Determine what to do with pointer receivers to make it similar o Java.
// - Add analytics:
//   - Add cyclomatic complexity per method.
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)

func abstractProject(proj *reader.Project) jsonData {
	projData := jsonMap{
		`duckTyping`: `true`,
		`language`:   `go`,
	}
	pkgData := jsonList{}
	for _, pkg := range proj.PreOrder() {
		pkgData.add(abstractPackage(pkg))
	}
	projData.addNotEmpty(`packages`, pkgData)
	return projData
}

func abstractPackage(pkg *packages.Package) jsonData {
	pkgData := jsonMap{}.
		addNotEmpty(`path`, pkg.PkgPath).
		addNotEmpty(`imports`, utils.SortedKeys(pkg.Imports))
	for _, f := range pkg.Syntax {
		addFile(pkg, f, pkgData)
	}
	return pkgData
}

func pos(pkg *packages.Package, pos token.Pos) string {
	return pkg.Fset.Position(pos).String()
}

func addFile(pkg *packages.Package, f *ast.File, pkgData jsonMap) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			addGenDecl(pkg, d, pkgData)
		case *ast.FuncDecl:
			abstractFuncDecl(pkg, d, pkgData)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pos(pkg, decl.Pos())))
		}
	}
}

func addGenDecl(pkg *packages.Package, decl *ast.GenDecl, pkgData jsonMap) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			abstractTypeSpec(pkg, s, pkgData)
		case *ast.ValueSpec:
			abstractValueSpec(pkg, s, pkgData)
		default:
			panic(fmt.Errorf(`unexpected specification: %s`, pos(pkg, spec.Pos())))
		}
	}
}

func abstractTypeSpec(pkg *packages.Package, spec *ast.TypeSpec, pkgData jsonMap) {
	tv, has := pkg.TypesInfo.Types[spec.Type]
	if !has {
		panic(fmt.Errorf(`type specification not found in types info: %s`, pos(pkg, spec.Type.Pos())))
	}
	data := jsonMap{
		`name`: spec.Name.Name,
		`type`: convertType(tv.Type),
	}
	pkgData.append(`types`, data)
}

func abstractValueSpec(pkg *packages.Package, spec *ast.ValueSpec, pkgData jsonMap) {
	for _, name := range spec.Names {
		if name.Name != `_` {
			tv, has := pkg.TypesInfo.Defs[name]
			if !has {
				panic(fmt.Errorf(`value specification not found in types info: %s`, pos(pkg, spec.Type.Pos())))
			}
			data := jsonMap{
				`name`: name.Name,
				`type`: convertType(tv.Type()),
			}
			pkgData.append(`values`, data)
		}
	}
}

func abstractFuncDecl(pkg *packages.Package, decl *ast.FuncDecl, pkgData jsonMap) {
	obj := pkg.TypesInfo.Defs[decl.Name]
	data := jsonMap{
		`name`:      decl.Name.Name,
		`signature`: convertSignature(obj.Type().(*types.Signature), false),
	}

	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(fmt.Errorf(`function declaration has unexpected receiver fields: %s`, pos(pkg, decl.Pos())))
		}
		tv := pkg.TypesInfo.Types[decl.Recv.List[0].Type]
		data.add(`receiver`, convertType(tv.Type))
	}

	// TODO: Add cyclomatic complexity and other information

	pkgData.append(`methods`, data)
}

func convertType(t types.Type) jsonData {
	switch t2 := t.(type) {
	case *types.Array:
		return convertArray(t2)
	case *types.Basic:
		return convertBasic(t2)
	case *types.Chan:
		return convertChan(t2)
	case *types.Interface:
		return convertInterface(t2)
	case *types.Map:
		return convertMap(t2)
	case *types.Named:
		return convertNamed(t2)
	case *types.Pointer:
		return convertPointer(t2)
	case *types.Signature:
		return convertSignature(t2, true)
	case *types.Slice:
		return convertSlice(t2)
	case *types.Struct:
		return convertStruct(t2)
	case *types.TypeParam:
		return convertTypeParam(t2)
	default:
		panic(fmt.Errorf(`unhandled type, %T: %s`, t, t))
	}
}

func convertArray(t *types.Array) jsonData {
	return jsonMap{
		`kind`: `list`,
		`elem`: convertType(t.Elem()),
	}
}

func convertBasic(t *types.Basic) jsonData {
	return t.Name()
}

func convertChan(t *types.Chan) jsonData {
	return jsonMap{
		`kind`: `chan`,
		`elem`: convertType(t.Elem()),
	}
}

func convertInterface(t *types.Interface) jsonData {
	t = t.Complete()
	if t.NumMethods() == 0 {
		return `any`
	}

	data := jsonMap{
		`kind`: `interface`,
	}
	for i := range t.NumMethods() {
		data.append(`methods`, convertFunc(t.Method(i)))
	}
	return data
}

func convertMap(t *types.Map) jsonData {
	return jsonMap{
		`kind`: `map`,
		`key`:  convertType(t.Key()),
		`elem`: convertType(t.Elem()),
	}
}

func convertNamed(t *types.Named) jsonData {
	return t.String()
}

func convertFunc(t *types.Func) jsonData {
	return jsonMap{
		`name`:      t.Name(),
		`signature`: convertSignature(t.Type().(*types.Signature), false),
	}
}

func convertPointer(t *types.Pointer) jsonData {
	return jsonMap{
		`kind`: `pointer`,
		`type`: convertType(t.Elem()),
	}
}

func convertSignature(t *types.Signature, showKind bool) jsonData {
	// Don't output receiver or receiver type here.
	data := jsonMap{}.
		addNotEmpty(`variadic`, t.Variadic()).
		addNotEmpty(`params`, convertTuple(t.Params())).
		addNotEmpty(`returns`, convertTuple(t.Results())).
		addNotEmpty(`typeParams`, convertTypeParamList(t.TypeParams()))
	if showKind {
		data.add(`kind`, `signature`)
	}
	return data
}

func convertSlice(t *types.Slice) jsonData {
	return jsonMap{
		`kind`: `list`,
		`type`: convertType(t.Elem()),
	}
}

func convertStruct(t *types.Struct) jsonData {
	data := jsonMap{
		`kind`: `struct`,
	}
	for i := range t.NumFields() {
		data.append(`fields`, convertVar(t.Field(i)))
	}
	return data
}

func convertTuple(t *types.Tuple) jsonData {
	data := jsonList{}
	for i := range t.Len() {
		data.add(convertVar(t.At(i)))
	}
	return data
}

func convertTypeParam(t *types.TypeParam) jsonData {
	return jsonMap{
		`kind`:       `typeParam`,
		`index`:      t.Index(),
		`constraint`: convertType(t.Constraint()),
	}
}

func convertTypeParamList(t *types.TypeParamList) jsonData {
	data := jsonList{}
	for i := range t.Len() {
		if p := t.At(i); p.Index() >= 0 {
			data.add(convertTypeParam(p))
		}
	}
	return data
}

func convertVar(t *types.Var) jsonData {
	data := jsonMap{
		`type`: convertType(t.Type()),
	}
	data.addNotEmpty(`name`, t.Name())
	return data
}
