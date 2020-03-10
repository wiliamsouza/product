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

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Create product.",
	Long:  `Create product using gRPC endpoint.`,
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

		product, err := c.CreateProduct(
			ctx,
			&grpcv1.CreateProductRequest{
				Title:        viper.GetString("title"),
				Description:  viper.GetString("description"),
				PriceInCents: viper.GetInt32("price"),
			})
		if err != nil {
			span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
			log.Fatalf("could not create product: %v", err)
		}

		log.Print(product)
	},
}

func init() {
	var (
		title       string
		description string
		price       int
	)
	pdtCreateCmd.AddCommand(productCmd)
	pdtCreateCmd.PersistentFlags().StringVarP(&title, "title", "", "", "Product title")
	pdtCreateCmd.PersistentFlags().StringVarP(&description, "description", "", "", "Product description")
	pdtCreateCmd.PersistentFlags().IntVarP(&price, "price", "", 0, "Product price")
	err := viper.BindPFlag("title", pdtCreateCmd.PersistentFlags().Lookup("title"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("description", pdtCreateCmd.PersistentFlags().Lookup("description"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("price", pdtCreateCmd.PersistentFlags().Lookup("price"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
}
