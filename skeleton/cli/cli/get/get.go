package getcmd

import (
	"context"
	"flag"

	"myapp/cli/root"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	root  *root.Config
	Count int
	Tail  bool
}

func New(rcfg *root.Config) *ffcli.Command {
	cfg := Config{root: rcfg}

	fs := flag.NewFlagSet("example get", flag.ContinueOnError)
	fs.BoolVar(&cfg.Tail, "tail", false, "do not exit, emit new items as they arrive")
	fs.IntVar(&cfg.Count, "count", 0, "how many existing items to retrieve")
	root.Config.CommonFlags(fs)

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "get [flags] <stream>",
		ShortHelp:  "retrieve items from an example api",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}
}

func (c *Config) Exec(ctx context.Context, args []string) error {
	// do thing here
	return nil
}
