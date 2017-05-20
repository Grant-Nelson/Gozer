package types

var _ Type = (*StringType)(nil)
var _ IndexableType = (*StringType)(nil)

// StringType is type of a string.
type StringType struct{}

// String is the set of all strings of 8-bit bytes, conventionally
// but not necessarily representing UTF-8-encoded text. A string may be empty,
// but not nil. Values of string type are immutable.
func String() *StringType {
	return &StringType{}
}

// Subtype gets the indexable subtype from this type.
func (t *StringType) Subtype() Type {
	return UInt8()
}

// String gets the name for this type.
func (t *StringType) String() string {
	if t == nil {
		return nilStr
	}
	return "string"
}
