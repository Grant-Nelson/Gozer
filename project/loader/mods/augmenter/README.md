# Augmenter

The augmenter modifier reads in specially marked up `*.go` files
that specify what changes need to be made to files for a package.

- [Augmenter](#augmenter)
  - [Adding](#adding)
    - [Adding an Import](#adding-an-import)
    - [Adding a Func](#adding-a-func)
    - [Adding a Var or Const](#adding-a-var-or-const)
    - [Adding a Data Type](#adding-a-data-type)
    - [Adding an Interface](#adding-an-interface)
    - [Adding a Field](#adding-a-field)
  - [Deleting](#deleting)
    - [Deleting a Func](#deleting-a-func)
    - [Deleting a Var or Const](#deleting-a-var-or-const)
    - [Deleting a Data Type or Interface](#deleting-a-data-type-or-interface)
    - [Deleting a Field](#deleting-a-field)
  - [Replacing](#replacing)
    - [Replacing a Signature](#replacing-a-signature)
  - [Rename](#rename)
    - [ReplaceRecv](#replacerecv)
  - [Ignoring](#ignoring)

## Adding

Adding new code into a package can be done with `//gozer:add`.
When adding something new to the code, the existing code is checked for a match
code to ensure that the code isn't accidentally overwriting existing code.
Any comments or directives, other than the `//gozer:add`, attached to the code
being added, will be added as well.

### Adding an Import

An import can be added with or without an alias.

```Go
//gozer:add
import "foo/bar"
```

### Adding a Func

A function can be added with or without a receiver.
The function should have a body or have a `go:linkname` directive.
A new `init()` function may be added.

```Go
//gozer:add
func (x *X) Foo(y int, z string) { /*...*/ }
```

### Adding a Var or Const

A variable can be added with or without initialization.
A constant can be added with an initialization.
A `var()` or `const()` group can be added.

```Go
//gozer:add
var x int
```

```Go
//gozer:add
var (
    x int
    y int
)
```

```Go
var (
    //gozer:add
    x int
    //gozer:add
    y int
)
```

In the prior example, `x` and `y` may not be added in the same group
so things like `iota` may not work correctly. To use `iota` correctly
add the vars as part of declaration, like the example in the middle.

### Adding a Data Type

A package-level type, struct, alias, etc can be added.

```Go
//gozer:add
type Foo struct { /*...*/ }
```

### Adding an Interface

A package-level interface can be added.

```Go
//gozer:add
type Foo interface { /*...*/ }
```

### Adding a Field

A package-level struct can have a field added to it.
The structure must exist originally or added in a previous modifier.
A function signature can be added to a package-level interface too.
The type may not have any other gozer directive on it.
Multiple different kinds of gozer directives can be defined in the same type.

```Go
type Foo struct {
    //gozer:add
    Bar int `tag`
}
```

```Go
type Foo interface {
    //gozer:add
    Bar(x int, y string)
}
```

## Deleting

A type, interface, function, field, var, or const can be deleted
with `gozer:delete`. The name of the code to delete must match
the original name.

### Deleting a Func

When deleting a function with a receiver, the receiver name and
if the receiver is a pointer must match te function. The type parameters
and signature does not have to match. The function may be a stub without
a function body.

```Go
//gozer:delete
func Foo()
```

### Deleting a Var or Const

A var or const can be deleted. They do not have to have the original's
identifier and the actual type of the var or const does not have to match.
Vars and consts can be deleted even if they were originally defined
as part of a multi-assignment.

```Go
//gozer:delete
var x int
```

A deletion may be applied to a group or multi-assignment to delete
all of the vars and consts within that group. For example, the following
delete both `x` and `y`.

```Go
//gozer:delete
var (
    x int
    y int
)
```

```Go
//gozer:delete
var x, y int
```

```Go
var (
    //gozer:delete
    x int
    //gozer:delete
    y int
)
```

Note that if the deleted var or const is in the middle of `iota`
then it will cause the iota to offset.

### Deleting a Data Type or Interface

Deleting a type can be done as long as the name matches.
The type may be any underlying type but if deleting an interface
then the underlying type must be an integer too.

```Go
//gozer:delete
type Foo struct{}
```

```Go
//gozer:delete
type (
    Foo struct{}
    Bar struct{}
)
```

```Go
type 
    //gozer:delete
    Foo struct{}
    //gozer:delete
    Bar struct{}
)
```

When deleting a data type, all the functions with that type as the
receiver can be deleted too with `gozer:deleteAll`.

### Deleting a Field

A field may be deleted from a structure or interface.

```Go
type Foo interface {
    //gozer:delete
    Bar()
}
```

## Replacing

Replacing is similar to deleting and adding where it matches based on name
and basic type kind and it must provide the new function, type, var, const,
field, or interface. Replace with `gozer:replace`.

To change the basic type kind, e.g. replace a struct with a function,
delete the original and add the new type separately.

```Go
//gozer:replace
type Foo struct { x int; y int }
```

May not use replace on a multi-var or multi-type declaration.
Only one spec can be replaced by one replace at a time.

### Replacing a Signature

A function can have the signature changed without modifying any of the code
with `gozer:replaceSig`. The function name must match with the package-level
function to replace the signature. If the function has a receiver, then it must
match the original receiver so that the original can be found.
The generic parameter can be modified with the signature.

```Go
//gozer:replaceSig
func Foo[T any](x T)
```

## Rename

An type, interface, function, field, var, or const can be renamed with
`gozer:rename <newName>`. `<newName>` must be a valid id.
The code must provide the existing name so that it can be matched.

```Go
//gozer:rename Bar
type Foo struct {}
```

For functions, a rename can be paired with a `gozer:replaceSig` to
change the name and signature of a function while keeping the body unchanged.

For imports, a rename can be done with a simple replace since imports
will use the import path to find a matching import to replace and
the replaced one can have a different alias.

### ReplaceRecv

The receiver of a function can be replaced with `gozer:replaceRecv (*)<newType>`.
This can be paired with `gozer:rename` to replace the receiver and function
at the same time. The new receiver type includes if the type should be
a pointer or not with a `*` preceding the new type. Replacing the receiver
type is part of renaming since the names for functions include the receiver
type, e.g. `(*Foo).Bar`. To replace the name of the receiver within the scope
of the body, e.g. the `x` in `(x *Foo)`, use `gozer:replaceSig` to change the
parameter names with `gozer:replaceRecv` to change the type.

```Go
//gozer:replaceRecv *Foo2
func (f *Foo) Bar(x int, y int)
```

## Ignoring

An import, type, function, variable, and constant can be ignored
with `gozer:ignore`. Ignored items will have no affect in the augmenter
but can be added to the augmenter code to reduce the number of warnings
and errors displayed in a code editor. For example, when adding a new
function with a receiver in the original code, as stubbed out type
can be added and ignored so that the editor doesn't complained about
the missing type.

If imports doesn't have any directives, they will default to ignore.
