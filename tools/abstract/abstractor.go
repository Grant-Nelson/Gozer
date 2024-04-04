package abstract

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/Gozer/reader"
)

func abstractProject(proj *reader.Project) jsonData {
	pkgData := jsonList{}
	for _, pkg := range proj.PreOrder() {
		pkgData.add(abstractPackage(pkg))
	}
	projData := jsonMap{}
	projData.addNotEmpty(`packages`, pkgData)
	return projData
}

func abstractPackage(pkg *packages.Package) jsonData {
	pkgData := jsonMap{}
	pkgData.addNotEmpty(`path`, pkg.PkgPath)
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
	pkgData.append(`types`, convertType(tv.Type))
}

func abstractValueSpec(pkg *packages.Package, spec *ast.ValueSpec, pkgData jsonMap) {
	// TODO: Implement

	//	defs := pkg.Source().TypesInfo.Defs
	//	for _, name := range spec.Names {
	//		def := defs[name]
	//		fmt.Printf(">>>(Value) %s\n", def.String())
	//	}
}

func abstractFuncDecl(pkg *packages.Package, decl *ast.FuncDecl, pkgData jsonMap) {

	// TODO: Implement

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
	case *types.Pointer:
		return convertPointer(t2)
	case *types.Signature:
		return convertSignature(t2)
	case *types.Slice:
		return convertSlice(t2)
	case *types.Struct:
		return convertStruct(t2)
	default:
		panic(fmt.Errorf(`unhandled type, %T: %s`, t, t))
	}
}

func convertArray(t *types.Array) jsonData {
	return jsonMap{
		`kind`: `list`,
		`type`: convertType(t.Elem()),
	}
}

func convertBasic(t *types.Basic) jsonData {
	return t.Name()
}

func convertChan(t *types.Chan) jsonData {
	panic(fmt.Errorf(`chan is unimplemented: %s`, t))
}

func convertInterface(t *types.Interface) jsonData {
	t = t.Complete()
	methods := jsonList{}
	for i := range t.NumMethods() {
		methods.add(convertFunc(t.Method(i)))
	}
	return jsonMap{
		`kind`:    `interface`,
		`methods`: methods,
	}
}

func convertMap(t *types.Map) jsonData {
	panic(fmt.Errorf(`map is unimplemented: %s`, t))
}

func convertFunc(t *types.Func) jsonData {
	visibility := `internal`
	if t.Exported() {
		visibility = `public`
	}
	return jsonMap{
		`name`:       t.Name,
		`visibility`: visibility,
		`signature`:  convertType(t.Type()),
	}
}

func convertPointer(t *types.Pointer) jsonData {
	return jsonMap{
		`kind`: `pointer`,
		`type`: convertType(t.Elem()),
	}
}

func convertSignature(t *types.Signature) jsonData {
	panic(fmt.Errorf(`signature is unimplemented: %s`, t))
}

func convertSlice(t *types.Slice) jsonData {
	return jsonMap{
		`kind`: `list`,
		`type`: convertType(t.Elem()),
	}
}

func convertStruct(t *types.Struct) jsonData {
	panic(fmt.Errorf(`struct is unimplemented: %s`, t))
}
