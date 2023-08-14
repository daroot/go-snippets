package main

// This file is boilerplate.
// The only thing you should need to edit is the "app" import line.

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	// The cli package must offer:
	//   Run(ctx context.Context, args, env []string, in io.Reader, out, err io.Writer) error
	cli "myapp/cli"
)

func main() {
	// These things change "global" module level state.
	// Do them here in main(), not in any library or setup code,
	// to avoid any possible inconsistency in tests and the like.
	// Avoid doing anything else.  Any code here is very hard to test reliably.
	//
	// Ensure we always use UTC for consistency.
	time.Local = time.UTC
	// Try to catch any logging from downstream libraries using either log.Print
	// or module level slog.Info etc.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// if we call os.Exit(1) directly from main() our signal defer will not run.
	// this pattern ensures we both run our signal stop() and exit with an
	// error code if necessary.
	var exitCode int
	defer func() { os.Exit(exitCode) }()

	// Create a global context root and when we get a terminating signal,
	// that should cancel our whole context tree.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Now run our app with global state set and signal handling initialized.
	if err := cli.Run(ctx, os.Args, os.Environ(), os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		exitCode = 1
		return
	}
}
