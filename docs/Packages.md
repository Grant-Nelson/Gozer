# Packages

Gozer will transpile each package of Go into its own module file
for the target language (i.e. `fmt` will become `fmt.ts`).
Any imports into the package will also be imported into the module.

The main ideas are:

- Allow several different locations in the target language to use different
  packages but any dependencies to those to be shared. For example if
  one location imports `cat` and one imports `dog` and both `cat` and `dog`
  import `fmt`.
- All packages eventually depend on the `runtime` package. The `runtime`
  defines the schedular. The `runtime` may also have registers to handle
  things like resolving `//go:linkname` directives.
- Any function that is exported from a package will have an
  [externalized function](#externalized-functions) that can be called from
  the target language.
- Any function that is exported from a package will have an
  [externalized type](#externalized-types).
- The `init` functions will be called once per package after that package is
  loaded. All imported packages will be initialized before the packages
  importing them.
- When compiling into an executable with a main method for an application or
  for testing. The main package will be loaded last and have an automatically
  kicked off main method that is run after linking is resolved and `init`s have
  been called.
- When compiling with the same build flags, packages which have not changed and
  had no dependencies that were changed will not have to be recomplied.
- When compiling for tests a custom test package for `X_test` packages will be
  created and a special build of the package being tested will be built.
  Those packages will be named in a way that they do not get picked up as
  "precompiled" packages when performing another build. Additionally a
  test runner main package is added into the build that will kick of the tests
  based on the command line arguments.

Known issues with design:

- Allowing multiple executables and packages to use the same package
  dependencies may cause problems:
  - Two separate transpilations both try to implement a stub via
    a `//go:linkname` directive. We will probably have the first link made to
    the stub sticky so that it can not be redefined. Even if the executable
    exits its main thread there may still be links to code that were made by
    that executable when another executable is started and is sharing the same
    packages.
  - May cause conflicts in naming. For example if two executables define a
    package with the same name, (e.g. `cat`) but with different `cat` pacakges
    per executable, then the package names will collide.
    This is already a problem even if they were run in their own instances
    since the cached module name could collide.
  - May cause conflicting configurations in global values. For example the
    [`testing.Testing() bool`](https://pkg.go.dev/testing#Testing) is different
    if an executable is for testing or not. With two executables being run,
    one could be for tests and one not.
  - Different build flags are used to transpile the modules. This is like
    a versioning problem since the cache of the modules would have to know
    about the different build flags via a difference in the module names
    to get the correct module build.

Most of these are problems with how the code is being used and the developer
attempting to use the packages in an odd way needs come up with a different
way to handle these issues such as compiling to a single executable only.
These problems are typical problems in any code that allows multiple modules
to be create and used together. There are existing tools to help deal with
diamond dependencies on different versions, etc that can be used to help
deal with the unexpected cases.

We assume that typically there will be either only one executable,
one test executable, or a set of libraries being used. If libraries are
being used, then we assume they were either built together or were built using
the same build flags and source code versions without any conflicting names.

## Package information

Each package will export a pacakge information object.
The package information will contain all top-level exported types and functions.
Since these are all part of the package name space we know the identifiers for
them will not collide. This will also contain some methods to support linking
for `//go:linkname` directives (these may change as we are implementing the code
and figuring out how they should work).
This will not include variables and constants.
When another package imports this package it will only have to get the
package information when it is using types from that package or calling into
the package. The functions on the type package are those designed for internal
use and be quickly called.

Variables and constants will be exported from the package separate from the
pacakge information. Additionally any top-level externalized functions will be
exported from the package.

## Externalized Functions

Any function that is exported from a package will have a function that
can be called by the target language. Calling this method is not as fast
as when transpiled methods call eachother so these externalized
methods should only be called from outside of the transpiled code.

When an externalized method is called it will select the approprate
[type dictionaries](./TypesAndGenerics.md#type-dictionaries), if needed,
to handle the given arguments. This call will kick off a new
[main thread](#main-threads) to run that function call.
The returns will always be put into a promise. If the last result value is
an `error` then that will be used as the error case in the promise.
All panics that would leave the function should be recovered and returned from
the promise.

Calling these functions should feel simple and as native to the target
langauge as possible. However the types will not be automaticvally
internalized/externalized so may have to have a ["to native" function](./TypesAndGenerics.md#casting-to-target-language-types)
called on them. This is desirable in most cases since the code consuming
the transpiled code may want to use all 64-bits of a `uint64` instead of being
limited to whatever size the native integer is.

It may be possible in the future to add directives into the Go source code
(e.g. `//gozer:retype y number`) to indicate what native types to use as
parameters and result values. These would automatically internalize and
externalize from the target language type. Generic type parameters will
automatically allow this because of how generics work, but even concrete types
(e.g. the `y` in `foo[T any](x T, y uint64)`) would have to be internalized
into a `Uint64` instance to be able to be passed in as a parameter.
Without the above directive, the transpiled exported code would look like
`foo<T>(x: T, y: Uint64)` but with the driective it would transpile to
something like `foo<T>(x: T, y: number)`. Obviously we would have to check
that such a conversion can be made automatically. Also it may require
that result types must be named in order to be retyped.

## Externalized Types

Any exported named types will be accessable and have any exported methods
exposed as an externalized function. Those externalized methods will be called
from the exported type, meaning if the exported type is null then the call
may fail even though in Go it would have worked with a nil receiver.
If the user wants to call the value more safely, they'll have to use the
type information for that type that contains the receiverless functions
then pass the type into those.

The fields of an exposed type will be the internalized representation
that can be externalized with ["to native" functions](./TypesAndGenerics.md#casting-to-target-language-types).

Although we want to make these types as eazy to use from the target language,
the main focus should be that the transpiled code runs fast and is intended
to mostly be used internally. This may cause some problems when creating
bindings to modules defined in the target langauge such as React since
some of the resulting externalized types may not work as expected when
copied by React. Because of the potential complication, these designs may
need to be adjusted later to make this kind of work easier. It may be as simple
as adding a wrapper JS object that protects the internalized data from being
mangaled by things like React. At this point we don't know.

## Main Threads

The thread for an externalized function call are defaulted as a "main thread"
such that when it exits, all other threads that were started because of it will
also be killed. If any of the threads it started get killed with a panic then
the main thread is killed and the panic from the other thread is returned from
the promise.

To not have this default threading, inside of the method call before transpilation
a schedular method (see [`setMainThread` in Non-breaking Flow Controls](./Blocks.md#non-breaking-flow-controls))
may be called to indicate that the current thread should
not be treated as a main thread which allows it to exit while leaving any
threads it creates running and if any thread panics and there is no main
thread for it, then it will log the panic and exit by itself.

A thread can also call `setMainThread`to become it's own main thread
thus taking it out of the main thread it may have been in and to become the
main thread for any new threads it creates.

When the application is compiled with a main method, the main method will
automatically be invoked once all packages have been loaded and it will
be invoked as a main thread. Since the schedular package a shared dependency
amongs all transpiled packages, this mean that two or more applications
can be started and each will run in the same schedular but act as if they
are running by itself because of how the main thread is being managed.

When the application is compiled to run tests, examples, and benchmarks then
it will be given it's own main method to kick off the tests and that main
method will start a main thread.
