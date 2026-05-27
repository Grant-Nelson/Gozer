# Variable Passing Between Blocks

*Design document for tracking*
*variable flow across block boundaries*

## Problem Statement

When a function is split into blocks, variables defined in one block may be needed in subsequent blocks. We must track:

1. Which variables are **live** at each block boundary
2. Which variables must be **passed** from one block to another
3. How **return values** from function calls become available in follow blocks

## Example

```go
func example(n int) int {
    x := n + 1          // defines x
    y := expensive(x)   // call splits here; x is live
    z := x + y          // needs x (from before) and y (return value)
    return z
}
```

After blocking:

```
Block 0 (entry):
  Params: [n]
  Body:
    x := n + 1
    call expensive(x), follow=Block1, followArgs=[x]

Block 1 (after call):
  Params: [x, $ret0]  // x passed in, $ret0 is return value
  Body:
    y := $ret0
    z := x + y
    return z
```

## Data Structures

### Block.Params

```go
type Block struct {
    // Params are variables passed into this block.
    // Each ident must have a types.Info entry for its Object.
    Params []*ast.Ident
    
    // ... other fields
}
```

### BlockRef.Args

```go
type BlockRef struct {
    Block *Block
    // Args are expressions passed to Block.Params.
    // len(Args) must equal len(Block.Params)
    Args []ast.Expr
}
```

### Synthetic Identifiers for Return Values

When a function call returns values, we need synthetic identifiers:

```go
type SyntheticIdent struct {
    Name string      // e.g., "$ret0", "$ret1"
    Type types.Type  // from function signature
}
```

These must be registered with `types.Info` so type lookups work.

## Algorithm: Liveness Analysis

### Step 1: Identify Variables at Split Points

At each split point (call, channel op, etc.), record:
- Variables defined before the split
- Variables used after the split

### Step 2: Compute Live Variables

A variable is **live** at a split if:
1. It is defined before the split, AND
2. It is used after the split (directly or transitively)

### Step 3: Populate Params and Args

For each edge from Block A to Block B:
1. B.Params = variables live at entry to B
2. For each `GotoBlockStmt` or `FuncCallStmt` in A targeting B:
   - Set Args to match B.Params

## Implementation Phases

### Phase 1: Basic Variable Tracking

Track which variables are defined in each block:

```go
type blockVars struct {
    defined map[*ast.Ident]bool  // variables defined in this block
    used    map[*ast.Ident]bool  // variables used in this block
}
```

### Phase 2: Split Point Analysis

At each split, compute forward liveness:

```go
func (fbb *funcBlockBuilder) computeLiveVars(splitIndex int) []*ast.Ident {
    // Walk statements after splitIndex
    // Collect all identifier uses
    // Filter to those defined before splitIndex
}
```

### Phase 3: Propagate Through Control Flow

Handle cases where variables flow through multiple blocks:

```
Block 0: x := 1
         if cond goto Block1 else goto Block2
Block 1: y := 2
         goto Block3
Block 2: y := 3
         goto Block3
Block 3: z := x + y  // x from Block0, y from Block1 or Block2
```

Block 3 needs `x` and `y` in Params. Both Block 1 and Block 2 must pass them.

### Phase 4: Handle Return Values

Function calls inject return values into the follow block:

```go
func (fbb *funcBlockBuilder) handleCallReturn(call *ir.FuncCallStmt) {
    sig := fbb.info().TypeOf(call.Fun).(*types.Signature)
    results := sig.Results()
    
    for i := 0; i < results.Len(); i++ {
        // Create synthetic ident for $ret{i}
        // Add to follow block's Params
    }
}
```

## Edge Cases

### Shadowing

```go
x := 1
if cond {
    x := 2  // shadows outer x
    foo(x)  // uses inner x
}
bar(x)  // uses outer x
```

Solution: Use `types.Info.Uses` and `types.Info.Defs` to distinguish by Object, not name.

### Address-Taken Variables

```go
x := 1
p := &x
foo()  // split point
*p = 2
```

If `x` is address-taken, it cannot be passed by value. Either:
1. Wrap in a cell/box type
2. Pass pointer explicitly

For initial implementation, assume no address-taken locals cross splits.

### Multiple Return Values

```go
a, b := multiReturn()
```

Follow block needs `$ret0` and `$ret1` in Params.

### Blank Identifier

```go
_, b := multiReturn()
```

Don't create param for blank identifier's return position.

## Integration with Blocker

### When Splitting

```go
func (fbb *funcBlockBuilder) splitCurBlock(nextBlk *Block) (Stmt, *GotoBlockStmt) {
    // ... existing code ...
    
    // NEW: Compute live variables
    live := fbb.computeLiveVars(fbb.stmtIndex)
    nextBlk.Params = live
    gotoStmt.Block.Args = fbb.identExprs(live)
    
    return stmt, gotoStmt
}
```

### For Function Calls

```go
func (fbb *funcBlockBuilder) remodelCallExpr(...) {
    // ... existing split code ...
    
    // Add return value params to follow block
    fbb.addReturnParams(follow, callExpr)
    
    // Combine live vars with return params
    call.Follow.Args = append(liveArgs, returnArgs...)
}
```

## Testing Strategy

### Unit Tests

Test variable passing with:

```go
func Test_Blocker_VarPassing_Simple(t *testing.T) {
    pkg := blockIrcFunc(t,
        `package main`,
        `func foo(n int) int {`,
        `    x := n + 1`,
        `    y := bar(x)`,
        `    return x + y`,
        `}`,
        `func bar(a int) int { return a * 2 }`)
    
    // Verify Block 1 has Params [x, $ret0]
    // Verify FuncCallStmt has Follow.Args [x]
}
```

### Integration Tests

Use `testApps/fib` once blocking is complete:

```bash
go run . build ./testApps/fib
```

## Future Considerations

### Optimization: Unused Variables

Don't pass variables that aren't used:

```go
x := 1
y := 2
foo()  // split
return y  // x not needed in follow
```

### Optimization: Phi Nodes

For complex control flow, consider SSA-style phi nodes:

```
Block 3:
  x = phi(Block1: x1, Block2: x2)
```

This is deferred until basic variable passing works.
