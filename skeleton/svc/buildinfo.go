package myapp

import (
	"runtime/debug"
)

// debug.BuildInfo cannot know which git tag is appropriate,
// so this var is meant to be filled in by
// `go build -ldflags -X=<package>.Version=<Value>`.
var (
	// Version is the git tag of this build (v1.2.3)
	Version = "unknown"
)

// BuildInfo contains vcs-related info about when the binary was built.
type BuildInfo struct {
	Version   string
	Commit    string
	Date      string
	GoVersion string
}

// Does not seem to include an equivalent of version tag as of go 1.18 ... :(
func getBuildInfo() BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{}
	}

	b := BuildInfo{
		Version:   Version,
		GoVersion: info.GoVersion,
	}

	dirty := false
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
