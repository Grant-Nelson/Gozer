package golang

import (
	"go/ast"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs/cMethod"
	"github.com/Snow-Gremlin/Gozer/constructs/cObject"
)

func funcReceiverIdent(funcDecl *ast.FuncDecl) string {
	if funcDecl == nil || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return ``
	}
	recv := funcDecl.Recv.List[0].Type
	for {
		switch r := recv.(type) {
		case *ast.IndexListExpr:
			recv = r.X
		case *ast.IndexExpr:
			recv = r.X
		case *ast.StarExpr:
			recv = r.X
		case *ast.Ident:
			return r.Name
		default:
			panic(terror.New(`unexpected type in receiver of function`).
				With(`type`, recv).
				With(`function`, funcDecl))
		}
	}
}

func (con *converter) addFunc(funcDecl *ast.FuncDecl) {
	m := cMethod.New()
	m.SetName(funcDecl.Name.Name)
	//directives := con.readDirectives(funcDecl.Doc)
	// "go:linkname"

	mSet := con.p.Methods()
	if recvName := funcReceiverIdent(funcDecl); len(recvName) > 0 {
		obj, exists := con.p.Objects().TryGetByName(recvName)
		if !exists {
			// Create placeholder object to add this function to
			obj = cObject.New()
			obj.SetName(recvName)
			con.p.Objects().Add(obj)
		}
		mSet = obj.Methods()
	}
	mSet.Add(m)

	// TODO: Add more method information.
	// - func unique receiver name
	// - func unique type parameters names
	// - parameters and return types
	// - func body
}
