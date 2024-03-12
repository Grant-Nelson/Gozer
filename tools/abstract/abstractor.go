package abstract

import (
	"fmt"
	"sort"

	"github.com/Snow-Gremlin/Gozer/reader"
)

type abstractor struct {
}

func (a *abstractor) abstractPackage(args *reader.ConvertPackageArgs) {
	fmt.Println(`Package:`, args.Package.Name())

	parts := []string{}
	for id, obj := range args.Info.Defs {
		part := id.Name + `: ` + obj.String()
		parts = append(parts, part)
	}
	sort.Strings(parts)

	for _, part := range parts {
		fmt.Println(`   ` + part)
	}
}
