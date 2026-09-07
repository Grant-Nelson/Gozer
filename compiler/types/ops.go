package types

// AddOps indicates that addition (numeric) or concatenation (string)
// operations are available.
//
//	| Operator        | Types         | Comment                                               |
//	|-----------------|---------------|-------------------------------------------------------|
//	| `(x + y) => z`  | `x, y, z = T` | addition / concatenation                              |
//	| `(x += y) => z` | `x, y, z = T` | in place addition / concatenation returns resulting x |
type AddOps interface {
	Type
	HasAddOps()
}

// ArithmeticOps indicates that arithmetic operations are available.
// Arithmetic operations always have additions operations.
//
//	| Operator        | Types         | Comment                                                      |
//	|-----------------|---------------|--------------------------------------------------------------|
//	| `(-x) => y `    | `x, y = T`    | negation unary operation                                     |
//	| `(x - y) => z`  | `x, y, z = T` | subtraction binary operation                                 |
//	| `(x * y) => z`  | `x, y, z = T` | multiplication binary operation                              |
//	| `(x / y) => z`  | `x, y, z = T` | division binary operation                                    |
//	| `(x -= y) => z` | `x, y, z = T` | in place subtraction binary operation returns resulting x    |
//	| `(x *= y) => z` | `x, y, z = T` | in place multiplication binary operation returns resulting x |
//	| `(x /= y) => z` | `x, y, z = T` | in place division binary operation returns resulting x       |
//	| `(x++) => y`    | `x, y = T`    | increment unary operation returns resulting x                |
//	| `(x--) => y`    | `x, y = T`    | decrement unary operation returns resulting x                |
type ArithmeticOps interface {
	AddOps
	HasArithmeticOps()
}

// ModuloOps indicates that modulo operations are available.
// Bitwise operations always have additions and arithmetic operations.
//
//	| Operator        | Types         | Comment                                                      |
//	|-----------------|---------------|--------------------------------------------------------------|
//	| `(x % y) => z`  | `x, y, z = T` | modulo binary operation                                      |
//	| `(x %= y) => z` | `x, y, z = T` | in place modulo binary operation returns resulting x         |
type ModuloOps interface {
	ArithmeticOps
	HasModuloOps()
}

// BitwiseOps indicates that bitwise operations are available.
// Bitwise operations always have additions and arithmetic operations.
//
//	| Operator         | Types                     | Comment                                            |
//	|------------------|---------------------------|----------------------------------------------------|
//	| `(^x) => y`      | `x, y = T`                | bitwise invert unary operation                     |
//	| `(x & y) => z`   | `x, y, z = T`             | bitwise AND binary operation                       |
//	| `(x | y) => z`   | `x, y, z = T`             | bitwise OR binary operation                        |
//	| `(x ^ y) => z`   | `x, y, z = T`             | bitwise XOR binary operation                       |
//	| `(x &^ y) => z`  | `x, y, z = T`             | bitwise AND NOT binary operation                   |
//	| `(x &= y) => z`  | `x, y, z = T`             | in place bitwise AND returns resulting x           |
//	| `(x |= y) => z`  | `x, y, z = T`             | in place bitwise OR returns resulting x            |
//	| `(x ^= y) => z`  | `x, y, z = T`             | in place bitwise XOR returns resulting x           |
//	| `(x &^= y) => z` | `x, y, z = T`             | in place bitwise AND NOT returns resulting x       |
//	| `(x << y) => z`  | `x, z = T`, `y = integer` | left shift operation                               |
//	| `(x >> y) => z`  | `x, z = T`, `y = integer` | right shift operation                              |
//	| `(x <<= y) => z` | `x, z = T`, `y = integer` | in place left shift operation returns resulting x  |
//	| `(x >>= y) => z` | `x, z = T`, `y = integer` | in place right shift operation returns resulting x |
type BitwiseOps interface {
	ArithmeticOps
	HasBitwiseOps()
}

// TODO: Move to sliceable
// ByteSliceOps indicates that byte slice operations are available.
//
//	| Operator         | Argument Types        | Comment                                            |
//	|------------------|-----------------------|----------------------------------------------------|
//	| `[]byte(x) => y` | `x = T`, `y = []byte` |
//	| `copyTo(x, y)`   | `x = T`, `y = []byte` |
type ByteSliceOps interface {
	Type
	HasByteSliceOps()
}

// TODO: Comment
// See [https://pkg.go.dev/builtin#len]
type LenOps interface {
	Type
	HasLenOps()
}

// TODO: Comment
// See [https://pkg.go.dev/builtin#cap]
type CapOps interface {
	LenOps
	HasCapOps()
}

// TODO: Comment
type IsNilOps interface {
	Type
	HasIsNilOps()
}

// TODO: Comment
type ComparableOps interface {
	Type
	HasComparableOps()
}

// TODO: Comment
type OrderableOps interface {
	ComparableOps
	HasOrderableOps()
}

// TODO: Comment
type DerefOps interface {
	Type
	HasDerefOps()
}

// TODO: Comment
// 	| `new(x) => y`  | `x => T`, `y => *T` | creates a pointer to the given value                 |
type RefOps interface {
	Type
	HasRefOps()
}

// TODO: Comment
type MakeOps interface {
	Type
	HasMakeOps()
}

// TODO: Comment
type Make3Ops interface {
	Type
	HasMake3Ops()
}

// TODO: Comment
type GetIndexOps interface {
	Type
	HasGetIndexOps()
}

// TODO: Comment
type TryGetIndexOps interface {
	Type
	HasTryGetIndexOps()
}

// TODO: Comment
type RefIndexOps interface {
	Type
	HasRefIndexOps()
}

// TODO: Comment
// Add clear
type SetIndexOps interface {
	Type
	HasSetIndexOps()
}

// TODO: Comment
type Slice2Ops interface {
	Type
	HasSlice2Ops()
}

// TODO: Comment
type Slice3Ops interface {
	Type
	HasSlice3Ops()
}

// TODO: Comment
type RecvOps interface {
	Type
	HasRecvOps()
}

// TODO: Comment
type SendOps interface {
	Type
	HasSendOps()
}

// TODO: Comment
type Range1Ops interface {
	Type
	HasRange1Ops()
}

// TODO: Comment
type Range2Ops interface {
	Type
	HasRange2Ops()
}

// TODO: Comment
type ComplexOps interface {
	Type
	ComplexType() Complex
	HasComplexOps()
}

// TODO: Comment
type RealImagOps interface {
	Type
	RealImagType() Float
	HasRealImagOps()
}
