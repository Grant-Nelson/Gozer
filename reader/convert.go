package reader

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/Gozer/constructs"
)

func convert(proj constructs.IProject, fileSet *token.FileSet, tPack *types.Package, info *types.Info, files []*ast.File) {
	pkg := constructs.NewPackage(tPack.Name())
	proj.Packages().Append(pkg)

	conv := &converter{
		proj:    proj,
		pkg:     pkg,
		fileSet: fileSet,
		tPack:   tPack,
		info:    info,
	}

	for _, in := range tPack.Imports() {

		inP, has := proj.AtPath(in.Path())
	}

	for _, file := range files {
		conv.convertFile(file)
	}
}

type converter struct {
	proj    constructs.IProject
	pkg     constructs.IPackage
	fileSet *token.FileSet
	tPack   *types.Package
	info    *types.Info
}

func (conv *converter) convertFile(file *ast.File) {

}
