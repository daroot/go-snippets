package envflag_test

import (
	"bytes"
	"flag"
	"testing"

	"github.com/matryer/is"

	"importfromprojectlocally/envflag"
)

func TestParse(t *testing.T) {
	t.Parallel()

	type tEnv []string
	type tArgs []string
	type target struct {
		Foo string
		Baz string
		Wuz string
	}

	testCases := map[string]struct {
		env       tEnv
		args      tArgs
		result    target
		shoulderr bool
	}{
		"basic": {
			tEnv{"FOO=bar", "WUZ_VALUE=a bear"},
			tArgs{"-baz", "quux"},
			target{"bar", "quux", "a bear"},
			false,
		},
		"no env value": {
			tEnv{"FOO="},
			tArgs{"-baz", "quux", "-wuz-value", "a bear"},
			target{"", "quux", "a bear"},
			false,
		},
		"no flag value": {
			tEnv{"FOO=bar"},
			tArgs{"-baz"},
			target{},
			true, // missing -baz argument
		},
		"flag overrides env": {
			tEnv{"FOO=bar"},
			tArgs{"-foo", "buz"},
			target{Foo: "buz"}, // not bar
			false,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			errbuf := &bytes.Buffer{}
			flagtarget := target{}

			fs := flag.NewFlagSet("test-"+name, flag.ContinueOnError)
			fs.StringVar(&flagtarget.Foo, "foo", "", "test foo")
			fs.StringVar(&flagtarget.Baz, "baz", "", "test baz")
			fs.StringVar(&flagtarget.Wuz, "wuz-value", "", "test hyphenated env")
			fs.SetOutput(errbuf)

			err := envflag.Parse(fs, tc.args, tc.env)
			if !tc.shoulderr {
				is.NoErr(err)
				is.Equal(errbuf.String(), "")
				is.Equal(flagtarget, tc.result)
			} else {
				is.True(errbuf.String() != "")
			}
		})
	}
}
