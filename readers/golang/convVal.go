package golang

import (
	"fmt"
	"go/ast"
)

func (con *converter) addValueSpec(vSpec *ast.ValueSpec, declDirectives []string) {
	directives := append(con.readDirectives(vSpec.Doc, vSpec.Comment), declDirectives...)

	fmt.Println(`Value Directives`, directives)
	// TODO: Implement

}
