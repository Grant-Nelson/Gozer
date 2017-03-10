package constructs

var _ Type = (*ListType)(nil)

// ListType for storing the types of lists.
type ListType struct {

	// Element gets the element type for this list.
	Element Type
}

// List creates a new list type with the given element type
func List(element Type) *ListType {
	return &ListType{
		Element: element,
	}
}

// String gets the name for this type.
func (t *ListType) String() string {
	if t == nil {
		return nilStr
	}
	return "List<" + ToString(t.Element) + ">"
}
