package testgolden_test

import (
	"errors"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"

	"importfromprojectlocally/testgolden"
)

type simpleStruct struct {
	Int    int
	Float  float64
	String string
}

const expectedSimpleGoldenBytes = `{
  "Int": 1,
  "Float": 3.141592653589793,
  "String": "Foo"
}`

// testcase is a generic wrapper for a test case with an arbitrary value type.
type testcase[T any] struct {
	value    T
	expected []byte
	validate func(t *testing.T, path string, expect []byte, opts ...cmp.Option) error
	cmpopts  []cmp.Option
}

func (tc testcase[T]) testCompare(t *testing.T, name, path string) bool {
	return testgolden.Compare(t, name, path, tc.value, tc.cmpopts...)
}

func (tc testcase[T]) validateGolden(t *testing.T, path string) error {
	return tc.validate(t, path, tc.expected, tc.cmpopts...)
}

// because Go can't do map[string]testcase without an instantiated type,
// we need to use an interface allow all varieties of testcase[x] to live in one map.
type compareTester interface {
	testCompare(*testing.T, string, string) bool
	validateGolden(*testing.T, string) error
}

func TestCompare(t *testing.T) {
	t.Parallel()

	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		if mkdErr := os.MkdirAll("testdata", 0o0750); mkdErr != nil {
			t.Fatal("unable to make testdata directory", err)
		}
	} else if err != nil {
		t.Fatal("unable to stat testdata directory", err)
	}

	testCases := map[string]compareTester{
		// Basic single json int value
		"int": testcase[int]{
			1,
			[]byte(`1`),
			matchGoldenContents,
			nil,
		},

		// Basic single json string value
		"string": testcase[string]{
			"foo",
			[]byte(`"foo"`),
			matchGoldenContents,
			nil,
		},

		// Struct with simple data types
		"simple struct": testcase[simpleStruct]{
			simpleStruct{1, math.Pi, "Foo"},
			[]byte(expectedSimpleGoldenBytes),
			matchGoldenContents,
			nil,
		},

		// Use a comparer that doesn't inspect one field.
		// golden file contents are different,
		// but Compare must pass
		"struct with comparer": testcase[simpleStruct]{
			simpleStruct{1, math.E, "Foo"},
			[]byte(expectedSimpleGoldenBytes),
			expectGoldenDiff,
			[]cmp.Option{cmp.Comparer(func(a, b simpleStruct) bool {
				return a.Int == b.Int && a.String == b.String
			})},
		},

		// Things to future test:
		// structs with unexported fields
		// complex structs
		// unwritable testdata dir
		// diff in golden returns failure
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			// Don't use t.TempDir(), because we explicitly need to be in testdata/ dir
			file, err := os.CreateTemp("testdata", "sample-*.json")
			is.NoErr(err) // create temporary golden file
			path := file.Name()
			defer os.Remove(path)

			testgolden.DoUpdate(true)
			result := tc.testCompare(t, name, path)
			is.True(result) // compare with Update to write golden

			is.NoErr(tc.validateGolden(t, path)) // golden file contents

			testgolden.DoUpdate(false)
			result = tc.testCompare(t, name+" recompare", path)
			is.True(result) // re-compare with no update
		})
	}
}

func matchGoldenContents(t *testing.T, path string, expected []byte, opts ...cmp.Option) error {
	t.Helper()
	loadedbytes := testgolden.Load(t, path)
	if diff := cmp.Diff(expected, loadedbytes, opts...); diff != "" {
		return fmt.Errorf("golden file does not match expectation: %v", diff)
	}
	return nil
}

func expectGoldenDiff(t *testing.T, path string, expected []byte, opts ...cmp.Option) error {
	t.Helper()
	loadedbytes := testgolden.Load(t, path)
	if diff := cmp.Diff(expected, loadedbytes, opts...); diff == "" {
		return errors.New("golden file is unexpectedly the same")
	}
	return nil
}
