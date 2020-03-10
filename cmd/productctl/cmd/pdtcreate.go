package cmd

import (
	"github.com/spf13/cobra"
)

// pdtCreateCmd represents the create command
var pdtCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create commands.",
	Long:  `This is used to provide create for all gRPC endpoints.`,
}

func init() {
	clientCmd.AddCommand(pdtCreateCmd)
}
