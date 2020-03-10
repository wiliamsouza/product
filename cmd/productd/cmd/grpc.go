package cmd

import (
	"database/sql"
	"log"
	"net"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	grpcServer "wiliam.dev/product/grpc"
	promotionv1 "wiliam.dev/product/grpc/client/promotion/v1alpha1"
	productv1 "wiliam.dev/product/grpc/v1beta1"
	"wiliam.dev/product/postgresql"
	"wiliam.dev/product/usecase"
)

// grpcCmd represents the serve command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start product gRPC server",
	Long:  `Listen on the provide network address.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Tracing
		tracingExporter, err := jaeger.NewExporter(jaeger.Options{
			Endpoint: viper.GetString("tracerEndpoint"),
			Process: jaeger.Process{
				ServiceName: viper.GetString("tracerServiceName"),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		trace.RegisterExporter(tracingExporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		l, err := net.Listen("tcp", viper.GetString("grpclistenAddress"))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		driverName, err := ocsql.Register("postgres", ocsql.WithAllTraceOptions())
		if err != nil {
			log.Fatalf("unable to register our ocsql driver: %v\n", err)
		}

		db, err := sql.Open(driverName, viper.GetString("dsn"))
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}

		// enable periodic recording of sql.DBStats
		dbstatsCloser := ocsql.RecordStats(db, 5*time.Second)

		defer func() {
			dbstatsCloser()
			db.Close()
		}()

		dbx := sqlx.NewDb(db, "postgres")

		if err = dbx.Ping(); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		dataStore := postgresql.NewProductDataStore(dbx)

		useCase := usecase.NewProductUseCase(dataStore)

		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		}
		conn, err := grpc.Dial(viper.GetString("promotionConnectAddress"), opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		promotionClient := promotionv1.NewPromotionAPIClient(conn)
		withPromotionProduct := usecase.NewPromotionUseCase(useCase, promotionClient)

		server := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))

		productServer := grpcServer.NewProductAPIServer(withPromotionProduct)

		productv1.RegisterProductAPIServer(server, productServer)

		log.Printf("serving: %s\n", viper.GetString("grpclistenAddress"))
		err = server.Serve(l)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	serveCmd.AddCommand(grpcCmd)

	var (
		rpcListenAddress  string
		dataSourceName    string
		tracerEndpoint    string
		tracerServiceName string
		connectAddress    string
	)

	grpcCmd.PersistentFlags().StringVarP(&rpcListenAddress, "grpclistenAddress", "", "localhost:13666", "gRPC listen address")
	grpcCmd.PersistentFlags().StringVarP(&dataSourceName, "dsn", "", "postgres://postgres:swordfish@127.0.0.1:5432/product?sslmode=disable", "Database data source name")
	grpcCmd.PersistentFlags().StringVarP(&tracerEndpoint, "tracerEndpoint", "", "http://localhost:14268", "Tracing exporter endpoint")
	grpcCmd.PersistentFlags().StringVarP(&tracerServiceName, "tracerServiceName", "", "product-grpc", "Tracing exporter service name")
	httpCmd.PersistentFlags().StringVarP(&connectAddress, "promotionConnectAddress", "", "localhost:13666", "gRPC address to connect")

	err := viper.BindPFlag("grpclistenAddress", grpcCmd.PersistentFlags().Lookup("grpclistenAddress"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("dsn", grpcCmd.PersistentFlags().Lookup("dsn"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("tracerEndpoint", grpcCmd.PersistentFlags().Lookup("tracerEndpoint"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("tracerServiceName", grpcCmd.PersistentFlags().Lookup("tracerServiceName"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("promotionConnectAddress", httpCmd.PersistentFlags().Lookup("promotionConnectAddress"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
}
