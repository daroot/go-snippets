package httptools_test

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"

	"importfromprojectlocally/httptools"
	"importfromprojectlocally/slogext"
)

// TestDoesServe
// TestSlowShutdown

func testHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"I'm doing science and I'm still alive."}`))
}

func TestGetsContextCancel(t *testing.T) {
	is := is.New(t)
	logbuf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(logbuf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mux := http.HandlerFunc(testHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = slogext.Add(ctx, logger)

	done := make(chan struct{}, 1)
	started := make(chan struct{}, 1)

	go func() {
		close(started)                         // indicate the Serve go routine got scheduled and is now running
		is.NoErr(httptools.Serve(ctx, 0, mux)) // Serve shuts down cleanly, context.Canceled did not propagate
		close(done)                            // now finished
	}()

	go func() {
		// run in a go routine to wait for the started signal
		// without interfering with the overall test timeout.
		<-started // wait for Serve to be running
		cancel()  // Do the thing; ensure it propagates
	}()

	select {
	case <-time.After(time.Millisecond * 50):
		t.Fatal("timed out waiting for server shutdown")
	case <-done:
	}

	logtext := logbuf.String()
	is.True(strings.Contains(logtext, "HTTP service starting"))                // successfully started
	is.True(strings.Contains(logtext, "HTTP service shutting down on cancel")) // got a cancel signal
	is.True(strings.Contains(logtext, "HTTP service stopped"))                 // succesfully stopped
}

func TestListenAndServePropagatesError(t *testing.T) {
	is := is.New(t)
	logbuf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(logbuf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mux := http.HandlerFunc(testHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = slogext.Add(ctx, logger)

	done := make(chan struct{}, 1)
	started := make(chan struct{}, 1)
	usedport := 0

	go func() {
		listener, err := net.Listen("tcp", ":0")
		is.NoErr(err)
		usedport = listener.Addr().(*net.TCPAddr).Port
		close(started) // signal that listener has started
		<-done         // wait until after attempting to start Serve
		listener.Close()
	}()

	go func() {
		<-started // wait for the listener routine to schedule and open the port.
		err := httptools.Serve(ctx, usedport, mux)
		t.Log(err)
		close(done)
		is.True(err != nil) // ports conflicted
	}()

	select {
	case <-time.After(time.Millisecond * 50):
		t.Fatal("timed out waiting for server shutdown")
	case <-done:
	}

	is.True(strings.Contains(logbuf.String(), "http.Server returned abnormally, stopping app"))
	is.True(strings.Contains(logbuf.String(), "bind: address already in use"))
}
