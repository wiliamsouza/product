package cmd

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
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

		// Metrics
		metricsExporter, err := prometheus.NewExporter(prometheus.Options{
			Namespace: viper.GetString("metricNamespace"),
		})
		if err != nil {
			log.Fatal(err)
		}
		view.SetReportingPeriod(60 * time.Second)
		view.RegisterExporter(metricsExporter)
		if err = view.Register(ocgrpc.DefaultServerViews...); err != nil {
			log.Fatalf("error registering default server views: %v", err)
		} else {
			log.Print("registered default gRPC server metrics views")
		}

		go func() {
			mux := http.NewServeMux()
			// enable OpenCensus zPages
			zpages.Handle(mux, "/debug")

			// Enable ocsql metrics with OpenCensus
			ocsql.RegisterAllViews()
			mux.Handle("/metrics", metricsExporter)
			if err = http.ListenAndServe(viper.GetString("metricListenAdrress"), mux); err != nil {
				log.Fatalf("failed to run metrics scrape endpoint: %v", err)
			}
		}()

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
		}
		conn, err := grpc.Dial("127.0.0.1:50051", opts...)
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
		rpcListenAddress    string
		dataSourceName      string
		tracerEndpoint      string
		tracerServiceName   string
		metricNamespace     string
		metricListenAdrress string
	)

	grpcCmd.PersistentFlags().StringVarP(&dataSourceName, "dsn", "", "postgres://cataloging:cataloging@127.0.0.1:5432/cataloging?sslmode=disable", "Database data source name")
	grpcCmd.PersistentFlags().StringVarP(&rpcListenAddress, "grpclistenAddress", "", "localhost:13666", "gRPC listen address")
	grpcCmd.PersistentFlags().StringVarP(&tracerEndpoint, "tracerEndpoint", "", "http://localhost:14268", "Tracing exporter endpoint")
	grpcCmd.PersistentFlags().StringVarP(&tracerServiceName, "tracerServiceName", "", "product-grpc", "Tracing exporter service name")
	grpcCmd.PersistentFlags().StringVarP(&metricNamespace, "metricNamespace", "", "product", "Metrics exporter namespace")
	grpcCmd.PersistentFlags().StringVarP(&metricListenAdrress, "metricListenAdrress", "", "localhost:8888", "Metrics listen address")

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
	err = viper.BindPFlag("metricNamespace", grpcCmd.PersistentFlags().Lookup("metricNamespace"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
	err = viper.BindPFlag("metricListenAdrress", grpcCmd.PersistentFlags().Lookup("metricListenAdrress"))
	if err != nil {
		log.Fatalf("failed to bind flag: %v", err)
	}
}
