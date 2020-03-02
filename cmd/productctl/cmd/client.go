package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	grpcv1 "wiliam.dev/product/grpc/v1beta1"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interact with product gRPC server",
	Long:  `List products.`,
	Run: func(cmd *cobra.Command, args []string) {
		address := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

		// Tracing
		exporter, err := jaeger.NewExporter(jaeger.Options{
			Endpoint: viper.GetString("clientTracerEndpoint"),
			Process: jaeger.Process{
				ServiceName: viper.GetString("clientTracerServiceName"),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		trace.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		log.Print("jaeger initialization completed.")

		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		}
		conn, err := grpc.Dial(address, opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := grpcv1.NewProductAPIClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ctx, span := trace.StartSpan(ctx, viper.GetString("clientTracerServiceName"))
		defer span.End()

		r, err := c.ListProducts(ctx, &grpcv1.ListProductsRequest{UserId: "UUID"})
		if err != nil {
			span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
			log.Fatalf("could not list products: %v", err)
		}

		for _, p := range r.Products {
			log.Printf("%s", p.Title)
		}
	},
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
	clientCmd.PersistentFlags().StringVarP(&tracerServiceName, "clientTtracerServiceName", "", "product-grpc-client", "Tracing exporter service name")
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
