# Gozer Agent Guide

This document provides context for AI agents working on the Gozer codebase.

## Project Summary

**Gozer** is a Go-to-other-languages transpiler. It enables writing code once in Go and transpiling it to target languages (currently TypeScript). The transpiled code supports pseudo-multithreaded execution in single-threaded environments via a scheduler and block-based execution model.

### Key Concepts

- **Blocks**: Functions are decomposed into statement blocks that can be scheduled cooperatively
- **Scheduler**: Controls pseudo-multithreaded execution by swapping between blocks
- **Atomic Functions**: Functions marked `//gozer:atomic` run without scheduler interruption

## Architecture: Phased Pipeline

Gozer is designed as a **phased pipeline** where each phase has a well-defined input/output contract. Phases should be worked on with **minimal overlap** - changes primarily occur within a phase, with cross-phase changes limited to shared data structures.

```
┌─────────┐    ┌─────────┐    ┌──────────┐    ┌──────────┐
│ Loader  │───▶│ Modeler │───▶│ Compiler │───▶│  Output  │
└─────────┘    └─────────┘    └──────────┘    └──────────┘
     │              │               │
     ▼              ▼               ▼
  AST + Types    IR Blocks      Target Code
```

### Phase 1: Loader (`project/loader/`)

**Input**: Go source files, patterns, configuration  
**Output**: `*project.Project` with parsed AST and type information

Responsibilities:
- Load packages via `golang.org/x/tools/go/packages`
- Apply modifiers (augmenter, type checker, cache, package dropper)
- Parse files into Go AST with full type information

**Key Files**: `loader.go`, `mods/` directory

### Phase 2: Modeler (`project/modeler/`)

**Input**: `*project.Project` with AST  
**Output**: `*project.Package` with IR (`pkg.Ir`)

Responsibilities:
- Convert Go AST to intermediate representation (IR)
- Break functions into schedulable blocks (blocker)
- Apply remodelers for target-specific transformations
- Track variable flow between blocks (inputs/outputs)

**Key Files**: `modeler.go`, `ir/`, `remodel/`

### Phase 3: Compiler (`targets/`)

**Input**: `*project.Package` with IR  
**Output**: Target language source code

Responsibilities:
- Emit target language code from IR
- Generate scheduler integration code
- Handle target-specific idioms

**Key Files**: `typescript.go` (others to be added)

## Shared Data Structures

Cross-phase communication happens through these structures:

| Structure | Location | Purpose |
|-----------|----------|---------|
| `Project` | `project/project.go` | Container for all packages |
| `Package` | `project/package.go` | Single Go package with AST and IR |
| `ir.Package` | `project/modeler/ir/package.go` | IR for a package |
| `ir.Func` | `project/modeler/ir/func.go` | Function as collection of blocks |
| `ir.Block` | `project/modeler/ir/block.go` | Statement block with inputs/outputs |
| `ir.Stmt` | `project/modeler/ir/stmt.go` | IR statements |

## Current Development Priority

### Focus: Blocker and Variable Flow

The **immediate priority** is completing the blocker (`project/modeler/remodel/blocker/`) to handle simple applications with:
- Function calls (including recursive calls)
- Basic types (int, bool, string, etc.)
- Control flow (if, for, goto, labels, break, continue)

**Explicitly deferred** (for now):
- Type declarations (`addTypeSpec`)
- Value declarations (`addValueSpec`)  
- Standard library usage
- Channels, goroutines, select

This allows validating the scheduler/blocking approach before adding complexity.

### Block Variable Passing Design

When a block is split (e.g., at a function call), we must track:

1. **Block Inputs** (`Block.Params`): Variables the block needs to receive
2. **Block Outputs** (`BlockRef.Args`): Variables passed when jumping to another block

Example transformation:
```go
// Original
func foo(n int) int {
    x := n + 1
    y := bar(x)    // <- split point
    return x + y
}

// After blocking
// Block 0: entry
//   x := n + 1
//   call bar(x), follow=Block1, followArgs=[x]
//
// Block 1: after call (params: [x, $ret0])
//   y := $ret0
//   return x + y
```

The blocker must:
1. Identify live variables at split points
2. Determine which variables are needed in the follow block
3. Populate `Block.Params` and `BlockRef.Args` accordingly

## Hierarchical Plan Structure

Plans in this project are organized hierarchically. A plan may have subplans that must be completed in order.

### Example: Complete IR Modeler

```
Plan: Complete IR Modeler
├── Subplan: Variable passing for blocks
│   ├── Task: Track live variables at block boundaries
│   ├── Task: Populate Block.Params from predecessors
│   └── Task: Populate BlockRef.Args for successors
├── Subplan: Finish function blocking
│   ├── Task: Handle call expressions in all contexts
│   ├── Task: Handle range statements
│   ├── Task: Handle switch/select statements
│   └── Task: Handle logical AND/OR short-circuit
└── Subplan: Add types and values
    ├── Task: Implement addTypeSpec
    └── Task: Implement addValueSpec
```

## Directory Overview

| Directory | Purpose |
|-----------|---------|
| `avail/` | Utility libraries (args, AST tools, errors, logging) |
| `cmd/` | Development utilities (esb, serve) |
| `docs/` | Design documentation |
| `experiments/` | Hand-written output tests |
| `project/` | Core pipeline (loader, modeler, IR) |
| `targets/` | Target language implementations |
| `testApps/` | Test applications |
| `tools/` | CLI tools (build, run, test, list) |

## Testing

The blocker has unit tests in `project/modeler/remodel/blocker/blocker_test.go`. Run with:

```bash
go test ./project/modeler/remodel/blocker/...
```

Tests use inline Go source and verify the resulting block structure.

## Key Documentation

- `docs/Blocks.md` - Block representation and scheduler design
- `docs/Patterns.md` - Code patterns and transformations
- `docs/Future.md` - Future feature plans
- `project/loader/mods/augmenter/README.md` - Augmenter system

## Code Style Notes

- Use `faults.ErrGroup` for collecting multiple errors
- Use `logger.Logger` for hierarchical verbose output
- IR statements mirror Go AST but are mutable and block-aware
- Directives use `//gozer:` prefix (e.g., `//gozer:atomic`)
