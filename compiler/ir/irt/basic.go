package irt

type Basic uint16

var _ Type = Basic(0)

const (
	Invalid Basic = 0

	isUnsized = 0x00_01
	is1Bit    = 0x00_02
	is8Bit    = 0x00_04
	is16Bit   = 0x00_08
	is32Bit   = 0x00_10
	is64Bit   = 0x00_20
	is128Bit  = 0x00_40

	isBoolean  = 0x01_00
	isInteger  = 0x02_00
	isUnsigned = 0x04_00
	isFloat    = 0x08_00
	isComplex  = 0x10_00
	isString   = 0x20_00
	isUntyped  = 0x40_00
	isPointer  = 0x80_00

	isOrdered   = isInteger | isFloat | isString
	isNumeric   = isInteger | isFloat | isComplex
	isConstType = isBoolean | isNumeric | isString

	Bool = is1Bit | isBoolean

	Int   = isUnsized | isInteger
	Int8  = is8Bit | isInteger
	Int16 = is16Bit | isInteger
	Int32 = is32Bit | isInteger
	Int64 = is64Bit | isInteger

	Uint   = isUnsized | isInteger | isUnsigned
	Uint8  = is8Bit | isInteger | isUnsigned
	Uint16 = is16Bit | isInteger | isUnsigned
	Uint32 = is32Bit | isInteger | isUnsigned
	Uint64 = is64Bit | isInteger | isUnsigned

	Uintptr = isUnsized | isPointer | isInteger | isUnsigned

	Float32 = is32Bit | isFloat
	Float64 = is64Bit | isFloat

	Complex64  = is64Bit | isComplex
	Complex128 = is128Bit | isComplex

	String        = isUnsized | isString
	UnsafePointer = isUnsized | isPointer

	UntypedBool    = isBoolean | isUntyped
	UntypedInt     = isInteger | isUntyped
	UntypedRune    = isInteger | isUntyped | is32Bit
	UntypedFloat   = isFloat | isUntyped
	UntypedComplex = isComplex | isUntyped
	UntypedString  = isString | isUntyped
	UntypedNil     = isUntyped
)

func (b Basic) Name() string     { return b.String() }
func (b Basic) Underlying() Type { return b }

func (b Basic) String() string {
	switch b {
	case Bool:
		return `bool`
	case Int:
		return `int`
	case Int8:
		return `int8`
	case Int16:
		return `int16`
	case Int32:
		return `int32`
	case Int64:
		return `int64`
	case Uint:
		return `uint`
	case Uint8:
		return `uint8`
	case Uint16:
		return `uint16`
	case Uint32:
		return `uint32`
	case Uint64:
		return `uint64`
	case Uintptr:
		return `uintptr`
	case Float32:
		return `float32`
	case Float64:
		return `float64`
	case Complex64:
		return `complex64`
	case Complex128:
		return `complex128`
	case String:
		return `string`
	case UnsafePointer:
		return `Pointer`
	case UntypedBool:
		return `untyped bool`
	case UntypedInt:
		return `untyped int`
	case UntypedRune:
		return `untyped rune`
	case UntypedFloat:
		return `untyped float`
	case UntypedComplex:
		return `untyped complex`
	case UntypedString:
		return `untyped string`
	case UntypedNil:
		return `untyped nil`
	default: // Invalid
		return `invalid type`
	}
}

func (b Basic) IsValid() bool {
	switch b {
	case UntypedBool, Bool,
		UntypedInt, Int, Int8, Int16, Int32, Int64,
		UntypedRune, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr,
		UntypedFloat, Float32, Float64,
		UntypedComplex, Complex64, Complex128,
		UntypedString, String,
		UntypedNil, UnsafePointer:
		return true
	default: // Invalid
		return false
	}
}

func (b Basic) IsUnsized() bool   { return b&isUnsized != 0 }
func (b Basic) Is8Bit() bool      { return b&is8Bit != 0 }
func (b Basic) Is16Bit() bool     { return b&is16Bit != 0 }
func (b Basic) Is32Bit() bool     { return b&is32Bit != 0 }
func (b Basic) Is64Bit() bool     { return b&is64Bit != 0 }
func (b Basic) Is128Bit() bool    { return b&is128Bit != 0 }
func (b Basic) IsBoolean() bool   { return b&isBoolean != 0 }
func (b Basic) IsInteger() bool   { return b&isInteger != 0 }
func (b Basic) IsUnsigned() bool  { return b&isUnsigned != 0 }
func (b Basic) IsFloat() bool     { return b&isFloat != 0 }
func (b Basic) IsComplex() bool   { return b&isComplex != 0 }
func (b Basic) IsString() bool    { return b&isString != 0 }
func (b Basic) IsUntyped() bool   { return b&isUntyped != 0 }
func (b Basic) IsOrdered() bool   { return b&isOrdered != 0 }
func (b Basic) IsNumeric() bool   { return b&isNumeric != 0 }
func (b Basic) IsConstType() bool { return b&isConstType != 0 }
