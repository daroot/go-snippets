// Package testbuffer contains a syncronized buffer for using in tests that
// write to an io.Writer-based log (zerolog, slog, zap, etc) but need to handle
// t.Parallel() subtests or which use goroutines to write to the log.
package testbuffer

import (
	"bytes"
	"strings"
	"sync"
)

// LogBuf is a syncronized io.Writer.
// It's meant to prevent any race detection on loggers in background go routines
// if you'd get using a raw bytes.Buffer.
//
// Typical usage pattern with zerolog looks something like:
//
//	logbuf := testbuffer.New()
//	log := zerolog.New(logbuf)
//	<do test things that write to log>
//	if !logbuf.Contains("some test value") { t.Error("missing expected log") }
type LogBuf struct {
	m *sync.RWMutex
	*bytes.Buffer
}

// New returns an initialized log buffer.
func New() *LogBuf {
	ml := LogBuf{}
	ml.Buffer = &bytes.Buffer{}
	ml.m = &sync.RWMutex{}
	return &ml
}

// Write satisfies io.Writer.
func (ml *LogBuf) Write(p []byte) (int, error) {
	ml.m.Lock()
	defer ml.m.Unlock()
	return ml.Buffer.Write(p)
}

// String satisfies fmt.Stringer
func (ml *LogBuf) String() string {
	ml.m.RLock()
	defer ml.m.RUnlock()
	return ml.Buffer.String()
}

// Reset resets the log buffer.
func (ml *LogBuf) Reset() {
	ml.m.Lock()
	defer ml.m.Unlock()
	ml.Buffer.Reset()
}

// Contains is a handy shortcut for checking if the string s is in the LogBuf
func (ml *LogBuf) Contains(s string) bool {
	ml.m.RLock()
	defer ml.m.RUnlock()
	return strings.Contains(ml.String(), s)
}
