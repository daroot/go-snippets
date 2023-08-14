package slogext_test

import (
	"errors"
	"importfromprojectlocally/slogext"
	"testing"
)

func TestSlogError(t *testing.T) {
	tl := newTestLogger(t)
	testerr := errors.New("my test error")
	tl.logger.Error("something failed", slogext.Error(testerr))
	tl.HasLogged(`error="` + testerr.Error() + `"`)
}
