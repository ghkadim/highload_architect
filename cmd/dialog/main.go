package main

import (
	"log"
	"net/http"

	"github.com/ghkadim/highload_architect/internal/config"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
	"github.com/ghkadim/highload_architect/internal/server/dialog/openapi"
	"github.com/ghkadim/highload_architect/internal/service/dialog"
	"github.com/ghkadim/highload_architect/internal/session"
	"github.com/ghkadim/highload_architect/internal/tarantool"
)

func main() {
	l := logger.Init(config.Get("DEBUG", false))
	defer func() { _ = l.Sync() }()

	logger.Info("Server starting")

	session_ := session.NewSession("")

	dialogSvc := dialog.NewService(
		tarantool.NewStorage(
			config.Get("TARANTOOL_ADDRESS", ""),
			http.Client{},
		),
	)

	apiController := server.NewServer(
		openapi.NewRouter(openapi.NewController(dialogSvc, session_)),
	)

	log.Fatal(http.ListenAndServe(":8080", apiController))
}
