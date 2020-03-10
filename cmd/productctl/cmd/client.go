package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client for product.",
	Long:  `This is used for operation in products endpoints.`,
}

func init() {
	var (
		host              string
		port              int
		tracerEndpoint    string
		tracerServiceName string
	)
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringVarP(&host, "host", "", "localhost", "Host name or IP address")
	clientCmd.PersistentFlags().IntVarP(&port, "port", "", 13666, "Port number to listen")
	clientCmd.PersistentFlags().StringVarP(&tracerEndpoint, "clientTracerEndpoint", "", "http://localhost:14268", "Tracing exporter endpoint")
	clientCmd.PersistentFlags().StringVarP(&tracerServiceName, "clientTracerServiceName", "", "product-grpc-client", "Tracing exporter service name")
	err := viper.BindPFlag("host", clientCmd.PersistentFlags().Lookup("host"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("port", clientCmd.PersistentFlags().Lookup("port"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("clientTracerEndpoint", clientCmd.PersistentFlags().Lookup("clientTracerEndpoint"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("clientTracerServiceName", clientCmd.PersistentFlags().Lookup("clientTracerServiceName"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
}
