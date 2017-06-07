package types

var _ Type = (*InterfaceType)(nil)
var _ NamedType = (*InterfaceType)(nil)
var _ SubtypableType = (*InterfaceType)(nil)

// InterfaceType for storing the types of interfaces.
type InterfaceType struct {

	// Name is the name of the interface.
	Name string

	// Functions is the set of functions for this interface.
	Functions *FunctionSet
}

// Interface creates a new interface type.
func Interface() *InterfaceType {
	return &InterfaceType{
		Name:      "",
		Functions: &FunctionSet{},
	}
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *InterfaceType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this interface.
func (t *InterfaceType) Find(name string) (Type, bool) {
	if t == nil {
		return nil, false
	}
	return t.Functions.Find(name)
}

// AddFunction adds a function to this interface.
// If a function by that name already exists, that function is returned.
func (t *InterfaceType) AddFunction(name string) *FunctionType {
	if t == nil {
		return nil
	}
	t2, _ := t.Functions.AddNew(name)
	return t2
}

// String gets the name for this type.
func (t *InterfaceType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the structure for this type.
func (t *InterfaceType) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "interface"
	if len(t.Name) > 0 {
		name = t.Name
	}
	if t.Functions.Len() <= 0 {
		return name + "{}"
	}
	return name + "{\n  " + t.Functions.FullString() + "\n}"
}
