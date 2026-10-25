# Avail

- [args](./args): a tool for parsing nested command line arguments into a tagged struct, including tools, flags,
  positional arguments, and generated help output.
- [assert](./assert): a collection of generic assertions that panic with a detailed fault when a check fails.
- [astTools](./astTools): a collection of helpers for Go ASTs, such as iterating declarations, identifiers, and nodes,
  reading directives and comments, and parsing packages from source strings.
- [crumb](./crumb): a debugging tool for dropping "bread crumbs" that print the file, line, and function where they're
  reached.
- [faults](./faults): an error type with a source location and key/value data, plus a concurrency-safe error group for
  collecting multiple errors.
- [iterator](./iterator): a wrapper around `iter.Seq` that adds LINQ-style methods for filtering, transforming,
  combining, and consuming sequences.
- [logger](./logger): a nil-safe logger with hierarchical indentation and timed log groups for verbose output.
- [predicate](./predicate): a collection of generic predicate functions and combinators like `And`, `Or`, and `Not`
  that are designed to work well with [iterator](./iterator).
- [source](./source): composable converters for rebasing, overriding, or skipping source file paths and their data,
  plus a helper to read them.
- [stack](./stack): a generic linked-list stack that reuses dead nodes, designed for depth-first traversals.
- [typeTools](./typeTools): a collection of helpers for `go/types`, including type checks and determining which
  operations a type (or type constraint) supports via [typeOp](./typeTools/typeOp).
- [wad](./wad): a simple and very fast binary serializer for writing and reading values and registered `Waddable` types
  in a fixed order designed to be used for caching or transmission. It is sensitive to code changes so not intended
  for long term persistence, although could be used like that if additional versioning is added into the serialization.

## Agent's Review

### typeTools

- `ops.go` still has `crumb.DropMsg` calls throughout, so every call prints to stdout.
- `unionTermsOps`: `Ops: ops` sits in the same composite literal as the `adj(...)` calls that modify `ops`. The evaluation order
  isn't specified, so the removals done by `adj` are probably lost.
- `mapTypeOps`: `Range2` is set to `Key()` instead of `Elem()`, and it includes `Cap`, which maps don't support.
- `arrayTypeOps` includes `IsNil` and `Make` and is missing `Comparable`. `structTypeOps` includes `IsNil`.
- `pointerTypeOps` and `unsafePointerTypeOps` include `Orderable`, but Go pointers aren't ordered.
- `chanTypeOps` is missing `Make` and `Ref`. Several other types are also missing `Ref`: interfaces, structs, and the
  comparable-only cases.
- `signatureTypeOps`: the Seq2 case gives `Range2` without `Range`, which contradicts the `Op` docs. `func(func() bool)` isn't
  handled.
- The `comparable` constraint, which has no terms, falls through to an empty `OpTypes`, so `Comparable` is lost.
- `floatTypeOps` sets `Complex` redundantly before the switch.
- In `is.go`, `Is[T]` and `IsPtrOf` don't call `Underlying()` or `types.Unalias`, but `IsBool` and its siblings do. So `IsSlice` returns
  false for a named or aliased slice while `IsUint8` sees through it.
- `Deref` returns nil for a non-pointer instead of the type itself.
- `anyType` panics with `GODEBUG=gotypesalias=0`, because the cast to `*types.Alias` fails.
- `typeOp`: `Add = 1 << iota` starts at bit 1, so bit 0 is unused but still inside `mask`. `Op(1).Valid()` is true and its `String()`
  is "".
- `typeOp` docs: `=>` should be `>=` in `Orderable`; the types for `new` in `Ref` are reversed; `Valid` mentions a `RealImag` op that
  doesn't exist.

### wad

- `read` uses `r.Read`, not `io.ReadFull`. Short reads (from `bufio`, gzip, network streams) corrupt the stream. It also panics when
  `n == len` and the reader returns `io.EOF`, and on a zero-length read at end of stream, such as a trailing empty string from a
  `bytes.Reader`.
- `typeName` has a leftover debug `fmt.Printf("Name: ...")`.
- `typeName` uses `t.String()`, so two packages with the same name and the same type name collide. Using `PkgPath` would avoid
  that.
- Registry indexes depend on registration (init) order, so data isn't stable across builds or when imports change. The
  header has no magic number to detect wrong input.
- `Writer.Wad` panics when `T` is an interface and `v` is nil, because `reflect.ValueOf(nil)` has no type. Nil maps, slices and
  funcs aren't encoded as nil.
- Reading 64-bit `int`/`uint` data on a 32-bit build silently truncates. `Bytes` and slice lengths are silently truncated to
  `uint32` when writing, and when reading, a corrupt length can allocate up to 4 GB.
- `nativeCaster` casts `*[8]byte` to `*uint64`/`*float64`, which isn't alignment-guaranteed. `-race`/checkptr flags this, and it can
  fault on architectures that require alignment.
- `defaultCaster.bufU` is never used.
- `Pos` is written and read without its `FileSet`, so it's only meaningful within the same process.

### predicate

- `InRange` doesn't check `min <= max`, and the ordering predicates are always false for NaN.
