# Composure

## Nullability

Given a nullable type (pointer, slice, chan, map), `X`, when using `X` after an
assignment all following uses will have the same null state until the next assignment.
This means that only right before the first usage of `X`, the null state of `X`
should be checked, then all following usage until next assignment doesn't have
to check the null state. For this documentation "null" means typed or untyped nil
in Go, meaning the data is null or the pointer address (`uintptr` value) is 0.

For example, given the following Go program:

1. First assume that all usages (including reads, writes, defines) of `X` needs
    to have a null state check and apply "maybe null" composure flag to all usages
2. Then for each assignment or definition of `X`:
    1. If assigned to "null" then set the composure to "is null"
        1. The first read of `X` than becomes a "dereference nil error" panic
        2. All following statements and expressions that would be handled after the
            panic can be removed
    2. If assigned to "not null" then set the composure to "not null"
        1. All following reads of `X` until the next assignment can be made without
            checking null state
    3. If assigned to a "maybe null" leave the composure as "maybe null"
        1. The first read of `X` stays labelled as "maybe null" so that when
            the code is outputted a null state check is performed right before
            reading `X` and panic a "deference nil error" panic
        2. For all following reads of `X` until the next assignment can be made
            without checking null state and may be set to composure "not null"
3. If every returned result is either all "null" or all "not null" then the
   returned result will be set to that common composure, however, if they are
   mixed or "maybe null" then the returned result will have the "maybe null"
   composure

```Go
X := foo()
x[4] := X[2]
X.Bar()
```

In the example code, if the `foo()` call returns "null" then `X[2]` will
be replaced with a panic for arrays (e.g., `X *[5]T`) and slices (e.g., `X []T`)
but not always for a map (e.g., `X map[int]V`) depending on the operator.
Maps index read ops can always be set to "not null" since the map can read
a default value. Map index writes should null state check.

TODO: Create pseudo version

TODO: `foo()` call returns "not null"

TODO: `foo()` call returns "maybe null"

When calling a function like `Bar()` note that the null state of `X` is not
checked until the receiver is used inside of `X` since Go allows null receivers.

TODO: `X == null` in a constant assignment of a bool which could then cause something
    like an if-statement to reduce. In a read of a "maybe null" `X` for this condition
    will result in the if-statement body or else (or after depending on terminating
    branches) being "null", "not null", or "maybe null" if the `X == null` is in
    an OR which may or may not cause `X` to resolve.

## Numeric Affirmation

When doing an op like `x++` or `x*y` we need to check what composure ranges
are possible to determine when overflows are possible or in a comparison.

TODO: When to enforce N-bit size or sign.

## Boolean Affirmation

TODO: `if false {` or `if true {`

## Index Boundaries

TODO: Check for out-of-bounds only when needing too.
