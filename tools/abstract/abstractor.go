package abstract

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/Gozer/reader"
	"github.com/Snow-Gremlin/Gozer/tools/abstract/models"
)

func abstract(projIn *reader.Project) models.ProjectModel {
	proj := models.NewProject(projIn)

	// Add all packages under the projects root path.
	rootPath := proj.Path() + `/`
	for _, pkg := range proj.Source().PreOrder() {
		if strings.HasPrefix(pkg.PkgPath, rootPath) {
			proj.AddPackage(pkg)
		}
	}

	proj.Packages().Enumerate().Foreach(func(pkg models.PackageModel) {
		for _, f := range pkg.Source().Syntax {
			handleFile(pkg, f)
		}
	})

	return proj
}

func handleFile(pkg models.PackageModel, f *ast.File) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			handleGenDecl(pkg, d)
		case *ast.FuncDecl:
			addFuncDecl(pkg, d)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func handleGenDecl(pkg models.PackageModel, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			handleTypeSpec(pkg, s)
		case *ast.ValueSpec:
			handleValueSpec(pkg, s)
		default:
			panic(fmt.Errorf(`unexpected specification: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func handleTypeSpec(pkg models.PackageModel, spec *ast.TypeSpec) {
	defs := pkg.Source().TypesInfo.Defs
	def := defs[spec.Name]

	n, ok := def.Type().(*types.Named)
	if !ok {
		panic(fmt.Errorf(`unexpected type for object, %T: %s`, def.Type(), def))
	}
	addType(pkg, n.Underlying())
}

func handleValueSpec(pkg models.PackageModel, spec *ast.ValueSpec) {
	// TODO: Implement

	/*
		defs := pkg.Source().TypesInfo.Defs
		for _, name := range spec.Names {
			def := defs[name]
			fmt.Printf(">>>(Value) %s\n", def.String())
		}
	*/
}

func addFuncDecl(pkg models.PackageModel, decl *ast.FuncDecl) {

	// TODO: Implement

}

func addType(pkg models.PackageModel, t types.Type) models.TypeModel {
	switch t2 := t.(type) {
	case *types.Basic:
		fmt.Println(">>>(Basic)", t2)
	case *types.Struct:
		fmt.Println(">>>(Struct)", t2)
	case *types.Interface:
		fmt.Println(">>>(Interface)", t2)
	case *types.Signature:
		fmt.Println(">>>(Signature)", t2)
	case *types.Pointer:
		panic(fmt.Errorf(`pointer is unimplemented: %s`, t2))
	case *types.Slice:
		panic(fmt.Errorf(`slice is unimplemented: %s`, t2))
	case *types.Map:
		panic(fmt.Errorf(`map is unimplemented: %s`, t2))
	case *types.Chan:
		panic(fmt.Errorf(`channel is unimplemented: %s`, t2))
	default:
		panic(fmt.Errorf(`unhandled type, %T: %s`, t, t))
	}
	return nil
}
