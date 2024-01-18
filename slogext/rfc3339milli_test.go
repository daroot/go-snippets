package slogext_test

import (
	"bytes"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"

	"importfromprojectlocally/slogext"

	"github.com/matryer/is"
)

const rfc3339MilliRegexp = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,3})?(Z|[+-]\d{2}:\d{2})`

func TestRFC3339Millis(t *testing.T) {
	testCases := map[string]struct {
		date     time.Time
		expected string
	}{
		"base": {
			time.Date(2023, 9, 16, 11, 37, 23, 42000000, time.UTC),
			"2023-09-16T11:37:23.042Z",
		},
		"milli value is zero padded": {
			time.Date(2023, 9, 16, 11, 37, 23, 0, time.UTC),
			"2023-09-16T11:37:23.000Z",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			attr := slog.Attr{Key: slog.TimeKey, Value: slog.TimeValue(tc.date)}
			a := slogext.RFC3339Millis([]string{}, attr)
			actual := strings.TrimSpace(a.Value.String())
			is.Equal(tc.expected, actual) // time value matches expected
		})
	}
}

func TestRFC3339MillisJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	log := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{ReplaceAttr: slogext.RFC3339Millis}))
	log.Info("test")
	if ok, err := regexp.Match(`"time":"`+rfc3339MilliRegexp+`"`, buf.Bytes()); !ok || err != nil {
		t.Log(buf.String())
		t.Error("time field format didn't match", rfc3339MilliRegexp, err)
	}
}
