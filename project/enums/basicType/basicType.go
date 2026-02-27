package basicType

// BasicType is a basic type, not a composite type.
//
// Aliased types, i.e. `rune` and `byte`, are not defined.
// Simply use the type they alias to.
// Untyped constants are treated like the largest sized equivalent
// and casts should be added as needed.
//
// https://go.dev/ref/spec#Types
type BasicType int

const (
	undefined BasicType = iota

	Uint8  // unsigned  8-bit integer (0 to 255)
	Uint16 // unsigned 16-bit integer (0 to 65535)
	Uint32 // unsigned 32-bit integer (0 to 4294967295)
	Uint64 // unsigned 64-bit integer (0 to 18446744073709551615)
	Uint   // unsigned integer, either 32 or 64 bits

	Int8  // signed  8-bit integer (-128 to 127)
	Int16 // signed 16-bit integer (-32768 to 32767)
	Int32 // signed 32-bit integer (-2147483648 to 2147483647)
	Int64 // signed 64-bit integer (-9223372036854775808 to 9223372036854775807)
	Int   // signed integer, either 32 or 64 bits

	Float32 // IEEE 754 32-bit floating-point number
	Float64 // IEEE 754 64-bit floating-point number

	Complex64  // complex number with a float32 real and imaginary part
	Complex128 // complex number with a float64 real and imaginary part

	Bool    // a boolean value, i.e. true or false
	String  // a string value, UTF-8, with underlying `[]byte`
	Uintptr // an unsigned integer large enough to store a pointer value
)

const max = Uintptr

var names []string

func init() {
	nameMap := map[BasicType]string{
		undefined: `UndefinedBasicType`,

		Uint8:  `uint8`,
		Uint16: `uint16`,
		Uint32: `uint32`,
		Uint64: `uint64`,
		Uint:   `uint`,

		Int8:  `int8`,
		Int16: `int16`,
		Int32: `int32`,
		Int64: `int64`,
		Int:   `int`,

		Float32: `float32`,
		Float64: `float64`,

		Complex64:  `complex64`,
		Complex128: `complex128`,

		Bool:    `bool`,
		String:  `string`,
		Uintptr: `uintptr`,
	}
	names = make([]string, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
	}
}

func (bt BasicType) Valid() bool {
	// undefined is not valid
	return bt > undefined && bt <= max
}

func (bt BasicType) String() string {
	if bt >= undefined && bt <= max {
		return names[bt]
	}
	return `UnknownBasicType`
}
