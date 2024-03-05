package golang

import (
	"fmt"
	"go/ast"
)

func (con *converter) addTypeSpec(tSpec *ast.TypeSpec, declDirectives []string) {
	directives := append(con.readDirectives(tSpec.Doc, tSpec.Comment), declDirectives...)
	name := tSpec.Name.Name
	//aliased := tSpec.Assign == token.NoPos

	fmt.Printf("%s: %+v\n", name, tSpec)
	fmt.Println(`Type Directives`, directives)
	// TODO: Implement
}
