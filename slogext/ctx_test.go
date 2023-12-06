package slogext_test

import (
	"bytes"
	"context"
	"importfromprojectlocally/slogext"
	"log/slog"
	"reflect"
	"regexp"
	"testing"
)

const timeRegexp = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}(Z|[+-]\d{2}:\d{2})`

type testLogger struct {
	t      *testing.T
	logger *slog.Logger
	buf    *bytes.Buffer
}

func (tl testLogger) HasLogged(re string) {
	tl.t.Helper()
	matched, err := regexp.MatchString(re, tl.buf.String())
	if err != nil {
		tl.t.Fatal(err)
	}
	if !matched {
		tl.t.Errorf("did not get expected log output %s, got: %s", re, tl.buf.String())
	}
}

func newTestLogger(t *testing.T) *testLogger {
	buf := &bytes.Buffer{}
	return &testLogger{
		t:      t,
		buf:    buf,
		logger: slog.New(slog.NewTextHandler(buf, nil)),
	}
}

func TestEmptyCtx(t *testing.T) {
	log := slogext.From(context.Background())
	handler := log.Handler()

	if _, ok := handler.(slogext.DiscardHandler); !ok {
		t.Error("From did not return the expected logger")
	}
}

func TestFrom(t *testing.T) {
	tl := newTestLogger(t)
	ctx := slogext.Add(context.Background(), tl.logger)

	log2 := slogext.From(ctx)
	if !reflect.DeepEqual(tl.logger, log2) {
		t.Error("Add/From cycle did not return the expected logger")
	}

	log2.Info("did it")
	tl.HasLogged(`time=` + timeRegexp + ` level=INFO msg="did it"`)
}

func TestDefault(t *testing.T) {
	tl := newTestLogger(t)

	// restore defaults when done to ensure any other tests don't get surprise default loggers.
	prevDefault := slogext.GetContextDefault()
	defer func() {
		slogext.SetContextDefault(prevDefault)
	}()

	slogext.SetContextDefault(tl.logger)

	log2 := slogext.From(context.Background())
	if !reflect.DeepEqual(tl.logger, log2) {
		t.Error("From empty context did not return the SetDefault logger")
	}

	log2.Warn("something", "count", 3)
	tl.HasLogged(`time=` + timeRegexp + ` level=WARN msg=something count=3`)
}
