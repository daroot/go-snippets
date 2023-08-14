package root

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"

	"myapp/slogext"
)

type Config struct {
	LogLevel *slog.LevelVar

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	BuildInfo BuildInfo `json:"Build"`
	// Client *some.Client
}

// CommonFlags adds the flags that are common to all commands to a flagset.
func (cfg *Config) CommonFlags(fs *flag.FlagSet) {
	fs.TextVar(cfg.LogLevel, "level", &slog.LevelVar{}, "logging level")
}

func (cfg *Config) WithLoggerCtx(ctx context.Context) context.Context {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	log.Debug("debug logging on.")

	log.Debug("initializing root context with slog")
	return slogext.Add(ctx, log)
}

func New(in io.Reader, out io.Writer, err io.Writer) (*ffcli.Command, *Config) {
	cfg := Config{
		Stdin:     in,
		Stdout:    out,
		Stderr:    err,
		BuildInfo: getBuildInfo(),
	}
	fs := flag.NewFlagSet("myapp", flag.ContinueOnError)
	cfg.CommonFlags(fs)

	return &ffcli.Command{
		ShortUsage: "myapp [flags] <subcommand>",
		ShortHelp:  "do myapp things",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}, &cfg
}
