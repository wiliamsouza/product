package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands.",
	Long:  `This is used to provide list for all gRPC endpoints.`,
}

func init() {
	clientCmd.AddCommand(listCmd)
}
