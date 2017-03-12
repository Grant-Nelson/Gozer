package constructs

var _ Construct = (Type)(nil)

// Type is the interface for all types.
type Type interface {

	// String gets the identifier name for this type.
	String() string
}

// FindSubtype finds the member, function type, of subtype
// with the given name inside the given type.
// If no subtype by that name is found, false is returned.
func FindSubtype(t Type, name string) (Type, bool) {
	if t != nil {
		switch t2 := t.(type) {
		case *StructureType:
			return t2.Find(name)
		case *ClassType:
			return t2.Find(name)
		case *InterfaceType:
			return t2.Find(name)
		case *PackageType:
			return t2.Find(name)
		}
	}
	return nil, false
}
