package cmd

import (
	"context"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	productv1 "wiliam.dev/product/grpc/v1beta1"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start product HTTP server",
	Long:  `Listen on the provide network address.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Tracing
		exporter, err := jaeger.NewExporter(jaeger.Options{
			Endpoint: viper.GetString("tracerEndpoint"),
			Process: jaeger.Process{
				ServiceName: viper.GetString("tracerServiceName"),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		trace.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		}
		conn, err := grpc.Dial(viper.GetString("connectAddress"), opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		mux := runtime.NewServeMux()
		err = productv1.RegisterProductAPIHandler(ctx, mux, conn)
		if err != nil {
			log.Fatal(err)
		}
		nmux := &ochttp.Handler{Handler: mux}
		err = http.ListenAndServe(viper.GetString("listenAddress"), nmux)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serveCmd.AddCommand(httpCmd)

	var (
		connectAddress    string
		listenAddress     string
		tracerEndpoint    string
		tracerServiceName string
	)

	httpCmd.PersistentFlags().StringVarP(&connectAddress, "connectAddress", "", "localhost:13666", "gRPC address to connect")
	httpCmd.PersistentFlags().StringVarP(&listenAddress, "listenAddress", "", "localhost:23666", "HTTP listen address")
	httpCmd.PersistentFlags().StringVarP(&tracerEndpoint, "tracerEndpoint", "", "http://localhost:14268", "Tracing exporter endpoint")
	httpCmd.PersistentFlags().StringVarP(&tracerServiceName, "tracerServiceName", "", "product-http", "Tracing exporter service name")

	err := viper.BindPFlag("connectAddress", httpCmd.PersistentFlags().Lookup("connectAddress"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("listenAddress", httpCmd.PersistentFlags().Lookup("listenAddress"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("tracerEndpoint", httpCmd.PersistentFlags().Lookup("tracerEndpoint"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("tracerServiceName", httpCmd.PersistentFlags().Lookup("tracerServiceName"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
}
