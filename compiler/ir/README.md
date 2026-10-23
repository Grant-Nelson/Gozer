# Intermediate Representation

The intermediate representation (IR) will be able to represent Go code
and target languages. The IR node interfaces have exported methods
so that custom IR nodes specific to a target language can be added into
the IR as needed.

The IR will not store Go's AST exactly. It will resolve type nodes
into the types that they represent and replace nodes like identifiers
with specific nodes based on the type that the identifier represents.
This means that the original Go AST will not be able to be recreated
from the IR. However, the IR will contain positional information to
be able to properly identify where in the GO code different nodes
came from.

When the code is first moved over to the IR, it will be similar to the
source Go code. Then by running several conversions, optimizations, and
analysis the representation will be put into a form that can be
converted into the target language. Depending on the target language,
different modifications can be made to the code.
Regardless of the modifications, the original functionality of the code
should be maintained as much as possible so that result executes as
close as possible to the original Go code.

## Block Representation

See [Blocks](../../docs/Blocks.md) documentation
