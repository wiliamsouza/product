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

// productsCmd represents the products command
var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "List products.",
	Long:  `List products using gRPC endpoint.`,
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

		log.Printf("%d products.\n", len(r.Products))
		for _, p := range r.Products {
			log.Printf("%s", p.Title)
		}
	},
}

func init() {
	listCmd.AddCommand(productsCmd)
}
