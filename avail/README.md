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
