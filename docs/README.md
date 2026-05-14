# Developer Documentation

- [Developer Documentation](#developer-documentation)
  - [Repo Structure](#repo-structure)
  - [Backend](#backend)
    - [Augmenter](#augmenter)
  - [Frontend](#frontend)

## Repo Structure

- **avail**: S collection of available tool library
  - **args**: A tool for reading command line arguments
  - **assert**: A collection of assertions to add into code to quickly check assumptions
  - **astTools**: Useful tools for reading and changing the Go AST
  - **crumb**: A tool to assist in debugging by printing the location of the crumb when it is reached
  - **faults**: A specialized error creator that add additional data and stack traces to an error
  - **iterator**: Adds additional functionality to the Go iter
  - **logger**: Adds hierarchical output for displaying states
  - **predicate**: A collection of predicate functions typically used in iterators
  - **source**: A tool for reading source files that allows redirecting paths,
      providing overrides of file data, skipping files, etc
- **cmd**: A collection of applications that can be used during development
  - **esb**: An app to convert TS into JS and optionally minify the output
  - **serve**: An app that serves TS on a browser locally using a default HTML page
- **docs**: A collection of documentation, diagrams, and images
- **experiments**: A collection of self-contained hand written outputs to quickly
    test different shapes of the output that the transpilation could match
- **project**: // TODO: Finish
- **testApps**: // TODO: Finish
- **tools**: // TODO: Finish

## Backend

The backend is the part of the code that loads a Go project and dependencies,
and performs some preparation of project for the [frontend](#frontend).
The following preparations are performed on a project.

### Augmenter

The augmenter is part of loading a project.
The Go files for a project and the dependencies are modified by the augmenter
so that code patterns that will not work for the target language
can be removed and replaced with something that will work.
For example, some pointers will not function in all languages,
such as Javascript, so need to be replaced via an augmenter.

The augmenter loads augmentation code that defines how to modify the code
being augmented. The augmentation code is annotated with directives.
For more information see [Augmenter](../project/loader/mods/augmenter/README.md).

## Frontend

The frontend takes a prepared project from the [backend](#backend) and
converts it into the target language.
There are different frontends for each target language.
