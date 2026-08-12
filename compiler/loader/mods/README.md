# Modifiers

## Augmenter

TODO: Fill out

## Cache

This modifier will attempt to shortcut the loading process of a package
by loading precompiled serialized information, i.e. load from a cache.

[See cache for more information](./cache/README.md)

## TypeChecker

TODO: Fill out

## TODO

- Add modifier to adjust the source file list before parsing
- Add modifier to seek specific AST patterns to replace with another AST branch.
  For example, used to replace one cast deep inside a large method without
  having to override the whole method.
- Add modifier to simplify code if needed:
  - Simplify constants (except concatenated strings that have separate variables)
- Get actual useful build flags for a package. This is to help when loading
  a cache manifest to know which build flags affect a package build
- Need post processing for determining things like inheritance
