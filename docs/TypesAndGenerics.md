# Types and Generics

The resulting transpiled code need to function as closely as we can to
the source Go code. This means that some of the quirks of Go must be carried
through to the output.

## Interfaces

For example, interfaces in many languages are different from Go's interfaces.
Go's interfaces are objects that contain type information

TODO: Finish

## Pointers

TODO: Finish

## Generics

Gozer will implement generics based on ideas from the Golang
[Generics implementation - GC Shape Stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-gcshape.md)
design document. Obviously most target languages won't have garbage collect (GC)
since most target languages like TS are based on JS.
This means we don't have to worry about footprint sizes for GC,
see [Footprints](#footprints) for more information.
However, the stenciling design gives us some useful ideas for our design.
The main idea for how we will handle generics is to create a single generic
method with additional parameters for type dictionaries when needed.

### No Monomorphic Instantiations

We do not want to make monomorphic instantiations from generic types since that
will make it difficult to perform the partial compiles of packages.
For example, if we have a `List[T any]` defined in a `list` package, and
the `cat` package depends on `list` to define a `List[Cat]` instantiation of
`List[T any]`. If `List` accesses information that is not exported from the
package we can't define the `List[Cat]` instantiation anywhere but inside
of `list`, meaning to do monomorphic instantiation we'd have to compile all
packages at the same time and not be able to compile packages like `list`
without knowing the packages that import it. This is undesirable so we
must handle generics without creating monomorphic instantiations.

### Type Dictionaries

There are a finite number of basic types, such as `int8`, `int`, `uint64`, `bool`,
`float64`, and `complex64`. When these types are used in type constraints
on type parameters for generic declarations, they indicate what can possibly
be done with that type, i.e. what operators can be called on that type.
We don't care if the type does or does not call one of the operators since
it may call another function or function pointer parameter that does call one
of those types. We base what is possible on the type constraints without
needing to inspect the body of a generic function or methods of a generic type.

- We can't run the target operator on all of those types. For example, the TS
  binary addition operator, `x + y`, can not be run on an `uint64` since TS
  integers are limited to the size of Javascript integers which is less than
  64 bits. We will have to create an object in TS to represent a `uint64` with
  a high and low bit (same for `int64`, `complex64`, and `complex128`).
- Since TS will not allow operator overloading we can not perform a `+` on our
  `uint64` implementation in TS. Instead we define a function `add(x, y uint64)`
  that will properly perform the addition.
- Since we do not want to make multiple instantiations from a generic function
  that accepts a `int` and `int64`, on instantiation that calls `add(x, y int64)`
  and one that performs a `x + y`, we make one instantiation that calls an `add`
  function that is specific to the type parameter being passed in.
- By adding type dictionaries as an initial parameter to a function, similar
  to a closure captured variable or a receiver, we do not have to box and unbox
  basic types like Java does, e.g. `int` boxed in `Integer`.

The basic type dictionaries have the following interfaces for operations:

| Operator Group | Operators | Applies To |
|:--|:-:|:--|
| Arithmetic | `+` (binary/unary), `-` (binary/unary), `*`, `/` | integers, floats, complex |
| Addition / Concatenation | `+` (binary) | strings, integers, floats, complex |
| Remainder | `%` | integers only |
| Bitwise | `&` (binary), `\|`, `^`, `&^` | integers only |
| Shift | `<<`, `>>` | integers only |
| Comparison | `==`, `!=` | comparable types |
| Ordering | `<`, `<=`, `>`, `>=` | ordered types |
| Logical | `&&`, `\|\|`, `!` | booleans only |
| Reference | `&` (unary) | pointers only |

If we have a type like `type Foo int` with `func (f Foo) inc()` then that
type will implement Arithmetic, Addition, Remainder, Bitwise, Shift,
Comparison, Ordering, and `interface { inc() }`. That means it can be passed
in as a type argument for a generic with the constraints `~int`, `~int|~string`,
`interface { inc() }`, `any`, etc.

It is possible with this setup that a type will implement an interface for
the generic type that was not defined in the constraints. To help prevent this
all operation functions are defined with a preceding `$` or some character not
allowed in identifiers for Go but is allowed in the target language.
We also rely on Go's type checker to insure the types are valid.
If the function is exported and called directly from TS with the wrong type,
that is fine since the function will still work if the interface is implemented.

TODO: Finish

This would allow us to create additional type definitions specific to the target
language. For example TS could have a `num` type definition defined so that
when calling an exported function from TS the type definition can be looked up
and used.

If we have a generic function like `foo[T int|float](x T)` then 

TODO: Finish

## Footprints

We will still keep some footprint information for reflection and some cast

TODO: Finish
