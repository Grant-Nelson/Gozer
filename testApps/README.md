# Test Apps

These are small test applications used to perform integration test

## Adding A Test App

1. Create a new folder for the test app
2. Add a main.go file to that folder
3. Add the build constraint `//go:build testApp` to all *.go files
4. Additional build constraints for different Go versions may be added
5. To check an expected output (TODO: once a test runner is finished)
   add a comment block starting with `// Output:` at the end of the main.go.
   The expected output follows prefixed with `//<single space>` that is ignored.

## Running Test Apps

> TODO: Add a test app runner that checks the output against an expected
> output similar to the Go repo known issue tests, then update these steps.
> The test runner needs to run the application in Go and check the output
> if it has one and then run the same test in TS via node.js and run other
> targets. Use build constraints to skip over specific build targets.

To run these test apps manually:

1. Build the gozer compiler main
2. Build the test app
3. Inspect the resulting output

Example:

- On windows: `go build -o gozer .\main.go ; .\gozer.exe build -v .\testApps\fib\main.go`
- On mac/linux: `go build -o gozer main.go && gozer build -v ./testApp/fib/main.go`
