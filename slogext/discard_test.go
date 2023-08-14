package slogext_test

import (
	"context"
	"importfromprojectlocally/slogext"
	"log/slog"
	"testing"

	"github.com/matryer/is"
)

func TestDiscardHandler(t *testing.T) {
	is := is.New(t)

	h := slogext.DiscardHandler{}
	ctx := context.Background()

	is.Equal(h.Enabled(ctx, slog.LevelInfo), false)                 // enabled is never true
	is.Equal(h.Handle(ctx, slog.Record{}), nil)                     // handle never fails
	is.Equal(h.WithAttrs([]slog.Attr{slog.Any("thing", "foo")}), h) // does not change internals
	is.Equal(h.WithGroup("blah"), h)                                // does not change internals
}
