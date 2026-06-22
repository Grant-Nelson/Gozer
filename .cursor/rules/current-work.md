# Current Work Focus

## Immediate Goal

Validate the scheduler/blocking approach by getting simple applications working.

## Target: Simple Applications

Applications that work with:
- Functions (including recursive calls)
- Basic types (int, bool, string)
- Control flow (if, for, goto, break, continue, labels)
- Function calls in various contexts

**NOT including** (deferred):
- Type declarations
- Value (const/var) declarations at package level  
- Standard library
- Channels, goroutines, select
- Interfaces
- Generics

## Current Blockers (pun intended)

The `blocker.go` file has several incomplete implementations marked with:
- `crumb.DropMsg("Unimplemented")`
- `panic(faults.New("unimplemented"))`
- `// TODO:` comments

### High Priority Incomplete Items

1. **Variable passing between blocks** - See `docs/VariablePassing.md`
   - `Block.Params` population
   - `BlockRef.Args` population
   - Synthetic return value identifiers

2. **Call expressions in context** (`remodelCallExpr`)
   - Works when call is the only thing in statement
   - Fails when call is inside binary expression, assignment RHS, etc.
   - Lines 566-589 in `blocker.go`

3. **Range statements** (`remodelRangeStmt`)
   - Currently panics
   - Lines 307-328 in `blocker.go`

### Medium Priority

4. **Logical AND/OR** (`remodelLogicalAndExpr`, `remodelLogicalOrExpr`)
   - Short-circuit evaluation needs special blocking
   - Lines 558-564 in `blocker.go`

5. **Receive expressions** (`remodelReceiveExpr`)
   - Channel receive is a blocking operation
   - Line 531-533 in `blocker.go`

6. **Fallthrough** (`remodelFallThroughBranchStmt`)
   - Switch case fallthrough
   - Lines 459-470 in `blocker.go`

## Test Status

Run tests: `go test ./project/modeler/remodel/blocker/...`

Tests with incomplete expected output (marked `// TODO: Finish`):
- `Test_Blocker_ImplicitReturns_TerminatingFor`
- `Test_Blocker_ImplicitReturns_TerminatingSwitch`
- `Test_Blocker_ForRange`
- `Test_Blocker_Call_AB`
- `Test_Blocker_Call_ABInBinary`
- `Test_Blocker_Call_ABWithDefine`
- `Test_Blocker_Call_Recursive`

## Next Steps

1. Design and implement variable tracking (liveness analysis)
2. Complete `remodelCallExpr` for all contexts
3. Add tests that verify variable passing
4. Get `testApps/fib` working end-to-end
