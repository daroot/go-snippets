// Package consterr provides const-compatible string errors
package consterr

// Err allows compile time constant defined error values,
// rather than needing any var that runs at init time.
//
// Err does not implement Unwrap, Is or As as used by [errors];
// it is for "fundamental" error types.
//
//nolint:errname // Err _is_ the type, not XXXError or similar.
type Err string

// Error satistfies the error interface.
func (e Err) Error() string { return string(e) }
