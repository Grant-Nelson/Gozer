package constructs

var _ Type = (*ClassType)(nil)

// ClassType for storing the types of class.
type ClassType struct {

	// Data is the type of the data in this class.
	Data Type

	// Interface is the interface for this class.
	Interface *InterfaceType
}

// Class creates a new class type.
func Class() *ClassType {
	return &ClassType{
		Data:      nil,
		Interface: Interface(),
	}
}

// Equals determins if the given type is the same as this type.
func (t *ClassType) Equals(other interface{}) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*ClassType)
	if !ok {
		return false
	}
	if !Equals(t.Data, tother.Data) {
		return false
	}
	if !Equals(t.Interface, tother.Interface) {
		return false
	}
	return true
}

// String gets the name for this type.
func (t *ClassType) String() string {
	if t == nil {
		return nilStr
	}
	return "{\n" +
		indent(ToString(t.Data), "  ") + "\n" +
		indent(ToString(t.Interface), "  ") + "\n" +
		"}"
}
