package types

// nilStr is the string to use for nil values.
const nilStr = "nil"

// Type is the interface for all types.
type Type interface {

	// String gets the identifier name for this type.
	String() string
}

// ToString creates a string for the given construct.
func ToString(t Type) string {
	if t == nil {
		return nilStr
	}
	return t.String()
}

// LookupType gets the type for the given type name.
func LookupType(typeName string) Type {
	switch typeName {
	case "bool":
		return Bool()
	case "byte":
		return Byte()
	case "complex64":
		return Complex64()
	case "complex128":
		return Complex128()
	case "float32":
		return Float32()
	case "float64":
		return Float64()
	case "int":
		return Int()
	case "int8":
		return Int8()
	case "int16":
		return Int16()
	case "int32":
		return Int32()
	case "int64":
		return Int64()
	case "rune":
		return Rune()
	case "string":
		return String()
	case "uint":
		return UInt()
	case "uint8":
		return UInt8()
	case "uint16":
		return UInt16()
	case "uint32":
		return UInt32()
	case "uint64":
		return UInt64()
	case "variant":
		return Variant()
	case "void":
		return Void()
	default:
		return nil
	}
}
