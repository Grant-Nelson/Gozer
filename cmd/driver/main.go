package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/cmd/driver/internal/driver"
)

func main() {
	ctx, cancel := withInterrupt(context.Background())
	defer cancel()

	if err := run(ctx, os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, `driver failed: %v`, err)
	}
}

func run(ctx context.Context, patterns []string, in io.Reader, out io.Writer) error {
	request := &packages.DriverRequest{}
	if err := json.NewDecoder(in).Decode(request); err != nil {
		return fmt.Errorf(`unable to unmarshal request: %w`, err)
	}

	response, err := driver.New(ctx, patterns, request).Drive()
	if err != nil {
		return fmt.Errorf(`error running driver: %w`, err)
	}

	if err := json.NewEncoder(out).Encode(response); err != nil {
		return fmt.Errorf(`unable to marshal response: %w`, err)
	}
	return nil
}

func withInterrupt(parentCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parentCtx)
	ch := make(chan os.Signal, 1)
	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			signal.Stop(ch)
		}
	}()
	signal.Notify(ch, os.Interrupt)
	return ctx, cancel
}
