// Package main runs an application
package main

// This file is boilerplate.
// The only thing that should need editing is the "app" import line.

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	// The app package must offer:
	// func Run(ctx context.Context, args, env []string, input io.ReadCloser, output, errout io.WriteCloser) error
	app "myapp/cli"
)

func main() {
	// Do things that change module or process level state in main(),
	// rather than any library or init() code,
	// to minimize inconsistency in tests and the like.
	// Avoid doing anything else,
	// as any code here is hard to test reliably.

	// Ensure UTC for timestamp consistency.
	time.Local = time.UTC
	// Try to catch any logging from downstream libraries using either log.Print
	// or module level slog.Info etc.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// calling os.Exit(1) directly from main() prevents the signal defer from running.
	// Ensure that signal stop() runs, but os.Exit still gets the appropriate return code.
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	// Create a global context root that cancels every descended context when a
	// termination signal is received.
	sigctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(sigctx, os.Args, os.Environ(), os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		exitCode = 1
	}
}
