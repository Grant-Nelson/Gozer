package ir

// Directive is an additional directives added to a declaration,
// statement, or package to tell the compiler to perform a special action.
//
// See https://go.dev/doc/comment#directives
// See https://pkg.go.dev/cmd/compile
type Directive interface {
	Node
	DirectiveNode()
}
