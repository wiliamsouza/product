package cmd

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	version    = "unknown-version"
	commit     = "unknown-commit"
	date       = "unknown-built-time"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Long:  `Version prints the version, as reported by main.Version.`,
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintf(w, "Version:\t%s\n", version)
			fmt.Fprintf(w, "Commit:\t%s\n", commit)
			fmt.Fprintf(w, "Built date:\t%s\n", date)
			fmt.Fprintf(w, "Go version:\t%s\n", runtime.Version())
			fmt.Fprintf(w, "Os:\t%s\n", runtime.GOOS)
			fmt.Fprintf(w, "Arch:\t%s\n", runtime.GOARCH)
			w.Flush()
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
