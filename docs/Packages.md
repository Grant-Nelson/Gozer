# Packages

*Design document for package transpilation*
*and module structure*

Gozer will transpile each package of Go into its own module file
for the target language (i.e., `fmt` will become `fmt.ts`).
Any imports into the package will also be imported into the module.

- [Packages](#packages)
  - [Main Ideas](#main-ideas)
  - [Known Issues](#known-issues)
  - [Package Loading Lifecycle](#package-loading-lifecycle)
  - [Runtime Package](#runtime-package)
  - [Package Information](#package-information)
  - [Externalized Functions](#externalized-functions)
    - [Type Dictionary Selection](#type-dictionary-selection)
    - [Promise Handling](#promise-handling)
    - [Future: Parameter Retyping](#future-parameter-retyping)
  - [Externalized Types](#externalized-types)
  - [Main Threads](#main-threads)
  - [Caching and Build Flags](#caching-and-build-flags)

## Main Ideas

- **Shared dependencies**: Allow several different locations in the target language
  to use different packages but share common dependencies. For example, if
  one location imports `cat` and one imports `dog` and both `cat` and `dog`
  import `fmt`, there is only one `fmt` module loaded.

- **Runtime package**: All packages eventually depend on the `runtime` package.
  The `runtime` defines the scheduler. The `runtime` may also have registries
  to handle things like resolving `//go:linkname` directives.

- **Externalized functions**: Any function that is exported from a package will
  have an [externalized function](#externalized-functions) that can be called
  from the target language.

- **Externalized types**: Any type that is exported from a package will have an
  [externalized type](#externalized-types).

- **Init functions**: The `init` functions will be called once per package after
  that package is loaded. All imported packages will be initialized before the
  packages importing them (dependency order).

  *TODO: Check that all `init`s can be called like this and we don't have to*
  *defer them until linking is finished or something like that.*

- **Main method**: When compiling into an executable with a main method for an
  application or for testing, the main package will be loaded last and have an
  automatically kicked off main method that is run after linking is resolved
  and `init`s have been called.

- **Incremental compilation**: When compiling with the same build flags, packages
  which have not changed and had no dependencies that were changed will not
  have to be recompiled.

- **Test compilation**: When compiling for tests, a custom test package for
  `X_test` packages will be created and a special build of the package being
  tested will be built. Those packages will be named in a way that they do not
  get picked up as "precompiled" packages when performing another build.
  However, if the same test is being run multiple times, the built test packages
  could be reused without recompiling to make calls like `gozer test ./...` only
  have to rebuild tests for packages that have changed.

  Additionally a test runner main package is added into the build that will
  kick off the tests based on the command line arguments. Building tests will
  be similar to how Go compiles tests where each package and `X_test` package
  pair will be built and run before moving to another package to build and run.
  Since a test may change package variables, each test package being run will
  need to reset and reload the modules. This can be done simply by compiling a
  test for a package, running it in node.js, shutting down node.js, then
  continue onto the next test package.

## Known Issues

Allowing multiple executables and packages to use the same package
dependencies may cause problems:

- **Conflicting linkname directives**: Two separate transpilations both try to
  implement a stub via a `//go:linkname` directive. We will probably have the
  first link made to the stub sticky so that it cannot be redefined. Even if
  the executable exits its main thread there may still be links to code that
  were made by that executable when another executable is started and is
  sharing the same packages.

- **Name collisions**: May cause conflicts in naming. For example, if two
  executables define a package with the same name (e.g., `cat`) but with
  different `cat` packages per executable, then the package names will collide.
  This is already a problem even if they were run in their own instances
  since the cached module name could collide.

- **Global configuration conflicts**: May cause conflicting configurations in
  global values. For example, [`testing.Testing() bool`](https://pkg.go.dev/testing#Testing)
  is different if an executable is for testing or not. With two executables
  being run, one could be for tests and one not.

- **Build flag variations**: Different build flags are used to transpile the
  modules. This is like a versioning problem since the cache of the modules
  would have to know about the different build flags via a difference in the
  module names to get the correct module build.

Most of these are problems with how the code is being used. The developer
attempting to use the packages in an odd way needs to come up with a different
way to handle these issues, such as compiling to a single executable only.
These problems are typical problems in any code that allows multiple modules
to be created and used together. There are existing tools to help deal with
diamond dependencies on different versions, etc., that can be used to help
deal with the unexpected cases.

We assume that typically there will be either only one executable,
one test executable, or a set of libraries being used. If libraries are
being used, then we assume they were either built together or were built using
the same build flags and source code versions without any conflicting names.

## Package Loading Lifecycle

When a package is used, the following sequence occurs:

```text
1. Load Dependencies (recursive)
   └── For each import, load that package first

2. Load Module
   └── Import the .ts/.js file for this package

3. Initialize Package
   ├── Resolve any //go:linkname directives that can be resolved
   ├── Run all init() functions in source order
   └── Package-level variables are initialized

4. Ready for Use
   └── Exported functions and types are accessible
```

For an executable with a `main` function:

```text
1. Load package with normal package life cycle
2. Start main() as a main thread
3. Wait for main thread to exit (or keep-alive)
```

## Runtime Package

The `runtime` package is the foundational dependency for all transpiled packages.
It provides:

- **Scheduler**: Thread management, block execution, cost tracking
  (see [Blocks.md](./Blocks.md))
- **Type system support**: Go interface wrappers, type assertions,
  type dictionaries and information for basic and built-in types
  (e.g. `int`, `uint64`, `float64`,  `map`, `chan`, `slice`), etc.
  (see [TypesAndGenerics.md](./TypesAndGenerics.md))
- **Linkname registry**: Resolution of `//go:linkname` directives
- **Panic/recover**: Panic propagation and recovery mechanism

```typescript
// Example runtime imports in a transpiled package
import { scheduler, GoInterface, Slice, Channel } from './runtime';
```

## Package Information

Each package will export a package information object.
The package information will contain all top-level exported types and functions.
Since these are all part of the package namespace we know the identifiers for
them will not collide.

```typescript
// Example: fmt package information
export const $pkg = {
    name: "fmt", // package name
    path: "fmt", // package path

    // Type information for exported types
    Stringer: {
        $name: "Stringer",
        $kind: interfaceKind,
        $signatures: [
            { name: "String", results: [ stringKind ] },
        ],
        $zero: () => { ... },
    },
    
    // Functions (receiverless, for internal fast calls)
    // and type information for the function.
    Println: {
        $name: "Println", // go name
        $kind: functionKind,
        $signature: { ... },
        $invoke: function(args: any[], dicT: TypeDict): BlockReturn { ... },
    },
    Sprintf: { ... },
    
    // Linkname support
    $linkname: (name: string, impl: Function) => { ... },
    $getLink: (name: string) => { ... },
};
```

This will also contain some methods to support linking for `//go:linkname`
directives (these may change as we are implementing the code and figuring
out how they should work).

Package information will **not** include variables and constants directly.
Variables and constants will be exported from the package separately from the
package information. Additionally, any top-level externalized functions will be
exported from the package.

When another package imports this package it will only have to get the
package information when it is using types from that package or calling into
the package. The functions on the package information are those designed for
internal use and to be quickly called.

## Externalized Functions

Any function that is exported from a package will have a function that
can be called by the target language. Calling this method is not as fast
as when transpiled functions call each other, so these externalized
functions should only be called from outside of the transpiled code.

```typescript
// Externalized function - used by target language code
export async function Println(...args: any[]): Promise<void> {
    // 1. Select type dictionaries based on argument types
    // 2. Create a new main thread
    // 3. Execute blocks
    // 4. Return promise with results
}
```

### Type Dictionary Selection

When an externalized method is called it will select the appropriate
[type dictionaries](./TypesAndGenerics.md#type-dictionaries), if needed,
to handle the given arguments. This is done by inspecting the runtime types
of the arguments and looking up the corresponding type dictionary.

```typescript
// Example: generic function externalization
// Go: func Max[T constraints.Ordered](a, b T) T

export async function Max<T>(a: T, b: T, dicT?: Ordered<T>): Promise<T> {
    // If dicT not provided, determine from argument types
    if (!dicT) {
        dicT = lookupOrderedDict(typeof a);
    }
    // ... execute with selected dictionary
}
```

### Promise Handling

The call will kick off a new [main thread](#main-threads) to run that function.
The returns will always be put into a promise:

- **Success case**: Promise resolves with return values
- **Error case**: If the last result value is an `error`, that is used as the
  rejection reason
- **Panic case**: All panics that would leave the function are recovered and
  returned as promise rejections

```typescript
// Go: func ReadFile(name string) ([]byte, error)

// Externalized as:
export async function ReadFile(name: string): Promise<Uint8Array> {
    return new Promise((resolve, reject) => {
        const thread = scheduler.newMainThread();
        thread.onComplete = (data: Uint8Array, err: GoError | null) => {
            if (err !== null) {
                reject(err);
            } else {
                resolve(data);
            }
        };
        thread.onPanic = (panicValue: any) => {
            reject(new GoPanic(panicValue));
        };
        // Start execution...
    });
}
```

### Future: Parameter Retyping

Calling these functions should feel simple and as native to the target
language as possible. However, the types will not be automatically
internalized/externalized so may have to have a
["to native" function](./TypesAndGenerics.md#casting-to-target-language-types)
called on them. This is desirable in most cases since the code consuming
the transpiled code may want to use all 64-bits of a `uint64` instead of being
limited to whatever size the native integer is.

It may be possible in the future to add directives into the Go source code
(e.g., `//gozer:retype y number`) to indicate what native types to use as
parameters and result values. These would automatically internalize and
externalize from the target language type. Generic type parameters will
automatically allow this because of how generics work, but even concrete types
(e.g., the `y` in `foo[T any](x T, y uint64)`) would have to be internalized
into a `Uint64` instance to be able to be passed in as a parameter.

Without the above directive, the transpiled exported code would look like
`foo<T>(x: T, y: Uint64)` but with the directive it would transpile to
something like `foo<T>(x: T, y: number)`. Obviously we would have to check
that such a conversion can be made automatically. Also, it may require
that result types must be named in order to be retyped.

## Externalized Types

Any exported named types will be accessible and have any exported methods
exposed as an externalized function. Those externalized methods will be called
from the exported type, meaning if the exported type is null then the call
may fail even though in Go it would have worked with a nil receiver.

If the user wants to call the method more safely, they'll have to use the
type information for that type that contains the receiverless functions
then pass the type into those.

```typescript
// Go type with method
// type Buffer struct { ... }
// func (b *Buffer) Write(p []byte) (int, error)

// Externalized type
export class Buffer {
    private $data: Uint8Array;
    
    // Externalized method - may fail if this is null
    async Write(p: Uint8Array): Promise<number> {
        // ... calls internal blocks
    }
}

// Safe alternative using type info
const n = await BufferTypeInfo.Write(maybeNullBuffer, data);
```

The fields of an exposed type will be the internalized representation
that can be externalized with
["to native" functions](./TypesAndGenerics.md#casting-to-target-language-types).

Although we want to make these types as easy to use from the target language,
the main focus should be that the transpiled code runs fast and is intended
to mostly be used internally. This may cause some problems when creating
bindings to modules defined in the target language such as React since
some of the resulting externalized types may not work as expected when
copied by React. Because of the potential complication, these designs may
need to be adjusted later to make this kind of work easier. It may be as simple
as adding a wrapper JS object that protects the internalized data from being
mangled by things like React. At this point we don't know.

## Main Threads

The thread for an externalized function call is defaulted as a "main thread"
such that when it exits, all other threads that were started because of it will
also be killed. If any of the threads it started get killed with a panic then
the main thread is killed and the panic from the other thread is returned from
the promise.

To not have this default threading, inside of the method call before transpilation
a scheduler method (see [`setMainThread`](./Blocks.md#non-breaking-flow-controls))
may be called to indicate that the current thread should not be treated as a
main thread which allows it to exit while leaving any threads it creates running
and if any thread panics and there is no main thread for it, then it will log
the panic and exit by itself.

A thread can also call `setMainThread(true)` to become its own main thread
thus taking it out of the main thread it may have been in and to become the
main thread for any new threads it creates.

When the application is compiled with a main method, the main method will
automatically be invoked once all packages have been loaded and it will
be invoked as a main thread. Since the scheduler package is a shared dependency
among all transpiled packages, this means that two or more applications
can be started and each will run in the same scheduler but act as if they
are running by themselves because of how the main thread is being managed.

When the application is compiled to run tests, examples, and benchmarks then
it will be given its own main method to kick off the tests and that main
method will start a main thread.

## Caching and Build Flags

When compiling packages, the build flags affect the output. To support
incremental compilation with different configurations:

- **Module naming**: Cached modules include a hash of relevant build flags
  in their filename or path (e.g., `fmt_abc123.ts`)
- **Dependency tracking**: If a dependency is recompiled (different flags or
  source changes), dependents must also be recompiled
- **Build flag examples**:
  - Target architecture (`GOARCH`)
  - Build tags (`-tags`)
  - Test mode (`-test`)

```text
cache/
├── fmt_release.ts       # Normal build
├── fmt_test.ts          # Test build (may have different behavior)
└── fmt_debug.ts         # Debug build with extra logging
```

The exact caching strategy will be refined during implementation, but the
key principle is that different build configurations produce isolated outputs
that don't conflict with each other.
