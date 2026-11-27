# Composure

This is similar to `prove` from Go, but will also test for type confirmation,
such as when faking a 16-bit number with 32-bit numbers (or a JS number),
then when does an op like `x + y` require the result to be truncated back into 16-bit.
See: [cmd/compile/internal/ssa/prove](https://cs.opensource.google/go/go/+/master:src/cmd/compile/internal/ssa/prove.go)

## How `prove` Works Under the Hood

In Go's compiler (`cmd/compile/internal/ssa/prove.go`), the `prove` pass is a
**dominator-based conditional fact propagation pass** over the SSA control-flow graph (CFG).
As it walks blocks in dominator order, it learns facts from branches, comparisons,
and operations, narrows value domains, removes redundant branches and
bounds/nil checks, and cleans up unreachable edges.

`prove` combines three main mechanisms to deduce what is true at a given point
in execution:

1. **Per-Value Domain Intervals (`limit` / range bounds):** Each integer/pointer
   `Value` carries lower and upper bounds (both signed and unsigned perspectives,
   plus non-nil / nil state) that tighten as facts are learned.

2. **Relational Constraints (`poset` / partial order):** `prove` tracks relative
   ordering relations between SSA values (such as `x < y` or `x <= y`) so a bound
   on one variable propagates to related variables or slice lengths.

3. **Dominator Tree Traversal with Scoped Facts:** Walking the CFG along the
   dominator tree lets facts learned on a true/false branch apply to dominated
   blocks, while pushing/popping or scoping facts so sibling branches do not
   leak into each other.

### Designing a `golang.org/x/tools/go/ssa` Analyzer

In `golang.org/x/tools/go/ssa`, SSA values are distinct from instructions
(`ssa.Instruction`s may produce a `ssa.Value` or merely perform a side effect/check).
To record state *at each instruction*, maintain a **state snapshot before
(and/or after) each `ssa.Instruction**`, keyed by instruction or `(block, instrIndex)`,
mapping each relevant `*ssa.Value` to an abstract domain entry.

#### 1. Tracking Pointer Nilness (`Nil`, `NotNil`, `MaybeNil`)

Represent pointer/interface/slice/map/func/chan nil-ness as a three-state lattice:

- **`NotNil`**: Guaranteed non-nil at this point.
- **`IsNil`**: Guaranteed nil at this point.
- **`MaybeNil`**: Unknown or could be either (top/default state).

##### Sources of nilness in `go/ssa`

- **Guaranteed `NotNil` producers:** `*ssa.Alloc`, `*ssa.MakeSlice`, `*ssa.MakeMap`,
  `*ssa.MakeChan`, `*ssa.MakeClosure`, `*ssa.FieldAddr`, and `*ssa.IndexAddr`
  once the base pointer has already been dereferenced or checked.

- **Dereferences & assertions:** Dereferencing a pointer (`*ssa.UnOp` with
  `token.MUL`, or `*ssa.FieldAddr`) implies the pointer is `NotNil` after the
  operation (and triggers a panic if it was `IsNil` before). A successful type
  assertion (`x.(T)` or `ok == true` in a comma-ok assertion) also establishes
  non-nil for the asserted value when non-nil is required.

- **Branch refinement (`*ssa.If`):** When a block ends with `*ssa.If` testing
  `p == nil` or `p != nil` (via `*ssa.BinOp` with `EQL` or `NEQ` against a nil
  `*ssa.Const` or known nil/not-nil), propagate `IsNil` down the matching successor
  and `NotNil` down the opposite successor—when that successor is only reached
  from that branch (or dominates/is dominated cleanly).

- **Phi joins (`*ssa.Phi`):** At a join point, combine incoming predecessor
  states via lattice join: `NotNil ⊔ NotNil = NotNil`, `Nil ⊔ Nil = Nil`,
  and anything else becomes `MaybeNil`.

#### 2. Detecting Addition and Multiplication Overflow Risk

In Go, signed integer overflow is defined (two's complement wraparound),
while unsigned overflow wraps modulo $2^N$, and both are common sources of
subtle logic or sizing bugs. Because `go/ssa` represents integers with
`types.Basic` (giving bit width and signed-ness), track a
**signed interval `[sMin, sMax]`** and an **unsigned interval `[uMin, uMax]`**
for each integer `Value` based on its bit width $W \in \{8, 16, 32, 64\}$.

Given intervals $[a_{\text{lo}}, a_{\text{hi}}]$ for `x` and
$[b_{\text{lo}}, b_{\text{hi}}]$ for `y`:

- **Addition (`*ssa.BinOp` with `token.ADD`):**
  - Compute exact mathematical range $[a_{\text{lo}} + b_{\text{lo}},\, a_{\text{hi}} + b_{\text{hi}}]$
    using wider arithmetic (e.g., `math/big.Int` so the analyzer itself never overflows).
  - **Signed overflow possible** if the exact sum range extends outside
    $[\text{MinSigned}(W), \text{MaxSigned}(W)]$.
  - **Unsigned overflow possible** if the exact sum range extends above
    $\text{MaxUnsigned}(W)$ (or below `0` for subtraction).
  - If no overflow occurs, clamp the resulting value's interval to the
    target width's representable range.
- **Multiplication (`*ssa.BinOp` with `token.MUL`):**
  - Multiplication has four corner products: $a_{\text{lo}} b_{\text{lo}}$,
    $a_{\text{lo}} b_{\text{hi}}$, $a_{\text{hi}} b_{\text{lo}}$, and
    $a_{\text{hi}} b_{\text{hi}}$.
  - Using `big.Int`, take the minimum and maximum of those four products as the exact
    product interval $[p_{\text{lo}}, p_{\text{hi}}]$.
  - Flag **overflow possible** if $p_{\text{lo}} < \text{Min}(W)$ or
    $p_{\text{hi}} > \text{Max}(W)$ for the operation's signed-ness.

#### 3. Identifying Upper and Lower Bounds Checks on Indexing and Slicing

In Go, indexing or slicing (`s[i]`, `a[i]`, `str[i]`, `s[low:high]`, `s[low:high:max]`)
panics if indices violate bounds:

- **Indexing (`*ssa.IndexAddr` for arrays/slices, `*ssa.Index` for strings):**
  - Index `i` must satisfy $0 \le i < \text{len}$.
  - **Lower-bound check needed** if `i`'s known lower bound can be less than `0`
    (for signed indices, or if unknown).
  - **Upper-bound check needed** if `i`'s known upper bound is not proven strictly
    less than the known length (or if length is unknown).
- **Slicing (`*ssa.Slice`):**
  - For a 2-index slice `x[low:high]`, Go requires
    $0 \le \text{low} \le \text{high} \le \text{cap}$ (or $\le \text{len}$
    for strings). For a 3-index slice `x[low:high:max]`, Go also requires
    $\text{high} \le \text{max} \le \text{cap}$.
  - Check **lower bound** (`low >= 0`), **monotonicity** (`low <= high`, and
    `high <= max` when present), and **upper bound** (`high <= cap` or `max <= cap`)
    against the tracked intervals and relational facts.

### Practical Architecture Tip for `go/ssa`

Loops with back-edges are where naive forward passes loop infinitely or
lose precision. Either:

1. Run a **worklist dataflow pass** over blocks with a widening/clamping
   step on intervals (widening toward $\pm\infty$ or type min/max when a
   bound grows across a back-edge, then narrowing)
2. Or seed loop induction bounds from common loop header comparisons
   (`i < n`, `i <= len(s)`) before propagating facts along dominators.

## Gozer Implementation Notes

### Nullability

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

### Numeric Affirmation

When doing an op like `x++` or `x*y` we need to check what composure ranges
are possible to determine when overflows are possible or in a comparison.

TODO: When to enforce N-bit size or sign.

### Boolean Affirmation

TODO: `if false {` or `if true {`

### Index Boundaries

TODO: Check for out-of-bounds only when needing too.
