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
		if !strings.HasPrefix(pkg.PkgPath, rootPath) {
			proj.AddPackage(pkg)
		}
	}

	proj.Packages().Enumerate().Foreach(func(pkg models.PackageModel) {
		for _, f := range pkg.Source().Syntax {
			addFile(proj, pkg, f)
		}
	})

	return proj
}

func addFile(proj models.ProjectModel, pkg models.PackageModel, f *ast.File) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			addGenDecl(proj, pkg, d)
		case *ast.FuncDecl:
			addFuncDecl(proj, pkg, d)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func addGenDecl(proj models.ProjectModel, pkg models.PackageModel, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			addTypeSpec(proj, pkg, s)
		case *ast.ValueSpec:
			addValueSpec(proj, pkg, s)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pkg.PosPath(decl.Pos())))
		}
	}
}

func addTypeSpec(proj models.ProjectModel, pkg models.PackageModel, spec *ast.TypeSpec) {
	// TODO: Add methods
}

func addValueSpec(proj models.ProjectModel, pkg models.PackageModel, spec *ast.ValueSpec) {
	// TODO: Add methods
}

func addFuncDecl(proj models.ProjectModel, pkg models.PackageModel, decl *ast.FuncDecl) {
	// TODO: Add methods
}
