package slogxray

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/matryer/is"
)

type testString string

func (s testString) String() string { return string(s) }

func TestAdapter(t *testing.T) {
	testCases := map[string]struct {
		slevel   slog.Level
		xlevel   xraylog.LogLevel
		expected string
	}{
		"both debug": {
			slog.LevelDebug,
			xraylog.LogLevelDebug,
			`{"level":"DEBUG","msg":"blah"}`,
		},
		"both info": {
			slog.LevelInfo,
			xraylog.LogLevelInfo,
			`{"level":"INFO","msg":"blah"}`,
		},
		"both warn": {
			slog.LevelWarn,
			xraylog.LogLevelWarn,
			`{"level":"WARN","msg":"blah"}`,
		},
		"both error": {
			slog.LevelError,
			xraylog.LogLevelError,
			`{"level":"ERROR","msg":"blah"}`,
		},
		"slog greater": {
			slog.LevelWarn,
			xraylog.LogLevelInfo,
			``,
		},
		"xray greater": {
			slog.LevelDebug,
			xraylog.LogLevelInfo,
			`{"level":"INFO","msg":"blah"}`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			buf := &bytes.Buffer{}
			log := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{
				Level: tc.slevel,
				// remove time from output so results are consistent test to test.
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return a
				},
			}))
			a := New(log)

			a.Log(tc.xlevel, testString("blah"))
			t.Log("got", strings.TrimSpace(buf.String()), "want", tc.expected)
			is.Equal(strings.TrimSpace(buf.String()), tc.expected)
		})
	}
}
