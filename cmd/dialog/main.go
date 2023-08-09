package main

import (
	"net/http"
	"sync"

	"github.com/ghkadim/highload_architect/internal/app/service/dialog"
	"github.com/ghkadim/highload_architect/internal/config"
	grpcController "github.com/ghkadim/highload_architect/internal/dialog/controller/grpc"
	"github.com/ghkadim/highload_architect/internal/dialog/controller/openapi"
	"github.com/ghkadim/highload_architect/internal/dialog/tarantool"
	"github.com/ghkadim/highload_architect/internal/grpc"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
	"github.com/ghkadim/highload_architect/internal/session"
)

func main() {
	l := logger.Init(config.Get("DEBUG", false))
	defer func() { _ = l.Sync() }()

	logger.Info("Server starting")

	session_ := session.NewSession(
		config.Get("SESSION_KEY", "secret"),
	)

	dialogSvc := dialog.NewService(
		tarantool.NewStorage(
			config.Get("TARANTOOL_ADDRESS", ""),
			http.Client{},
		),
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
