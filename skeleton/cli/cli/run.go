package myappcli

import (
	"context"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3/ffcli"

	getcmd "myapp/cli/get"
	putcmd "myapp/cli/put"
	"myapp/cli/root"
	versioncmd "myapp/cli/version"
)

// Run executes the app.
// It satisifies the interface expected by cmd/main.go
// It sets up config from any arguments and environment.
// Then it executes the application using the provided io streams.
func Run(ctx context.Context, args, env []string, stdin io.Reader, stdout, stderr io.Writer) error {
	rootCmd, cfg := root.New(stdin, stdout, stderr)

	rootCmd.Subcommands = []*ffcli.Command{
		getcmd.New(cfg),
		putcmd.New(cfg),
		versioncmd.New(cfg),
	}

	if err := rootCmd.Parse(args[1:]); err != nil {
		return fmt.Errorf("parsing cli: %w", err)
	}

	// Do all global initialization here.
	// clients for apis, loggers, etc
	ctx = cfg.WithLoggerCtx(ctx)

	return rootCmd.Run(ctx)
}
