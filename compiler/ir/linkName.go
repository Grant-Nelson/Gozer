package ir

import "go/token"

// LinkName is an alternative name of a symbol for a var, const, type, or func
// as defined by the `//go:linkname localName [importPath.name]` directive.
// See <https://pkg.go.dev/cmd/compile#hdr-Linkname_Directive>
//
// If a function is a stub the link may be defined as either a push or pull link.
// If a var or const value, the two values will be connected with the other
// and the local value will be used a pointer with a get/set method.
type LinkName struct {
	// LinkPos is the location that this linkname is being defined.
	LinkPos token.Pos

	// LocalName is the name of the local symbol being linked to or from.
	LocalName string

	RemotePath string

	RemoteName string
}

var _ Node = (*Func)(nil)

func (n *LinkName) Pos() token.Pos { return n.LinkPos }

func (n *LinkName) String() string {
	return n.LocalName + ` => ` + n.RemotePath + `.` + n.RemoteName
}
