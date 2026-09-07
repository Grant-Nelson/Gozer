package types

// Type is a Gozer IR type.
//
// 	| Operator       | Argument Types      | Comment                                              |
// 	|----------------|---------------------|------------------------------------------------------|
// 	| `zero() => x`  | `x => T`            | creates a new zero value of this type                |
// 	| `new() => x`   | `x => *T`           | creates a new zero value and returns a pointer to it |
// 	| `copy(x) => y` | `x, y => T`         | shallow copy of value                                |
//
type Type interface {
	IsType()
}

type FixedSize interface {
	Type
	BitSize() int
}
