package constructs

var _ Type = (*BaseType)(nil)

// BaseType for storing the types which are low-level built-in types.
type BaseType struct {
	typeName string
}

// newBaseType creates a new base type.
func newBaseType(typeName string) *BaseType {
	return &BaseType{
		typeName: typeName,
	}
}

// Equals determins if the given type is the same as this type.
func (t *BaseType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*BaseType)
	if !ok {
		return false
	}
	if t.typeName != tother.typeName {
		return false
	}
	return true
}

// String gets the name for this type.
func (t *BaseType) String() string {
	if t == nil {
		return nilStr
	}
	return t.typeName
}

//==========================================================================

// Void is the non-type return value.
func Void() *BaseType {
	return newBaseType("void")
}

// Bool is the set of boolean values, true and false.
func Bool() *BaseType {
	return newBaseType("bool")
}

// Byte is an alias for uint8 and is equivalent to uint8 in all ways.
// It is used, by convention, to distinguish byte values
// from 8-bit unsigned integer values.
func Byte() *BaseType {
	return newBaseType("byte")
}

// Float32 is the set of all IEEE-754 32-bit floating-point numbers.
func Float32() *BaseType {
	return newBaseType("float32")
}

// Float64 is the set of all IEEE-754 64-bit floating-point numbers.
func Float64() *BaseType {
	return newBaseType("float64")
}

// Int is a signed integer type that is at least 32 bits in size.
// It is a distinct type, however, and not an alias for, say, int32.
func Int() *BaseType {
	return newBaseType("int")
}

// Int16 is the set of all signed 16-bit integers.
// Range: -32768 through 32767.
func Int16() *BaseType {
	return newBaseType("int16")
}

// Int32 is the set of all signed 32-bit integers.
// Range: -2147483648 through 2147483647.
func Int32() *BaseType {
	return newBaseType("int32")
}

// Int64 is the set of all signed 64-bit integers.
// Range: -9223372036854775808 through 9223372036854775807.
func Int64() *BaseType {
	return newBaseType("int64")
}

// Int8 is the set of all signed 8-bit integers.
// Range: -128 through 127.
func Int8() *BaseType {
	return newBaseType("int8")
}

// Rune is an alias for int32 and is equivalent to int32 in all ways.
// It is used, by convention, to distinguish character values from integer values.
func Rune() *BaseType {
	return newBaseType("rune")
}

// String is the set of all strings of 8-bit bytes, conventionally
// but not necessarily representing UTF-8-encoded text. A string may be empty,
// but not nil. Values of string type are immutable.
func String() *BaseType {
	return newBaseType("string")
}

// UInt is an unsigned integer type that is at least 32 bits in size.
// It is a distinct type, however, and not an alias for, say, uint32.
func UInt() *BaseType {
	return newBaseType("uint")
}

// UInt16 is the set of all unsigned 16-bit integers.
// Range: 0 through 65535.
func UInt16() *BaseType {
	return newBaseType("uint16")
}

// UInt32 is the set of all unsigned 32-bit integers.
// Range: 0 through 4294967295.
func UInt32() *BaseType {
	return newBaseType("uint32")
}

// UInt64 is the set of all unsigned 64-bit integers.
// Range: 0 through 18446744073709551615.
func UInt64() *BaseType {
	return newBaseType("uint64")
}

// UInt8 is the set of all unsigned 8-bit integers.
// Range: 0 through 255.
func UInt8() *BaseType {
	return newBaseType("uint8")
}
