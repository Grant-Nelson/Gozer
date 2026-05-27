# Block Costs

*Design document for block cost calculation*
*and scheduler optimization*

## Overview

Block costs provide a lightweight mechanism for the scheduler to decide when to
check for thread swapping. Instead of checking wall-clock time after every block
(which is expensive), blocks report their computational "cost" and the scheduler
accumulates this cost per thread. Only when the accumulated cost exceeds a
threshold does the scheduler perform the more expensive time check and potential
thread swap.

## Motivation

Consider a tight for-loop copying values between slices:

```go
for i := 0; i < len(src); i++ {
    dst[i] = src[i]
}
```

This compiles to a loop body block that executes repeatedly. Each iteration is
cheap (one assignment, one increment, one comparison). Without cost-based scheduling:

- **100 iterations**: Should complete without interruption - total work is trivial
- **1,000,000 iterations**: Should periodically yield to allow other threads to run

With costs, the scheduler can:
1. Let small loops complete without overhead
2. Batch iterations of large loops, checking time only periodically
3. Avoid the expense of time checks on every block return

## Cost Model

This is an initial idea of the base cost and is subject to change as we
experiment with ideas and tune the optimization for the schedular.

### Basic Rules

| Construct | Base Cost |
|-----------|-----------|
| Simple statement (assignment, increment) | 1 |
| Binary expression | 1 |
| Unary expression | 1 |
| Function call (to atomic function) | 1 |
| Conditional (if without else) | 1 |
| Conditional (if with else) | 1 |
| Index expression | 1 |
| Selector expression | 1 |

### Compound Costs

Costs are additive within a statement:

```go
x := a + b*c + d[i]
// Cost: 1 (assign) + 1 (add) + 1 (mul) + 1 (add) + 1 (index) = 5
```

### Block Cost Calculation

A block's cost is the sum of all statement costs in its body:

```go
// Block cost = 3
x := 1        // cost 1
y := x + 2    // cost 2 (assign + add)
```

Costs are **calculated statically** after block creation but
**returned at runtime** since the block must actually execute to accumulate cost.

## Scheduler Integration

### Thread State

Each thread maintains an accumulated cost:

```typescript
interface Thread {
    id: number;
    accumulatedCost: number;
    // ... other fields
}
```

### Block Return Value

Blocks return their cost along with flow control information:

```typescript
interface BlockReturn {
    cost: number;
    next: BlockRef | ReturnValue | CallInfo;
}
```

### Scheduler Logic

```typescript
const cost_threshold = 400;

function runThread(thread: Thread): void {
    while (true) {
        const result = executeBlock(thread.currentBlock);
        thread.accumulatedCost += result.cost;
        
        if (thread.accumulatedCost <= cost_threshold) {
            // Continue to next block without yielding
            ...
        }
        thread.accumulatedCost = 0;
        // Yield to the scheduler's slower determination of thread swapping
        ...
    }
}
```

### Threshold Tuning

The cost threshold balances:
- **Too low**: Frequent time checks, high overhead
- **Too high**: Long delays before yielding, poor responsiveness

Suggested starting point: **400** (roughly equivalent to a simple loop body
running 400 times).

The threshold may be configurable per-application or auto-tuned based on
observed performance. Auto-tuning would look at how much time passed
since the last time check and compare that duration against the swap duration
(the amount of time to to allow a thread to run before trying to swap).
If the duration between time checks is multiples of the swap duration then
we can adjust the threshold lower so that the time is checked more often.
Alternitively, if the duration between time checks is sigificantly less than
the swap duration, the threshold can be increased thus allowing the block to
keep running longer before checking the time. The threshold itself should
probably remain constant and a scalar is adjusted to prevent drifting too far
from the original threshold and to damplen how quicky the threshold is adjusted
to prevent the threshold from bouncing around.

## Cost Calculation Implementation

### When to Calculate

Costs are calculated in the modeler as a remodeler implementation that is
run after all blocks are created by the blocker:

```pseudo
Loader → Modeler { Blocker → Cost Calculator } → Compiler
```

This allows the cost calculator to see the final block structure. Since
the cost calculator is an implementation of a remodeler, it can be added
into the sequence of the modeler and moved around based on the target language
or removed from the modeler steps if the target language doesn't use the costs.
If costs aren't being used, the cost can be left at zero.

### IR Integration

Add cost to the Block structure:

```go
type Block struct {
    Index  int
    Hint   string
    Body   []Stmt
    Cost   int  // Calculated static cost
    // ... other fields
}
```

### Cost Visitor

Walk the block body and sum costs:

```go
func calculateBlockCost(b *Block) int {
    cost := 0
    for _, stmt := range b.Body {
        cost += statementCost(stmt)
    }
    return cost
}

func statementCost(s Stmt) int {
    switch s := s.(type) {
    case *AssignStmt:
        return 1 + exprListCost(s.Rhs)
    case *ExprStmt:
        return exprCost(s.X)
    case *IfStmt:
        return 1 + exprCost(s.Cond)
        // Note: if body is in same block, add body cost
    // ... other cases
    }
}

func exprCost(e ast.Expr) int {
    switch e := e.(type) {
    case *ast.BinaryExpr:
        return 1 + exprCost(e.X) + exprCost(e.Y)
    case *ast.UnaryExpr:
        return 1 + exprCost(e.X)
    case *ast.CallExpr:
        return 1 + exprListCost(e.Args)
    case *ast.Ident, *ast.BasicLit:
        return 0  // Simple references are free
    // ... other cases
    }
}
```

## Advanced Considerations

### High-Cost Block Splitting (Future)

If a block's static cost exceeds a maximum (e.g., 1000), consider splitting it
even without a natural split point. This prevents scenarios like unrolled loops
from blocking the scheduler. `xi := compute(i)` in the following represents a
statement that won't cause a block split to occur and not a call to some
`compute` function:

```go
// Unrolled loop - no natural split points
x1 := compute(1)
x2 := compute(2)
x3 := compute(3)
// ... 500 more statements
x500 := compute(500)
```

Without splitting, this block runs to completion. With cost-based splitting,
insert yield points:

```go
// Block 1 (cost ~250)
x1 := compute(1)
// ... 
x250 := compute(250)
goto Block2

// Block 2 (cost ~250)
x251 := compute(251)
// ...
```

**Implementation note**: This requires the same variable scope handling
as with any other block split. All live variables at the split point must be
passed to the next block.

### Dynamic Cost Adjustment

Some operations may have variable runtime cost:

- String concatenation (depends on string length)
- Slice operations (depends on slice size)
- Map operations (depends on map size)

However, we shouldn't be computing costs at runtime since that adds to the
overhead and costs are only to reduce the schedular overhead as quickly and
simply as possible. Unless we determine later that it is needed, we should
treat a block's cost as static and constant.

For v1, use pessimistic static estimates. Future versions could:
1. Use runtime size information to adjust costs
2. Profile actual execution times and feed back to cost model
   (see [Threshold Tuning](#threshold-tuning))

### Inlined If-Statement Costs

When an if-statement body remains in the same block (not split out), its cost
should be included:

```go
// Block with inline if
x := 1                    // cost 1
if x > 0 {                // cost 1 (condition)
    y := x + 1            // cost 2
}
z := x + y                // cost 2
// Total block cost: 6
```

However, the if-body only executes conditionally. Options:
1. **Pessimistic**: Always include body cost
2. **Optimistic**: Exclude body cost
3. **Average**: Include half the body cost

Recommendation: Use pessimistic (always include) for simplicity and safety.

### Atomic Function Costs

When calling an atomic function (one that doesn't yield), the atomic function's
body executes without scheduler involvement. However, the cost is still
returned so that it does contribute to the thread's accumulated cost eventhough
no yield can occur within it. This prevents the block calling the atomic
from taking too much time especially since a yeild can not occur within the
atomic. The atomic may take a variable amount of time if it contains a for-loop
so the cost will only be a best estimate for v1
(see [Dynamic Cost Adjustment](#dynamic-cost-adjustment)).

```go
//gozer:atomic
func fastCompute(x int) int {
    // This entire body runs without yielding
    return x * x + x
}

// In a non-atomic function:
y := fastCompute(10)
```

### Zero-Cost Blocks

Some blocks may have zero cost:
- Empty blocks (generated during restructuring)
- Blocks with only goto statements

Zero cost blocks should be removed as part of a cleanup remodeler being
run in the modeler after the blocker is run, but it may be run after costs have
been computed. If a zero cost block still exists after the removal then the
cost should at minimum be 1.

The only time the costs should be zero is if cost calculations are being skipped
and the schedular isn't using costs.

Zero cost blocks like a block containing a goto can be removed in many cases
since if there are only goto's that jump to the block containing only a goto,
then that block containing only a goto can be removed and the gotos that
were jumping to that block can be instead rewritten to goto the block that
the removed block was going to.

The reason we need to clamp the block to 1 at minimum is so that we don't
accidently cause a starvation of other threads by jumping between blocks with
no cost, such as a deep recursive call returning back up the call stack.
We would want to yeild occationally during all those returns if needed.

## Example: Loop Costing

```go
func sumSlice(s []int) int {
    total := 0
    for i := 0; i < len(s); i++ {
        total += s[i]
    }
    return total
}
```

After blocking:

```pseudo
Block 0 (entry, cost: 2)
    total := 0          // cost 1
    i := 0              // cost 1
    goto Block 1

Block 1 (loop body, cost: 5)
    if !(i < len(s))    // cost 3 (if + compare + len)
        goto Block 3
    total += s[i]       // cost 2 (add + index)
    goto Block 2

Block 2 (loop post, cost: 1)
    i++                 // cost 1
    goto Block 1

Block 3 (after loop, cost: 1)
    return total        // cost 1 (return is flow control)
```

With a threshold of 400 and loop body cost of 5:
- Loop can run ~80 iterations before triggering a time check
- For less than 80-element slice: completes without interruption
- For 10,000-element slice: checks time roughly every 80 iterations

## Implementation Phases

### Phase 1: Basic Cost Calculation

- Add `Cost` field to `ir.Block`
- Implement `CalculateCosts` that is a RemodelFactory and RemodelFuncExt
  - Creating this remodeler should take options that can be added to later
    for controlling the cost of different statements and/or controlling
    if some kinds of statements are counted or not.
  - Implement a `calculateBlockCost(Block)` for basic statments and
    sets the code for the given block.
  - The `RemodelFunc` iterates over all the blocks in a function and
    calls `calculateBlockCost` on each.

### Phase 2: Scheduler Integration

- This phase must be done after the TS schedular is mostly fleshed out
- Add `accumulatedCost` to thread state
- Implement threshold-based time checking
- Make threshold configurable

### Phase 3: Refinement

- Tune default threshold based on benchmarks
- Handle edge cases (zero-cost blocks, very high-cost blocks)
- Consider inline if-statement cost strategies

### Phase 4: Advanced Features (Future)

- High-cost block splitting
- Dynamic cost adjustment
- Profiling feedback loop
