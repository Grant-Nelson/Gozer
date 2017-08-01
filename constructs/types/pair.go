package types

var _ Type = (*PairType)(nil)
var _ IndexableType = (*PairType)(nil)

// PairType for storing the types of pairs.
type PairType struct {

	// Key is the key type for this pair.
	Key Type

	// Value gets the value type for this pair.
	Value Type
}

// Pair creates a new pair type for the given key and value types.
func Pair(key Type, value Type) *PairType {
	return &PairType{
		Key:   key,
		Value: value,
	}
}

// ElementType gets the indexable subtype from this type.
func (t *PairType) ElementType() Type {
	if t == nil {
		return nil
	}
	return t.Value
}

// String gets the name for this type.
func (t *PairType) String() string {
	if t == nil {
		return nilStr
	}
	return "pair[" + ToString(t.Key) + "]" + ToString(t.Value)
}
