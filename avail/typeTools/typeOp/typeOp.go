package typeOp

import "strings"

type Op uint64

const (
	// None is the base for all operators and still have operations.
	//
	//	| Operator          | Types    | Comment                                              |
	//	|-------------------|----------|------------------------------------------------------|
	//	| `zero() => x`     | `x T`    | creates a new zero value of this type                |
	//	| `new() => x`      | `x *T`   | creates a new zero value and returns a pointer to it |
	//	| `shallow(x) => y` | `x, y T` | shallow copy of value                                |
	None Op = 0

	// Add indicates that addition (numeric) or concatenation (string)
	// operations are available.
	//
	//	| Operator        | Types       | Comment                           |
	//	|-----------------|-------------|-----------------------------------|
	//	| `(x + y) => z`  | `x, y, z T` | addition / concatenation          |
	//	| `(x += y) => z` | `x, y, z T` | in place addition / concatenation |
	Add Op = 1 << iota

	// Arith indicates that arithmetic operations are available.
	// Arithmetic operations should always have additions operations too.
	//
	//	| Operator    | Types       | Comment                                  |
	//	|-------------|-------------|------------------------------------------|
	//	| `x = -y`    | `x, y T`    | negation unary operation                 |
	//	| `x = y - z` | `x, y, z T` | subtraction binary operation             |
	//	| `x = y * z` | `x, y, z T` | multiplication binary operation          |
	//	| `x = y / z` | `x, y, z T` | division binary operation                |
	//	| `x -= y`    | `x, y T`    | in place subtraction binary operation    |
	//	| `x *= y`    | `x, y T`    | in place multiplication binary operation |
	//	| `x /= y`    | `x, y T`    | in place division binary operation       |
	//	| `x++`       | `x T`       | increment unary operation                |
	//	| `x--`       | `x T`       | decrement unary operation                |
	Arith

	// Mod indicates that modulo operations are available.
	// Mod operations should always have additions and arithmetic operations.
	//
	//	| Operator    | Types       | Comment                          |
	//	|-------------|-------------|----------------------------------|
	//	| `x = y % z` | `x, y, z T` | modulo binary operation          |
	//	| `x %= y`    | `x, y T`    | in place modulo binary operation |
	Mod

	// Bitwise indicates that bitwise operations are available.
	// Bitwise operations should always have additions and arithmetic operations.
	//
	//	| Operator     | Types                   | Comment                          |
	//	|--------------|-------------------------|----------------------------------|
	//	| `x = ^y`     | `x, y T`                | bitwise invert unary operation   |
	//	| `x = y & z`  | `x, y, z T`             | bitwise AND binary operation     |
	//	| `x = y | z`  | `x, y, z T`             | bitwise OR binary operation      |
	//	| `x = y ^ z`  | `x, y, z T`             | bitwise XOR binary operation     |
	//	| `x = y &^ z` | `x, y, z T`             | bitwise AND NOT binary operation |
	//	| `x = y << z` | `x, y T, z untyped int` | left shift operation             |
	//	| `x = y >> z` | `x, y T, z untyped int` | right shift operation            |
	//	| `x &= y`     | `x, y T`                | in place bitwise AND             |
	//	| `x |= y`     | `x, y T`                | in place bitwise OR              |
	//	| `x ^= y`     | `x, y T`                | in place bitwise XOR             |
	//	| `x &^= y`    | `x, y T`                | in place bitwise AND NOT         |
	//	| `x <<= y`    | `x T, y untyped int`    | in place left shift operation    |
	//	| `x >>= y`    | `x T, y untyped int`    | in place right shift operation   |
	Bitwise

	// Len indicates that the length operator is available.
	// See [https://pkg.go.dev/builtin#len]
	//
	//	| Operator     | Types        | Comment                                          |
	//	|--------------|--------------|--------------------------------------------------|
	//	| `x = len(y)` | `x int, y T` | determines the length of a string, slice, or map |
	Len

	// Cap indicates that the capacity operator is available.
	// The capacity operator should have the length operator too.
	// See [https://pkg.go.dev/builtin#cap]
	//
	//	| Operator     | Types        | Comment                                   |
	//	|--------------|--------------|-------------------------------------------|
	//	| `x = cap(y)` | `x int, y T` | determines the capacity of a slice or map |
	Cap

	// IsNil indicates that nil checking is available.
	//
	//	| Operator       | Types         | Comment                            |
	//	|----------------|---------------|------------------------------------|
	//	| `x = y == nil` | `x bool, y T` | determines the value is nil or not |
	IsNil

	// Comparable indicates that comparison operators are available.
	// The comparable operators may or may not be paired with a nil checking operator.
	//
	//	| Operator     | Types            | Comment                                      |
	//	|--------------|------------------|----------------------------------------------|
	//	| `x = y == z` | `x bool, y, z T` | determines the value is equal to another     |
	//	| `x = y != z` | `x bool, y, z T` | determines the value is not equal to another |
	Comparable

	// Orderable indicates that orderable operators are available.
	// Orderable operations should always have comparable operations too.
	//
	//	| Operator     | Types            | Comment                                                  |
	//	|--------------|------------------|----------------------------------------------------------|
	//	| `x = y < z`  | `x bool, y, z T` | determines the value is less than another                |
	//	| `x = y <= z` | `x bool, y, z T` | determines the value is less than or equal to another    |
	//	| `x = y > z`  | `x bool, y, z T` | determines the value is greater than another             |
	//	| `x = y => z` | `x bool, y, z T` | determines the value is greater than or equal to another |
	Orderable

	// Ref indicates that referencing operations are available.
	//
	//	| Operator     | Types       | Comment                                            |
	//	|--------------|-------------|----------------------------------------------------|
	//	| `x = &y`     | `x *T, y T` | creates a pointer to the given value               |
	//	| `x = new(y)` | `x T, y *T` | creates a pointer to a new copy of the given value |
	Ref

	// Make indicates that make operations are available.
	// The make can build slices, maps, and channels.
	//
	// If the first parameter is zero or not given, the default is used.
	// The second parameter is only defined for slices by Go so shouldn't be called by any type other
	// than a slice. However, for other types the second parameter can simply be ignored.
	//
	// See [https://pkg.go.dev/builtin#make]
	//
	//	| Operator     | Types               | Comment                                    |
	//	|--------------|---------------------|--------------------------------------------|
	//	| `make()`     |                     | creates a type with default length         |
	//	| `make(x)`    | `x unsigned int`    | creates a type with a length               |
	//	| `make(x, y)` | `x, y unsigned int` | creates a slice with a length and capacity |
	Make

	// GetIndex indicates that index getter operations are available.
	// Index getter operations should always have the length operation.
	// The index getter operations for strings, slices, and arrays use untyped ints,
	// and for maps use the map's key type.
	//
	// Go does not define the two results operation for a non-map, so it will not be called in normal
	// Go code however, may be used to not panic on out-of-bounds for slices and arrays.
	//
	//	| Operator      | Types                   | Comment                                             |
	//	|---------------|-------------------------|-----------------------------------------------------|
	//	| `x = y[z]`    | `x E, y T, z K`         | gets the value at the given index or key            |
	//	| `x, y = z[w]` | `x E, y bool, z T, w K` | gets the value and exists at the given index or key |
	GetIndex

	// SetIndex indicates that index setter operations are available.
	// Index setter operations should always have index getter operations and the length operation.
	// This is separate from the index getters operations because strings are immutable
	// and therefore will not have the setter operations on them.
	//
	// Any type that can have indices set on it can also have all the indices cleared.
	// See [https://pkg.go.dev/builtin#clear]
	//
	//	| Operator   | Types           | Comment                                             |
	//	|------------|-----------------|-----------------------------------------------------|
	//	| `x[y] = z` | `x T, y K, z E` | sets the value at the given index or key            |
	//	| `clear(x)` | `x T`           | clears all the values from the slice, array, or map |
	SetIndex

	// RefIndex indicates that the element pointer operations are available.
	// For slices, the pointer is on the internal array such that if the slice is grown and
	// a new array is allocated, the pointer remains pointing at the original array.
	// This is not available for maps.
	//
	//	| Operator    | Types            | Comment                                        |
	//	|-------------|------------------|------------------------------------------------|
	//	| `x = &y[z]` | `x *E, y T, z K` | gets a pointer to the value at the given index |
	RefIndex

	// Slice indicates that non-capacity slice operations are available.
	// Slice operations should always have index getter operations and the length operation.
	// Being able to slice also indicates that the type may be used as an array of the element type.
	//
	// This works for slices, arrays, and strings. For arrays and slices, the slice returns as `[]E`,
	// however, for strings the slice will return a substring of type `string` yet will still be
	// able to be used as `[]byte` when needed.
	//
	// The `copyTo` function can be used in the builtin functions `append` and `copy`.
	//
	//	| Operator       | Types                           | Comment                                                                |
	//	|----------------|---------------------------------|------------------------------------------------------------------------|
	//  | `x = y[z:w]`   | `x ~[]E, y T, z, w untyped int` | creates a slice of the type                                            |
	//	| `x = []E(y)`   | `x = []E, y = T`                | converts the slice, array, or string into a slice of its element types |
	//	| `copyTo(x, y)` | `x = T, y = []E`                | copies this value over the given section of the given slice            |
	Slice

	Slice3 // s[x:y:z] for (slice, array)

	Range // for _=range x, for y=range x

	Range2 // for _,_=range x, for y,_=range x, for y,z=range x

	Recv // x<-y, x,y<-z

	// See [https://pkg.go.dev/builtin#close]
	Send // x->y, close(y)

	Complex // complex(x,y)

	RealImag // real(x), imag(x)
)

func (op Op) All(other Op) bool { return op&other == other }
func (op Op) Any(other Op) bool { return op&other != None }

func (op Op) String() string {
	if op == None {
		return "none"
	}

	parts := []string{}
	add := func(other Op, name string) {
		if op.All(other) {
			parts = append(parts, name)
		}
	}

	add(Add, `Add`)
	add(Arith, `Arith`)
	add(Bitwise, `Bitwise`)
	add(ByteSlice, `ByteSlice`)
	add(Cap, `Cap`)
	add(Clear, `Clear`)
	add(Comparable, `Comparable`)
	add(Complex, `Complex`)
	add(Deref, `Deref`)
	add(GetIndex, `GetIndex`)
	add(GetIndex2, `GetIndex2`)
	add(IsNil, `IsNil`)
	add(Len, `Len`)
	add(Make, `Make`)
	add(Make3, `Make3`)
	add(Mod, `Mod`)
	add(Orderable, `Orderable`)
	add(Range, `Range`)
	add(Range2, `Range2`)
	add(RealImag, `RealImag`)
	add(Recv, `Recv`)
	add(Ref, `Ref`)
	add(RefIndex, `RefIndex`)
	add(Send, `Send`)
	add(SetIndex, `SetIndex`)
	add(Slice, `Slice`)
	add(Slice3, `Slice3`)

	return strings.Join(parts, `|`)
}
