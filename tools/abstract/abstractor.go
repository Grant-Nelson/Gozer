package abstract

import (
	"strings"

	"github.com/Snow-Gremlin/Gozer/reader"
	"github.com/Snow-Gremlin/Gozer/tools/abstract/models"
)

func abstract(proj *reader.Project) models.ProjectModel {
	basePath := proj.Packages[0].PkgPath
	projOut := models.NewProject(basePath)
	for _, pkg := range proj.PreOrder() {
		if strings.HasPrefix(pkg.PkgPath, basePath) {
			pkgOut := models.NewPackage(pkg)
			projOut.Packages().Append(pkgOut)
		}
	}

	return projOut
}
