package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*KeyValueExp)(nil)

// KeyValueExp is a map definition typed value.
type KeyValueExp struct {

	// Keys is the key of a map entry.
	Key Expression

	// Values is the value of a map entry.
	Value Expression
}

// KeyValue creates new composite literal expression.
func KeyValue(key Expression, value Expression) *KeyValueExp {
	return &KeyValueExp{
		Key:   key,
		Value: value,
	}
}

// ReturnType is the type of the composite literal.
func (e *KeyValueExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return types.Pair(e.Key.ReturnType(), e.Value.ReturnType())
}

// String gets the string for the composite literal.
func (e *KeyValueExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.Key) + ": " + ToString(e.Value)
}
