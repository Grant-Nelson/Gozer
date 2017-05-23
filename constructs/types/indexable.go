package types

var _ Type = (IndexableType)(nil)

// IndexableType is a type add-on for types which are indexable
// such as strings, lists, and maps.
type IndexableType interface {

	// ElementType gets the element type from the indexable type.
	ElementType() Type

	// String gets the identifier name for this type.
	String() string
}

// GetElementType gets the element type from an indexable type.
// If no element type exists, false is returned.
func GetIndexableType(t Type) (Type, bool) {
	if t != nil {
		if t2, ok := t.(IndexableType); ok {
			return t2.ElementType(), true
		}
	}
	return nil, false
}
