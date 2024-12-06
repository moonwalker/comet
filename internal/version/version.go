package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	version = "dev"
	commit  = "HEAD"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			commit = kv.Value
		}
	}
}

func Info() string {
	return fmt.Sprintf("%s, build %.7s", version, commit)
}

func InfoEx() string {
	return fmt.Sprintf("%s on %s/%s", Info(), runtime.GOOS, runtime.GOARCH)
}
