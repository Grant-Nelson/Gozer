# Types and Generics

*Design document for type representation*
*and generic implementation*

The resulting transpiled code needs to function as closely as we can to
the source Go code. This means that some of the quirks of Go must be carried
through to the output.

- [Types and Generics](#types-and-generics)
  - [Go Interfaces](#go-interfaces)
    - [Interface Structure](#interface-structure)
    - [Typed Nil vs Untyped Nil](#typed-nil-vs-untyped-nil)
    - [Method Invocation on Interfaces](#method-invocation-on-interfaces)
    - [Type Assertions and Type Switches](#type-assertions-and-type-switches)
    - [Method Sets and Interface Satisfaction](#method-sets-and-interface-satisfaction)
  - [Type Information](#type-information)
  - [Pointers](#pointers)
    - [Pointer to Variable](#pointer-to-variable)
    - [Pointer to Array Element](#pointer-to-array-element)
    - [Pointer to Slice, Array, or Map](#pointer-to-slice-array-or-map)
    - [Pointer to Struct Field](#pointer-to-struct-field)
    - [Pointer to Named Type](#pointer-to-named-type)
    - [Pointer Casting Limitations](#pointer-casting-limitations)
  - [Generics](#generics)
    - [No Monomorphic Instantiations](#no-monomorphic-instantiations)
    - [Type Dictionaries](#type-dictionaries)
    - [Operator Groups](#operator-groups)
    - [Type Dictionary Examples](#type-dictionary-examples)
    - [Indexer Operations](#indexer-operations)
    - [Mapper Operations](#mapper-operations)
    - [Make and Zero Values](#make-and-zero-values)
    - [Target Language Type Dictionaries](#target-language-type-dictionaries)
  - [Integration with Blocks](#integration-with-blocks)

## Go Interfaces

In many languages, "interface" refers to a collection of method signatures only
(as in OO languages like C# or Java). However, Go interfaces are objects that
contain type information. In these documents, we use "Go interface(s)" when
discussing Go's specific interface semantics.

### Interface Structure

Go interfaces will be transpiled into two parts:

1. **An interface type** - A collection of method signatures (the traditional OO concept)
2. **A wrapper object** - Contains:
   - A reference to the concrete type information
   - A reference to the underlying object

```typescript
interface Stringer {
    String(): string;
}

class GoInterface<T> {
    typeInfo: TypeInfo | null;  // null for untyped nil
    value: T | null;            // null for nil value
}
```

### Typed Nil vs Untyped Nil

The type and object may both be nil (untyped nil), or just the object
may be nil (typed nil). This allows for a non-nil Go interface that
contains a nil object.

```go
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

Output:

```text
Bar: false => Error: runtime error: invalid memory address or nil pointer dereference
Foo: false => Baz: false
Bar: true => Baz: false
Foo: true => Baz: true
Bar: true => Baz: true
```

Key observations:
- `var b Bar` is an untyped nil - comparing to nil returns false, calling methods panics
- `var f *Foo` is a typed nil - comparing to nil returns false, but methods can be called
- `b = f` assigns a typed nil to the interface - now `b != nil` is true
- The method `Baz()` can check if its receiver is nil

### Method Invocation on Interfaces

When a function is called on a Go interface, it will call the receiverless
instance of the function on the contained type with the contained object.
The receiverless function has the receiver moved to be the first argument.
This allows nil receivers to be passed correctly.

```go
// Original method
func (f *Foo) Baz() { ... }

// Transpiled receiverless form
func Foo_Baz(f *Foo) { ... }

// Interface call
b.Baz()  // becomes: Foo_Baz(b.value)
```

### Type Assertions and Type Switches

When casting from a Go interface to another Go interface or to a concrete type,
the Go interface is unwrapped and the internal type is used to determine if
the cast can be completed.

```go
// Type assertion
if s, ok := b.(Stringer); ok {
    // s is the Stringer interface wrapping the same value
}

// Type switch
switch v := b.(type) {
case *Foo:
    // v is *Foo
case Stringer:
    // v is Stringer
default:
    // v is the original interface type
}
```

Transpilation must:
1. Extract the type info from the Go interface
2. Check if the type implements the target interface (for interface assertions)
3. Check if the type matches exactly (for concrete type assertions)
4. Rewrap the value in the new interface type if successful

### Method Sets and Interface Satisfaction

A type satisfies an interface if its method set includes all interface methods.
The method set depends on whether we have a value or pointer:

- `T`'s method set includes methods with receiver `T`
- `*T`'s method set includes methods with receiver `T` or `*T`

```go
type Counter struct { n int }
func (c Counter) Value() int { return c.n }   // receiver T
func (c *Counter) Inc() { c.n++ }              // receiver *T

// Counter implements interface{ Value() int }
// *Counter implements interface{ Value() int; Inc() }
// Counter does NOT implement interface{ Inc() } - Inc has *Counter receiver
```

## Type Information

The type information will be an immutable object typically defined once as a
singleton. The type information will contain:

1. **Receiverless function definitions** - When called, these push a `CallNode`
   onto the scheduler call stack with any closure, receiver, type parameters,
   and type dictionaries, then call the first block for that function with the
   arguments passed in. See [Blocks.md](./Blocks.md) for scheduler details.

2. **Field and method type information** - Allows reflection and method values.
   The type information for a method defines the signature such that the method
   can be used as a function pointer. See [VariablePassing.md](./VariablePassing.md)
   for how method receivers become bound variables.

3. **String representation** - Information to show what Go would output when
   printing the type (e.g., using `%T` or `%#v`).

The type information may also extend a [type dictionary](#type-dictionaries)
for types like `type Foo int` that may also have methods attached to them.

How the type information looks depends on what is needed by Go and the target
language to appropriately define a type whilst leveraging as much of the
target language's type system as possible.

## Pointers

In the languages we are targeting (mostly sugar on JavaScript), there are no
real pointers. We will implement pointers as objects that reference
the value being pointed at.

### Pointer to Variable

If the pointer is taken of a variable, the object will have getter and
setter methods to read and write that variable.

```typescript
class VarPointer<T> {
    get(): T;
    set(value: T): void;
}
```

### Pointer to Array Element

If the pointer is taken on the element of an array, the object will contain
a getter and setter for that element inside the array as well as the index and
a reference to the array itself so that the pointer can be cast back into a
slice of the array when needed.

```typescript
class ArrayElemPointer<T> {
    array: T[];
    index: number;
    get(): T;
    set(value: T): void;
    toSlice(): Slice<T>;  // slice from index to end
}
```

The `toSlice` may also include an optional length so that it can be used
by methods such as Go's [`unsafe.Slice`](https://pkg.go.dev/unsafe#Slice).

Additionally, if the pointer to the element has methods such as when the
slice is filled with non-pointer named types, (e.g. `[]Cat`), then the
pointer needs to have the methods for `*Cat` able to be called on it too.

### Pointer to Slice, Array, or Map

If the pointer is taken of a slice, array, or map, it will need to be able
get or set that value as a whole as well as provide indexers to access
elements as specified by Go's pointers to a slice and map.

### Pointer to Struct Field

If the pointer is taken of a field in a structure, the pointer must have
a getter and setter to modify that field including a reference to the
structure itself.

Special rules exist for these so that a pointer to the first field will be
equal to the pointer to the structure itself and allow casting from the
first field back into the structure or vice versa. If the first field is
an embedded type, then the first field in that embedded type is equal to
the pointer to the embedded type and equal to the structure, and so on.

```go
type Outer struct {
    Inner  // embedded, first field
    x int
}
type Inner struct {
    y int  // first field of Inner
}

var o Outer
// These are all equivalent pointers:
// &o == &o.Inner == &o.Inner.y (as unsafe.Pointer)
```

### Pointer to Named Type

A pointer to a named type will be able to get or set that type as well as
call the methods for that type as defined by Go.

### Pointer Casting Limitations

For some pointers we will be able to cast to other types (e.g., the field
pointer above), however some pointers we will not be able to cast between.

Depending on the target language, we may be able to cast between a `[]byte` and
a `[]uint32` but we will not be able to cast between different footprints.
This goes beyond just footprint sizes but includes the field sizes too. If we
have two types with the same fields, then we can cast between them by using
the same data but providing different type information. This means the structure
will stay the same but the methods could be changed.

We cannot support casting from a type like `type color struct { value uint32 }`
into `type argb struct { b, g, r, a byte }` even though they have the same sized
footprint. It may be possible by creating a wrapper object that at runtime does
a lot of bit manipulation to make that work but that is overkill and very
complicated.

We will also not be able to alter our slices to capture memory outside of an
array like we can in Go. The slice will have to point to an array (or nil) and
the size and capacity will be bound to that array size and only allowed to
change when the underlying array is copied to a larger array.

## Generics

Gozer will implement generics based on ideas from the Go
[Generics implementation - GC Shape Stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-gcshape.md)
design document. Obviously most target languages won't have garbage collection (GC)
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
without knowing the packages that import it.

This is undesirable so we must handle generics without creating monomorphic
instantiations (e.g. `List_int`, `List_Cat`)

### Type Dictionaries

There are a finite number of basic types, such as `int8`, `int`, `uint64`, `bool`,
`float64`, and `complex64`. When these types are used in type constraints
on type parameters for generic declarations, they indicate what can possibly
be done with that type, i.e., what operators can be called on that type.

We don't care if the type does or does not call one of the operators since
it may call another function or function pointer parameter that does call one
of those operators. We base what is possible on the type constraints without
needing to inspect the body of a generic function or methods of a generic type.

Key observations:

- We can't run the target language's operator on all types. For example, the TS
  binary addition operator, `x + y`, can not be run on a `uint64` since TS
  integers are limited to the size of JavaScript integers which is less than
  64 bits. We will have to create an object in TS to represent a `uint64` with
  a high and low part (same for `int64`, `complex64`, and `complex128`).

- Since TS will not allow operator overloading we can not perform a `+` on our
  `uint64` implementation in TS. Instead we define a function `$add(x, y uint64)`
  that will properly perform the addition.

- Since we do not want to make multiple instantiations from a generic function
  that accepts an `int` and `int64`, one instantiation that calls `$add(x, y int64)`
  and one that performs `x + y`, we make one instantiation that calls an `$add`
  function that is specific to the type parameter being passed in.

- By adding type dictionaries as an initial parameter to a function, similar
  to a closure captured variable or a receiver
  (see [BoundFunc](./VariablePassing.md#boundfunc-for-closures)),
  we do not have to box and unbox basic types like Java does, e.g., `int` boxed in `Integer`.

- Most type dictionaries should be immutable singletons so that the same
  instance can be reused, like those for basic types. However, the type
  dictionaries can be extended for types like `type Foo int` that can have
  methods added to them.

- For pointers, slices, arrays, channels, and maps the type dictionary itself
  will be generic to allow them to be created for different element and key types.

- The type dictionaries will be used to help reflect determine how it can
  modify a `reflect.Value`.

### Operator Groups

The basic type dictionaries have the following interfaces for operations:

| Operator Group | Operators | Applies To |
|:--|:-:|:--|
| AllTypes | (copy and zero) | all types |
| Addition / Concatenation | `+` (binary) | strings, integers, floats, complex |
| Arithmetic | `+` (binary/unary), `-` (binary/unary), `*`, `/` | integers, floats, complex |
| Remainder | `%` | integers only |
| Bitwise | `&` (binary), `\|`, `^`, `&^` | integers only |
| Shift | `<<`, `>>` | integers only |
| Comparison | `==`, `!=` | comparable types (basic types, pointers, channels, interfaces, structs/arrays of comparable types) |
| Ordering | `<`, `<=`, `>`, `>=` | ordered types (integers, floats, strings) |
| Logical | `&&`, `\|\|`, `!` | booleans only |
| Reference | `&` (unary) | addressable values |
| Dereference | `*` (unary) | pointers only |
| Indexer | `[]` | slices, arrays, strings |
| Mapper | `[]` | maps |
| Receive | `<-c` | channels only |
| Sender | `c<-` | channels only |

If we have a type like `type Foo int` with `func (f Foo) inc()` then that
type will implement Arithmetic, Addition, Remainder, Bitwise, Shift,
Comparison, Ordering, and `interface { inc() }`. That means it can be passed
in as a type argument for a generic with the constraints `~int`, `~int|~string`,
`interface { inc() }`, `any`, etc.

If we have a type constraint like `~int|~uint`, the type dictionary will have
to be an interface that implements several of the basic type dictionaries,
e.g., `interface{ Arithmetic; Remainder; Bitwise; ... }`. These combined
interfaces can either be predefined and reused or defined with each use
depending on what works best for the target language.

It is possible with this setup that a type will implement an interface for
the generic type that was not defined in the constraints. To help prevent misuse:

- All operation functions are defined with a preceding `$` or some character not
  allowed in identifiers for Go but is allowed in the target language.
- We also rely on Go's type checker to ensure the types are valid.
- If the function is exported and called directly from TS with the wrong type,
  that is fine since the function will still work if the interface is implemented.
  See [Target Language Type Dictionaries](#target-language-type-dictionaries) for more.

### Type Dictionary Examples

If we have a generic function like `foo[T int|float64](x T)` then it will
transpile into code similar to `foo<T>(x: T, dicT: Arithmetic<T>)` and it will
be called with `foo<int>(i, IntegerDic)` or `foo<float64>(f, Float64Dic)`.
Internal to the function when doing an operator on `x`, such as `y := -x`,
it will use the dictionary like `T y = dicT.$neg(x)`.

An example `float64` type dictionary:

```typescript
class Float64Dictionary implements Arithmetic<number> {
    $zero(): number { return 0.0; }
    $copy(x: number): number { return x; }
    $add(x: number, y: number): number { return x + y; }
    $neg(x: number): number { return -x; }
    $sub(x: number, y: number): number { return x - y; }
    $mul(x: number, y: number): number { return x * y; }
    $div(x: number, y: number): number { return x / y; }
}
```

An example `int64` type dictionary (requires wrapper for 64-bit integers):

```typescript
class Int64Dictionary implements Arithmetic<Int64> {
    $zero(): Int64 { return new Int64(0, 0); }
    $copy(x: Int64): Int64 { return new Int64(x.high, x.low); }
    $add(x: Int64, y: Int64): Int64 { return int64Add(x, y); }
    $neg(x: Int64): Int64 { return int64Neg(x); }
    // ... etc
}
```

If a type dictionary isn't needed, like for `foo[T any]` where the constraints
don't allow any operations to be performed on type `T`, then no type dictionary
will be added to the parameters. In most these cases, having the generic
type argument extending an interface will work.

### Indexer Operations

For the indexer methods, there should be defined `$len`, `$cap`, `$get`, `$set`,
and `$ref`. The get and set will be used for reading and writing to the slice
or array. The `$ref` will be used to create a pointer to the specified
element so that it can be used for code like `p := &a[i]` where `p` can be used
to read and write the element in the array.
As discussed in [Pointer to Array Element](#pointer-to-array-element) the pointer
returned from `$ref` (an "index pointer") will carry enough information to cast
back into a slice.

```typescript
interface Indexer<V> {
    $get(container: any, key: integer): V;
    $set(container: any, key: integer, value: V): void;
    $ref(container: any, key: integer): IndexPointer<V>;
    $len(container: any): number;
    $cap(container: any): number;
}
```

If the target language has types like Javascript's
[TypedArray](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/TypedArray)
for faster `Int8Array`, `Int32Array`, etc arrays, then those should be
used for specific arrays and slices (e.g. `[]int8`, `[]int32`).
Therefore in some cases we wouldn't use the generic type information for those
slices but specifically designed ones to improve performance.

Note that Go can use `int64` as an index, so the compiler will have to convert
those to integers before use. This means we will not get the full 64-bits of
indexing. However, that is fine since most Javascript (and Typescript)
arrays can only index up-to a JS integer but will error when higher than
32-bit values.

### Mapper Operations

For the mapper methods, there should be defined `$len`, `$cap`, `$get`, `$set`.
The get and set will be used for reading and writing to the map.

```typescript
interface Indexer<K, V> {
    $get(container: Map, key: K): V;
    $set(container: Map, key: K, value: V): void;
    $len(container: Map): number;
    $cap(container: Map): number;
}
```

### Make and Zero Values

For `make` calls, the `$zero` will be called for element initialization.
However, for slices, arrays, and maps there may be additional methods that
allow sizes and capacities that the `make` can be transpiled to call.

```typescript
interface Makeable<T> {
    $zero(): T;
    $make(len: number, cap?: number): T;  // for slices
}
```

This connects to named result initialization in
[VariablePassing.md](./VariablePassing.md#named-results-as-local-variables) -
named results use `$zero()` for their initial values.

### Target Language Type Dictionaries

The type dictionaries would allow any type that meets the constraints to call
the function even if the original Go constraints do not include it.
This would allow us to create additional type dictionaries specific to the target
language. For example, TS could have a `number` type dictionary defined.

Functions will be transpiled into a form that can be called by a block which also
takes the type dictionaries if needed. Those blocked functions will be created
to run quickly. If a function is exported, a function with the original name
will be exported that optionally takes a type dictionary. These functions
are intended to be called from the target language so will do some additional
work.

If the type dictionary is not given, then it will look at the type passed
into the type arguments to determine the type dictionary to use. Since TS cannot
determine if a given argument is a `uint` or `int` we can instead use types
that can be determined, like `number`. We do not need to be strict and enforce
a constraint like `uint|int` at runtime because in that case a `number` can
work inside of the function even though it was not written explicitly to work
for a `number`.

However, when looking up the dictionary we do have to make sure
that it implements the interface required for the function so that a `number`
cannot be passed into something requiring an Indexer type dictionary.

We will try to base this lookup off the type arguments instead of the
passed in arguments since some methods have no parameters and only a return
value. However as we implement this, if we determine it is too complicated
to base it off the type arguments then we can do the look up with parameters
and make the type dictionaries required for functions with only return values.

## Integration with Blocks

Type dictionaries and type information integrate with the block-based execution
model described in [Blocks.md](./Blocks.md):

1. **Type dictionaries as parameters** - Similar to bound variables in closures
   (see [BoundFunc](./VariablePassing.md#boundfunc-for-closures)), type dictionaries
   are passed as implicit first parameters to generic functions.

2. **Method calls through interfaces** - When calling a method on a Go interface,
   the type information provides the block index to call, and the scheduler
   handles the invocation.

3. **Reflect and type switches** - Runtime type information allows the scheduler
   to determine which blocks to execute for type assertions and switches.

```pseudo
// Generic function call with type dictionary
call foo[int](x), follow=Block1
// Transpiles to:
call foo(x, IntDic), follow=Block1

// Interface method call
b.Baz()
// Transpiles to:
call b.typeInfo.Baz(b.value), follow=Block1
```
