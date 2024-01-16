package myapp

import (
	"runtime/debug"
)

// This var is meant to be filled in by
// `go build -ldflags -X=<myapp package>.Version=<Value>`.
var (
	// Version is the git tag of this build (v1.2.3)
	// Go's debug.BuildInfo does not include a git version tag or branch
	// (since a hash can be part of multiple tags/branches)
	Version = "unknown"
)

// BuildInfo contains vcs-related info about when the binary was built.
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// GetBuildInfo extracts commit info from runtime/debug about the binary as to
// from what VCS (git) revision hash and related commit date it was built from.
// It will amend the hash with `-dirty` if the vcs tree had local uncommited
// changes.
func GetBuildInfo() BuildInfo {
	b := BuildInfo{
		Version: Version,
		Commit:  "unknown",
		Date:    "unknown",
	}
	dirty := false
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{}
	}
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			b.Commit = kv.Value
		case "vcs.time":
			b.Date = kv.Value
		case "vcs.modified":
			dirty = kv.Value == "true"
		}
	}
	if dirty {
		b.Commit += "-dirty"
	}
	return b
}
