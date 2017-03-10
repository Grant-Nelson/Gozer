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

// String gets the name for this type.
func (t *MapType) String() string {
	if t == nil {
		return nilStr
	}
	return "Map<" + ToString(t.Key) + ", " + ToString(t.Value) + ">"
}
