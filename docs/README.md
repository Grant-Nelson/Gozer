# Developer Documentation

- [Developer Documentation](#developer-documentation)
  - [Backend](#backend)
    - [Augmenter](#augmenter)
  - [Frontend](#frontend)

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
