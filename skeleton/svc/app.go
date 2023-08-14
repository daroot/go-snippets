// Package myapp contains the business application logic for myapp
package myapp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"myapp/slogext"
)

// app wraps all our configuration and instantiated state.
type myapp struct {
	cfg    *Config
	mux    *http.ServeMux
	things int
	stuff  string
}

// Run our app
func (a myapp) Run(ctx context.Context) error {
	log := slogext.From(ctx)
	log = log.With(slog.String("app", "myapp"), slog.Int("port", a.cfg.Port))

	log.Info("starting")

	log.Info("doing arbitrary things until shutdown signaled",
		slog.Int("things", a.things),
		slog.String("stuff", a.stuff),
	)

	// go func() { err := httptools.Serve(ctx, app.cfg.Port, app.mux) }

	select {
	case <-ctx.Done():
		log.Info("got shutdown signal")
		if err := ctx.Err(); !errors.Is(err, context.Canceled) {
			// context.Canceled isn't a "real" error, but anything else is.
			return err
		}
	case <-time.After(time.Second * 30):
		log.Info("ran for some time and finished my work")
	}

	log.Info("doing cleanup things")

	log.Info("exiting")
	return nil
}

// Create our app from a config
// It takes and returns a context
// so that it can choose to set up an initial context for Run(),
// and deal with any signal canceling/shutdown logic
func newApp(ctx context.Context, output io.Writer, cfg *Config) (context.Context, *myapp) {
	app := &myapp{things: 1, stuff: "foo"}

	log := slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	log.Debug("debug logging on.")

	log.Debug("setting default context logger")
	slogext.SetContextDefault(log)

	log.Debug("doing all the app setup things we should do")
	log.Debug("creating HTTP router")
	app.mux = http.NewServeMux()
	log.Debug("loading key and secrets")
	log.Debug("setting up network clients")
	log.Debug("setting up 3p SDKs")
	log.Debug("etc")

	log.Info("app setup complete", slog.Any("config", cfg))
	return ctx, app
}

// Run executes this app.
// It satisifies the interface expected by cmd/main.go
// It sets up the config from any arguments and environment,
// then executes the application logging to the provided output writer.
// args should be os.Args
// env should be os.Environ()
func Run(ctx context.Context, output io.Writer, args, env []string) error {
	cfg, err := NewConfig(args, env)
	if err != nil {
		return fmt.Errorf("couldn't load config: %w", err)
	}

	// newApp returns a configured app and the (modified if desired) context.
	ctx, app := newApp(ctx, output, cfg)

	return app.Run(ctx)
}
