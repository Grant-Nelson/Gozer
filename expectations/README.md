# Expectations

This is a collection of tests (or example code) that is run periodically
to ensure that the assumptions made about the Go compiler and standard library
have not changed.

If the Go compiler or standard libraries change in a way that causes
these to fail, we will have to adjust the Gozer code to take into account
the new changes, i.e. repair any code that expected the code to run differently.
