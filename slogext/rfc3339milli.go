package slogext

import (
	"log/slog"
)

// RFC3339Millis is a helper function for use in slog.HandlerOptions.ReplaceAttr.
// It replaces the default formatting of the record's time
// with a millisecond precision RFC3339 timestamp.
// The default JSON formatting includes down to nanosecond precision,
// which is often unnecessarily verbose.
// The TextHandler defaults to the same output as RFC3339Milli,
// but no HandlerOption is exposed to select this for JSON.
func RFC3339Millis(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{
			Key: slog.TimeKey,
			// This could use the same code as the internal slog.writeRFC3339Millis
			// for performance without extra memory allocations.
			Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.999Z07:00")),
		}
	}
	return a
}
