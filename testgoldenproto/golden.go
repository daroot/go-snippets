// Package testgoldenproto contain a test.Testing helper
// for comparing google protobuf-based objects to known good ("golden") outputs.
//
// It is separated from testgolden
// to isolate the numerous protobuf library dependencies.
// for projects which are not proto based.
package testgoldenproto

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	testgolden "importfromprojectlocally/testgolden"
)

// Compare is a shorthand testing.T wrapper
// around testgolden.Save and Load,
// for taking a test output proto.Message
// and either saving into or comparing with a known golden file.
//
// As with testgolden.Compare,
// if go test is run with -update-golden,
// the goldfile is first saved with the actual output,
// ensuring the test will pass comparison.
// if -update-golden flag is not passed,
// then if the actual output varies
// from the expect contents of goldfile,
// the test case will fail.
//
// Compare uses protojson for making golden files,
// so humans can look at git diff after running an `-update-golden`,
// and verify that the changes match expectations.
//
// As a result, proto objects and cmp.Diff with protocmp.Transform() are used,
// because protojson.Marshal output is deliberately not stable,
// and the raw json bytes from protojson can **never** be compared.
// See https://github.com/golang/protobuf/issues/1121 for details.
func Compare[T proto.Message](t *testing.T, testname string, goldfile string, actual T) bool {
	t.Helper()

	// When updating, first Save the new golden file;
	// this overwrites any previous golden file and
	// should ensure no changes in cmp.Diff below.
	if testgolden.DoingUpdate() {
		jbytes, err := protojson.MarshalOptions{UseProtoNames: true, Indent: "  "}.Marshal(actual)
		if err != nil {
			t.Errorf("unable to marshal got: %v", err)
			return false
		}
		testgolden.Save(t, goldfile, jbytes)
	}

	var expected T
	jbytes := testgolden.Load(t, goldfile)
	if err := protojson.Unmarshal(jbytes, expected); err != nil {
		t.Errorf("unable to unmarshal golden proto object: %v", err)
		return false
	}

	// do actual diff with protocmp.Transform.
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		t.Errorf("%v: actual differed from expected: %v", testname, diff)
		return false
	}
	return true
}
