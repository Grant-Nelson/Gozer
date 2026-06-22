# Blocker Development Rules

Applies to: `project/modeler/remodel/blocker/**`, `project/modeler/ir/**`

## Context

The blocker transforms Go functions into schedulable blocks. Each block represents a unit of execution that runs without scheduler interruption.

## Current State

The blocker handles:
- Labels and goto statements
- For loops (with init, cond, post)
- Break and continue (with labels)
- If statements
- Basic expressions

The blocker does NOT yet handle:
- Range statements (`remodelRangeStmt` - panics with unimplemented)
- Switch/type switch/select statements
- Logical AND/OR short-circuit evaluation
- Channel receive expressions
- Call expressions in non-statement context (e.g., in binary expressions)
- Fallthrough statements

## Variable Passing Priority

The **immediate work** is designing how blocks track variable flow:

### Block.Params

Variables the block receives as inputs. These are identifiers that must have `types.Info` entries.

### BlockRef.Args

Expressions passed when jumping to a block. These become the `Params` of the target block.

### Design Requirements

1. **At split points** (function calls, channel ops), determine which variables are:
   - Live (used after the split)
   - Modified before the split

2. **Follow blocks** need:
   - Variables live across the split
   - Return values from the call (as synthetic variables)

3. **Integration with types.Info**:
   - Use existing `*ast.Ident` for params to leverage type information
   - Create synthetic idents for return values with proper type objects

## Testing Approach

Add tests to `blocker_test.go` following the existing pattern:
1. Define inline Go source
2. Call `blockIrcFunc()` to process it
3. Use `stringForFunc()` to get string representation
4. Use `diffString()` to compare expected output

Incomplete tests are marked with `// TODO: Finish` in the expected output.

## Code Patterns

### Splitting a block

```go
nextBlk := fbb.fn.NewBlock(`hint`)
stmt, gotoStmt := fbb.splitCurBlock(nextBlk)
// stmt is the current statement (removed from block)
// gotoStmt jumps to nextBlk (can be modified)
```

### Creating a function call

```go
follow := fbb.fn.NewBlock(`follow call`)
_, gotoFollow := fbb.splitCurBlock(follow)
fbb.curStmtList[fbb.stmtIndex+1] = &ir.FuncCallStmt{
    Ast:    e,
    Fun:    e.Fun,
    Args:   e.Args,
    Follow: gotoFollow.Block,
}
```

### Looking up type info

```go
obj := fbb.info().Uses[ident]   // for identifier usage
obj := fbb.info().Defs[ident]   // for identifier definition
typ := fbb.info().TypeOf(expr)  // for expression type
```

## Debug Helpers

- `crumb.DropMsg()` - prints location when reached (for debugging flow)
- `fmt.Printf` statements marked `// TODO: Remove` exist for debugging
- Block hints help identify blocks in test output
