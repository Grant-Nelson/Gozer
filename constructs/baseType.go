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

// String gets the name for this type.
func (t *BaseType) String() string {
	if t == nil {
		return nilStr
	}
	return t.typeName
}

//==========================================================================

// Bool is the set of boolean values, true and false.
// https://golang.org/src/builtin/builtin.go?h=make#L14
func Bool() *BaseType {
	return newBaseType("bool")
}

// Byte is an alias for uint8 and is equivalent to uint8 in all ways.
// It is used, by convention, to distinguish byte values
// from 8-bit unsigned integer values.
// https://golang.org/src/builtin/builtin.go?h=make#L88
func Byte() *BaseType {
	return newBaseType("byte")
}

// Complex64 is the set of all IEEE-754 32-bit floating-point numbers
// followed by another IEEE-754 32-bit float-point number as an imaginary.
// https://golang.org/src/builtin/builtin.go?h=make#L62
func Complex64() *BaseType {
	return newBaseType("complex64")
}

// Complex128 is the set of all IEEE-754 64-bit floating-point numbers
// followed by another IEEE-754 64-bit float-point number as an imaginary.
// https://golang.org/src/builtin/builtin.go?h=make#L66
func Complex128() *BaseType {
	return newBaseType("complex128")
}

// Float32 is the set of all IEEE-754 32-bit floating-point numbers.
// https://golang.org/src/builtin/builtin.go?h=make#L119
func Float32() *BaseType {
	return newBaseType("float32")
}

// Float64 is the set of all IEEE-754 64-bit floating-point numbers.
// https://golang.org/src/builtin/builtin.go?h=make#L58
func Float64() *BaseType {
	return newBaseType("float64")
}

// Int is a signed integer type that is at least 32 bits in size.
// It is a distinct type, however, and not an alias for, say, int32.
// https://golang.org/src/builtin/builtin.go?h=make#L75
func Int() *BaseType {
	return newBaseType("int")
}

// Int8 is the set of all signed 8-bit integers.
// Range: -128 through 127.
// https://golang.org/src/builtin/builtin.go?h=make#L40
func Int8() *BaseType {
	return newBaseType("int8")
}

// Int16 is the set of all signed 16-bit integers.
// Range: -32768 through 32767.
// https://golang.org/src/builtin/builtin.go?h=make#L44
func Int16() *BaseType {
	return newBaseType("int16")
}

// Int32 is the set of all signed 32-bit integers.
// Range: -2147483648 through 2147483647.
// https://golang.org/src/builtin/builtin.go?h=make#L48
func Int32() *BaseType {
	return newBaseType("int32")
}

// Int64 is the set of all signed 64-bit integers.
// Range: -9223372036854775808 through 9223372036854775807.
// https://golang.org/src/builtin/builtin.go?h=make#L52
func Int64() *BaseType {
	return newBaseType("int64")
}

// Rune is an alias for int32 and is equivalent to int32 in all ways.
// It is used, by convention, to distinguish character values from integer values.
// https://golang.org/src/builtin/builtin.go?h=make#L92
func Rune() *BaseType {
	return newBaseType("rune")
}

// UInt is an unsigned integer type that is at least 32 bits in size.
// It is a distinct type, however, and not an alias for, say, uint32.
// https://golang.org/src/builtin/builtin.go?h=make#L79
func UInt() *BaseType {
	return newBaseType("uint")
}

// UInt8 is the set of all unsigned 8-bit integers.
// Range: 0 through 255.
// https://golang.org/src/builtin/builtin.go?h=make#L24
func UInt8() *BaseType {
	return newBaseType("uint8")
}

// UInt16 is the set of all unsigned 16-bit integers.
// Range: 0 through 65535.
// https://golang.org/src/builtin/builtin.go?h=make#L28
func UInt16() *BaseType {
	return newBaseType("uint16")
}

// UInt32 is the set of all unsigned 32-bit integers.
// Range: 0 through 4294967295.
// https://golang.org/src/builtin/builtin.go?h=make#L32
func UInt32() *BaseType {
	return newBaseType("uint32")
}

// UInt64 is the set of all unsigned 64-bit integers.
// Range: 0 through 18446744073709551615.
// https://golang.org/src/builtin/builtin.go?h=make#L36
func UInt64() *BaseType {
	return newBaseType("uint64")
}

// Variant is a stand-in for any type, but typically represents the same
// type for any given function invocation. This takes the place of
// the "Type" type in the Go documentation.
// https://golang.org/src/builtin/builtin.go?h=make#L106
func Variant() *BaseType {
	return newBaseType("variant")
}

// Void is the non-type return value.
func Void() *BaseType {
	return newBaseType("void")
}
