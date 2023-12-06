package myapp

import (
	"flag"
	"fmt"
	"log/slog"

	"myapp/consterr"
	"myapp/envflag"
)

const (
	// ErrNoArgs is returned if there is no command line args;
	// even if no args are on the command line, Args[0] is the command name.
	ErrNoArgs = consterr.Err("missing args, pass in os.Args")
)

// Config contains the service configuration as parsed from CLI flags, env
// variables, and embedded build information.
type Config struct {
	AppName  string
	LogLevel *slog.LevelVar
	Port     int

	BuildInfo BuildInfo `json:"Build"`
}

// NewConfig creates an application Config from command line flags and environment variables.
func NewConfig(args []string, env []string) (*Config, error) {
	if len(args) == 0 {
		return nil, ErrNoArgs
	}

	c := Config{
		AppName:   args[0],
		BuildInfo: getBuildInfo(),
	}

	// ContinueOnError as to never panic or os.Exit() except at the top level.
	fs := flag.NewFlagSet(c.AppName, flag.ContinueOnError)
	fs.TextVar(c.LogLevel, "log-level", &slog.LevelVar{}, "logging level (debug, info, warn, error)")
	fs.IntVar(&c.Port, "port", 8000, "network port to listen on")
	// Add other fields here

	// envflag wraps fs.Parse to also pull from equiv ENV variables if no cli arg is set
	if err := envflag.Parse(fs, args[1:], env); err != nil {
		return nil, fmt.Errorf("parsing config from environment: %w", err)
	}

	return &c, nil
}
