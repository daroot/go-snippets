package consterr_test

import (
	"importfromprojectlocally/consterr"
	"testing"
)

const testerr = consterr.Err("constant error is constant")

func TestConstErr(t *testing.T) {
	// Really, the test is that consterr.Err can be compiled as a const above.
	if testerr.Error() != "constant error is constant" {
		t.Fail()
	}
}
