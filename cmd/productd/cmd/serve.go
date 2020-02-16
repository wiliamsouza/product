package cmd

import (
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start http or grpc daemons",
	Long:  `This is used to control daemons start for gRPC and HTTP services.`,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
