package abstract

import (
	"fmt"
	"go/ast"
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
			addFile(pkg, f)
		}
	})

	return proj
}

func addFile(pkg models.PackageModel, f *ast.File) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			addGenDecl(pkg, d)
		case *ast.FuncDecl:
			addFuncDecl(pkg, d)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func addGenDecl(pkg models.PackageModel, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			addTypeSpec(pkg, s)
		case *ast.ValueSpec:
			addValueSpec(pkg, s)
		default:
			panic(fmt.Errorf(`unexpected specification: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func addTypeSpec(pkg models.PackageModel, spec *ast.TypeSpec) {
	// TODO: Add methods
}

func addValueSpec(pkg models.PackageModel, spec *ast.ValueSpec) {
	defs := pkg.Source().TypesInfo.Defs
	for _, name := range spec.Names {
		def := defs[name]
		if def.Exported() {
			fmt.Printf(">>> %s\n", def.String())
		}
	}
}

func addFuncDecl(pkg models.PackageModel, decl *ast.FuncDecl) {
	// TODO: Add methods
}
