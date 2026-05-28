# Types and Generics

The resulting transpiled code need to function as closely as we can to
the source Go code. This means that some of the quirks of Go must be carried
through to the output.

## Interfaces and Type information

For example, interfaces in many languages are different from Go's interfaces.
In fact, in these documents, "interface" is used as would be defined in an OO
language like C#, where the interface is a collection of signatures only.
However, Go interfaces are objects that contain type information so when
discussing Go interfaces we will try to always use "Go interface(s)".

Go interfaces will be transpiled into two parts. It will create an actual
interface as a collection of singatures and an object that impements the
interface and contains a reference to a type and an object.
The type and object may be nil (i.e. the `untyped nil`), or just the object
may be nil (i.e. a typed nil). This allows for a non-nil Go interface that
contains a nil object.

When casting from an Go interface to another Go interface or to another type,
the Go interface is unwrapped and the internal type is used to determine if
the cast can be completed.

When a function is called on a Go interface, it will call the receiverless
instance of the function on the contained type with the contained object.
The receiverless function where the receiver is moved to be the first argument
of the function. That way if the Go interface contains a nil object, it can still
be passed in as nil.

The type information will be an immutable object typically defined once as a
singleton. The type information will contain all the receiverless function
definitions that when called will push a `CallNode` onto the shedular call stack
with any closure, reciever, type parameters, and type dictionaries, then
will call the first block for that function with the arguments passed in.
The type information may also extend a [type dictionary](#type-dictionaries)
for structs like `type Foo int` that may also have methods attached to them.

The type information can also return the type information for any field or
method. The type information for a method will define the signature such that
the method can be used in a function pointer for when we are passing a method
into another function that has the function pointer as a parameter. How
the type information looks depends on what is needed by Go and the target
language to approparately define a type whilst leveraging as much of the
target language's type system as possible. This type information is used by
reflections and contains the information to show what Go would output when
printing the type (e.g. using `%T`).

The following is an example of how Go interfaces can be nil or non-nil while
still containing a nil object:

```Go
package main

import "fmt"

func main() {
  var b Bar
  fmt.Printf("Bar: %t => ", b != nil)
  try(func() { b.Baz() })

  var f *Foo
  fmt.Printf("Foo: %t => ", f != nil)
  try(func() { f.Baz() })

  b = f
  fmt.Printf("Bar: %t => ", b != nil)
  try(func() { b.Baz() })

  f = &Foo{}
  fmt.Printf("Foo: %t => ", f != nil)
  try(func() { f.Baz() })

  b = f
  fmt.Printf("Bar: %t => ", b != nil)
  try(func() { b.Baz() })
}

func try(handle func()) {
  defer func() {
    if r := recover(); r != nil {
      fmt.Printf("Error: %v\n", r)
    }
  }()
  handle()
}

type Bar interface{ Baz() }

type Foo struct{}

func (f *Foo) Baz() { fmt.Printf("Baz: %t\n", f != nil) }
```

Will output:

```text
Bar: false => Error: runtime error: invalid memory address or nil pointer dereference
Foo: false => Baz: false
Bar: true => Baz: false
Foo: true => Baz: true
Bar: true => Baz: true
```

## Pointers

In the languages we are targetting, mostly sugar on Javascript, there are no
real pointers. We will implement pointers as an object that references
the value being pointed at.

- If the pointer is taken of a variable then the object will have getter and
  setter methods to read and write that variable.
- If the pointer is taken on the element of an array, the object will contain
  a getter and setter for that element inside the array as well as the index and
  a reference to the array itself so that the pointer can be cast back into a
  slice of the array when needed.
- If the pointer is taken of a slice, array, or map. It will need to be able
  get or set that value as a whole as well as provide indexers to access
  elements as specified by Go's pointers to a slice and map.
- If the pointer is taken of a field in a structure, the pointer must have
  a getter and setter to modify that field including a reference to the
  structure itself. Special rules exist for these so that a pointer to
  the first field will be equal to the pointer to the structure itself and
  allow casting from the first field back into the structure or vice versa.
  If the first field is an embedded type, then the first field in that embedded
  type is equal to the pointer to the embedded type and equal to the structure,
  and so on.
- A pointer to a named type will be able to get or set that type as well as
  call the methods for that type as defined by Go.

For some pointers we will be able to cast to other types, for example the
above field pointer, however some pointers we will not be able to cast between.
Depending on the target language, we may be able to cast between a `[]byte` and
a `[]uint32` but we will not be able to cast between different footprints.
This goes beyond just footprint sizes but includes the field sizes too. If we
have two types with the same fields, then we can cast between them by using
the same data but providing different type information. This means the structure
will stay the same but the methods could be changed. We can not support
casting from a type like `type color struct { value uint32 }` into
`type argb struct { b, g, r, a byte }` even through they have the same sized
footprint. It may be possible by creating a wrapper object tha at runtime does
a lot of bit manipulation to make that work but that is overkill and very
complicated. We will also not be able to alter our slices to capture memory
outside of an array like we can in Go. The slice will have to point to an
array (or nil) and the size and capacity will be bound to that array size and
only allowed to change when the underlying array is copied to a larger array.

## Generics

Gozer will implement generics based on ideas from the Golang
[Generics implementation - GC Shape Stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-gcshape.md)
design document. Obviously most target languages won't have garbage collect (GC)
since most target languages like TS are based on JS.
This means we don't have to worry about footprint sizes for GC and can allow
several different sizes of types to be used in one implementation.
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
  `uint64` implementation in TS. Instead we define a function `$add(x, y uint64)`
  that will properly perform the addition.
- Since we do not want to make multiple instantiations from a generic function
  that accepts a `int` and `int64`, on instantiation that calls `$add(x, y int64)`
  and one that performs a `x + y`, we make one instantiation that calls an `$add`
  function that is specific to the type parameter being passed in.
- By adding type dictionaries as an initial parameter to a function, similar
  to a closure captured variable or a receiver, we do not have to box and unbox
  basic types like Java does, e.g. `int` boxed in `Integer`.
- Most type dictionaries should be immutable singletons so that the same
  instance can be reused, like those for basic types. However the type
  dictionaries can be extended for types like `type Foo int` that can have
  methods added to them.
- For pointers, slices, arrays, and maps the type dictionary itself will be
  generic to allow then to be created for different element and key types.
- The type dictionaries will be used to help reflect determine how it can
  modify and a `reflect.Value`.

The basic type dictionaries have the following interfaces for operations:

| Operator Group | Operators | Applies To |
|:--|:-:|:--|
| AllTypes | (copy and zero) | all types |
| Addition / Concatenation | `+` (binary) | strings, integers, floats, complex |
| Arithmetic | `+` (binary/unary), `-` (binary/unary), `*`, `/` | integers, floats, complex |
| Remainder | `%` | integers only |
| Bitwise | `&` (binary), `\|`, `^`, `&^` | integers only |
| Shift | `<<`, `>>` | integers only |
| Comparison | `==`, `!=` | comparable types |
| Ordering | `<`, `<=`, `>`, `>=` | ordered types |
| Logical | `&&`, `\|\|`, `!` | booleans only |
| Reference | `&` (unary) | pointers only |
| Indexer | `[]` | maps, slices, arrays |

If we have a type like `type Foo int` with `func (f Foo) inc()` then that
type will implement Arithmetic, Addition, Remainder, Bitwise, Shift,
Comparison, Ordering, and `interface { inc() }`. That means it can be passed
in as a type argument for a generic with the constraints `~int`, `~int|~string`,
`interface { inc() }`, `any`, etc.

If we have a type constraint like `~int|~uint`, the type dictionary will have
to be an interface that impelments several of the basic type dictionaries,
e.g. `interface{ Arithmetic | Remainder | Bitwise | ...  }`. These unioned
interfaces can either be predefined and reused or defined with each use
depending on what works best for the target language.

It is possible with this setup that a type will implement an interface for
the generic type that was not defined in the constraints. To help prevent this:

- All operation functions are defined with a preceding `$` or some character not
  allowed in identifiers for Go but is allowed in the target language.
- We also rely on Go's type checker to insure the types are valid.
- If the function is exported and called directly from TS with the wrong type,
  that is fine since the function will still work if the interface is implemented.
  See [Target Language Type Dictionaries](#target-language-type-dictionaries) for more.

If we have a generic function like `foo[T int|float64](x T)` then it will
transpile into code similar to `foo<T>(x: T, dicT: Arithmetic<T>)` and it will
be called with `foo<int>(i, IntegerDic)` or `foo<float64>(f, Float64Dic)`.
Internal to the function when doing an operator on `x`, such as `y := -x`,
it will use the dictionary like, `T y = dicT.$neg(x)`.

An example `float64` type dictionary would look like:

```TS
class Float64Dictionary implements Arithmetic<float64> {
    $zero(): float64 { return 0.0; }
    $add(x: float64, y: float64): float64 { return x + y; }
    $neg(x: float64): float64 { return -x; }
    $sub(x: float64, y: float64): float64 { return x - y; }
    $mul(x: float64, y: float64): float64 { return x * y; }
    $div(x: float64, y: float64): float64 { return x / y; }
}
```

If a type dictionary isn't needed, like for `foo[T any]` where the constraints
do allow any operations to be performed on type `T`, then no type dictionary
will be added to the parameters. In most these cases, having the generic
type argument extending an interface will work.

For the indexer methods, there should be defined a `$get`, `$set`, and `$ref`.
The get and set will be used for reading and writing to the map, slice, or array.
The `$ref` will be used to create a pointer to the specified element so that
it can be for code like `p := &a[i]` where `p` can be used to read and write
the element in the array. The reference to maps will have to follow the same
methods to implement Go maps as close as possible while likely using a map
type in the target language.

For `make` calls, the `$zero` will be called, however for slices, arrays,
and maps there may be additional methods that allow sizes and capacities that
the `make` can be transpiled to call.

### Target Language Type Dictionaries

The type dictionaries would allow any type that meets the constraints to call
the function even if the original Go constraints do not implement.
This would allow us to create additional type dictionaries specific to the target
language. For example TS could have a `number` type dictionary defined.
Functions will be transpiled into a form that can be called by a block which also
takes the type dictionaries if needed. Those blocked functions will be created
to run quickly. If a function is exported, a function with the original name
will be exported that optionally takes a type dictionary. These functions
are intended to be called from the target language so will do some additional
work. If the type dictionary is not given, then it will look at the type passed
into the arguments to determine the type dictionary to use. Since TS can not
determine if a given argument is an `uint` or `int` we can instead use types
that can be determined, like `number`. We do not need to be strict and enforce
a constraint like `uint|int` at runtime because in that case a `number` can
work inside of the function even though it was not written explicitly to work
for a `number`. However, when looking up the dictionary we do have to make sure
that it implements the interface required for the function so that a `number`
can not be passed into something requiring an Indexer type dictionary.
