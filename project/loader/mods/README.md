# Modifiers

## TODO

- Add Modifiers
  - to preload cached packages then shortcut modified files
  - to store modified files in a cached package that saves when load is done
  - to simplify constants (except concatenated strings that have separate variables)
  - to remove defers into a `deferBlock` call
  - to remove Goto and labels (aka flatten)
  - to inject Jumps and labels to replace other flow-controls
  - to generate return structures for multiple returns
  - to replace multiple assignments with a `multiAssign` call
  - to flatten select statements and switches as needed
  - to adjust imports
- Need post processing for determining things like inheritance
