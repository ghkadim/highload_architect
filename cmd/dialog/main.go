package main

import (
	"context"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/ghkadim/highload_architect/internal/config"
	grpcController "github.com/ghkadim/highload_architect/internal/dialog/controller/grpc"
	"github.com/ghkadim/highload_architect/internal/dialog/controller/openapi"
	"github.com/ghkadim/highload_architect/internal/dialog/counter"
	"github.com/ghkadim/highload_architect/internal/dialog/service/dialog"
	"github.com/ghkadim/highload_architect/internal/dialog/tarantool"
	"github.com/ghkadim/highload_architect/internal/grpc"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
	"github.com/ghkadim/highload_architect/internal/session"
	"github.com/ghkadim/highload_architect/internal/trace"
)

func main() {
	l := logger.Init(config.Get("DEBUG", false))
	defer func() { _ = l.Sync() }()

	logger.Infof("Server starting")
	exporter := trace.NewJaegerExporter()
	defer func() { _ = exporter.Shutdown(context.Background()) }()
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	session_ := session.NewSession(
		config.Get("SESSION_KEY", "secret"),
	)

	counterCl, err := counter.NewClient(config.Get("COUNTER_ADDRESS", ""))
	if err != nil {
		logger.Fatalf("Failed to create counter client: %v", err)
	}

	dialogSvc := dialog.NewService(
		tarantool.NewStorage(
			config.Get("TARANTOOL_ADDRESS", ""),
			http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			},
		),
		counterCl,
	)

	httpServer := server.NewServer(
		openapi.NewRouter(openapi.NewController(dialogSvc, session_)),
	)

	grpcServer := grpc.NewServer(
		grpcController.NewController(dialogSvc),
	)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		httpServer.ListenAndServe(":8080")
	}()
	go func() {
		defer wg.Done()
		grpcServer.ListenAndServe(":8081")
	}()
	wg.Wait()
}
