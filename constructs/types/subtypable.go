package types

var _ Type = (SubtypableType)(nil)

// SubtypableType is a type add-on for types which contain other types
// such as classes and interfaces.
type SubtypableType interface {

	// Find finds the member, function type, of subtype
	// with the given name inside this type.
	// If no subtype by that name is found, false is returned.
	Find(name string) (Type, bool)

	// String gets the identifier name for this type.
	String() string
}

// FindSubtype finds the member, function type, of subtype
// with the given name inside the given type.
// If no subtype by that name is found, false is returned.
func FindSubtype(t Type, name string) (Type, bool) {
	if t != nil {
		if t2, ok := t.(SubtypableType); ok {
			return t2.Find(name)
		}
	}
	return nil, false
}
