package versioncmd

import (
	"context"
	"flag"
	"log/slog"

	"myapp/cli/root"
	"myapp/slogext"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	root *root.Config
}

func New(rcfg *root.Config) *ffcli.Command {
	cfg := Config{root: rcfg}

	fs := flag.NewFlagSet("example version", flag.ContinueOnError)
	root.Config.CommonFlags(fs)

	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "version",
		ShortHelp:  "get version information",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}
}

func (c *Config) Exec(ctx context.Context, args []string) error {
	log := slogext.From(ctx)

	log.Info("myapp version information", slog.Any("buildinfo", c.root.BuildInfo))

	return nil
}
