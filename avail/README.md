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

I told an agent to do a quick review of things that could be problems and this is what it found.

### assert

- `faults.New` doesn't format its message, so every assert message literally contains `%w` (for example, "%w equality check").
- `Zero` and `Nonzero` report "must be nil check" and "not nil check" (copy-paste).
- The ordered checks (`Lt`, `Gt`, `Leq`, `Geq`, `Range`) all pass silently when a value is NaN.
- `Range`, `StrRange`, `SliceRange` and `MapRange` don't validate that `low <= high`.
- The comment refers to `disableAsserts`, but the constant is `enableAsserts`.

### faults

- `ErrGroup.Unwrap`: `g.errs[:max:max]` with `max = len-1` drops the last error and panics on an empty group. This also makes
  `expand` lose the last error whenever groups are merged.
- `ErrGroup.Full()` returns true for a nil group, so `FullOrNil()` on a nil group returns a typed-nil `*ErrGroup` as a non-nil
  `error`.
- `crunch` reorders the caller's slice in place, because a variadic call with `errs...` shares the caller's backing array.
- `expand` drops the merged group's `remainder` (its discarded count).
- The zero value `ErrGroup{}` panics because `lock` is a nil pointer. The same applies to `logger.Logger`.
- `Error()` with exactly one error and `remainder > 0` omits the "discarded" note.
- The `Empty` doc comment says the opposite of what the function does.
- `Fault.Error()` never renders the inner errors, so `From(*ErrGroup)` produces just "error group" and all the detail is lost.
- `SetCurLoc` skips any directory ending in `/assert` or `/runtime`, which includes third-party packages such as testify's `assert`.
  It also skips the faults package's own tests, and it gives up after 10 frames.
- In `Recover`, the location is the deferred call site, not where the panic happened.
- `With` puts errors into `data`, which are printed with `%v`, so multi-line errors break the tab-indented layout. `bool` and named
  numeric types get stringified while built-in numerics don't.

### iterator

- `Limit`: it yields before checking `count`, so `Limit(n)` returns n+1 values.
- `LessThan`: off by one. It returns true when the count is exactly `count`, so it behaves like less-than-or-equal.
- `Zip`: when both iterators are exhausted, or only the second has values left, the loop never exits. Neither `stop1` nor `stop2`
  is called when `yield` returns false; they should be deferred.
- `Expand`: the constraint `TIn iter.Seq[TOut]` only matches exactly `iter.Seq[TOut]`, so `Iterator[Iterator[T]]` can't be used. It
  needs `~func(func(TOut) bool)`.
- The `SkipWhile` doc says "until the predicate returns true", but the code skips while it returns true.
- `Reduce` and `Aggregate` call `r(next, acc)`, which is the reverse of the documented "first and second value" order and matters
  for non-commutative functions. `Reduce` also can't signal an empty sequence.
- Doc names don't match: "Index" for `Indexed`, "ToStringF" for `ToStringsF`, "Appends" for `Append`.

### stack

- `Trim` doesn't actually free memory while any live node shares a chunk with the tombs, because nodes are slices of one
  allocated array.
- `Pop` and `Peek` can't distinguish an empty stack from a stored zero value.
- `PushSeq` corrupts the stack if the sequence reads or mutates the same stack (for example, `s.PushSeq(s.PopAll())`).
- `Push(a, b)` gives a different order than `PushOne(a).PushOne(b)`, and the "reverse order" wording in the doc is ambiguous.

### logger

- Indentation is shared state, so concurrent `LogGroup` and `Indent` calls from different goroutines produce wrong indentation
  even though the writes are locked.
- Calling the returned dedent or group-end function twice dedents twice.
- Blank lines get the prefix, which leaves trailing whitespace.
- Write errors are ignored.

### crumb

- `Out` is typed `*os.File`, so crumbs can't be redirected to an arbitrary `io.Writer` as the package doc suggests.
- `runtime.FuncForPC` gives the wrong function name for inlined callers; `runtime.CallersFrames` handles that correctly.
- Typos: "there way" should be "their way", "an it" should be "and it".

### source

- `PathRebase` uses a raw string prefix, so `/a/b` also matches `/a/bc/x`. With an empty `oldBase`, it concatenates with no
  separator.
- `Skipper` and `Override` match uncleaned paths, so `./a` and `a` are treated as different.
- `ToReader` passes an `io.Reader` through without any close semantics, while reading files fully into memory.
- Doc typos: "id" should be "is", "te" should be "the".

### astTools

- `NodesWithStack`: when `Skip` is set, `ast.Inspect` never calls `f(nil)` for that node, so it's never popped and the stack is
  corrupted from then on. `Nodes` has the same missing close event, which contradicts its doc.
- `Idents` mutates the shared `ds.Node` for each name in a `ValueSpec`. Values the caller kept from earlier names all end up
  pointing at the last ident, and `Start`/`End` change meaning.
- `Idents` yields the quoted import path instead of the import's name. It also yields method names and `init` as if they were
  package-level identifiers, and filters `_` only for values.
- `ParsePackage`: empty `inputFiles` means no patterns, which loads the current directory's package. The working directory is
  captured once from the process, and error strings are capitalized.
- `RemoveDirectives` matches by `Pos`, so synthetic comments at `NoPos` all match each other.
- Doc typos: "belongs too" should be "belongs to", "the directly containing" should be "the directory containing".

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
