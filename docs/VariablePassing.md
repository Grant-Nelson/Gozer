# Variable Passing Between Blocks

*Design document for tracking*
*variable flow across block boundaries*

## SSA Analysis

Creating blocks is very similar to control flow graphs (CFG) such as
SSA analysis that uses phi nodes (`φ()`) to deal with single assignments.
Blocks is very much based on this idea and we can use the same analysis
to help with variable passing. However, Go already looks for unused varables
with its own SSA analysis so we don't have to add this complication if not need.

For complex control flow, SSA-style phi nodes looks like:

```pseudo
Block 3:
  x = φ(Block1: x1, Block2: x2)
```

Since our blocks pass variables through arguments that have indices to indicate
which parameter the argument is for, we can do the above without the phi node
by having Block1 x1 put into parameter index i, and Block2 x2 also put into
index i, then when assigning x we simply read from index i. The parameter
indices are similar to but not the same as a phi node since we do not have
registers that we can leave a variable in. So we don't have to try to resolve
where in memory Block1 x1 and Block2 x2 live in relation to x in memory.

We are not compiling to run on metal (i.e. compile into an executable that
is run by a processor) but instead we are targetting another high-level language
so we are forced to pack an argument list and unpack it at runtime. We should try
to optimize as best as possible but remember, we are not compiling to a LLVM
or assembly. We are instead creating blocks to produce the best schedualling
in a single threaded environment but contain the target language statements
for simplicity and speed. We will not have virtual registers.

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

```pseudo
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

### BoundFunc for Closures

Function pointers with captured state:

```go
type BoundFunc struct {
    // Func is the function being referenced.
    Func *Func
    
    // Bound are values captured at function creation time.
    // These are prepended to the function's parameters when called.
    // For methods, Bound[0] is the receiver.
    // For closures, Bound contains captured variables.
    Bound []ast.Expr
}
```

Most function references have empty `Bound`. Closures and method values
populate `Bound` when the function pointer is created.

### NamedResult for Function Returns

Track named results for proper initialization and bare returns:

```go
type NamedResult struct {
    Ident      *ast.Ident  // The named result identifier
    Type       types.Type  // Result type
    NeedsZero  bool        // Must be zero-initialized (after analysis)
    Captured   bool        // Captured by closure/defer
}
```

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

```pseudo
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

### Phase 5: Named Results

Handle named results as initialized local variables:

1. At function entry, insert zero-initialization for each named result
2. Track which named results are captured by closures/defers
3. For bare `return` statements, expand to explicit return of named results
4. Apply optimization to skip zero-init where provably unnecessary

### Phase 6: Closures and Bound Variables

Handle closures and method values:

1. Detect closures that need multiple blocks (contain blocking operations)
2. Extract such closures as separate functions
3. Identify captured variables via `types.Info`
4. Create `BoundFunc` with captured values
5. For method values (`obj.Method`), bind receiver as first bound variable
6. Handle capture-by-reference for modified variables (cell pattern)

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

## Closures and Captured Variables

When a closure (anonymous function) or defer needs to be broken into multiple blocks,
it must be extracted from the block that defines it. The captured variables require
special handling.

### Bound Variables

Functions can have **bound variables** - values set when the function is created,
not when it's called. Most functions have no bound variables. Closures use bound
variables for captured state.

```go
type BoundFunc struct {
    Func   *Func           // The function to call
    Bound  []any           // Captured/bound values
}
```

When the function is invoked, bound variables are available alongside parameters.

### Closure Example

```go
func outer(n int) func() int {
    x := n * 2
    return func() int {  // closure captures x
        return x + 1
    }
}
```

The returned closure captures `x`. When creating the function pointer:
- Create `BoundFunc` with the inner function
- Bind `x`'s current value into `Bound`

When the closure is called later, `x` comes from `Bound`, not from any block's params.

### Closures Spanning Multiple Blocks

If a closure body requires blocking (contains a function call, channel op, etc.),
the closure becomes its own function with multiple blocks:

```go
func outer(n int) func() int {
    x := n * 2
    return func() int {
        y := expensive(x)  // <- causes block split in closure
        return x + y
    }
}
```

The closure is extracted as a separate function:
- Bound variables: `[x]`
- Entry block receives `x` from bound state
- Blocks pass `x` through params as usual

### Method Receivers as Bound Variables

When taking a method as a function pointer, the receiver becomes a bound variable:

```go
type Counter struct { value int }
func (c *Counter) Inc() { c.value++ }

fn := counter.Inc  // fn is BoundFunc{Func: Inc, Bound: [counter]}
```

This is handled the same way as closure captures - the receiver is bound at the
time the function pointer is created.

## Defers and Named Results

### Named Results as Local Variables

Named results (e.g., `func foo() (ok bool, err error)`) should be treated as
local variables that are:
1. Declared at function entry
2. Initialized to their type's zero value
3. Available for assignment throughout the function
4. Returned implicitly on bare `return` statements

```go
func example() (result int, err error) {
    // Equivalent to:
    // var result int    // = 0
    // var err error     // = nil
    
    result = 42
    return  // returns (42, nil)
}
```

### Defer with Recover Pattern

Named results are essential for the defer-recover pattern:

```go
func safe() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()
    
    riskyOperation()
    return nil
}
```

Here:
1. `err` is a named result, initialized to `nil`
2. The defer closure captures `err` (as a bound variable)
3. If panic occurs, recover assigns to `err`
4. The function returns the modified `err`

The closure must capture `err` by reference (or via a cell/box) so modifications
are visible to the return.

### Zero Value Returns

Named results enable zero-value returns:

```go
func zero[T any]() (v T) {
    return  // returns zero value of T
}
```

The named result `v` is initialized to `T`'s zero value. The bare `return`
returns whatever `v` currently holds.

Note: Handling generic instantiation will be be added after variable passing.
However, we need to prepare to allow for parameters with variable types.

### Implementation

At function entry, before user code:

```pseudo
Block 0 (entry):
  Params: [user params...]
  Body:
    // Initialize named results
    result := 0      // zero value for int
    err := nil       // zero value for error
    // ... user code follows
```

For bare returns, collect current values of named results:

```pseudo
return  // becomes: return result, err
```

### Optimization: Skip Zero Initialization

Using liveness/phi-node analysis, we can skip zero initialization if:
1. The named result is **always assigned** before any use, AND
2. The named result is **not captured** by any closure/defer, AND  
3. The named result is **not used in a bare return**

```go
// Can skip zero-init: always assigned before use, no closure capture
func compute(n int) (result int) {
    result = n * 2
    return result  // explicit return, not bare
}

// Cannot skip: used in bare return
func maybeZero(n int) (result int) {
    if n > 0 {
        result = n
    }
    return  // bare return - might return uninitialized if we skipped
}

// Cannot skip: captured by defer
func withDefer() (err error) {
    defer func() { log(err) }()
    // ...
}
```

## Address-Taken Variables and Closures

When a variable's address is taken OR it's captured by a closure that may
modify it, the variable cannot be passed by value between blocks.

### Cell/Box Pattern

Wrap the variable in a cell that is passed by reference:

```go
type Cell[T any] struct {
    Value T
}
```

```go
x := 1
p := &x
foo()  // split
*p = 2
use(x)  // needs updated value
```

Becomes:

```pseudo
Block 0:
  xCell := &Cell{Value: 1}
  p := &xCell.Value
  call foo(), follow=Block1, followArgs=[xCell]

Block 1:
  Params: [xCell]
  *p = 2  // p still points into xCell
  use(xCell.Value)
```

### Closure Capture by Reference

When a closure modifies a captured variable:

```go
func counter() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}
```

The closure captures `n` by reference (via cell):
- `nCell := &Cell{Value: 0}`
- Closure binds `nCell`
- `n++` becomes `nCell.Value++`

## Future Considerations

### Optimization: Unused Variables

Don't pass variables that aren't used:

```go
x := 1
y := 2
foo()  // split
return y  // x not needed in follow
```

### Closure Inlining

If a closure is:
1. Called immediately (IIFE pattern)
2. Called exactly once
3. Doesn't escape the function

Consider inlining it rather than creating bound state.
