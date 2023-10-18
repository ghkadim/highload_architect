package main

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/ghkadim/highload_architect/internal/config"
	grpcController "github.com/ghkadim/highload_architect/internal/counter/controller/grpc"
	"github.com/ghkadim/highload_architect/internal/counter/controller/openapi"
	"github.com/ghkadim/highload_architect/internal/counter/service/counter"
	"github.com/ghkadim/highload_architect/internal/counter/tarantool"
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

	service := counter.NewService(
		tarantool.NewStorage(
			config.Get("TARANTOOL_USER", "storage"),
			config.Get("TARANTOOL_PASSWORD", "passw0rd"),
			config.Get("TARANTOOL_ADDRESS", "localhost:3301"),
		),
	)

	httpServer := server.NewServer(
		openapi.NewRouter(openapi.NewController(service, session_)),
	)

	grpcServer := grpc.NewServer(
		grpcController.NewController(service),
	)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		httpServer.ListenAndServe(config.Get("HTTP_LISTEN_ADDRESS", ":8080"))
	}()
	go func() {
		defer wg.Done()
		grpcServer.ListenAndServe(config.Get("GRPC_LISTEN_ADDRESS", ":8081"))
	}()
	wg.Wait()
}
