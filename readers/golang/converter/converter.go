package loader

import (
	"go/build"

	"github.com/Snow-Gremlin/Gozer/constructs"
)

func Convert(pkg *build.Package) (*constructs.CPackage, error) {
	cpkg := &constructs.CPackage{
		Name:    pkg.Name,
		Path:    pkg.ImportPath,
		Imports: getImports(pkg),
	}

	return cpkg, nil
}

func getImports(pkg *build.Package) []*constructs.CImport {
	is := make([]*constructs.CImport, len(pkg.Imports))
	for i, p := range pkg.Imports {
		is[i] = &constructs.CImport{Path: p}
	}
	return is
}
