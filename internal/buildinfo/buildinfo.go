// Package buildinfo carries version metadata injected at build time via ldflags.
package buildinfo

import "fmt"

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func String() string {
	return fmt.Sprintf("oym-conventions %s (commit %s, built %s)", Version, GitCommit, BuildTime)
}
