package slogext

import "log/slog"

// Error is a convenience function creating an slog.Attr
// with a fixed "error" key, and the given error value,
// rather than having to repeatedly use slog.Any("error", error)
func Error(err error) slog.Attr {
	return slog.Attr{Key: "error", Value: slog.AnyValue(err)}
}
