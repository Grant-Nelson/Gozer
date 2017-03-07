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

// Equals determins if the given type is the same as this type.
func (t *ListType) Equals(other interface{}) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*ListType)
	if !ok {
		return false
	}
	if !Equals(t.Element, tother.Element) {
		return false
	}
	return true
}

// String gets the name for this type.
func (t *ListType) String() string {
	if t == nil {
		return nilStr
	}
	return "List<" + ToString(t.Element) + ">"
}
