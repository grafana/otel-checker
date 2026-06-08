package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is the build-time injected version string. Set via:
//
//	go build -ldflags "-X main.version=v0.1.0"
//
// When unset, falls back to whatever Go's runtime/debug.ReadBuildInfo()
// exposes — typically meaningful for `go install` from a tagged release,
// "(devel)" or a pseudo-version for local builds.
var version = ""

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Short:        "Print the otel-checker version",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		Run: func(cc *cobra.Command, _ []string) {
			fmt.Fprintf(cc.OutOrStdout(), "otel-checker %s\n", currentVersion())
		},
	}
}

func currentVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}
