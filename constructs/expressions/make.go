package expressions

import (
	"strings"

	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*MakeExp)(nil)

// MakeExp is a make expression.
type MakeExp struct {

	// Type is the type being made
	Type types.Type

	// Length is the optional initial length of the slice to make.
	Length Expression

	// Capacity is the optional initial capacity of the slice to make.
	Capacity Expression
}

// Make creates new make expression.
func Make() *MakeExp {
	return &MakeExp{
		Type:     types.Void(),
		Length:   nil,
		Capacity: nil,
	}
}

// ReturnType is the return type from the called make.
func (e *MakeExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Type
}

// String gets the string for the call.
func (e *MakeExp) String() string {
	if e == nil {
		return nilStr
	}
	parts := make([]string, 0, 3)
	parts = append(parts, types.ToString(e.Type))
	if e.Length != nil {
		parts = append(parts, ToString(e.Length))
	}
	if e.Capacity != nil {
		if e.Length == nil {
			parts = append(parts, "0")
		}
		parts = append(parts, ToString(e.Capacity))
	}
	return "make(" + strings.Join(parts, ", ") + ")"
}
