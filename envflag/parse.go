// Package envflag provides supplemental environment parsing for [flag.FlagSet.Parse]
package envflag

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Parse supplements [flag.FlagSet.Parse] to pull from env variables.
// flag names are mapped to env variable names via
// uppercase with _ separators in place of `[.-/]`
// ie, flag.Name `foo-bar` maps to environment variable `FOO_BAR`
//
// fs must have flag.ContinueOnError set,
// args is expected to be os.Args[1:],
// and environ is expected to be the return of os.Environ().
func Parse(fs *flag.FlagSet, args, environ []string) error {
	if fs.ErrorHandling() != flag.ContinueOnError {
		return errors.New("flag set does not have ContinueOnError")
	}
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parsing arg flags: %w", err)
	}

	// turn os.environ format into a useful map.
	env := map[string]string{}
	for _, entry := range environ {
		k, v, found := strings.Cut(entry, "=")
		if !found {
			return fmt.Errorf("unexpected environment entry %v", entry)
		}
		env[k] = v
	}

	// replacer would be a global const if Go allowed such things.
	// a sync.OnceValue(func replacer() *strings.Replacer {...})
	// would also work, but sync.OnceValue also can't be a const.
	// Given the lack of options to have immutable global state,
	// and the low likelihood this code will be performance load bearing
	// so just inline it and ignore the issue.
	replacer := strings.NewReplacer(
		"-", "_",
		".", "_",
		"/", "_",
	)
	envName := func(s string) string { return replacer.Replace(strings.ToUpper(s)) }

	// args override env, so figure out which flags were manually specified
	hadArg := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		hadArg[envName(f.Name)] = true
	})

	verr := error(nil)
	fs.VisitAll(func(f *flag.Flag) {
		key := envName(f.Name)
		if verr != nil || hadArg[key] {
			// if we had a manual flag or error, skip doing more
			return
		}

		value := env[key]
		if value == "" {
			return
		}

		for _, v := range strings.Split(value, ",") {
			if err := fs.Set(f.Name, v); err != nil {
				verr = fmt.Errorf("setting flag %v from env var %v: %w", f.Name, key, err)
				return
			}
		}
	})
	return verr
}
