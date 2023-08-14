package slogext

import (
	"context"
	"log/slog"
)

//nolint:gochecknoglobals // deliberate package state
var (
	defaultContextLogger *slog.Logger
)

type ctxKey struct{}

// From returns an *slog.Logger associated to a context.
// If the context does not contain an existing logger,
// if slogext.SetContextDefault has been called, returns that logger,
// otherwise return an *slog.Logger which silently Discards all output.
//
// slogext.From(ctx).Info("did a thing", slog.Int("count", 3))
func From(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	} else if l = defaultContextLogger; l != nil {
		return l
	}
	return slog.New(DiscardHandler{})
}

// Add creates a new context from an existing one with an *slog.Logger associated.
//
// logger = slog.New(...)
// logctx = slogext.Add(r.Context(), logger.With("request-id", "1234"))
// doMyfunc(logctx, ...)
func Add(ctx context.Context, l *slog.Logger) context.Context {
	// zerolog ensures it does not store a disabled logger.  Do we care?
	return context.WithValue(ctx, ctxKey{}, l)
}

// GetContextDefault returns the *slog.Logger set by SetContextDefault,
// or nil if that has not been called.
func GetContextDefault() *slog.Logger {
	return defaultContextLogger
}

// SetContextDefault makes l the *slog.Logger returned by From when
// a context has no previously associated logger.
//
//	slogext.SetContextDefault(mylogger)
//	if logger := slogext.From(context.Background()); logger == mylogger {
//	  logger.Info("successfully updated default context logger")
//	}
func SetContextDefault(l *slog.Logger) {
	defaultContextLogger = l
}
