package constructs

var _ Type = (*MapType)(nil)

// MapType for storing the types of maps.
type MapType struct {

	// Key is the key type for this map.
	Key Type

	// Value gets the value type for this map.
	Value Type
}

// Map creates a new map type for the given key and value types.
func Map(key Type, value Type) *MapType {
	return &MapType{
		Key:   key,
		Value: value,
	}
}

// Equals determins if the given type is the same as this type.
func (t *MapType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*MapType)
	if !ok {
		return false
	}
	if !t.Key.Equals(tother.Key) {
		return false
	}
	if !t.Value.Equals(tother.Value) {
		return false
	}
	return true
}

// String gets the name for this type.
func (t *MapType) String() string {
	return "Map<" + t.Key.String() + ", " + t.Value.String() + ">"
}
