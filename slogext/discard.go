package slogext

import (
	"context"
	"log/slog"
)

// DiscardHandler is an slog.Handler which simply drops all prospective output.
type DiscardHandler struct{}

// Enabled reports it is not enabled at any level.
func (d DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool { return false }

// Handle handles the Record by doing nothing.
func (d DiscardHandler) Handle(_ context.Context, _ slog.Record) error { return nil }

// WithAttrs returns the handler with no updates.
func (d DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler { return d }

// WithGroup returns the handler with no updates.
func (d DiscardHandler) WithGroup(_ string) slog.Handler { return d }
