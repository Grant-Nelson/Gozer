package ir

import (
	"go/token"
	"go/types"
	"strings"
)

// Parameter represents a single parameter that can be passed into a block.
type Param struct {
	// ParamObj is the object for this parameter.
	// This may be synthesized if the parameter was not named.
	ParamObj *types.Var
}

var _ Node = (*Param)(nil)

func (n *Param) Pos() token.Pos       { return n.ParamObj.Pos() }
func (n *Param) Type() types.Type     { return n.ParamObj.Type() }
func (n *Param) Object() types.Object { return n.ParamObj }
func (n *Param) Blank() bool          { return n.ParamObj.Name() == `_` }

func (n *Param) String() string {
	result := toString(n.ParamObj)
	result, _ = strings.CutPrefix(result, `var `)
	return result
}
