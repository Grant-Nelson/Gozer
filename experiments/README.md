# Experiments

These are a set of experiments that were created to try out different
ideas while prototyping and designing Gozer.

## Running Experiments

1. Run `go build main.go` in the Gozer/cmd/serve folder
2. Run `../../cmd/serve/main.exe` from the experiment folder
3. Open the browser to the URL the server says it is serving to,
   typically `http://localhost:8080/`

## List of Experiments

- [exp001](./exp001/): The first experiment were the application is simply
  a bunch of goto's between blocks for a single function using hand-written
  blocks and an experimental simple scheduler.
- [exp002](./exp002/): This experiment add to the prior experiment to include
  calling another method and returning from the method.
