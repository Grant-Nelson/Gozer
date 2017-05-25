package types

import "github.com/grant-nelson/Gozer/common"

var _ Type = (*ClassType)(nil)
var _ SubtypableType = (*ClassType)(nil)

// ClassType for storing the types of class.
type ClassType struct {

	// Parent is the parent type for this function.
	// It is a package or nil.
	Parent Type

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
		Parent:    nil,
		Name:      "",
		Data:      nil,
		Interface: Interface(),
	}
	c.Interface.Parent = c
	return c
}

// Find looks up a subtype to this class.
func (t *ClassType) Find(name string) (Type, bool) {
	if t2, exists := t.Interface.Find(name); exists {
		return t2, true
	}
	if structType, ok := t.Data.(*StructureType); ok {
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

	result := ""
	if str := ToString(t.Data); (len(str) > 0) && (str != "nil") {
		result += "  " + common.Indent(str, "  ") + "\n"
	}
	if str := ToString(t.Interface); (len(str) > 0) && (str != "interface{}") {
		result += "  " + common.Indent(str, "  ") + "\n"
	}

	if len(result) <= 0 {
		return "{}"
	}
	return "{\n" + result + "}"
}
