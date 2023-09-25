// Package testgolden contains [test.Testing] helpers
// for comparing a current test result
// against a previously saved "golden" result,
//
// If the path to any golden file
// does not contain a `testdata` directory segment,
// then `testdata` is prepended to the initial path.
// to ensure golden files do not interfere with the normal go toolchain,
// per `go help packages` and the [go package docs].
//
// A normal workflow is as follows:
//   - Write the initial tests using Compare,
//     which will result in a failing test, as no golden file exists.
//   - Run the specific test using `go test -update-golden -run <regex matching test name>`,
//     to create a new golden file.
//   - Validate that the contents of the golden file match expectations.
//   - Commit the new golden file to source control.
//   - When future code changes what the expected output should be,
//     the existing test calling Compare will fail.
//   - Re-run the specific failing test with `-update-golden` again,
//   - Use a VCS diff on the golden file to inspect the changes
//     and verify they are expected.
//   - Commit the updated golden file.
//
// [go package docs]: https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns
package testgolden

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// doUpdate provides the '-update-golden' flag for `go test`
//
//nolint:gochecknoglobals // this must be package global for flag.Parse to work correctly
var doingUpdate = flag.Bool("update-golden", false, "update test cases that use testgolden.Compare")

// DoingUpdate returns if go tests were called with the `-update-golden` flag.
// Useful for writing helpers which use a Compare style workflow,
// but offer their own diff/comparison method,
// such as the companion testgoldenproto package.
func DoingUpdate() bool {
	return *doingUpdate
}

// DoUpdate sets the state of whether Compare will write goldenfiles
// prior to testing for equality.
func DoUpdate(newstate bool) {
	*doingUpdate = newstate
}

// ensure path p comes only from within a testdata/ subdir,
// aND IF one is not already in p, use the local testdata/ subdir.
func normalizeTestPath(p string) string {
	if !strings.Contains(p, "testdata") {
		return filepath.Join("testdata", p)
	}
	return p
}

// make a testdata dir for path p if it does not exist,
// using filepath.Dir in case the user asked for a path
// in a different testdata dir
// that is relative to where the current test is running.
func ensureTestDirForPath(p string) error {
	testDir := filepath.Dir(p)
	_, err := os.Stat(testDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(testDir, 0o0750)
	}
	return err
}

// Load fetches bytes previously saved to goldfile.
// Load performs t.Fail (via Errorf) if goldfile is not present,
// but not t.FailNow, in case the test needs to clean up resources.
func Load(t *testing.T, goldfile string) []byte {
	t.Helper()

	goldfile = filepath.Clean(normalizeTestPath(goldfile))

	expected, err := os.ReadFile(goldfile)
	if err != nil {
		t.Errorf("unable to load golden file %s, %v", goldfile, err)
		return []byte{}
	}
	return expected
}

// Save stores expected bytes to compare for future uses of Load.
// If the goldfile cannot be saved (or a testdata dir not created),
// the test will be t.FailNow()'d (via t.Fatalf),
// to try and minimize any other tests doing updates.
func Save(t *testing.T, goldfile string, expected []byte) {
	t.Helper()

	goldfile = filepath.Clean(normalizeTestPath(goldfile))

	if err := ensureTestDirForPath(goldfile); err != nil {
		t.Fatalf("unable to ensure testdata directory %v", err)
	}

	err := os.WriteFile(goldfile, expected, 0o0640) //#nosec G306 -- group read is acceptable
	if err != nil {
		t.Fatalf("unable to write golden file %s, %v", goldfile, err)
	}
}

// Compare is a [testing.T] wrapper
// for the workflow of taking an actual result
// optionally writing an that value as an updated expected result,
// and then comparing actual against the expected results.
//
// actual must be a type that is (un)marshalable via encoding/json
//
// If -update-golden flag is set,
// the expected goldfile is saved with the current actual output
// and Compare should always return true.
// If -update-golden flag is not used,
// then Compare will diff the actual vs expected,
// failing the test if the results are not identical.
//
// Compare returns a bool indicating if the outputs matched successfully,
// for the case of using wrappers like `is` or `testify`,
// or where the test needs to optionally do additional verifications or cleanup.
//
// The optional cmpopts can be used to pass [cmp.Diff] any set of [cmp.Comparer] or
// [cmp.Transform] functions necessary to compare the actual to expected result.
func Compare[T any](t *testing.T, testname string, goldfile string, actual T, cmpopts ...cmp.Option) bool {
	t.Helper()

	// if our update-golden flag is set,
	// first we save the actual output as the new expected,
	// which should additionallly ensure
	// that expected will match actual for this run.
	if DoingUpdate() {
		actualbytes, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			t.Fatalf("%v: failed to json encode actual (type %T): %v", testname, actual, err)
		}

		Save(t, goldfile, actualbytes)
	}

	expectedbytes := Load(t, goldfile)

	var expected T
	if err := json.Unmarshal(expectedbytes, &expected); err != nil {
		t.Errorf("%v: unable to recreate expected value from json bytes: %v", testname, err)
		return false
	}

	if diff := cmp.Diff(expected, actual, cmpopts...); diff != "" {
		t.Errorf("%v: actual (type %T) differed from expected (type %T) : %v", testname, actual, expected, diff)
		return false
	}
	return true
}
