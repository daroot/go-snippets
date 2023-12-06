package httptools

import (
	"context"
	"errors"
	"fmt"
	golog "log"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"importfromprojectlocally/slogext"
)

// adapter to write stdlib log via log/slog as errors, for use in http.Server reporting.
type slogErrWriter struct {
	*slog.Logger
}

func (sw slogErrWriter) Write(p []byte) (int, error) {
	sw.Error(string(p))
	return 0, nil
}

// Serve an http.Handler (which may be a mux like http.ServeMux, chi.Router, gorilla.Mux, etc)
// on a given port with reasonable defaults.
// Run until the supplied context is canceled, then try to shutdown gracefully.
// Return the http.Server or shutdown error if unable to start or exit gracefully.
// Will return nil, not http.ErrServerClosed, if the shutdown is graceful.
func Serve(ctx context.Context, port int, mux http.Handler) error {
	log := slogext.From(ctx).With("component", "http", "port", port)

	// The ancestor ctx chain should have a signal.NotifyContext() in it,
	// so when it cancels on a signal, everything downstream does too.
	// localcancel is used to ensure http.Server.Shutdown runs
	// in the event that ListenAndServe exits for any reason other than
	// the global signal cancellation, such as a port conflict.
	ctx, localcancel := context.WithCancel(ctx)

	//nolint:gomnd // the whole point of this is to encode magic defaults
	srv := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      30 * time.Second, // Absolute maximum possible request response time.
		IdleTimeout:       300 * time.Second,
		ErrorLog:          golog.New(slogErrWriter{log.With(slog.String("component", "http.Server"))}, "", 0),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx // all requests inherit from the global context
		},
	}

	log.Info("HTTP service starting")

	wg := &sync.WaitGroup{}
	var (
		srverr  error
		shuterr error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if srverr = srv.ListenAndServe(); srverr != nil {
			if !errors.Is(srverr, http.ErrServerClosed) {
				log.Error("http.Server returned abnormally, stopping app", slogext.Error(srverr))
				localcancel() // shutdown waiter will not block forever.
				return
			}
			log.Debug("closing http.Server on clean shutdown")
			srverr = nil // ErrServerClosed is not an error for the caller.
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for the context to cancel (via the local or upstream cancel())
		<-ctx.Done()
		log.Info("HTTP service shutting down on cancel")
		// Needs a clean context so it's not pre-canceled during shutdown.
		//nolint:gomnd // again, deliberately encode magic default
		shutctx, shutcancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer shutcancel()
		shuterr = srv.Shutdown(shutctx)
	}()

	wg.Wait()

	if srverr != nil {
		return srverr
	}
	if shuterr != nil {
		log.Error("during shutdown", slogext.Error(shuterr))
	}

	log.Info("HTTP service stopped")
	return shuterr
}
