package types

import "github.com/grant-nelson/Gozer/common"

var _ Type = (*ClassType)(nil)
var _ NamedType = (*ClassType)(nil)
var _ SubtypableType = (*ClassType)(nil)

// ClassType for storing the types of class.
type ClassType struct {

	// Name is the name of the class.
	Name string

	// Data is the type of the data in this class.
	// This should not be a package, interface, or class.
	Data Type

	// Interface is the interface for this class.
	Interface *InterfaceType
}

// Class creates a new class type.
func Class() *ClassType {
	c := &ClassType{
		Name:      "",
		Data:      nil,
		Interface: Interface(),
	}
	return c
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *ClassType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this class.
func (t *ClassType) Find(name string) (Type, bool) {
	if t == nil {
		return nil, false
	}
	if t2, exists := t.Interface.Find(name); exists {
		return t2, true
	}
	if structType, ok := t.Data.(SubtypableType); ok {
		t2, exists := structType.Find(name)
		return t2, exists
	}
	return nil, false
}

// String gets the name for this type.
func (t *ClassType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the structure for this type.
func (t *ClassType) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "class"
	if len(t.Name) > 0 {
		name = t.Name
	}
	result := ""
	if str := ToString(t.Data); (len(str) > 0) && (str != nilStr) {
		result += "  " + common.Indent(str, "  ") + "\n"
	}
	if str := ToString(t.Interface); (len(str) > 0) && (str != nilStr) && (str != "interface{}") {
		result += "  " + common.Indent(str, "  ") + "\n"
	}
	if len(result) <= 0 {
		return name + "{}"
	}
	return name + "{\n" + result + "}"
}
