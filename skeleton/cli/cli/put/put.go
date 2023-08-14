package putcmd

import (
	"context"
	"flag"

	"myapp/cli/root"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	root *root.Config
}

func New(rcfg *root.Config) *ffcli.Command {
	cfg := Config{root: rcfg}

	fs := flag.NewFlagSet("example put", flag.ContinueOnError)
	cfg.root.CommonFlags(fs)

	return &ffcli.Command{
		Name:       "put",
		ShortUsage: "put [flags] <destination>",
		ShortHelp:  "send items to an example api",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}
}

func (c *Config) Exec(ctx context.Context, args []string) error {
	// use cfg.root.Client to do thing here
	return nil
}
