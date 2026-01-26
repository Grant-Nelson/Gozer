# Intermediate Representation

*This documentation is design notes and ideas.*
*They may not reflect the final implementation.*

The intermediate representation (IR) will be able to represent the
Go code almost like it was originally but it is not guaranteed to be
bidirectional. It will contain additional information and constructs
that Go does not provide, such as block information, type information,
and value range constraints.

When the code is first moved over to the IR, it will be very similar
to the source Go code. Then by running several conversions, optimizations, and
analysis the representation will be put into a form that can be
converted into the target language. Depending on the target language,
different modifications can be made to the code. Regardless of the modifications
the original functionality of the code should be maintained as much as possible
so that result executes as similarly to if the Go was run by itself.
